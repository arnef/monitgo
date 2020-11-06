package config

import (
	"fmt"
	"io/ioutil"

	"git.arnef.de/monitgo/database"
	"git.arnef.de/monitgo/monitor"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Nodes    []monitor.Node
	Telegram *Bot
	InfluxDB *database.InfluxDB `yaml:"influxdb"`
}

type Bot struct {
	Token string
}

func Get(path string) Config {
	fmt.Printf("🔧️ loading config from %s\n", path)
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var config Config
	yaml.Unmarshal(raw, &config)
	return config
}
