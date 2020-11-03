package watch

import (
	"encoding/json"
	"fmt"
	"time"

	"git.arnef.de/monitgo/bot"
	"git.arnef.de/monitgo/monitor"
	"github.com/urfave/cli/v2"
)

func Cmd(ctx *cli.Context) error {
	sleep := ctx.Uint64("interval")
	watcher := new(sleep)
	if !ctx.Bool("no-bot") {
		go watcher.bot.Listen()
	}
	watcher.start()
	return nil
}

type watcher struct {
	sleep        uint64
	bot          bot.Bot
	lastResponse string
}

func new(sleep uint64) watcher {
	return watcher{
		sleep:        sleep,
		lastResponse: "",
		bot:          bot.New(),
	}
}

func (w *watcher) start() {
	fmt.Printf("👀 watcher runs every %d seconds\n", w.sleep)
	w.run()
	for range time.Tick(time.Duration(w.sleep) * time.Second) {
		go w.run()
	}
}

func (w *watcher) run() {
	stats := monitor.GetStatus()
	resp, err := json.Marshal(&stats)
	if err != nil {
		panic(err)
	}
	response := string(resp)
	if response != w.lastResponse {
		var lastResponse map[string]monitor.Status
		json.Unmarshal([]byte(w.lastResponse), &lastResponse)
		message := ""
		for i, s := range stats {
			if prev, ok := lastResponse[i]; ok {
				if prev.Error != "" && s.Error == "" {
					message += fmt.Sprintf("✅ *%s*\nresolved: _%s_\n", s.Name, prev.Error)
				}
				if len(prev.Data) > 0 {
					resolved := ""
					for _, i := range prev.Data {
						errorResolved := true
						for _, i2 := range s.Data {
							if i.ID == i2.ID {
								errorResolved = false
							}
						}
						if errorResolved {
							resolved += fmt.Sprintf("_%s_ up again\n", i.Name)

						}
					}
					if resolved != "" {
						message += fmt.Sprintf("🚀 *%s*\n%s", s.Name, resolved)
					}
				}
				message += "\n"
			}
			// something is wrong lets fire a telegram message
			if s.Error != "" {
				message += fmt.Sprintf("❗️ *%s*\n_%s_", s.Name, s.Error)
			} else if len(s.Data) > 0 {
				message += fmt.Sprintf("🔥️ *%s*\n", s.Name)
				for _, d := range s.Data {
					message += fmt.Sprintf("_%s_ down\n", d.Name)
				}
			}
		}
		if message != "" {
			w.bot.Broadcast(message)
		}
	}
	w.lastResponse = response
}
