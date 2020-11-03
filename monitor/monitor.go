package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.arnef.de/monitgo/config"
	"git.arnef.de/monitgo/docker"
)

type Status struct {
	Name  string
	Error string
	Data  []docker.Stats
	host  string
}

func GetStatus(nodes []config.Node) map[string]Status {
	stati := make([]Status, len(nodes))
	for i, node := range nodes {
		stati[i].host = node.Host
		stati[i].Name = node.Name
		url := fmt.Sprintf("http://%s:%d/stats", node.Host, node.Port)
		resp, err := http.Get(url)
		if err != nil {
			stati[i].Error = err.Error()
		} else {
			defer resp.Body.Close()
			var stats = struct {
				Data  []docker.Stats
				Error *string
			}{}
			err = json.NewDecoder(resp.Body).Decode(&stats)
			if err != nil {
				stati[i].Error = err.Error()
			} else {
				if stats.Error != nil {
					stati[i].Error = *stats.Error
				} else {
					for _, row := range stats.Data {
						if row.MemUsage == 0 {
							stati[i].Data = append(stati[i].Data, row)
						}
					}
				}
			}
		}
	}
	result := make(map[string]Status)
	for _, s := range stati {
		result[s.host] = s
	}
	return result
}
