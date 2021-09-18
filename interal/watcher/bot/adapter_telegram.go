package bot

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/tucnak/telebot.v2"
)

type TelegramBotConfig struct {
	Token  string
	ChatID int
	Admin  []int
}

func NewTelegramBot(cfg *TelegramBotConfig) Bot {
	tb := &telegrambot{}

	t, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.Token,
		Poller: &telebot.LongPoller{Timeout: 30 * time.Second},
	})

	if err == nil {
		tb.bot = t
		tb.chatID = telebot.ChatID(cfg.ChatID)
	} else {
		log.Error(err)
	}

	return tb
}

type telegrambot struct {
	bot    *telebot.Bot
	chatID telebot.ChatID
}

func (t *telegrambot) Send(msg Message) {
	if t.bot != nil && t.chatID > 0 {
		_, err := t.bot.Send(t.chatID, strings.ReplaceAll(msg.HTML, "<br>", "\n"), telebot.ModeHTML)
		if err != nil {
			log.Error(err)
		}
	}
}

func (t *telegrambot) Listen(commandHandler CommandHandler) {
	if t.bot != nil {
		log.Info("telegram bot listening for commands")
		defer t.bot.Stop()
		t.bot.Handle("/uptime", func(msg *telebot.Message) {
			t.Send(commandHandler(UptimeCommand))
		})
		t.bot.Handle("/status", func(msg *telebot.Message) {
			t.Send(commandHandler(StatusCommand))
		})
		t.bot.Handle("/start", func(msg *telebot.Message) {
			t.chatID = telebot.ChatID(msg.Chat.ID)
			t.Send(MessageFromHTML(fmt.Sprintf("Add this chat id to your monitgo config<br><code>%d</code>", msg.Chat.ID)))
		})
		t.bot.Start()
	}
}
