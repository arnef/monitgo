package node

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type RemoteNode struct {
	Name     string
	Endpoint string
}

type Node interface {
	GetCPUs() (int, error)
}

func (r *RemoteNode) GetCPUs() (int, error) {
	lscpu, err := r.exec("lscpu --json")
	fmt.Println(lscpu, err)
	return 0, errors.New("not implemented")
}

func (r *RemoteNode) exec(cmd string) (string, error) {
	resp, err := http.Post(r.Endpoint+"/xc", "text/plain; charset=utf-8", bytes.NewReader([]byte(cmd)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	out, err := ioutil.ReadAll(resp.Body)
	return string(out), err
}
