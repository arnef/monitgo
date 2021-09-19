package database

import (
	"fmt"

	"github.com/arnef/monitgo/pkg"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	log "github.com/sirupsen/logrus"
)

type InfluxDBConfig struct {
	Host         string
	Port         uint
	Database     string
	Username     string
	Password     string
	Organization string
}

func NewInfluxDB(cfg *InfluxDBConfig) Database {
	log.Debug(cfg)
	return &influxdb{
		client: influxdb2.NewClient(
			fmt.Sprintf("http://%s:%d", cfg.Host, cfg.Port),
			fmt.Sprintf("%s:%s", cfg.Username, cfg.Password),
		),
		database:     cfg.Database,
		organization: cfg.Organization,
	}
}

type influxdb struct {
	client       influxdb2.Client
	database     string
	organization string
}

func (db *influxdb) OnSnapshot(snap []pkg.NodeSnapshot) {
	if db.client != nil {
		writeAPI := db.client.WriteAPI(db.organization, db.database)

		for _, node := range snap {
			if node.Error == nil {

				writeAPI.WritePoint(influxdb2.NewPoint("mem_usage", map[string]string{
					"name": node.Name,
					"host": node.Name,
				}, map[string]interface{}{
					"total":      int64(node.MemoryUsage.TotalBytes),
					"used":       int64(node.MemoryUsage.UsedBytes),
					"percentage": node.MemoryUsage.Percentage(),
				}, node.Timestamp))

				writeAPI.WritePoint(influxdb2.NewPoint("disk_usage", map[string]string{
					"name": node.Name,
				}, map[string]interface{}{
					"total":      int64(node.DiskUsage.TotalBytes),
					"used":       int64(node.DiskUsage.UsedBytes),
					"percentage": node.DiskUsage.Percentage(),
				}, node.Timestamp))

				writeAPI.WritePoint(influxdb2.NewPoint("cpu_load", map[string]string{
					"name":    node.Name,
					"average": "1",
				}, map[string]interface{}{
					"value": node.CPU,
				}, node.Timestamp))

				// node container points

				for _, container := range node.Container {
					if container == nil {
						continue
					}
					writeAPI.WritePoint(influxdb2.NewPoint("container", map[string]string{
						"id":   container.ID,
						"name": container.Name,
						"host": node.Name,
					}, map[string]interface{}{
						"cpu":       container.CPU,
						"mem_usage": int64(container.MemoryUsage.UsedBytes),
						"net_in":    int64(container.Network.TotalRxBytes),
						"net_out":   int64(container.Network.TotalTxBytes),
					}, container.Timestamp))
				}
			}
		}
		writeAPI.Flush()
	}

}
