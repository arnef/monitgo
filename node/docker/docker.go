package docker

import (
	"time"
)

// Stats docker
type Stats struct {
	ID   string
	Name string
	// CPU percentage
	CPU float64
	// MemUsage in bytes
	Memory  MemoryStats
	Network map[string]NetworkStats
}

type MemoryStats struct {
	UsedBytes  uint64
	TotalBytes uint64
}

type NetworkStats struct {
	TotalRxBytes uint64
	TotalTxBytes uint64
}

var (
	diff      map[string]Stats
	timestamp time.Time
)
