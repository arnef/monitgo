package monitor

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func GetStatus(nodes []Node) Data {
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

func getNodeStatus(node Node) (string, Status) {
	url := fmt.Sprintf("http://%s:%d/stats", node.Host, node.Port)
	resp, err := http.Get(url)
	var errorVal string
	if err != nil {
		errorVal = err.Error()
		return node.Host, Status{
			Name:  node.Name,
			Error: &errorVal,
		}
	}
	var status Status
	err = json.NewDecoder(resp.Body).Decode(&status)
	if err != nil {
		errorVal = err.Error()
		return node.Host, Status{
			Name:  node.Name,
			Error: &errorVal,
		}
	}
	status.Name = node.Name
	return node.Host, status
}
