package bot

import (
	"encoding/json"
	"fmt"

	"git.arnef.de/monitgo/monitor"
	"git.arnef.de/monitgo/node/docker"
)

func (b *Bot) analyze(stats monitor.Data) string {

	resp, err := json.Marshal(&stats)
	if err != nil {
		panic(err)
	}
	response := string(resp)
	if response != b.lastResponse {
		var lastResponse monitor.Data
		json.Unmarshal([]byte(b.lastResponse), &lastResponse)
		message := ""
		for i, s := range stats {
			if prev, ok := lastResponse[i]; ok {
				if prev.Error != nil && s.Error == nil {
					message += fmt.Sprintf("✅ *%s*\nresolved: _%s_\n", s.Name, prev.Error)
				}
				if len(prev.Container) > 0 {
					resolved := ""
					for _, i := range prev.Container {
						errorResolved := true
						for _, i2 := range s.Container {
							if i.ID == i2.ID {
								errorResolved = false
							}
						}
						if errorResolved {
							resolved += fmt.Sprintf("_%s_ is up again\n", i.Name)

						}
					}
					if resolved != "" {
						message += fmt.Sprintf("🚀 *%s*\n%s", s.Name, resolved)
					}
				}
				message += "\n"
			}
			// something is wrong lets fire a telegram message
			if s.Error != nil {
				message += fmt.Sprintf("❗️ *%s*\n_%s_\n", s.Name, s.Error)
			} else {
				filtered := containerWithNoMemUsage(s.Container)
				if len(filtered) > 0 {
					message += fmt.Sprintf("🔥️ *%s*\n", s.Name)
					for _, d := range filtered {
						message += fmt.Sprintf("_%s_ is down\n", d.Name)
					}
				}
			}
		}
		b.lastResponse = response
		return message
	}
	return ""
}

func containerWithNoMemUsage(data []docker.Stats) []docker.Stats {
	var filteredData []docker.Stats

	for _, d := range data {
		if d.MemUsage == 0 {
			filteredData = append(filteredData, d)
		}
	}

	return filteredData
}
