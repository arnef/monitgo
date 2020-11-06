package docker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"git.arnef.de/monitgo/utils"
)

// Stats docker
type Stats struct {
	ID       string
	Name     string
	CPU      float64
	MemUsage float64
	NetIn    float64
	NetOut   float64
	BlockIO  float64
}

func GetStats() ([]Stats, error) {
	out, err := docker("stats", "--no-stream", "--all", "--format",
		"{\"ID\": \"{{ .Container }}\", \"Name\": \"{{ .Name }}\", \"CPU\": \"{{ .CPUPerc }}\", \"MemUsage\": \"{{ .MemUsage }}\", \"NetIO\": \"{{ .NetIO }}\", \"BlockIO\": \"{{ .BlockIO }}\" }")

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	stats := make([]Stats, len(lines)-1)
	for i := range stats {
		var raw = struct {
			ID       string
			Name     string
			CPU      string
			MemUsage string
			NetIO    string
			BlockIO  string
		}{}
		err := json.Unmarshal([]byte(lines[i]), &raw)
		netsplit := strings.Split(raw.NetIO, "/")
		stats[i] = Stats{
			ID:       raw.ID,
			Name:     raw.Name,
			CPU:      utils.MustParsePercentage(raw.CPU),
			MemUsage: utils.MustParseMegabyte(raw.MemUsage),
			NetIn:    utils.MustParseMegabyte(netsplit[0]),
			NetOut:   utils.MustParseMegabyte(netsplit[1]),
			BlockIO:  utils.MustParseMegabyte(raw.BlockIO),
		}
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
