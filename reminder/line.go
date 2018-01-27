package reminder

import (
	"encoding/json"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type LineConfig struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func NewLineConfig() *LineConfig {
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", os.Getenv("CHANNEL_ID"))
	values.Set("client_secret", os.Getenv("CHANNEL_SECRET"))

	url := "https://api.line.me/v2/oauth/accessToken"

	req, err := http.NewRequest(
		"POST",
		url,
		strings.NewReader(values.Encode()),
	)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	config := new(LineConfig)

	if err := json.NewDecoder(res.Body).Decode(config); err != nil {
		log.Fatal(err)
	}

	return config
}

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

func ReceiveEvent(req *http.Request) ([]linebot.Event, error) {
	config := NewLineConfig()
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), config.AccessToken)
	if err != nil {
		return []linebot.Event{}, err
	}
	received, err := bot.ParseRequest(req)
	return received, err
}
