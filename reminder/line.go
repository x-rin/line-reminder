package reminder

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"net/http"
)

// LineAPI - linebotの使用するメソッドを定義
type LineAPI interface {
	PushMessage(to string, message linebot.Message) (*linebot.BasicResponse, error)
	ReplyMessage(replyToken string, message linebot.Message) (*linebot.BasicResponse, error)
	GetProfile(userID string) (*linebot.UserProfileResponse, error)
	ParseRequest(r *http.Request) ([]linebot.Event, error)
}

type lineAPI struct {
	client *linebot.Client
}

// NewLineAPI - LineAPIを生成
func NewLineAPI(config *Config) (LineAPI, error) {
	client, err := linebot.New(config.Channel.Secret, config.Access.Token)
	if err != nil {
		return nil, err
	}
	return &lineAPI{
		client: client,
	}, nil
}

func (la *lineAPI) PushMessage(to string, message linebot.Message) (*linebot.BasicResponse, error) {
	return la.client.PushMessage(to, message).Do()
}

func (la *lineAPI) ReplyMessage(replyToken string, message linebot.Message) (*linebot.BasicResponse, error) {
	return la.client.ReplyMessage(replyToken, message).Do()
}

func (la *lineAPI) GetProfile(userID string) (*linebot.UserProfileResponse, error) {
	return la.client.GetProfile(userID).Do()
}

func (la *lineAPI) ParseRequest(r *http.Request) ([]linebot.Event, error) {
	return la.client.ParseRequest(r)
}
