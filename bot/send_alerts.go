package bot

import (
	a "git.arnef.de/monitgo/alerts"
)

func (b *Bot) SendAlerts(alerts a.Alerts) {
	raw, html := b.alertsToMessage(alerts)
	b.Broadcast(raw, html)
}

func (b *Bot) SaveStatus(alerts a.Alerts) {
	b.statusAlerts = alerts
}

// func isErrorAlert(alert a.Alert) bool {
// 	return alert.State == a.Error || alert.State == a.Down || alert.State == a.Warning
// }
