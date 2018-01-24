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
	config := lineConfig{}
	log.Println("TOKEN:" + config.AccessToken)
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), config.AccessToken)
	if err != nil {
		log.Fatal(err.Error())
	}

	if _, err := bot.PushMessage(os.Getenv("GROUP_ID"), linebot.NewTextMessage(os.Getenv("REPORT_MESSAGE"))).Do(); err != nil {
		log.Fatal(err.Error())
	}


	source := c.PostForm("id")
	log.Println("source: " + source)
	statusKey := strings.ToUpper(source) + "_STATUS"
	os.Setenv(statusKey, "success")
	c.JSON(http.StatusOK, gin.H{
		"status": os.Getenv(statusKey),
	})
}
