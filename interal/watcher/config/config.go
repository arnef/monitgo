package config

import (
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Nodes []NodeConfig
}

type NodeConfig struct {
	Name string
	Host string
	Port int
}

func FromPath(path string) (*Config, error) {
	log.Infoln("read in config from %s", path)
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(raw, &config)
	return &config, err
}

func (c *Config) String() string {
	return fmt.Sprintf("Nodes: %v\n", c.Nodes)
}
