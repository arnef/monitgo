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
}

func mapNode2Generic(data []pkg.NodeSnapshot) map[string]GenericSnaphot {
	dataMap := make(map[string]GenericSnaphot)
	for _, d := range data {
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
		if d != nil {
			dataMap[d.Name] = GenericSnaphot{
				Name:        d.Name,
				Error:       d.Error,
				Timestamp:   d.Timestamp,
				CPU:         d.CPU,
				MemoryUsage: &d.MemoryUsage,
				Network:     &d.Network,
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

func (c *GenericSnaphot) highDiskUsageOccurred(prev *GenericSnaphot) bool {
	return c.DiskUsage.Percentage() > 80 && (prev == nil || prev.DiskUsage.Percentage() <= 80)
}
func (c *GenericSnaphot) highDiskUsageResolved(prev *GenericSnaphot) bool {
	return c.DiskUsage.Percentage() <= 80 && prev != nil && prev.DiskUsage.Percentage() > 80
}

func (c *GenericSnaphot) highMemoryUsageOccurred(prev *GenericSnaphot) bool {
	return c.MemoryUsage.Percentage() > 80 && (prev == nil || prev.MemoryUsage.Percentage() <= 80)
}
func (c *GenericSnaphot) highMemoryUsageResolved(prev *GenericSnaphot) bool {
	return c.MemoryUsage.Percentage() <= 80 && prev != nil && prev.MemoryUsage.Percentage() > 80
}

func (c *GenericSnaphot) highCPUUsageOccurred(prev *GenericSnaphot) bool {
	return c.CPU > 80 && (prev == nil || prev.CPU <= 80)
}
func (c *GenericSnaphot) highCPUUsageResolved(prev *GenericSnaphot) bool {
	return c.CPU <= 80 && prev != nil && prev.CPU > 80
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
	return snap.MemoryUsage.TotalBytes == 0 && snap.MemoryUsage.UsedBytes == 0
}
