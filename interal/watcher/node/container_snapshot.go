package node

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	"github.com/arnef/monitgo/pkg"
	"github.com/docker/docker/api/types"
)

const MAX_ROUTINES = 10

func (n *Node) container(snapshot *pkg.NodeSnapshot) {
	sem := make(chan int, MAX_ROUTINES)

	client, err := n.DockerClient()
	if err != nil {
		snapshot.Error = err
		return
	}
	defer client.Close()
	ctx := context.Background()
	containerList, err := client.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		snapshot.Error = err
		return
	}

	snapshot.Container = make([]*pkg.ContainerSnapshot, len(containerList))
	wg := sync.WaitGroup{}
	for i := range containerList {
		// Blocks if MAX_ROUTINES reached
		sem <- 1
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			container := containerList[i]
			if _, ignore := container.Labels["monitgo.ignore"]; !ignore {
				cs := pkg.ContainerSnapshot{}
				cs.Timestamp = snapshot.Timestamp
				cs.ID = container.ID
				cs.Name = strings.TrimPrefix(strings.Join(container.Names, ","), "/")

				resp, err := client.ContainerStats(ctx, container.ID, false)
				if err != nil {
					cs.Error = err
				} else {
					var stats types.StatsJSON
					err = json.NewDecoder(resp.Body).Decode(&stats)
					if err != nil {
						cs.Error = err
					} else {
						cs.MemoryUsage = pkg.Usage{
							TotalBytes: stats.MemoryStats.MaxUsage,
							UsedBytes:  stats.MemoryStats.Usage - stats.MemoryStats.Stats["cache"],
						}
						cs.CPU = calculateCPUPercentUnix(stats.PreCPUStats.CPUUsage.TotalUsage, stats.PreCPUStats.SystemUsage, &stats)

						var rxBytes uint64
						var txBytes uint64

						for _, net := range stats.Networks {
							rxBytes += net.RxBytes
							txBytes += net.TxBytes
						}
						cs.Network = pkg.Network{
							TotalRxBytes: rxBytes,
							TotalTxBytes: txBytes,
						}
						cs.State = pkg.ContainerStateType(container.State)
					}
				}

				snapshot.Container[i] = &cs
			}
			// makes place for new action
			<-sem
		}(i)
	}
	wg.Wait()

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
