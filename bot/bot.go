package bot

import (
	"fmt"

	"git.arnef.de/monitgo/config"
	"git.arnef.de/monitgo/monitor"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/urfave/cli/v2"
)

type Bot struct {
	chatIDs []int64
	api     *tgbotapi.BotAPI
}

func New() Bot {
	config := config.Get()
	api, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		panic(err)
	}
	return Bot{
		api:     api,
		chatIDs: []int64{1225509414},
	}
}

func (b *Bot) Send(chatID int64, message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ParseMode = tgbotapi.ModeMarkdown
	b.api.Send(msg)
}

func (b *Bot) Broadcast(message string) {
	for _, chatID := range b.chatIDs {
		b.Send(chatID, message)
	}
}

func (b *Bot) listen() {
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
	}
}

func (b *Bot) start(msg tgbotapi.Update) {

	b.chatIDs = append(b.chatIDs, msg.Message.Chat.ID)

	message := fmt.Sprintf("Hey %s! I will now keep you up to date!\n/help", msg.Message.From)
	b.Send(msg.Message.Chat.ID, message)
}

func (b *Bot) status(msg tgbotapi.Update) {
	status := monitor.GetStatus()
	message := ""
	for _, s := range status {
		if s.Error != "" {
			message += fmt.Sprintf("❗️ *%s*\n_%s_", s.Name, s.Error)
		} else if len(s.Data) > 0 {
			message += fmt.Sprintf("🔥️ *%s*\n", s.Name)
			for _, d := range s.Data {
				message += fmt.Sprintf("_%s_ down\n", d.Name)
			}
		} else {
			message += fmt.Sprintf("✅️ *%s*\n", s.Name)
		}
	}
	b.Send(msg.Message.Chat.ID, message)
}

func (b *Bot) help(msg tgbotapi.Update) {
	message := "Available commands:\n/start - Subscribe\n/status - Print the current status"
	b.Send(msg.Message.Chat.ID, message)
}

func Cmd(ctx *cli.Context) error {
	bot := New()
	bot.listen()
	return nil
}
