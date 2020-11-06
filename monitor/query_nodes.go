package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func GetStatus(nodes []Node) map[string]Status {
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
		result[s.key] = s
	}
	return result
}

func getNodeStatus(node Node) Status {
	url := fmt.Sprintf("http://%s:%d/stats", node.Host, node.Port)
	resp, err := http.Get(url)
	if err != nil {
		return Status{
			key:   node.Host,
			Name:  node.Name,
			Error: err,
		}
	}
	var status Status
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		return Status{
			key:   node.Host,
			Name:  node.Name,
			Error: err,
		}
	}

	status.key = node.Host
	status.Name = node.Name
	return status
	// status := Status{
	// 	host: node.Host,
	// 	Name: node.Name,
	// }
	// var errVal string
	// if err != nil {
	// 	errVal = err.Error()
	// } else {
	// 	defer resp.Body.Close()
	// 	var stats = struct {
	// 		Data  []docker.Stats
	// 		Error *string
	// 	}{}
	// 	err = json.NewDecoder(resp.Body).Decode(&stats)
	// 	if err != nil {
	// 		errVal = err.Error()
	// 	} else {
	// 		if stats.Error != nil {
	// 			errVal = *stats.Error
	// 		} else {
	// 			for _, row := range stats.Data {
	// 				// if row.MemUsage == 0 {
	// 				status.Data = append(status.Data, row)
	// 				// }
	// 			}
	// 		}
	// 	}
	// }
	// if errVal != "" {
	// 	status.Error = &errVal
	// }

	// return status
}
