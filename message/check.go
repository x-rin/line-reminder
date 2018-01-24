package message

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/go/src/pkg/net/http"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"os"
)

func (m *message) Check(c *gin.Context) {
	//TODO: ACCESS_TOKENを自動で更新できるような仕組みがほしい
	//有効期限がきれると自分たちで手動で環境変数をセットしなおさなきゃらない。
	bot, _ := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_TOKEN"))
	received, err := bot.ParseRequest(c.Request)
	if err != nil {
		log.Println(err)
	}

	for _, event := range received {
		log.Println("Group ID: " + event.Source.GroupID)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
