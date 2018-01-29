package line_reminder

import (
	"encoding/json"
	"github.com/line/line-bot-sdk-go/linebot"
	"net/http"
	"os"
)

func (l *lineReminder) GetWebHook(req *http.Request) (string, error) {
	received, err := l.client.ReceiveEvent(req)
	if err != nil {
		return "", err
	}

	var status = "false"

	for _, event := range received {
		//log.Println("groupId: " + event.Source.GroupID)
		//log.Println("userId: " + event.Source.UserID)
		msg, _ := ExtractMessage(event)
		if msg == os.Getenv("REPORT_MESSAGE") {
			status = SetStatus(event.Source.UserID, "true")
			err := l.client.ReplyMessage(event.ReplyToken, os.Getenv("REPLY_SUCCESS"))
			if err != nil {
				return "", err
			}
		}
	}

	return status, nil
}

func ExtractMessage(event linebot.Event) (string, error) {
	textMsg := new(linebot.TextMessage)
	byteMsg, _ := event.Message.MarshalJSON()
	if err := json.Unmarshal(byteMsg, textMsg); err != nil {
		return "", err
	}
	return textMsg.Text, nil
}
