package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"git.arnef.de/monitgo/log"
	"git.arnef.de/monitgo/utils"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

func GetStats() (map[string]Stats, error) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	// defer cli.Close()
	// defer ctx.Done()
	defer func() { fmt.Println("GetStats Done") }()

	containerList, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	cli.Close()
	ctx.Done()
	statsList := make([]*Stats, len(containerList))
	wg := sync.WaitGroup{}
	// wg := goccm.New(10)
	var statsError error
	for i := range containerList {
		wg.Add(1)
		// wg.Wait()
		go func(i int) {
			defer wg.Done()
			if _, ignore := containerList[i].Labels["monitgo.ignore"]; !ignore {
				ctx := context.Background()
				cli, err := client.NewEnvClient()
				if err != nil {
					log.Debug(err)
					statsError = err
					return
					// return nil, err
				}
				defer cli.Close()
				defer ctx.Done()
				container := containerList[i]
				log.Debug("get container stats %v", container.Names)
				resp, err := cli.ContainerStats(ctx, container.ID, false)
				if err != nil {
					statsError = err
					log.Debug(err)
					return
				}
				defer resp.Body.Close()

				var stats types.StatsJSON

				err = json.NewDecoder(resp.Body).Decode(&stats)

				if err != nil {
					statsError = err
					log.Debug(err)
					return
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
				statsList[i] = &Stats{
					ID:   id,
					Name: stats.Name[1:],
					CPU:  utils.Round(cpu),
					Memory: MemoryStats{
						TotalBytes: stats.MemoryStats.MaxUsage,
						UsedBytes:  stats.MemoryStats.Usage - stats.MemoryStats.Stats["cache"],
					},
					Network: network,
				}
			}
		}(i)

	}
	// wg.WaitAllDone()
	wg.Wait()
	statsMap := make(map[string]Stats)
	for _, s := range statsList {
		if s != nil {
			statsMap[s.Name] = *s
		}
	}
	if statsError != nil {
		log.Debug(statsError)
	}
	return statsMap, statsError
}
