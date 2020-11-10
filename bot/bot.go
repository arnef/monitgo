package bot

import (
	"fmt"
	"time"

	"git.arnef.de/monitgo/alerts"
	"git.arnef.de/monitgo/config"
	"github.com/hako/durafmt"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	chatIDs    []int64
	api        *tb.Bot
	config     config.Config
	startTime  time.Time
	lastAlerts alerts.Alerts
}

func New(config config.Config) Bot {

	// api, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	api, err := tb.NewBot(tb.Settings{
		Token:  config.Telegram.Token,
		Poller: &tb.LongPoller{Timeout: 30 * time.Second},
	})
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

// func (b *Bot) reply(chatID int64, message string) {
// 	msg := tgbotapi.NewMessage(chatID, message)
// 	msg.ParseMode = tgbotapi.ModeHTML
// 	b.api.Send(msg)
// }

func (b *Bot) isAuthorized(chatID int64) bool {
	for _, id := range b.chatIDs {
		if id == chatID {
			return true
		}
	}
	return false
}

func (b *Bot) Broadcast(message string) {
	for _, chatID := range b.chatIDs {
		b.api.Send(tb.ChatID(chatID), message, tb.ModeHTML)
	}
}

func (b *Bot) Listen() {
	fmt.Println("🤖 telegram bot running")
	defer b.api.Stop()
	b.api.SetCommands([]tb.Command{
		{
			Text:        "uptime",
			Description: "Print current uptime",
		},
		{
			Text:        "status",
			Description: "Print the current status",
		},
	})

	b.api.Handle("/start", b.start)
	b.api.Handle("/uptime", b.uptime)
	b.api.Handle("/status", b.status)
	b.api.Handle("/help", func(m *tb.Message) {
		cmds, err := b.api.GetCommands()
		var message string
		if err != nil {
			fmt.Println(err)
			message = err.Error()
		} else {
			message = "Available commands:\n"
			for _, c := range cmds {
				message += fmt.Sprintf("/%s - %s\n", c.Text, c.Description)
			}
		}
		b.api.Send(m.Sender, message, tb.ModeHTML)
	})
	b.api.Start()
}

func (b *Bot) isAdmin(msg *tb.Message) bool {
	for _, id := range b.config.Telegram.Admin {
		if id == msg.Sender.ID {
			return true
		}
	}
	return false
}

func (b *Bot) start(msg *tb.Message) {
	inList := false
	var message string
	for _, id := range b.chatIDs {
		if id == msg.Chat.ID {
			inList = true
			message = fmt.Sprintf("Hey %s! This chat is already kept up to date!\n/help", msg.Sender.FirstName)
		}
	}
	if !inList {
		if b.isAdmin(msg) {
			b.chatIDs = append(b.chatIDs, msg.Chat.ID)
			b.persistChatIDs()
			message = fmt.Sprintf("Hey %s! I will now keep you up to date!\n/help", msg.Sender.FirstName)
		} else {
			message = fmt.Sprintf("Hey %s! You're not allowed to control this bot.", msg.Sender.FirstName)
		}
	}
	b.api.Send(msg.Sender, message, tb.ModeHTML)

}

func (b *Bot) status(msg *tb.Message) {
	message := b.alertsToMessage()
	if message == "" {
		message = "🎉️ No alerts right now!"
	}
	b.api.Send(msg.Sender, message, tb.ModeHTML)
}

func (b *Bot) uptime(msg *tb.Message) {
	uptime := fmt.Sprintf("<b>Monitgo Watcher</b>\nUptime: %s\n", durafmt.ParseShort(time.Since(b.startTime)))
	b.api.Send(msg.Sender, uptime, tb.ModeHTML)
}
