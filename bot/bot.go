package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

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
	bot := Bot{
		api:     api,
		chatIDs: []int64{},
	}
	bot.restoreChatIDs()
	return bot
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
	b.Send(msg.Message.Chat.ID, message)

}

func (b *Bot) persistChatIDs() {

	data, err := json.Marshal(b.chatIDs)
	if err != nil {
		logError(err)
		return
	}
	file, err := configFile()
	if err != nil {
		logError(err)
		return
	}

	if err := ioutil.WriteFile(file, []byte(data), 0600); err != nil {
		logError(err)
		return
	}
}

func configFile() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	dir := path.Join(configDir, "monitgo")

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, os.ModePerm); err != nil {
			return "", err
		}
	}

	file := path.Join(dir, "chat_ids.json")
	if _, err := os.Stat(file); os.IsNotExist(err) {
		f, err := os.Create(file)
		defer f.Close()
		if err != nil {
			return "", err
		}
	}
	return file, nil
}

func logError(err error) {
	fmt.Printf("🤖 ERROR: %s", err.Error())
}

func (b *Bot) restoreChatIDs() {
	file, err := configFile()
	if err != nil {
		logError(err)
		return
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		logError(err)
		return
	}
	if err := json.Unmarshal(data, &b.chatIDs); err != nil {
		logError(err)
		return
	}
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
	bot.Listen()
	return nil
}
