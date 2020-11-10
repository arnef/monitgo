package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"git.arnef.de/monitgo/node"
	"git.arnef.de/monitgo/node/docker"
	"git.arnef.de/monitgo/node/host"
)

var (
	prev      rawData
	timestamp time.Time
)

type rawData map[string]node.JsonStats

func GetStatus(nodes []NodeConfig) Data {
	var status Data = make(Data)
	now := time.Now()
	wg := sync.WaitGroup{}
	raw := make(rawData)
	for i := range nodes {
		wg.Add(1)
		go func(i int) {
			rawNode, err := queryNode(nodes[i])
			if err != nil {
				status[nodes[i].Host] = NewStatusError(err.Error())
			} else {
				raw[nodes[i].Host] = *rawNode
				status[nodes[i].Host] = processNode(nodes[i].Host, nodes[i].Name, *rawNode, now)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()

	initQuery := prev == nil

	prev = raw

	timestamp = now

	if initQuery {
		return nil
	}

	return status
}

func processNode(host string, name string, data node.JsonStats, now time.Time) Status {
	status := Status{Name: name, Error: data.Error}
	if data.Error != nil {
		return status
	}
	status.Host = HostStats{
		CPU: data.Host.CPU,
		Memory: UsageStats{
			TotalBytes: data.Host.Memory["Mem"].TotalBytes,
			UsedBytes:  data.Host.Memory["Mem"].UsedBytes,
			Percentage: float64(data.Host.Memory["Mem"].UsedBytes) * 100 / float64(data.Host.Memory["Mem"].TotalBytes),
		},
		Disk: calculateDiskUsage(data.Host.Disk),
	}
	status.Container = make(map[string]ContainerStats)

	for key, con := range data.Container {
		var rx, tx float64
		currentTotalRx, currentTotalTx := calculateNetwork(con.Network)

		if n, ok := prev[host]; ok {

			if p, ok := n.Container[key]; ok && p.Memory.UsedBytes > 0 && con.Memory.UsedBytes > 0 {
				duration := time.Since(timestamp).Seconds()
				prevTotalRx, prevTotalTx := calculateNetwork(p.Network)
				rx = float64(currentTotalRx-prevTotalRx) / duration
				tx = float64(currentTotalTx-prevTotalTx) / duration
			}
		}

		status.Container[key] = ContainerStats{
			Name: con.Name,
			CPU:  con.CPU,
			Memory: UsageStats{
				TotalBytes: con.Memory.TotalBytes,
				UsedBytes:  con.Memory.UsedBytes,
				Percentage: calculatePercentage(con.Memory.UsedBytes, con.Memory.TotalBytes),
			},
			Network: NetworkStats{
				RxBytesPerSecond: uint64(rx),
				TxBytesPerSecond: uint64(tx),
			},
		}
	}
	return status
}

func calculateDiskUsage(disks map[string]host.Usage) UsageStats {
	usage := UsageStats{}

	for _, disk := range disks {
		usage.UsedBytes += disk.UsedBytes
		usage.TotalBytes += disk.TotalBytes
	}
	usage.Percentage = calculatePercentage(usage.UsedBytes, usage.TotalBytes)
	return usage
}

func calculatePercentage(used uint64, total uint64) float64 {
	return float64(used) * 100 / float64(total)
}

func queryNode(nc NodeConfig) (*node.JsonStats, error) {
	url := fmt.Sprintf("http://%s:%d/stats", nc.Host, nc.Port)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	var status node.JsonStats

	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return nil, err
	}

	return &status, nil
}

func calculateNetwork(networks map[string]docker.NetworkStats) (uint64, uint64) {
	var totalTxBytes uint64
	var totalRxBytes uint64

	for _, net := range networks {
		totalRxBytes += net.TotalRxBytes
		totalTxBytes += net.TotalTxBytes
	}

	return totalRxBytes, totalTxBytes
}
