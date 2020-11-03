package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"git.arnef.de/monitgo/config"
	"git.arnef.de/monitgo/docker"
)

type Status struct {
	Name  string
	Error string
	Data  []docker.Stats
	host  string
}

func getNodeStatus(node config.Node) Status {
	url := fmt.Sprintf("http://%s:%d/stats", node.Host, node.Port)
	resp, err := http.Get(url)
	status := Status{
		host: node.Host,
		Name: node.Name,
	}
	if err != nil {
		status.Error = err.Error()
	} else {
		defer resp.Body.Close()
		var stats = struct {
			Data  []docker.Stats
			Error *string
		}{}
		err = json.NewDecoder(resp.Body).Decode(&stats)
		if err != nil {
			status.Error = err.Error()
		} else {
			if stats.Error != nil {
				status.Error = *stats.Error
			} else {
				for _, row := range stats.Data {
					if row.MemUsage == 0 {
						status.Data = append(status.Data, row)
					}
				}
			}
		}
	}

	return status
}

func GetStatus(nodes []config.Node) map[string]Status {
	stati := make([]Status, len(nodes))
	wg := sync.WaitGroup{}

	for i := range nodes {
		wg.Add(1)
		go func(i int) {
			stati[i] = getNodeStatus(nodes[i])
			wg.Done()
		}(i)
	}
	wg.Wait()
	result := make(map[string]Status)
	for _, s := range stati {
		result[s.host] = s
	}
	return result
}
