package alerts

import (
	"time"

	"github.com/arnef/monitgo/pkg"
)

type GenericSnaphot struct {
	Name        string
	Error       error
	Timestamp   time.Time
	CPU         float64
	MemoryUsage *pkg.Usage
	DiskUsage   *pkg.Usage
	Network     *pkg.Network
	State       pkg.ContainerStateType
}

func mapNode2Generic(data []pkg.NodeSnapshot) map[string]GenericSnaphot {
	dataMap := make(map[string]GenericSnaphot)
	for _, d := range data {
		d := d
		dataMap[d.Name] = GenericSnaphot{
			Name:        d.Name,
			Error:       d.Error,
			Timestamp:   d.Timestamp,
			CPU:         d.CPU,
			MemoryUsage: &d.MemoryUsage,
			DiskUsage:   &d.DiskUsage,
		}
	}
	return dataMap
}

func mapContainer2Generic(data []*pkg.ContainerSnapshot) map[string]GenericSnaphot {
	dataMap := make(map[string]GenericSnaphot)
	for _, d := range data {
		d := d
		if d != nil {
			dataMap[d.Name] = GenericSnaphot{
				Name:        d.Name,
				Error:       d.Error,
				Timestamp:   d.Timestamp,
				CPU:         d.CPU,
				MemoryUsage: &d.MemoryUsage,
				Network:     &d.Network,
				State:       d.State,
			}
		}
	}
	return dataMap
}

func (c *GenericSnaphot) errorOccurred(prev *GenericSnaphot) bool {

	return c.Error != nil && (prev == nil || prev.Error == nil)
}

func (c *GenericSnaphot) errorResolved(prev *GenericSnaphot) bool {
	return c.Error == nil && prev != nil && prev.Error != nil
}

func (c *GenericSnaphot) highDiskUsageOccurred(prev *GenericSnaphot, diskUsageDiskUsage float64) bool {
	return c.DiskUsage.Percentage() > diskUsageDiskUsage &&
		(prev == nil || prev.DiskUsage.Percentage() <= diskUsageDiskUsage)
}
func (c *GenericSnaphot) highDiskUsageResolved(prev *GenericSnaphot, diskUsageDiskUsage float64) bool {
	return c.DiskUsage.Percentage() <= diskUsageDiskUsage &&
		prev != nil && prev.DiskUsage.Percentage() > diskUsageDiskUsage
}

func (c *GenericSnaphot) highMemoryUsageOccurred(prev *GenericSnaphot, memoryUsageThreshold float64) bool {
	return c.MemoryUsage.Percentage() > memoryUsageThreshold &&
		(prev == nil || prev.MemoryUsage.Percentage() <= memoryUsageThreshold)
}
func (c *GenericSnaphot) highMemoryUsageResolved(prev *GenericSnaphot, memoryUsageThreshold float64) bool {
	return c.MemoryUsage.Percentage() <= memoryUsageThreshold &&
		prev != nil && prev.MemoryUsage.Percentage() > memoryUsageThreshold
}

func (c *GenericSnaphot) highCPUUsageOccurred(prev *GenericSnaphot, cpuUsageDiskUsage float64) bool {
	return c.CPU > cpuUsageDiskUsage &&
		(prev == nil || prev.CPU <= cpuUsageDiskUsage)
}
func (c *GenericSnaphot) highCPUUsageResolved(prev *GenericSnaphot, cpuUsageDiskUsage float64) bool {
	return c.CPU <= cpuUsageDiskUsage &&
		prev != nil && prev.CPU > cpuUsageDiskUsage
}

// container stuff

func (c *GenericSnaphot) WentDown(prev *GenericSnaphot) bool {
	return isDown(c) && (prev == nil || !isDown(prev))
}
func (c *GenericSnaphot) CameUp(prev *GenericSnaphot) bool {
	return !isDown(c) && prev != nil && isDown(prev)
}
func (c *GenericSnaphot) started(prev *GenericSnaphot) bool {
	return !isDown(c) && prev == nil
}

func isDown(snap *GenericSnaphot) bool {
	return snap.State != pkg.ContainerStateRunning
}
