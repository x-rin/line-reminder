package reminder

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

// LineService - LineAPIを使用するメソッドを定義
type LineService interface {
	GetNameByID(id string) (string, error)
	Send(groupID, message string) error
	Hear(request *http.Request) (linebot.Event, error)
	Reply(replyToken string, message string) error
}

type lineService struct {
	client *linebot.Client
}

// NewLineService - LineServiceを生成
func NewLineService(client *linebot.Client) LineService {
	return &lineService{
		client: client,
	}
}

func (ls *lineService) GetNameByID(id string) (string, error) {
	profile := ls.client.GetProfile(id)
	target, err := profile.Do()
	if err != nil {
		return "", err
	}
	return target.DisplayName, nil
}

func (ls *lineService) Send(groupID, message string) error {
	msgRequest := ls.client.PushMessage(groupID, ls.quickReplyMessage(message))
	if _, err := msgRequest.Do(); err != nil {
		return err
	}
	return nil
}

func (ls *lineService) Hear(request *http.Request) (linebot.Event, error) {
	received, err := ls.client.ParseRequest(request)
	if err != nil {
		return linebot.Event{}, err
	}
	var retEvent linebot.Event
	for _, event := range received {
		//log.Println("groupId: " + event.Source.GroupID)
		//log.Println("userId: " + event.Source.UserID)
		retEvent = *event
	}
	return retEvent, nil
}

func (ls *lineService) Reply(replyToken string, message string) error {
	msgRequest := ls.client.ReplyMessage(replyToken, ls.quickReplyMessage(message))
	if _, err := msgRequest.Do(); err != nil {
		return err
	}
	return nil
}

func (ls *lineService) quickReplyMessage(mainMessage string) linebot.SendingMessage {
	messageAction := linebot.NewMessageAction("飲みました", "飲みました")
	button := linebot.NewQuickReplyButton("https://www.aomori-ringo.or.jp/wp-content/uploads/2018/06/wasefuji.png", messageAction)
	quickReply := linebot.NewQuickReplyItems(button)
	textMessage := linebot.NewTextMessage(mainMessage)
	return textMessage.WithQuickReplies(quickReply)
}
