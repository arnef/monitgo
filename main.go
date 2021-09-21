package main

import (
	"fmt"
	"os"

	"github.com/arnef/monitgo/cmd/node"
	"github.com/arnef/monitgo/cmd/watch"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var version string = "1.0.0"

func main() {
	/// setup logger
	log.SetOutput(os.Stdout)

	err := (&cli.App{
		Version: version,
		Name:    "monitgo",
		Usage:   "Monitoring",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"vv"},
			},
		},
		Before: func(ctx *cli.Context) error {
			if ctx.Bool("debug") {
				log.SetLevel(log.DebugLevel)
				log.SetReportCaller(true)
			}
			return nil
		},
		Commands: []*cli.Command{
			&node.Command,
			&watch.Command,
		},
	}).Run(os.Args)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
