package bot

import (
	"time"

	"github.com/arnef/monitgo/pkg"
	"github.com/hako/durafmt"
	log "github.com/sirupsen/logrus"
)

type Bot interface {
	Send(msg Message)
	Listen(CommandHandler)
}

type Command string

const (
	UptimeCommand Command = "cmd_uptime"
	StatusCommand Command = "cmd_status"
)

type CommandHandler = func(cmd Command) Message

type Message struct {
	Plain    string
	HTML     string
	Markdown string
}

type BotManager struct {
	bots   []Bot
	uptime time.Time
	alerts []pkg.Alert
}

func NewManager() *BotManager {
	return &BotManager{
		uptime: time.Now(),
	}
}

func (b *BotManager) HandleAlerts(new []pkg.Alert, all []pkg.Alert) {
	msg := BuildMessage(new)
	b.alerts = all
	b.sendMessage(msg)
}

func (b *BotManager) sendMessage(msg Message) {
	for _, b := range b.bots {
		b.Send(msg)
	}
}

func (b *BotManager) RegisterBot(bot Bot) {
	log.Debug("Register Bot", bot)
	b.bots = append(b.bots, bot)
}

func (b *BotManager) onCommand(command Command) Message {
	switch command {
	case UptimeCommand:
		duration := time.Since(b.uptime)
		return MessageFromHTML("uptime: " + durafmt.ParseShort(duration).String())
	case StatusCommand:
		msg := BuildMessage(b.alerts)
		if len(msg.Plain) == 0 {
			return MessageFromHTML("🎉️ No alerts right now!")
		} else {
			return msg
		}
	}
	return MessageFromHTML("<span>unkown command</span>")
}

func (b *BotManager) Listen() {
	for _, bot := range b.bots {
		go bot.Listen(b.onCommand)
	}
}
