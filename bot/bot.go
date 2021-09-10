package bot

import (
	"fmt"
	"strings"
	"time"

	ntb "git.arnef.de/arnef/talkbot/pkg"
	"git.arnef.de/monitgo/alerts"
	"git.arnef.de/monitgo/config"
	"github.com/hako/durafmt"
	tb "gopkg.in/tucnak/telebot.v2"

	mb "maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type Bot struct {
	chatIDs   []int64
	telegram  *tb.Bot
	talk      *ntb.Bot
	matrix    *mb.Client
	config    config.Config
	startTime time.Time
	// lastAlerts alerts.Alerts
	statusAlerts alerts.Alerts
}

func New(config config.Config) Bot {
	var telegram *tb.Bot
	var talk *ntb.Bot
	var matrix *mb.Client
	var err error
	if config.Telegram != nil {
		t, err := tb.NewBot(tb.Settings{
			Token:  config.Telegram.Token,
			Poller: &tb.LongPoller{Timeout: 30 * time.Second},
		})
		if err != nil {
			panic(err)
		}
		telegram = t
	}
	if config.Talk != nil {
		t, err := ntb.NewBot(ntb.Settings{
			URL:      config.Talk.URL,
			Username: config.Talk.Username,
			BotUID:   config.Talk.BotID,
			Password: config.Talk.Password,
			ChatID:   config.Talk.ChatID,
		})
		if err != nil {
			panic(err)
		}
		talk = t
	}
	if config.Matrix != nil {
		matrix, err = mb.NewClient(config.Matrix.Homeserver, id.UserID(config.Matrix.UserID), config.Matrix.AccessToken)
		if err != nil {
			panic(err)
		}
	}

	bot := Bot{
		telegram:  telegram,
		talk:      talk,
		matrix:    matrix,
		chatIDs:   []int64{},
		config:    config,
		startTime: time.Now(),
	}
	bot.restoreChatIDs()
	return bot
}

// func (b *Bot) isAuthorized(chatID int64) bool {
// 	for _, id := range b.chatIDs {
// 		if id == chatID {
// 			return true
// 		}
// 	}
// 	return false
// }

func (b *Bot) Broadcast(raw string, message string) {
	if len(raw) == 0 {
		return
	}
	if b.telegram != nil {
		for _, chatID := range b.chatIDs {
			b.telegram.Send(tb.ChatID(chatID), message, tb.ModeHTML)
		}
	}
	if b.talk != nil {
		b.talk.Send(raw)
	}
	if b.matrix != nil {
		b.matrix.SendMessageEvent(
			id.RoomID(b.config.Matrix.RoomID),
			event.EventMessage,
			map[string]string{
				"msgtype":        "m.text",
				"body":           raw,
				"format":         "org.matrix.custom.html",
				"formatted_body": strings.ReplaceAll(message, "\n", "<br>"),
			},
		)

	}

}

func (b *Bot) Listen() {

	if b.telegram != nil {
		fmt.Println("🤖 telegram bot running")
		defer b.telegram.Stop()
		b.telegram.SetCommands([]tb.Command{
			{
				Text:        "uptime",
				Description: "Print current uptime",
			},
			{
				Text:        "status",
				Description: "Print the current status",
			},
		})

		b.telegram.Handle("/start", b.start)
		b.telegram.Handle("/uptime", func(msg *tb.Message) {
			_, uptime := b.uptime()
			b.send(msg, uptime)
		})
		b.telegram.Handle("/status", b.status)
		b.telegram.Handle("/help", func(m *tb.Message) {
			cmds, err := b.telegram.GetCommands()
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
			b.send(m, message)
		})
		go b.telegram.Start()
	}
	if b.talk != nil {
		fmt.Println("🤖 talk bot running")
		defer b.talk.Stop()
		b.talk.Handle("uptime", func(_ ntb.Message) {
			raw, _ := b.uptime()
			b.talk.Send(raw)
		})
		b.talk.Handle("status", func(_ ntb.Message) {
			b.status(nil)
		})
		go b.talk.Start()
	}

	if b.matrix != nil {
		syncer, ok := b.matrix.Syncer.(*mb.DefaultSyncer)
		if !ok {
			panic("cannot get default syncer")
		}
		fmt.Println("🤖 matrix bot running")

		syncer.OnEventType(event.EventMessage, func(_ mb.EventSource, evt *event.Event) {
			if evt.RoomID == id.RoomID(b.config.Matrix.RoomID) {
				msg := evt.Content.AsMessage().Body
				switch msg {
				case "!uptime":
					raw, message := b.uptime()
					// b.send(nil, uptime)
					b.matrix.SendMessageEvent(
						id.RoomID(b.config.Matrix.RoomID),
						event.EventMessage,
						map[string]string{
							"msgtype":        "m.text",
							"body":           raw,
							"format":         "org.matrix.custom.html",
							"formatted_body": strings.ReplaceAll(message, "\n", "<br>"),
						},
					)
				case "!status":
					b.status(nil)
				}

			}
		})
		defer b.matrix.StopSync()
		go b.matrix.Sync()
	}
}

func (b *Bot) send(msg *tb.Message, message string) {
	if len(message) == 0 {
		return
	}
	if b.telegram != nil && msg != nil {
		b.telegram.Send(tb.ChatID(msg.Chat.ID), message, tb.ModeHTML)
	}
	if b.talk != nil {
		b.talk.Send(message)
	}
	if b.matrix != nil {
		b.matrix.SendMessageEvent(
			id.RoomID(b.config.Matrix.RoomID),
			event.EventMessage,
			map[string]string{
				"msgtype":        "m.text",
				"body":           message,
				"format":         "org.matrix.custom.html",
				"formatted_body": strings.ReplaceAll(message, "\n", "<br>"),
			},
		)
	}
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
	b.send(msg, message)

}

func (b *Bot) status(msg *tb.Message) {
	raw, message := b.alertsToMessage(b.statusAlerts)
	if message == "" {
		message = "🎉️ No alerts right now!"
		raw = message
	}
	if msg != nil {
		b.send(msg, message)
	} else {
		b.send(nil, raw)
	}
}

func (b *Bot) uptime() (string, string) {
	uptimeHtml := fmt.Sprintf("<b>Monitgo Watcher</b>\nUptime: %s\n", durafmt.ParseShort(time.Since(b.startTime)))
	uptimeRaw := fmt.Sprintf("Monitgo Watcher\nUptime: %s\n", durafmt.ParseShort(time.Since(b.startTime)))
	return uptimeRaw, uptimeHtml
}
