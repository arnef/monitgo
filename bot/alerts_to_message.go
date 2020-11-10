package bot

import (
	"fmt"

	a "git.arnef.de/monitgo/alerts"
)

func (b *Bot) alertsToMessage() string {
	message := ""

	for host, alerts := range b.lastAlerts {
		if len(alerts) > 0 {
			message += fmt.Sprintf("<b>%s</b>\n", host)
			for _, alert := range alerts {
				switch alert.State {
				case a.Down:
					message += fmt.Sprintf("🔥️ <i>%s</i> is down\n", alert.Container)
				case a.Away:
					message += fmt.Sprintf("🗑️ <i>%s</i> removed\n", alert.Container)
				case a.Running:
					message += fmt.Sprintf("🚀️ <i>%s</i> is up again\n", alert.Container)
				case a.Error:
					message += fmt.Sprintf("❗️ %s", alert.Error)
				case a.ErrorResolved:
					message += fmt.Sprintf("✅️ <s>%s</s>\n", alert.Error)
				case a.Warning:
					message += fmt.Sprintf("⚠️ %s\n", alert.Warning)
				case a.WarningResolved:
					message += fmt.Sprintf("💚 <s>%s</s>", alert.Warning)

				}
			}
		}
	}
	return message
}
