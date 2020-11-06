package docker

import (
	"encoding/json"
	"strings"

	"git.arnef.de/monitgo/node/cmd"
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
	out, err := cmd.Exec("docker", "stats", "--no-stream", "--all", "--format",
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
