package api

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"os"
	"log"
	"net/http"
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
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
