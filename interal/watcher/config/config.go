package config

import (
	"fmt"
	"os"

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
	Telegram *bot.TelegramBotConfig
	// databases
	InfluxDB *database.InfluxDBConfig

	// Thresholds
	CPUUsageThreshold    *float64
	MemoryUsageThreshold *float64
	DiskUsageDiskUsage   *float64
}

func FromPath(path string) (*Config, error) {
	log.Infof("read in config from %s", path)
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(raw, &config)

	defaultThreshold := float64(80)
	if config.CPUUsageThreshold == nil {
		config.CPUUsageThreshold = &defaultThreshold
	}
	if config.DiskUsageDiskUsage == nil {
		config.DiskUsageDiskUsage = &defaultThreshold
	}
	if config.MemoryUsageThreshold == nil {
		config.MemoryUsageThreshold = &defaultThreshold
	}

	return &config, err
}

func (c *Config) String() string {
	return fmt.Sprintf("Nodes: %v\n", c.Nodes)
}
