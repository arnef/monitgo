package main

import (
	"fmt"
	"os"

	"git.arnef.de/monitgo/node"
	"git.arnef.de/monitgo/watch"
	"github.com/urfave/cli/v2"
)

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
					&cli.BoolFlag{
						Name:  "no-bot",
						Value: false,
						Usage: "don't start telegram bot",
					},
					&cli.PathFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "./config.yml",
					},
				},
				Action: watch.Cmd,
			},
		},
	}).Run(os.Args)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
