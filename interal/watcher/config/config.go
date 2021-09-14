package config

import (
	"fmt"
	"io/ioutil"

	"github.com/arnef/monitgo/interal/watcher/bot"
	"github.com/arnef/monitgo/interal/watcher/database"
	"github.com/arnef/monitgo/interal/watcher/node"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Nodes []node.Node
	// bots
	Matrix   *bot.MatrixBotConfig
	Talk     *bot.TalkBotConfig
	Telegram *bot.TelegramBotConfig
	// databases
	InfluxDB *database.InfluxDBConfig
}

func FromPath(path string) (*Config, error) {
	log.Infof("read in config from %s", path)
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
