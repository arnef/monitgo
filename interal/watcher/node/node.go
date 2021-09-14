package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

type Node struct {
	Name             string
	Port             int
	Host             string
	DockerAPIVersion string
	CPUs             int
	NoDocker         bool
}

func (n *Node) Validate() error {
	if n.Port == 0 {
		n.Port = 5000
	}
	if len(n.Name) == 0 {
		n.Name = fmt.Sprintf("%s:%d", n.Host, n.Port)
	}
	if n.CPUs == 0 {
		out, err := n.Exec("lscpu", "--json")
		if err != nil {
			return err
		}
		var data map[string][]map[string]string
		err = json.Unmarshal(out, &data)
		if err != nil {
			return err
		}

		for _, line := range data["lscpu"] {
			if line["field"] == "CPU(s):" {
				val, err := strconv.Atoi(line["data"])
				if err != nil {
					log.Debug(err)
					return err
				}
				n.CPUs = val
				break
			}
		}
	}
	return nil
}

func (n *Node) Exec(command string, args ...string) ([]byte, error) {
	url := fmt.Sprintf("http://%s:%d/xc", n.Host, n.Port)

	resp, err := http.Post(url, "text/plain; charset=utf-8", bytes.NewReader(
		[]byte(strings.Join(append([]string{command}, args...), " ")),
	))

	if err != nil {
		return nil, err
	}

	out := bytes.Buffer{}
	defer resp.Body.Close()
	_, err = io.Copy(&out, resp.Body)
	return out.Bytes(), err
}

func (n *Node) DockerClient() (*client.Client, error) {
	// var err error
	// if n.client == nil {
	// 	n.client, err = client.NewClient(fmt.Sprintf("http://%s:%d", n.Host, n.Port), n.DockerAPIVersion, nil, nil)
	// }
	// return n.client, err
	return client.NewClient(fmt.Sprintf("http://%s:%d", n.Host, n.Port), n.DockerAPIVersion, nil, nil)
}
