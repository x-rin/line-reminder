package reminder

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"os"
)

func PostMessage(message string) error {
	config := NewLineConfig()
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), config.AccessToken)
	if err != nil {
		return err
	}

	if _, err := bot.PushMessage(os.Getenv("GROUP_ID"), linebot.NewTextMessage(message)).Do(); err != nil {
		return err
	}
	return nil
}

func ReplyMessage(token string, message string) error {
	config := NewLineConfig()
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), config.AccessToken)
	if err != nil {
		return err
	}

	if _, err := bot.ReplyMessage(token, linebot.NewTextMessage(message)).Do(); err != nil {
		return err
	}
	return nil
}

func GetProfile(id string) (string, error) {
	config := NewLineConfig()
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), config.AccessToken)
	if err != nil {
		return "", err
	}

	res, err := bot.GetProfile(id).Do()
	if err != nil {
		return "", err
	}
	return res.DisplayName, nil
}
