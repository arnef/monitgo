package monitor

import (
	"fmt"
	"time"
)

func Init(nodes []Node, sleep uint64) {
	if monit == nil {
		monit = &monitor{
			nodes: nodes,
			sleep: sleep,
		}
	} else {
		panic("monitor service already initialized")
	}
}

func Register(d DataReceiver) {
	if monit != nil {
		monit.subscriber = append(monit.subscriber, d)
	}
}

func Start() error {
	if monit == nil {
		return fmt.Errorf("monitor service not initialized")
	}
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
	nodes      []Node
	sleep      uint64
}

func notifyAll(data Data) {
	for _, subscriber := range monit.subscriber {
		subscriber.Push(data)
	}
}
