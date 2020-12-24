package host

import (
	"git.arnef.de/monitgo/log"
)

type Stats struct {
	CPU    float64
	Memory map[string]Usage
	Disk   map[string]Usage
}

type Usage struct {
	TotalBytes uint64
	UsedBytes  uint64
}

func GetStats() (*Stats, error) {

	cpuLoad, err := getNormalizedLoad()
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	memUsage, err := getMemUsage()
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	diskUsage, err := getDiskUsage()
	if err != nil {
		log.Debug(err)
		return nil, err
	}

	stats := Stats{
		CPU:    cpuLoad,
		Memory: memUsage,
		Disk:   diskUsage,
	}

	return &stats, nil
}
