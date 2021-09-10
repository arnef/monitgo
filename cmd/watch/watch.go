package watch

import (
	"github.com/arnef/monitgo/interal/watcher"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
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
	Action: Action,
}

func Action(c *cli.Context) error {
	return watcher.Start(c.Path("config"), c.Int("interval"))
}
