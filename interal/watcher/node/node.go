package node

import (
	"bytes"
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

var clients = map[string]*client.Client{}

func (n *Node) Validate() error {
	if n.Port == 0 {
		n.Port = 5000
	}
	if len(n.Name) == 0 {
		n.Name = fmt.Sprintf("%s:%d", n.Host, n.Port)
	}
	if n.CPUs == 0 {
		out, err := n.Exec("lscpu")
		if err != nil {
			log.Error(err)
			return fmt.Errorf("could not detect cpus on %s. please add it to your config", n.Name)
		}
		for _, line := range strings.Split(string(out), "\n") {
			if strings.HasPrefix(line, "CPU(s):") {
				valStr := strings.TrimSpace(strings.Replace(line, "CPU(s):", "", 1))
				val, err := strconv.Atoi(valStr)

				if err != nil {
					log.Error(err)
					return err
				}
				n.CPUs = val
				break
			}
		}

		if n.CPUs == 0 {
			return fmt.Errorf("could not detect cpus on %s. please add it to your config", n.Name)
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

func CurrentClient(n *Node) (*client.Client, error) {
	if c, exists := clients[n.Name]; exists {
		return c, nil
	}

	c, err := n.DockerClient()
	if err != nil {
		return nil, err
	}

	clients[n.Name] = c

	return c, nil
}
