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
}

func GetStatus() []Status {
	config := config.Get()
	stati := make([]Status, len(config.Nodes))
	for i, node := range config.Nodes {
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
						if row.PIDs == "0" {
							stati[i].Data = append(stati[i].Data, row)
						}
					}
				}
			}
		}
	}
	return stati
}
