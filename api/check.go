package api

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"

	"encoding/json"
)

func Check(c *gin.Context) {
	config := NewLineConfig()
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), config.AccessToken)
	received, err := bot.ParseRequest(c.Request)
	if err != nil {
		log.Println(err)
	}

	for _, event := range received {
		log.Println("groupId: " + event.Source.GroupID)
		if event.Source.UserID == os.Getenv("TARGET_ID") {

			textMsg := new(TextMessage)
			byteMsg, _ := event.Message.MarshalJSON()
			if err := json.Unmarshal(byteMsg, textMsg); err != nil {
				//画像メッセージの場合もあるからただエラーを出力するだけにする
				log.Println(err.Error())
				c.JSON(http.StatusOK, gin.H{
					"status": "false",
				})
			}

			if textMsg.Text == os.Getenv("REPORT_MESSAGE") {
				err := PostMessage(os.Getenv("REPLY_SUCCESS"))
				if err != nil {
					log.Fatal(err.Error())
				}
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
