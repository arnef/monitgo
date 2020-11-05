package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Nodes    []Node
	Telegram *Bot
}

type Node struct {
	Name string
	Host string
	Port uint
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
