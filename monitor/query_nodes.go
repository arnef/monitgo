package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"git.arnef.de/monitgo/node"
)

func GetStatus(nodes []NodeConfig) Data {
	var status Data = make(Data)
	wg := sync.WaitGroup{}

	for i := range nodes {
		wg.Add(1)
		go func(i int) {
			key, value := getNodeStatus(nodes[i])
			status[key] = value
			wg.Done()
		}(i)
	}
	wg.Wait()
	return status
}

func getNodeStatus(nc NodeConfig) (string, Status) {
	url := fmt.Sprintf("http://%s:%d/stats", nc.Host, nc.Port)
	resp, err := http.Get(url)
	if err != nil {
		return nc.Host, NewStatusError(err.Error())
	}
	var status node.JsonStats
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return nc.Host, NewStatusError(err.Error())
	}

	hostDisk := UsageStats{}
	for _, du := range status.Host.Disk {
		hostDisk.TotalBytes += du.TotalBytes
		hostDisk.UsedBytes += du.UsedBytes
	}
	hostDisk.Percentage = float64(hostDisk.UsedBytes) * 100 / float64(hostDisk.TotalBytes)

	mem := UsageStats{
		TotalBytes: status.Host.Memory["Mem"].TotalBytes,
		UsedBytes:  status.Host.Memory["Mem"].UsedBytes,
	}
	mem.Percentage = float64(mem.UsedBytes) * 100 / float64(mem.TotalBytes)

	hostStats := HostStats{
		CPU:    status.Host.CPU,
		Disk:   hostDisk,
		Memory: mem,
	}

	containerStats := make(map[string]ContainerStats)
	for id, c := range status.Container {
		var totalTxBytes uint64
		var totalRxBytes uint64

		for _, net := range c.Network {
			totalRxBytes += net.TotalRxBytes
			totalTxBytes += net.TotalTxBytes
		}
		containerStats[id] = ContainerStats{
			Name: c.Name,
			CPU:  c.CPU,
			Memory: UsageStats{
				TotalBytes: c.Memory.TotalBytes,
				UsedBytes:  c.Memory.UsedBytes,
				Percentage: float64(c.Memory.UsedBytes) * 100 / float64(c.Memory.TotalBytes),
			},
		}
	}
	return nc.Host, Status{
		Name:      nc.Name,
		Host:      hostStats,
		Container: containerStats,
	}
}
