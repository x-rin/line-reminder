package api

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"os"
)

type TextMessage struct {
	Type	string 	`json: type`
	Text	string  `json: text`
}

func PostMessage(message string) error{
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
