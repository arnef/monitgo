package bot

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

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
