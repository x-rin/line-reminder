package reminder

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"os"
)

func GetWebHook(c *gin.Context) {
	config := NewLineConfig()
	received, err := config.ReceiveEvent(c.Request)
	if err != nil {
		log.Println(err)
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
			err := config.ReplyMessage(event.ReplyToken, os.Getenv("REPLY_SUCCESS"))
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}

	Response(c, status)
}
