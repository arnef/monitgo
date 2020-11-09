package docker

import "github.com/docker/docker/api/types"

func calculateNetwork(network map[string]types.NetworkStats) (uint64, uint64) {
	var rx, tx uint64

	for _, net := range network {
		rx += net.RxBytes
		tx += net.TxBytes
	}

	return rx, tx
}

func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) float64 {
	var (
		cpuPercent = 0.0
		// calculate the change for the cpu usage of the container in between readings
		cpuDelta = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		// calculate the change for the entire system between readings
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
		onlineCPUs  = float64(len(v.CPUStats.CPUUsage.PercpuUsage))
	)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * onlineCPUs * 100.0
	}
	return cpuPercent
}
