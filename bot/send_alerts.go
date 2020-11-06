package bot

import (
	a "git.arnef.de/monitgo/alerts"
)

func (b *Bot) SendAlerts(alerts a.Alerts) {
	b.lastAlerts = alerts
	message := b.alertsToMessage()
	b.Broadcast(message)
}
