package bot

import (
	"time"

	talk "git.arnef.de/arnef/talkbot/pkg"
	log "github.com/sirupsen/logrus"
)

type TalkBotConfig struct {
	URL      string
	Username string
	Password string
	UID      string
	Chat     string
}

func NewTalkBot(cfg *TalkBotConfig) Bot {
	if len(cfg.UID) == 0 {
		cfg.UID = cfg.Username
	}
	log.Debug(cfg)
	bot, err := talk.NewBot(talk.Settings{
		URL:      cfg.URL,
		Username: cfg.Username,
		BotUID:   cfg.UID,
		Password: cfg.Password,
		ChatID:   cfg.Chat,
	})
	tb := &talkbot{}
	if err == nil {
		tb.bot = bot
	} else {
		log.Error(err)
	}
	return tb
}

type talkbot struct {
	bot *talk.Bot
}

func (t *talkbot) Send(msg Message) {
	if t.bot != nil {
		err := t.bot.Send(msg.Plain)
		if err != nil {
			log.Error(err)
		}
	}
}

func (t *talkbot) Listen(commandHandler CommandHandler) {
	startTime := time.Now().Unix()
	if t.bot != nil {
		log.Info("talk bot listening for commands")
		defer t.bot.Stop()

		t.bot.Handle("uptime", func(m talk.Message) {
			if m.Timestamp > startTime {
				t.Send(commandHandler(UptimeCommand))
			}
		})
		t.bot.Handle("status", func(m talk.Message) {
			if m.Timestamp > startTime {
				t.Send(commandHandler(StatusCommand))
			}
		})
		t.bot.Start()
	}

}
