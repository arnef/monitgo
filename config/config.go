package config

import (
	"fmt"
	"io/ioutil"

	"git.arnef.de/monitgo/database"
	"git.arnef.de/monitgo/monitor"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Nodes    []monitor.NodeConfig
	Telegram *Bot
	Talk     *TalkBot           `yaml:"talk"`
	InfluxDB *database.InfluxDB `yaml:"influxdb"`
}

type TalkBot struct {
	URL      string `yaml:"url"`
	Username string `yaml:"username"`
	BotID    string `yaml:"uid"`
	Password string `yaml:"password"`
	ChatID   string `yaml:"chat"`
}
type Bot struct {
	Token string
	Admin []int
}

func Get(path string) Config {
	fmt.Printf("🔧️ loading config from %s\n", path)
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var config Config
	yaml.Unmarshal(raw, &config)
	fmt.Println(config)
	return config
}
