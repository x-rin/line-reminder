package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
)

func PostReport(c *gin.Context) {
	source := c.PostForm("id")
	reportErr := PostMessage(source + ": " + os.Getenv("REPORT_MESSAGE"))
	if reportErr != nil {
		log.Fatal(reportErr.Error())
	}
	statusKey := strings.ToUpper(source) + "_STATUS"
	os.Setenv(statusKey, "true")

	replyErr := PostMessage(os.Getenv("REPLY_SUCCESS"))
	if replyErr != nil {
		log.Fatal(reportErr.Error())
	}
	Response(c, os.Getenv(statusKey))
}
