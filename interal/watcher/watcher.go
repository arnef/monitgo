package watcher

import (
	"fmt"

	"github.com/arnef/monitgo/interal/watcher/config"
	"github.com/arnef/monitgo/interal/watcher/node"
	log "github.com/sirupsen/logrus"
)

func Start(configPath string, interval int) error {
	cfg, err := config.FromPath(configPath)
	if err != nil {
		return err
	}
	log.Debug(cfg)
	watcher := watcher{}
	for _, n := range cfg.Nodes {
		nodePort := n.Port
		if nodePort == 0 {
			nodePort = 5000
		}
		watcher.nodes = append(watcher.nodes, &node.RemoteNode{
			Name:     n.Name,
			Endpoint: fmt.Sprintf("http://%s:%d", n.Host, nodePort),
		})
	}
	return watcher.Run()
}

type watcher struct {
	nodes []node.Node
}

func (w *watcher) Run() error {

	for _, n := range w.nodes {
		fmt.Println(n.GetCPUs())
	}
	return nil
}
