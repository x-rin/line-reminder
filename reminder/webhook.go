package reminder

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"os"
	"encoding/json"
)

func GetWebHook(c *gin.Context) {
	received, err := ReceiveEvent(c.Request)
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
			status = setStatus(event.Source.UserID, "true")
			err := ReplyMessage(event.ReplyToken, os.Getenv("REPLY_SUCCESS"))
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}

	Response(c, status)
}
