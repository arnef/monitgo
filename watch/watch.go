package watch

import (
	"encoding/json"
	"fmt"
	"time"

	"git.arnef.de/monitgo/bot"
	"git.arnef.de/monitgo/config"
	"git.arnef.de/monitgo/monitor"
	"github.com/urfave/cli/v2"
)

func Cmd(ctx *cli.Context) error {
	sleep := ctx.Uint64("interval")
	watcher := new(sleep)
	watcher.noBot = ctx.Bool("no-bot")
	watcher.config = config.Get(ctx.Path("config"))
	if !watcher.noBot {
		watcher.bot = bot.New(watcher.config)
		go watcher.bot.Listen()
	}
	watcher.start()
	return nil
}

type watcher struct {
	sleep        uint64
	noBot        bool
	bot          bot.Bot
	lastResponse string
	config       config.Config
}

func new(sleep uint64) watcher {
	return watcher{
		sleep:        sleep,
		lastResponse: "",
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
	stats := monitor.GetStatus(w.config.Nodes)
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
			if s.Error != "" {
				message += fmt.Sprintf("❗️ *%s*\n_%s_", s.Name, s.Error)
			} else if len(s.Data) > 0 {
				message += fmt.Sprintf("🔥️ *%s*\n", s.Name)
				for _, d := range s.Data {
					message += fmt.Sprintf("_%s_ is down\n", d.Name)
				}
			}
		}
		if message != "" && !w.noBot {
			w.bot.Broadcast(message)
		}
	}
	w.lastResponse = response
}
