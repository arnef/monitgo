package docker

import (
	"fmt"
	"time"
)

// Stats docker
type Stats struct {
	ID   string
	Name string
	// CPU percentage
	CPU float64
	// MemUsage in bytes
	MemUsage uint64
	// NetRx in bytes
	NetRx uint64
	// NetTx in bytes
	NetTx uint64
}

var (
	diff      map[string]Stats
	timestamp time.Time
)

func GetStats() ([]Stats, error) {
	runningTime := time.Now()
	if diff == nil {
		query, err := getDockerStats()
		if err != nil {
			return nil, err
		}
		diff = query
		timestamp = runningTime
		return nil, fmt.Errorf("Not initialized")
	} else {
		current, err := getDockerStats()
		if err != nil {
			return nil, err
		}
		var diffStats []Stats
		for i, stat := range current {
			if _, ok := diff[i]; ok {
				duration := runningTime.Sub(timestamp).Seconds()
				diffStats = append(diffStats, Stats{
					ID:       stat.ID,
					Name:     stat.Name,
					CPU:      stat.CPU,
					MemUsage: stat.MemUsage,
					NetRx:    (current[i].NetRx - diff[i].NetRx) / uint64(duration),
					NetTx:    (current[i].NetTx - diff[i].NetTx) / uint64(duration),
				})
			}
		}

		diff = current
		timestamp = runningTime
		return diffStats, nil
	}
}
