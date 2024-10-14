package alerts

import (
	"fmt"

	"github.com/arnef/monitgo/interal/watcher/config"
	"github.com/arnef/monitgo/pkg"
)

type AlertManager struct {
	previous []pkg.NodeSnapshot
	handler  []pkg.AlertHandler
	config   *config.Config
}

func NewManager(config *config.Config) *AlertManager {
	if config == nil {
		panic("config needed")
	}

	return &AlertManager{
		config: config,
	}
}

func (a *AlertManager) RegisterAlertHandler(handler pkg.AlertHandler) {
	a.handler = append(a.handler, handler)
}

func (a *AlertManager) Notify(new []pkg.Alert, all []pkg.Alert) {
	for _, handler := range a.handler {
		handler(new, all)
	}
}

func (a *AlertManager) HandleSnaphsot(snapshot []pkg.NodeSnapshot) {

	alerts := a.generate(a.previous, snapshot)
	a.previous = snapshot

	all := a.generate(nil, snapshot)
	if len(alerts) > 0 {
		go a.Notify(alerts, all)
	}
}

func (a *AlertManager) generate(previous []pkg.NodeSnapshot, current []pkg.NodeSnapshot) []pkg.Alert {
	prev := mapNode2Generic(previous)
	cur := mapNode2Generic(current)
	alerts := alertlist{}
	for name, curNode := range cur {
		var prevNode *GenericSnaphot

		if p, ok := prev[name]; ok {
			prevNode = &p
		}
		alerts.key = name
		if curNode.errorOccurred(prevNode) {
			alerts.append(pkg.Error, curNode.Error.Error())
		} else {
			if curNode.errorResolved(prevNode) {
				alerts.append(pkg.ErrorResolved, prevNode.Error.Error())
			}
			if curNode.highCPUUsageOccurred(prevNode, *a.config.CPUUsageThreshold) {
				alerts.append(pkg.Warning, "High CPU usage")
			}
			if curNode.highCPUUsageResolved(prevNode, *a.config.CPUUsageThreshold) {
				alerts.append(pkg.WarningResolved, "High CPU usage")
			}
			if curNode.highDiskUsageOccurred(prevNode, *a.config.DiskUsageDiskUsage) {
				alerts.append(pkg.Warning, "High Disk Usage")
			}
			if curNode.highDiskUsageResolved(prevNode, *a.config.DiskUsageDiskUsage) {
				alerts.append(pkg.WarningResolved, "High Disk Usage")
			}
			if curNode.highMemoryUsageOccurred(prevNode, *a.config.MemoryUsageThreshold) {
				alerts.append(pkg.Warning, "High Memory usage")
			}
			if curNode.highMemoryUsageResolved(prevNode, *a.config.MemoryUsageThreshold) {
				alerts.append(pkg.WarningResolved, "High Memory usage")
			}

			prevC := map[string]GenericSnaphot{}
			if prevNode != nil {
				for _, n := range previous {
					if n.Name == prevNode.Name {
						prevC = mapContainer2Generic(n.Container)
						break
					}
				}
			}
			curC := map[string]GenericSnaphot{}
			for _, n := range current {
				if n.Name == curNode.Name {
					curC = mapContainer2Generic(n.Container)
				}
			}

			a.generateContainer(prevC, curC, &alerts)

		}
	}
	return alerts.value
}

func (a *AlertManager) generateContainer(previous map[string]GenericSnaphot, current map[string]GenericSnaphot, alerts *alertlist) {

	for id, curContainer := range current {
		var prevContainer *GenericSnaphot
		if val, ok := previous[id]; ok {
			prevContainer = &val
			delete(previous, id)
		}

		if curContainer.errorOccurred(prevContainer) {
			alerts.append(pkg.Error, fmt.Sprintf("[%s] %s", curContainer.Name, curContainer.Error))
		} else {
			if curContainer.errorResolved(prevContainer) {
				alerts.append(pkg.ErrorResolved, fmt.Sprintf("[%s] %s", curContainer.Name, prevContainer.Error))
			} else if curContainer.WentDown(prevContainer) {
				alerts.append(pkg.Down, curContainer.Name)
			} else if curContainer.CameUp(prevContainer) {
				alerts.append(pkg.Running, curContainer.Name)
			} else if len(previous) > 0 {
				if curContainer.started(prevContainer) {
					alerts.append(pkg.Started, curContainer.Name)
				}
			}
		}
	}

	for _, deleted := range previous {
		alerts.append(pkg.Away, deleted.Name)
	}
}

type alertlist struct {
	value []pkg.Alert
	key   string
}

func (a *alertlist) append(t pkg.AlertType, message string) {
	a.value = append(a.value, pkg.Alert{
		Key:     a.key,
		Type:    t,
		Message: message,
	})
}
