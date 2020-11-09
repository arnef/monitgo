package database

import (
	"fmt"
	"time"

	"git.arnef.de/monitgo/monitor"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type InfluxDB struct {
	Host         string
	Port         uint
	DatabaseName string `yaml:"database"`
	Username     string
	Password     string
	Organization string
}

var (
	client *influxdb2.Client
)

func Init(influxdb InfluxDB) {
	if client == nil {
		c := influxdb2.NewClient(
			fmt.Sprintf("http://%s:%d", influxdb.Host, influxdb.Port),
			fmt.Sprintf("%s:%s", influxdb.Username, influxdb.Password),
		)
		client = &c
		fmt.Println("💽️ Influx db initalized")
	}
}

func (db *InfluxDB) Push(data monitor.Data) {
	if client != nil {
		writeAPI := (*client).WriteAPI(db.Organization, db.DatabaseName)
		now := time.Now()

		for host, stats := range data {
			if stats.Error == nil {
				for _, container := range stats.Container {
					p := influxdb2.NewPoint("container",
						map[string]string{
							"id":   container.ID,
							"name": container.Name,
							"host": host,
						},
						map[string]interface{}{
							"cpu":       container.CPU,
							"mem_usage": container.MemUsage,
							"net_in":    container.NetRx,
							"net_out":   container.NetTx,
						},
						now,
					)
					writeAPI.WritePoint(p)
				}

				p := influxdb2.NewPoint("mem_usage", map[string]string{
					"name": stats.Name,
					"host": host,
				}, map[string]interface{}{
					"total":      stats.Host.MemUsage.Total,
					"used":       stats.Host.MemUsage.Used,
					"percentage": stats.Host.MemUsage.Percentage,
				}, now)
				writeAPI.WritePoint(p)

				p = influxdb2.NewPoint("disk_usage", map[string]string{
					"name": stats.Name,
					"host": host,
				}, map[string]interface{}{
					"total":      stats.Host.DiskUsage.Total,
					"used":       stats.Host.DiskUsage.Used,
					"percentage": stats.Host.DiskUsage.Percentage,
				}, now)
				writeAPI.WritePoint(p)

				for i, cpu := range stats.Host.CPULoad {
					average := "1"
					if i == 1 {
						average = "5"
					} else if i == 2 {
						average = "15"
					}
					p := influxdb2.NewPoint("cpu_load",
						map[string]string{
							"name":    stats.Name,
							"host":    host,
							"average": average,
						},
						map[string]interface{}{
							"value": cpu,
						},
						now,
					)
					writeAPI.WritePoint(p)
				}
			}
		}
		writeAPI.Flush()
	}
}
