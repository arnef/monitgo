package alerts

import (
	"encoding/json"

	"git.arnef.de/monitgo/monitor"
)

type AlertManager struct {
	prev   stateMap
	sender []AlertSender
}

type State int

const (
	Error         State = 0
	ErrorResolved State = 1
	Running       State = 2
	Down          State = 3
	Away          State = 4
)

type Alert struct {
	Error     *string
	Container string
	State     State
}

type Alerts map[string][]Alert

func (a *AlertManager) Register(sender AlertSender) {
	a.sender = append(a.sender, sender)
}

func (a *AlertManager) notify(alerts Alerts) {
	for _, sender := range a.sender {
		sender.SendAlerts(alerts)
	}
}

func (a *AlertManager) Push(data monitor.Data) {

	result := make(Alerts)

	state := buildState(data)
	if didChange(a.prev, state) {
		for host, data := range data {
			key := data.Name
			if data.Error != nil {
				if isHostError(a.prev, state, host) {
					result[key] = []Alert{
						{
							Container: "",
							State:     Error,
							Error:     data.Error,
						},
					}
				}
			} else {
				if isHostErrorResolved(a.prev, state, host) {
					result[key] = []Alert{
						{
							Container: "",
							State:     ErrorResolved,
							Error:     a.prev[host].Error,
						},
					}
				}
				for _, container := range data.Container {
					containerID := container.ID
					if isDown(a.prev, state, host, containerID) {
						result[key] = append(result[key], Alert{
							Container: container.Name,
							State:     Down,
						})
					} else if isUpAgain(a.prev, state, host, containerID) {
						result[key] = append(result[key], Alert{
							Container: container.Name,
							State:     Running,
						})
					}
				}
			}
		}
		a.prev = state
		a.notify(result)
	}
}

func isHostError(prev stateMap, cur stateMap, host string) bool {
	return cur[host].Error != nil && (prev == nil || prev[host].Error == nil)
}

func isHostErrorResolved(prev stateMap, cur stateMap, host string) bool {
	return cur[host].Error == nil && prev != nil && prev[host].Error != nil
}

func isDown(prev stateMap, cur stateMap, host string, containerID string) bool {
	return cur[host].Container[containerID] == Down && (prev == nil || prev[host].Container[containerID] != Down)
}

func isUpAgain(prev stateMap, cur stateMap, host string, containerID string) bool {
	return prev != nil && prev[host].Container[containerID] == Down && cur[host].Container[containerID] == Running
}

type AlertSender interface {
	SendAlerts(alerts Alerts)
}

type state struct {
	Error     *string
	Container map[string]State
}
type stateMap map[string]state

func buildState(data monitor.Data) stateMap {
	result := make(stateMap)
	for host, stats := range data {
		result[host] = state{
			Error:     stats.Error,
			Container: make(map[string]State),
		}
		for _, container := range stats.Container {
			state := Running
			if container.MemUsage == 0 {
				state = Down
			}
			result[host].Container[container.ID] = state
		}
	}

	return result
}

func didChange(prev stateMap, cur stateMap) bool {
	if prev == nil {
		return true
	}

	p, err := json.Marshal(prev)
	if err != nil {
		return true
	}
	c, err := json.Marshal(cur)
	if err != nil {
		panic(err)
	}
	return string(c) != string(p)

}
