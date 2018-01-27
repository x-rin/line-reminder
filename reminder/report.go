package reminder

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func PostReport(c *gin.Context) {
	id := c.PostForm("id")
	config := NewLineConfig()
	source, err := config.GetProfile(id)
	if err != nil {
		log.Fatal(err.Error())
	}

	reportErr := config.PostMessage(source + ": " + os.Getenv("REPORT_MESSAGE"))
	if reportErr != nil {
		log.Fatal(reportErr.Error())
	}

	status := SetStatus(id, "true")

	replyErr := config.PostMessage(os.Getenv("REPLY_SUCCESS"))
	if replyErr != nil {
		log.Fatal(reportErr.Error())
	}

	Response(c, status)
}
