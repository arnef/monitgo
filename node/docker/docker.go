package docker

import (
	"encoding/json"
	"strings"
	"time"

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

var (
	diff      map[string]Stats
	timestamp time.Time
)

func GetStats() ([]Stats, error) {
	runningTime := time.Now()
	if diff == nil {
		query, err := doGetStats()
		if err != nil {
			return nil, err
		}
		diff = query
		timestamp = runningTime
		time.Sleep(2 * time.Second)
	}

	runningTime = time.Now()
	current, err := doGetStats()
	if err != nil {
		return nil, err
	}
	var diffStats []Stats
	for i, stat := range current {
		if _, ok := diff[i]; ok {
			duration := runningTime.Sub(timestamp).Seconds()
			diffStats = append(diffStats, Stats{
				ID:       stat.ID,
				Name:     stat.Name,
				CPU:      stat.CPU,
				MemUsage: stat.MemUsage,
				NetIn:    (current[i].NetIn - diff[i].NetIn) / duration,
				NetOut:   (current[i].NetOut - diff[i].NetOut) / duration,
				BlockIO:  stat.BlockIO,
			})
		}
	}

	diff = current
	timestamp = runningTime
	return diffStats, nil
}

func doGetStats() (map[string]Stats, error) {
	out, err := cmd.Exec("docker", "stats", "--no-stream", "--all", "--format",
		"{\"ID\": \"{{ .Container }}\", \"Name\": \"{{ .Name }}\", \"CPU\": \"{{ .CPUPerc }}\", \"MemUsage\": \"{{ .MemUsage }}\", \"NetIO\": \"{{ .NetIO }}\", \"BlockIO\": \"{{ .BlockIO }}\" }")

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(out), "\n")
	stats := make(map[string]Stats)
	for i := range lines {
		if i < len(lines)-1 {
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
			stats[raw.ID] = Stats{
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
	}
	return stats, nil
}
