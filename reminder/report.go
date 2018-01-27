package reminder

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func PostReport(c *gin.Context) {
	id := c.PostForm("id")
	source, err := GetProfile(id)
	if err != nil {
		log.Fatal(err.Error())
	}

	reportErr := PostMessage(source + ": " + os.Getenv("REPORT_MESSAGE"))
	if reportErr != nil {
		log.Fatal(reportErr.Error())
	}

	status := SetStatus(id, "true")

	replyErr := PostMessage(os.Getenv("REPLY_SUCCESS"))
	if replyErr != nil {
		log.Fatal(reportErr.Error())
	}

	Response(c, status)
}
