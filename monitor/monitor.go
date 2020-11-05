package monitor

import (
	"fmt"
	"time"

	"git.arnef.de/monitgo/config"
)

func Init(nodes []config.Node, sleep uint64) {
	if monit == nil {
		monit = &monitor{
			nodes: nodes,
			sleep: sleep,
		}
	} else {
		panic("Monit already initialized")
	}
}

func Register(d DataReceiver) {
	if monit != nil {
		monit.subscriber = append(monit.subscriber, d)
	}
}

func Start() error {
	fmt.Printf("👀️ getting data every %d seconds\n", monit.sleep)
	query()
	for range time.Tick(time.Duration(monit.sleep) * time.Second) {
		query()
	}
	return nil
}

func query() {
	stats := GetStatus(monit.nodes)
	notifyAll(stats)
}

var (
	monit *monitor
)

type monitor struct {
	subscriber []DataReceiver
	nodes      []config.Node
	sleep      uint64
}

func notifyAll(data Data) {
	for _, subscriber := range monit.subscriber {
		subscriber.Push(data)
	}
}
