package main

import (
	"fmt"
	"os"

	"git.arnef.de/monitgo/alerts"
	"git.arnef.de/monitgo/bot"
	"git.arnef.de/monitgo/config"
	"git.arnef.de/monitgo/database"
	"git.arnef.de/monitgo/monitor"
	"git.arnef.de/monitgo/node"
	"github.com/urfave/cli/v2"
)

type Logger struct{}

func (l *Logger) Push(data monitor.Data) {
	fmt.Println(data)
}

func main() {
	err := (&cli.App{
		Version: "0.1.0",
		Name:    "monit",
		Usage:   "Monitoring Docker",
		Commands: []*cli.Command{
			{
				Name:  "node",
				Usage: "",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "host",
						Value: "127.0.0.1",
					},
					&cli.UintFlag{
						Name:  "port",
						Value: 5000,
					},
					&cli.BoolFlag{
						Name:    "dry-run",
						Aliases: []string{"d"},
						Value:   false,
					},
				},
				Action: node.Cmd,
			},
			{
				Name:  "watch",
				Usage: "",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:    "interval",
						Aliases: []string{"n"},
						Value:   60,
						Usage:   "interval in seconds",
					},
					&cli.PathFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "./config.yml",
					},
				},
				Action: func(ctx *cli.Context) error {
					conf := config.Get("./config.yml")

					monitor.Init(conf.Nodes, ctx.Uint64("interval"))

					am := alerts.AlertManager{}
					monitor.Register(&am)

					if conf.Telegram != nil {
						bot := bot.New(conf)
						go bot.Listen()
						am.Register(&bot)
					}

					if conf.InfluxDB != nil {
						database.Init(*conf.InfluxDB)
						monitor.Register(conf.InfluxDB)
					}

					return monitor.Start()
				},
			},
		},
	}).Run(os.Args)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
