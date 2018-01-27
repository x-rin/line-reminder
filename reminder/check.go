package reminder

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"

	"encoding/json"
	"strings"
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

		textMsg := new(linebot.TextMessage)
		byteMsg, _ := event.Message.MarshalJSON()
		if err := json.Unmarshal(byteMsg, textMsg); err != nil {
			//画像メッセージの場合もあるからただエラーを出力するだけにする
			log.Println(err.Error())
			c.JSON(http.StatusOK, gin.H{
				"status": "false",
			})
		}

		if textMsg.Text == os.Getenv("REPORT_MESSAGE") {
			statusKey := strings.ToUpper(event.Source.UserID) + "_STATUS"
			os.Setenv(statusKey, "true")
			err := ReplyMessage(event.ReplyToken, os.Getenv("REPLY_SUCCESS"))
			if err != nil {
				log.Fatal(err.Error())
			}
		}
	}

	Response(c, "ok")
}
