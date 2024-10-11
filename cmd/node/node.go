package node

import (
	"strings"

	"github.com/arnef/monitgo/interal/node-exporter/api"
	"github.com/urfave/cli/v2"
)

var Command = cli.Command{
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
		&cli.UintFlag{
			Name:  "port",
			Value: 5000,
		},
		&cli.PathFlag{
			Name:  "docker",
			Value: "/var/run/docker.sock",
		},
		&cli.StringFlag{
			Name:  "allowed",
			Value: "uptime,free,lscpu,df",
		},
	},
	Action: Action,
}

func Action(c *cli.Context) error {

	a := api.New(c.String("host"), c.Int("port"), c.Path("docker"), strings.Split(c.String("allowed"), ","))
	return a.Start()
}
