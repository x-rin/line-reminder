package message

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"os"
	"log"
	"net/http"
)

func Check(c *gin.Context) {
	var groupId string
	config := lineConfig{}
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), config.AccessToken)
	received, err := bot.ParseRequest(c.Request)
	if err != nil {
		log.Println(err)
	}

	for _, event := range received {
		groupId = event.Source.GroupID
	}

	log.Println("groupId: " + groupId)
	c.JSON(http.StatusOK, gin.H{
		"group_id": groupId,
	})
}
