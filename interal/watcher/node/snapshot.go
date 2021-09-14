package node

import (
	"time"

	"github.com/arnef/monitgo/interal/watcher/parser"
	"github.com/arnef/monitgo/pkg"
)

func (n *Node) Snapshot() pkg.NodeSnapshot {
	snapshot := pkg.NodeSnapshot{}
	snapshot.Name = n.Name
	snapshot.Timestamp = time.Now()

	n.uptime(&snapshot)
	n.memoryUsage(&snapshot)
	n.diskUsage(&snapshot)
	if !n.NoDocker && snapshot.Error == nil {
		n.container(&snapshot)
	}
	return snapshot
}

func (n *Node) uptime(snapshot *pkg.NodeSnapshot) {
	out, err := n.Exec("uptime")
	if err != nil {
		snapshot.Error = err
		return
	}
	load, err := parser.LoadAverage(string(out))
	if err != nil {
		snapshot.Error = err
		return
	}

	snapshot.CPU = load[1] / float64(n.CPUs)
}

func (n *Node) memoryUsage(snapshot *pkg.NodeSnapshot) {
	out, err := n.Exec("free", "--bytes")
	if err != nil {
		snapshot.Error = err
		return
	}
	total, used, err := parser.Free(string(out))
	if err != nil {
		snapshot.Error = err
		return
	}

	snapshot.MemoryUsage = pkg.Usage{TotalBytes: total, UsedBytes: used}
}

func (n *Node) diskUsage(snapshot *pkg.NodeSnapshot) {
	out, err := n.Exec("df", "--output=source,size,used")
	if err != nil {
		snapshot.Error = err
		return
	}
	total, used, err := parser.Df(string(out))
	if err != nil {
		snapshot.Error = err
		return
	}

	snapshot.DiskUsage = pkg.Usage{TotalBytes: total, UsedBytes: used}
}
