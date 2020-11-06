package bot

import (
	"fmt"
	"time"

	"git.arnef.de/monitgo/alerts"
	"git.arnef.de/monitgo/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hako/durafmt"
)

type Bot struct {
	chatIDs    []int64
	api        *tgbotapi.BotAPI
	config     config.Config
	startTime  time.Time
	lastAlerts alerts.Alerts
}

func New(config config.Config) Bot {
	api, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		panic(err)
	}
	bot := Bot{
		api:       api,
		chatIDs:   []int64{},
		config:    config,
		startTime: time.Now(),
	}
	bot.restoreChatIDs()
	return bot
}

func (b *Bot) reply(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = tgbotapi.ModeHTML
	b.api.Send(msg)
}

func (b *Bot) Broadcast(message string) {
	for _, chatID := range b.chatIDs {
		b.reply(chatID, message)
	}
}

func (b *Bot) Listen() {
	fmt.Println("🤖 telegram bot running")
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := b.api.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.Text[0] == '/' {
			b.handleCommand(update)
		}
	}
}

func (b *Bot) handleCommand(cmd tgbotapi.Update) {

	if cmd.Message.Text == "/start" {
		b.start(cmd)
	} else if cmd.Message.Text == "/status" {
		b.status(cmd)
	} else if cmd.Message.Text == "/help" {
		b.help(cmd)
	} else if cmd.Message.Text == "/alerts" {
		b.alerts(cmd)
	} else {
		b.help(cmd)
	}
}

func (b *Bot) start(msg tgbotapi.Update) {
	inList := false
	for _, id := range b.chatIDs {
		if id == msg.Message.Chat.ID {
			inList = true
		}
	}
	if !inList {
		b.chatIDs = append(b.chatIDs, msg.Message.Chat.ID)
		b.persistChatIDs()
	}
	message := fmt.Sprintf("Hey %s! I will now keep you up to date!\n/help", msg.Message.From)
	b.reply(msg.Message.Chat.ID, message)

}

func (b *Bot) alerts(msg tgbotapi.Update) {
	message := b.alertsToMessage()
	if message == "" {
		message = "🎉️ No alerts right now!"
	}
	b.reply(msg.Message.Chat.ID, message)
}

func (b *Bot) status(msg tgbotapi.Update) {
	uptime := fmt.Sprintf("<b>Monitgo Watcher</b>\nUptime: %s\n", durafmt.ParseShort(time.Since(b.startTime)))
	b.reply(msg.Message.Chat.ID, uptime)
}

func (b *Bot) help(msg tgbotapi.Update) {
	message := "Available commands:\n/start - Subscribe\n/status - Print the current status\n/alerts - Print current alerts"
	b.reply(msg.Message.Chat.ID, message)
}
