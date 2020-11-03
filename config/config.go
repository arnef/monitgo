package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Nodes    []Node
	Telegram Bot
}

type Node struct {
	Name string
	Host string
	Port uint
}

type Bot struct {
	Token string
}

func Get() Config {
	raw, err := ioutil.ReadFile("config.yml")
	if err != nil {
		panic(err)
	}
	var config Config
	yaml.Unmarshal(raw, &config)

	return config
}
