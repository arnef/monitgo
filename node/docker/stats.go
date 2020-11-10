package docker

import (
	"context"
	"encoding/json"
	"sync"

	"git.arnef.de/monitgo/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func GetStats() (map[string]Stats, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return nil, err
	}

	statsMap := make(map[string]Stats)
	wg := sync.WaitGroup{}
	var statsError error
	for i := range containerList {
		wg.Add(1)
		go func(i int) {
			container := containerList[i]
			resp, err := cli.ContainerStats(ctx, container.ID, false)
			if err != nil {
				statsError = err
				// panic(err)
				wg.Done()
				return
				// return nil, err
			}
			defer resp.Body.Close()

			var stats types.StatsJSON

			err = json.NewDecoder(resp.Body).Decode(&stats)

			if err != nil {
				statsError = err
				// panic(err)
				wg.Done()
				return
				// return nil, err
			}

			network := make(map[string]NetworkStats)
			for name, net := range stats.Networks {
				network[name] = NetworkStats{
					TotalRxBytes: net.RxBytes,
					TotalTxBytes: net.TxBytes,
				}
			}
			cpu := calculateCPUPercentUnix(stats.PreCPUStats.CPUUsage.TotalUsage, stats.PreCPUStats.SystemUsage, &stats)
			id := container.ID[:12]
			statsMap[id] = Stats{
				ID:   id,
				Name: stats.Name[1:],
				CPU:  utils.Round(cpu),
				Memory: MemoryStats{
					TotalBytes: stats.MemoryStats.MaxUsage,
					UsedBytes:  stats.MemoryStats.Usage - stats.MemoryStats.Stats["cache"],
				},
				Network: network,
			}
			wg.Done()
		}(i)

	}
	wg.Wait()

	return statsMap, statsError
}
