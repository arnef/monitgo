package monitor

import "git.arnef.de/monitgo/node/docker"

type DataReceiver interface {
	Push(data Data)
}

type Status struct {
	Name  string
	Error string
	Data  []docker.Stats
	host  string
}

type Data map[string]Status
