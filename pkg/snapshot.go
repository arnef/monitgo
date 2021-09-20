package pkg

import "time"

type ContainerStateType string

const (
	ContainerStateCreated    ContainerStateType = "created"
	ContainerStateRunning    ContainerStateType = "running"
	ContainerStatePaused     ContainerStateType = "paused"
	ContainerStateRestarting ContainerStateType = "restarting"
	ContainerStateRemoving   ContainerStateType = "removing"
	ContainerStateExisted    ContainerStateType = "exited"
	ContainerStateDead       ContainerStateType = "dead"
)

type SnapshotHandler = func(snap []NodeSnapshot)

type Snaphot struct {
	Name        string
	Error       error
	Timestamp   time.Time
	CPU         float64
	MemoryUsage Usage
}

type Usage struct {
	TotalBytes uint64
	UsedBytes  uint64
}

type Network struct {
	TotalRxBytes uint64
	TotalTxBytes uint64
}

func (u *Usage) Percentage() float64 {
	if u.TotalBytes == 0 {
		return -1
	}
	return float64(u.UsedBytes) * 100 / float64(u.TotalBytes)
}

type NodeSnapshot struct {
	DiskUsage Usage
	Container []*ContainerSnapshot
	Snaphot
}

type ContainerSnapshot struct {
	ID      string
	Network Network
	Snaphot
	State ContainerStateType
}
