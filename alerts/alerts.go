package alerts

import (
	"git.arnef.de/monitgo/monitor"
)

type AlertManager struct {
	prev   *monitor.Data
	sender []AlertSender
}

type State int

const (
	Error           State = 0
	ErrorResolved   State = 1
	Running         State = 2
	Down            State = 3
	Away            State = 4
	Warning         State = 5
	WarningResolved State = 6
)

type Alert struct {
	Error     string
	Warning   string
	Container string
	State     State
}

type Alerts map[string][]Alert

func (a *AlertManager) Register(sender AlertSender) {
	a.sender = append(a.sender, sender)
}

func (a *AlertManager) notify(alerts Alerts, status Alerts) {
	for _, sender := range a.sender {
		sender.SendAlerts(alerts)
		sender.SaveStatus(status)
	}
}

func (a *AlertManager) generate(data monitor.Data, prev *monitor.Data) Alerts {
	alerts := make(map[string][]Alert)
	for host, node := range data {
		if node.Error != nil {
			if a.errorOccoured(host, node, prev) {
				alerts[node.Name] = append(alerts[node.Name], Alert{
					Error: *node.Error,
					State: Error,
				})
			}
		} else {
			if yes, err := a.errorResolved(host, node, prev); yes {
				alerts[node.Name] = append(alerts[node.Name], Alert{
					Error: *err,
					State: ErrorResolved,
				})
			}
			if a.highCPUUsageOccurred(host, node.Host, prev) {
				alerts[node.Name] = append(alerts[node.Name], Alert{
					Warning: "High CPU usage",
					State:   Warning,
				})
			}
			if a.highCPUUsageResolved(host, node.Host, prev) {
				alerts[node.Name] = append(alerts[node.Name], Alert{
					Warning: "High CPU usage",
					State:   WarningResolved,
				})
			}

			if a.highDiskUsageOccurred(host, node.Host, prev) {
				alerts[node.Name] = append(alerts[node.Name], Alert{
					Warning: "High Disk usage",
					State:   Warning,
				})
			}
			if a.highDiskUsageResolved(host, node.Host, prev) {
				alerts[node.Name] = append(alerts[node.Name], Alert{
					Warning: "High Disk usage",
					State:   WarningResolved,
				})
			}

			if a.highMemoryUsageOccurred(host, node.Host, prev) {
				alerts[node.Name] = append(alerts[node.Name], Alert{
					Warning: "High Memory usage",
					State:   Warning,
				})
			}
			if a.highMemoryUsageResolved(host, node.Host, prev) {
				alerts[node.Name] = append(alerts[node.Name], Alert{
					Warning: "High Memory usage",
					State:   WarningResolved,
				})
			}

			for id, container := range node.Container {
				if a.containerWentDown(host, id, container, prev) {
					alerts[node.Name] = append(alerts[node.Name], Alert{
						Container: container.Name,
						State:     Down,
					})
				}
				if a.containerWentUpAgin(host, id, container, prev) {
					alerts[node.Name] = append(alerts[node.Name], Alert{
						Container: container.Name,
						State:     Running,
					})
				}
			}

			for _, name := range a.getTrashedContainer(host, node.Container, prev) {
				alerts[node.Name] = append(alerts[node.Name], Alert{
					Container: name,
					State:     Away,
				})
			}
		}
	}

	return alerts
}

func (a *AlertManager) Push(data monitor.Data) {

	sendableAlerts := a.generate(data, a.prev)
	a.prev = &data

	statusAlerts := a.generate(data, nil)
	a.notify(sendableAlerts, statusAlerts)
}

func (a *AlertManager) getTrashedContainer(host string, container map[string]monitor.ContainerStats, prev *monitor.Data) []string {
	var names []string

	if prev != nil {
		if p, ok := (*prev)[host]; ok {
			for id, con := range p.Container {
				if _, ok := container[id]; !ok {
					names = append(names, con.Name)
				}
			}
		}
	}

	return names
}

func (a *AlertManager) errorOccoured(key string, node monitor.Status, prev *monitor.Data) bool {
	if node.Error != nil {
		if prev == nil {
			return true
		}
		if val, ok := (*prev)[key]; ok {
			return val.Error == nil
		}
	}

	return false
}

func (a *AlertManager) errorResolved(key string, node monitor.Status, prev *monitor.Data) (bool, *string) {
	if node.Error == nil && prev != nil {
		if val, ok := (*prev)[key]; ok {
			if val.Error != nil {
				return true, val.Error
			}

		}
	}
	return false, nil
}

func (a *AlertManager) highDiskUsageOccurred(key string, host monitor.HostStats, prev *monitor.Data) bool {
	p := a.getPreviousHost(key, prev)
	return host.Disk.Percentage > 80 && (p == nil || p.Disk.Percentage <= 80)
}
func (a *AlertManager) highDiskUsageResolved(key string, host monitor.HostStats, prev *monitor.Data) bool {
	p := a.getPreviousHost(key, prev)
	return p != nil && p.Disk.Percentage > 80 && host.Disk.Percentage <= 80
}

func (a *AlertManager) highMemoryUsageOccurred(key string, host monitor.HostStats, prev *monitor.Data) bool {
	p := a.getPreviousHost(key, prev)
	return host.Memory.Percentage > 80 && (p == nil || p.Memory.Percentage <= 80)
}
func (a *AlertManager) highMemoryUsageResolved(key string, host monitor.HostStats, prev *monitor.Data) bool {
	p := a.getPreviousHost(key, prev)
	return p != nil && p.Memory.Percentage > 80 && host.Memory.Percentage <= 80
}

func (a *AlertManager) highCPUUsageOccurred(key string, host monitor.HostStats, prev *monitor.Data) bool {
	p := a.getPreviousHost(key, prev)
	return host.CPU > 80 && (p == nil || p.CPU <= 80)
}
func (a *AlertManager) highCPUUsageResolved(key string, host monitor.HostStats, prev *monitor.Data) bool {
	p := a.getPreviousHost(key, prev)
	return p != nil && p.CPU > 80 && host.CPU <= 80
}

func (a *AlertManager) containerWentDown(host string, id string, container monitor.ContainerStats, prev *monitor.Data) bool {
	p := a.getPreviousContainer(host, id, prev)
	return container.Memory.UsedBytes == 0 && (p == nil || p.Memory.UsedBytes > 0)
}

func (a *AlertManager) containerWentUpAgin(host string, id string, container monitor.ContainerStats, prev *monitor.Data) bool {
	p := a.getPreviousContainer(host, id, prev)
	return p != nil && p.Memory.UsedBytes == 0 && container.Memory.UsedBytes > 0
}

func (a *AlertManager) getPreviousHost(host string, prev *monitor.Data) *monitor.HostStats {
	if prev != nil {
		if h, ok := (*prev)[host]; ok {
			return &h.Host
		}
	}
	return nil
}

func (a *AlertManager) getPreviousContainer(host string, id string, prev *monitor.Data) *monitor.ContainerStats {
	if prev != nil {
		if _, ok := (*prev)[host]; ok {
			if val, ok := (*prev)[host].Container[id]; ok {
				return &val
			}
		}
	}
	return nil
}

type AlertSender interface {
	SendAlerts(alerts Alerts)
	SaveStatus(alerts Alerts)
}
