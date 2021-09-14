package bot

import (
	"fmt"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/arnef/monitgo/pkg"
	"github.com/k3a/html2text"
)

func BuildMessage(alerts []pkg.Alert) Message {

	groupedAlerts := map[string][]pkg.Alert{}

	for _, alert := range alerts {
		groupedAlerts[alert.Key] = append(groupedAlerts[alert.Key], alert)
	}

	message := ""

	for host, alerts := range groupedAlerts {
		if len(alerts) > 0 {
			if len(message) > 0 {
				message += "<br>"
			}
			message += fmt.Sprintf("<b>%s</b><br>", host)
			for _, alert := range alerts {
				switch alert.Type {
				case pkg.Down:
					message += fmt.Sprintf("🔥️ <i>%s</i> is down<br>", alert.Message)
				case pkg.Away:
					message += fmt.Sprintf("🗑️ <i>%s</i> removed<br>", alert.Message)
				case pkg.Running:
					message += fmt.Sprintf("🚀️ <i>%s</i> is up again<br>", alert.Message)
				case pkg.Error:
					message += fmt.Sprintf("❗️ %s<br>", alert.Message)
				case pkg.ErrorResolved:
					message += fmt.Sprintf("✅️ <s>%s</s><br>", alert.Message)
				case pkg.Warning:
					message += fmt.Sprintf("⚠️ %s<br>", alert.Message)
				case pkg.WarningResolved:
					message += fmt.Sprintf("💚 <s>%s</s><br>", alert.Message)
				default:
					message += fmt.Sprintf("[Unkown Type] %s<br>", alert.Message)
				}
			}
		}
	}
	return MessageFromHTML(message)
}

func MessageFromHTML(html string) Message {

	converter := md.NewConverter("", true, nil)
	markdown, _ := converter.ConvertString(html)
	plain := html2text.HTML2Text(html)

	return Message{
		Plain:    plain,
		Markdown: markdown,
		HTML:     html,
	}
}
