package watch

import (
	"encoding/json"
	"fmt"
	"time"

	"git.arnef.de/monitgo/bot"
	"git.arnef.de/monitgo/monitor"
	"github.com/urfave/cli/v2"
)

func Cmd(ctx *cli.Context) error {
	sleep := ctx.Uint64("interval")
	watcher := new(sleep)
	watcher.start()
	return nil
}

type watcher struct {
	sleep        uint64
	bot          bot.Bot
	lastResponse string
}

func new(sleep uint64) watcher {
	return watcher{
		sleep:        sleep,
		lastResponse: "",
		bot:          bot.New(),
	}
}

func (w *watcher) start() {
	w.run()
	for range time.Tick(time.Duration(w.sleep) * time.Second) {
		go w.run()
	}
}

func (w *watcher) run() {
	stats := monitor.GetStatus()
	resp, err := json.Marshal(&stats)
	if err != nil {
		panic(err)
	}
	response := string(resp)
	if response != w.lastResponse {
		message := ""
		for _, s := range stats {
			// something is wrong lets fire a telegram message
			if s.Error != "" {
				message += fmt.Sprintf("❗️ *%s*\n_%s_", s.Name, s.Error)
			} else if len(s.Data) > 0 {
				message += fmt.Sprintf("🔥️ *%s*\n", s.Name)
				for _, d := range s.Data {
					message += fmt.Sprintf("_%s_ down\n", d.Name)
				}
			}
		}
		if message != "" {
			w.bot.Broadcast(message)
		}
	} else {
		fmt.Println("No changes")
	}
	w.lastResponse = response
}
