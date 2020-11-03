package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Stats docker
type Stats struct {
	ID       string
	Name     string
	CPU      string
	MemUsage string
	NetIO    string
	BlockIO  string
	PIDs     string
}

func GetStats() ([]Stats, error) {

	out, err := docker("stats", "--no-stream", "--all", "--format",
		"{\"ID\": \"{{ .Container }}\", \"Name\": \"{{ .Name }}\", \"CPU\": \"{{ .CPUPerc }}\", \"MemUsage\": \"{{ .MemUsage }}\", \"NetIO\": \"{{ .NetIO }}\", \"BlockIO\": \"{{ .BlockIO }}\", \"PIDs\": \"{{ .PIDs }}\" }")

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	stats := make([]Stats, len(lines)-1)
	for i := range stats {
		err := json.Unmarshal([]byte(lines[i]), &stats[i])
		if err != nil {
			return nil, err
		}
	}
	return stats, nil
}

func docker(args ...string) ([]byte, error) {
	docker, err := exec.LookPath("docker")
	if err != nil {

		return nil, err
	}

	var outb, errb bytes.Buffer
	cmd := &exec.Cmd{
		Path:   docker,
		Stderr: &errb,
		Stdout: &outb,
		Args:   append([]string{docker}, args...),
	}
	err = cmd.Run()
	if err != nil {

		return nil, fmt.Errorf(errb.String())
	}

	return outb.Bytes(), err
}
