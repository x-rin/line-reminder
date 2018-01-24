package message

import (
	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"net/http"
	"os"
	"strings"
)

func PostReport(c *gin.Context) {
	source := c.PostForm("id")
	config := NewLineConfig()
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), config.AccessToken)
	if err != nil {
		log.Fatal(err.Error())
	}

	if _, err := bot.PushMessage(os.Getenv("GROUP_ID"), linebot.NewTextMessage(source+": "+os.Getenv("REPORT_MESSAGE"))).Do(); err != nil {
		log.Fatal(err.Error())
	}

	statusKey := strings.ToUpper(source) + "_STATUS"
	os.Setenv(statusKey, "true")
	c.JSON(http.StatusOK, gin.H{
		"status": os.Getenv(statusKey),
	})
}
