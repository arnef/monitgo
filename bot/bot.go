package bot

import (
	"encoding/json"
	"fmt"
	"time"

	"git.arnef.de/monitgo/config"
	"git.arnef.de/monitgo/monitor"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hako/durafmt"
)

type Bot struct {
	chatIDs      []int64
	api          *tgbotapi.BotAPI
	config       config.Config
	startTime    time.Time
	lastMessage  *string
	lastResponse string
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
	msg.ParseMode = tgbotapi.ModeMarkdown
	b.api.Send(msg)
}

func (b *Bot) statusToMessage() string {
	uptime := durafmt.ParseShort(time.Since(b.startTime))
	if b.lastResponse != "" {
		message := ""
		var data monitor.Data
		err := json.Unmarshal([]byte(b.lastResponse), &data)
		if err == nil {
			for _, s := range data {
				if s.Error != "" {
					message += fmt.Sprintf("❗️ *%s*\n_%s_\n", s.Name, s.Error)
				} else if len(s.Data) > 0 {
					message += fmt.Sprintf("🔥️ *%s*\n", s.Name)
					for _, d := range s.Data {
						message += fmt.Sprintf("_%s_ down\n", d.Name)
					}
				} else {
					message += fmt.Sprintf("✅️ *%s*\n", s.Name)
				}

			}
			return fmt.Sprintf("*Monitgo Watcher*\nUptime: %s\n\nNodes:\n%s", uptime, message)
		}
	}

	return fmt.Sprintf("*Monitgo Watcher*\nUptime: %s\n\nNot enought data!", uptime)

}

func (b *Bot) asyncSend(chatID int64, callable func() string) {
	msg := tgbotapi.NewMessage(chatID, "⏳")
	m, err := b.api.Send(msg)
	if err == nil {

		newMsg := tgbotapi.NewEditMessageText(chatID, m.MessageID, callable())
		newMsg.ParseMode = tgbotapi.ModeMarkdown
		b.api.Send(newMsg)
	} else {
		// TODO err message
	}
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
	} else {
		fmt.Printf("🤖 > %s unkown command\n", cmd.Message.Text)
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

func logError(err error) {
	fmt.Printf("🤖 ERROR: %s\n", err.Error())
}

func (b *Bot) status(msg tgbotapi.Update) {
	b.reply(msg.Message.Chat.ID, b.statusToMessage())
}

func (b *Bot) help(msg tgbotapi.Update) {
	message := "Available commands:\n/start - Subscribe\n/status - Print the current status"
	b.reply(msg.Message.Chat.ID, message)
}
