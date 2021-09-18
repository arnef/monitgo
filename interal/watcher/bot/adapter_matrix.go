package bot

import (
	"time"

	log "github.com/sirupsen/logrus"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/id"
)

type MatrixBotConfig struct {
	Homeserver  string
	UserID      string
	RoomID      string
	AccessToken string
}

func NewMatrixBot(cfg *MatrixBotConfig) Bot {
	mb := &matrixbot{}
	client, err := mautrix.NewClient(cfg.Homeserver, id.UserID(cfg.UserID), cfg.AccessToken)
	if err == nil {
		mb.client = client
		mb.room = id.RoomID(cfg.RoomID)
	} else {
		log.Error(err)
	}
	return mb
}

type matrixbot struct {
	client *mautrix.Client
	room   id.RoomID
}

func (m *matrixbot) Send(msg Message) {
	if m.client != nil {
		_, err := m.client.SendMessageEvent(m.room, event.EventMessage, map[string]string{
			"msgtype":        "m.text",
			"body":           msg.Plain,
			"format":         "org.matrix.custom.html",
			"formatted_body": msg.HTML,
		})
		if err != nil {
			log.Error(err)
		}
	}
}

func (m *matrixbot) Listen(commandHandler CommandHandler) {
	// only handle commands that was send after the bot was started
	startTime := time.Now().UnixNano() / 1000000
	if m.client != nil {
		syncer, ok := m.client.Syncer.(*mautrix.DefaultSyncer)
		if ok {
			syncer.OnEventType(event.EventMessage, func(_ mautrix.EventSource, evt *event.Event) {
				if evt.RoomID == m.room && evt.Timestamp > startTime {
					switch evt.Content.AsMessage().Body {
					case "!uptime":
						m.Send(commandHandler(UptimeCommand))
					case "!status":
						m.Send(commandHandler(StatusCommand))
					}
				}
			})
			log.Info("matrix bot listening for commands")
			defer m.client.StopSync()
			m.client.Sync()
		}

	}
}
