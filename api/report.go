package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
)

func PostReport(c *gin.Context) {
	envKey := c.PostForm("id")
	source, err := GetProfile(envKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	reportErr := PostMessage(source + ": " + os.Getenv("REPORT_MESSAGE"))
	if reportErr != nil {
		log.Fatal(reportErr.Error())
	}
	statusKey := strings.ToUpper(envKey) + "_STATUS"
	os.Setenv(statusKey, "true")

	replyErr := PostMessage(os.Getenv("REPLY_SUCCESS"))
	if replyErr != nil {
		log.Fatal(reportErr.Error())
	}
	Response(c, os.Getenv(statusKey))
}
