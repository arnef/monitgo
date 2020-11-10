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
	fmt.Println("[DEBUG] write data")
	if client != nil {
		writeAPI := (*client).WriteAPI(db.Organization, db.DatabaseName)
		now := time.Now()

		for host, stats := range data {
			if stats.Error == nil {
				for id, container := range stats.Container {
					p := influxdb2.NewPoint("container",
						map[string]string{
							"id":   id,
							"name": container.Name,
							"host": host,
						},
						map[string]interface{}{
							"cpu":       container.CPU,
							"mem_usage": int64(container.Memory.UsedBytes),
							"net_in":    int64(container.Network.RxBytesPerSecond),
							"net_out":   int64(container.Network.TxBytesPerSecond),
						},
						now,
					)
					writeAPI.WritePoint(p)
				}

				p := influxdb2.NewPoint("mem_usage", map[string]string{
					"name": stats.Name,
					"host": host,
				}, map[string]interface{}{
					"total":      int64(stats.Host.Memory.TotalBytes),
					"used":       int64(stats.Host.Memory.UsedBytes),
					"percentage": stats.Host.Memory.Percentage,
				}, now)
				writeAPI.WritePoint(p)

				p = influxdb2.NewPoint("disk_usage", map[string]string{
					"name": stats.Name,
					"host": host,
				}, map[string]interface{}{
					"total":      int64(stats.Host.Disk.TotalBytes),
					"used":       int64(stats.Host.Disk.UsedBytes),
					"percentage": stats.Host.Disk.Percentage,
				}, now)
				writeAPI.WritePoint(p)

				p = influxdb2.NewPoint("cpu_load",
					map[string]string{
						"name":    stats.Name,
						"host":    host,
						"average": "1",
					},
					map[string]interface{}{
						"value": stats.Host.CPU,
					},
					now,
				)
				writeAPI.WritePoint(p)

			}
		}
		writeAPI.Flush()
	}
}
