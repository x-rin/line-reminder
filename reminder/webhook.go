package reminder

import (
	"encoding/json"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
)

func (con *LineConfig) GetWebHook(req *http.Request) (string, error) {
	received, err := con.ReceiveEvent(req)
	if err != nil {
		return "", err
	}

	var status = "false"

	for _, event := range received {
		//log.Println("groupId: " + event.Source.GroupID)
		//log.Println("userId: " + event.Source.UserID)
		textMsg := new(linebot.TextMessage)
		byteMsg, _ := event.Message.MarshalJSON()
		if err := json.Unmarshal(byteMsg, textMsg); err != nil {
			//画像メッセージの場合もあるからただエラーを出力するだけにする
			log.Println(err.Error())
		}

		if textMsg.Text == os.Getenv("REPORT_MESSAGE") {
			status = SetStatus(event.Source.UserID, "true")
			err := con.ReplyMessage(event.ReplyToken, os.Getenv("REPLY_SUCCESS"))
			if err != nil {
				return "", err
			}
		}
	}

	return status, nil
}
