package watcher

import (
	"sync"
	"time"

	"github.com/arnef/monitgo/interal/watcher/alerts"
	"github.com/arnef/monitgo/interal/watcher/bot"
	"github.com/arnef/monitgo/interal/watcher/config"
	"github.com/arnef/monitgo/interal/watcher/database"
	"github.com/arnef/monitgo/interal/watcher/node"
	"github.com/arnef/monitgo/pkg"
	log "github.com/sirupsen/logrus"
)

func Start(configPath string, interval int) error {
	cfg, err := config.FromPath(configPath)
	if err != nil {
		return err
	}
	log.Debug(cfg)
	watcher := watcher{
		nodes: cfg.Nodes,
	}
	for i := range watcher.nodes {
		if err := watcher.nodes[i].Validate(); err != nil {
			return err
		}
	}
	log.Debug(watcher.nodes)

	alertManager := alerts.NewManager()
	watcher.registerSnapshotHandler(alertManager.HandleSnaphsot)

	botManager := bot.NewManager()
	alertManager.RegisterAlertHandler(botManager.HandleAlerts)

	if cfg.Matrix != nil {
		botManager.RegisterBot(bot.NewMatrixBot(cfg.Matrix))
	}
	if cfg.Talk != nil {
		botManager.RegisterBot(bot.NewTalkBot(cfg.Talk))
	}
	if cfg.Telegram != nil {
		botManager.RegisterBot(bot.NewTelegramBot(cfg.Telegram))
	}
	go botManager.Listen()

	if cfg.InfluxDB != nil {
		watcher.registerSnapshotHandler(database.NewInfluxDB(cfg.InfluxDB).OnSnapshot)
	}

	return watcher.run(interval)
}

type watcher struct {
	nodes           []node.Node
	snapshotHanlder []pkg.SnapshotHandler
}

func (w *watcher) registerSnapshotHandler(handler pkg.SnapshotHandler) {
	w.snapshotHanlder = append(w.snapshotHanlder, handler)
}

func (w *watcher) notifySnapshotHandler(snapshot []pkg.NodeSnapshot) {
	for _, handler := range w.snapshotHanlder {
		handler(snapshot)
	}
}

func (w *watcher) run(interval int) error {
	if err := w.doWork(); err != nil {
		return err
	}

	for range time.Tick(time.Duration(interval) * time.Second) {
		if err := w.doWork(); err != nil {
			return err
		}
	}

	return nil
}

func (w *watcher) doWork() error {
	current := make([]pkg.NodeSnapshot, len(w.nodes))
	wg := sync.WaitGroup{}
	for i := range w.nodes {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			current[i] = w.nodes[i].Snapshot()
		}(i)
	}
	wg.Wait()
	go w.notifySnapshotHandler(current)
	return nil
}
