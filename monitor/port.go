package monitor

import (
	"git.arnef.de/monitgo/node/docker"
	"git.arnef.de/monitgo/node/host"
)

type DataReceiver interface {
	Push(data Data)
}

type Status struct {
	Name      string
	Error     error
	Container []docker.Stats
	Host      host.HostStats
	key       string
}

type Node struct {
	Name string
	Host string
	Port uint
}

type Data map[string]Status
