package bot

import (
	"git.arnef.de/monitgo/monitor"
)

func (b *Bot) Push(data monitor.Data) {
	message := b.analyze(data)
	if message != "" {
		b.Broadcast(message)
	}
}
