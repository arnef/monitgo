package alerts

import (
	"time"

	"github.com/arnef/monitgo/pkg"
)

type genericSnaphot struct {
	Name        string
	Error       error
	Timestamp   time.Time
	CPU         float64
	MemoryUsage *pkg.Usage
	DiskUsage   *pkg.Usage
	Network     *pkg.Network
}

func mapNode2Generic(data []pkg.NodeSnapshot) map[string]genericSnaphot {
	dataMap := make(map[string]genericSnaphot)
	for _, d := range data {
		dataMap[d.Name] = genericSnaphot{
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

func mapContainer2Generic(data []*pkg.ContainerSnapshot) map[string]genericSnaphot {
	dataMap := make(map[string]genericSnaphot)
	for _, d := range data {
		if d != nil {
			dataMap[d.Name] = genericSnaphot{
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

func (c *genericSnaphot) errorOccurred(prev *genericSnaphot) bool {

	return c.Error != nil && (prev == nil || prev.Error == nil)
}

func (c *genericSnaphot) errorResolved(prev *genericSnaphot) bool {
	return c.Error == nil && prev != nil && prev.Error != nil
}

func (c *genericSnaphot) highDiskUsageOccurred(prev *genericSnaphot) bool {
	return c.DiskUsage.Percentage() > 80 && (prev == nil || prev.DiskUsage.Percentage() <= 80)
}
func (c *genericSnaphot) highDiskUsageResolved(prev *genericSnaphot) bool {
	return c.DiskUsage.Percentage() <= 80 && prev != nil && prev.DiskUsage.Percentage() > 80
}

func (c *genericSnaphot) highMemoryUsageOccurred(prev *genericSnaphot) bool {
	return c.MemoryUsage.Percentage() > 80 && (prev == nil || prev.MemoryUsage.Percentage() <= 80)
}
func (c *genericSnaphot) highMemoryUsageResolved(prev *genericSnaphot) bool {
	return c.MemoryUsage.Percentage() <= 80 && prev != nil && prev.MemoryUsage.Percentage() > 80
}

func (c *genericSnaphot) highCPUUsageOccurred(prev *genericSnaphot) bool {
	return c.CPU > 80 && (prev == nil || prev.CPU <= 80)
}
func (c *genericSnaphot) highCPUUsageResolved(prev *genericSnaphot) bool {
	return c.CPU <= 80 && prev != nil && prev.CPU > 80
}

// container stuff

func (c *genericSnaphot) wentDown(prev *genericSnaphot) bool {
	return c.MemoryUsage.UsedBytes == 0 && (prev == nil || prev.MemoryUsage.UsedBytes > 0)
}
func (c *genericSnaphot) cameUp(prev *genericSnaphot) bool {
	return c.MemoryUsage.UsedBytes > 0 && prev != nil && prev.MemoryUsage.UsedBytes == 0
}
func (c *genericSnaphot) started(prev *genericSnaphot) bool {
	return c.MemoryUsage.UsedBytes > 0 && prev == nil
}
