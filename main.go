package main

import (
	"fmt"
	"os"

	"git.arnef.de/monitgo/bot"
	"git.arnef.de/monitgo/node"
	"github.com/urfave/cli/v2"
)

func main() {
	(&cli.App{
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
				},
				Action: node.Cmd,
			},
			{
				Name:  "monitor",
				Usage: "",
				Action: func(ctx *cli.Context) error {

					// status := monitor.GetStatus()

					// fmt.Println(status)
					// bot := bot.New()
					// bot.Send(status)
					return fmt.Errorf("Not implemented")

				},
			},
			{
				Name:   "bot",
				Usage:  "",
				Action: bot.Cmd,
			},
		},
	}).Run(os.Args)
}
