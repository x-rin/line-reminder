package reminder

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func PostReminder(c *gin.Context) {
	id := c.PostForm("id")
	config := NewLineConfig()
	target, err := config.GetProfile(id)
	if err != nil {
		log.Fatal(err.Error())
	}

	rmdErr := config.PostMessage(target + ": " + os.Getenv("REMINDER_MESSAGE"))
	if rmdErr != nil {
		log.Fatal(rmdErr.Error())
	}

	status := SetStatus(id, "false")

	Response(c, status)
}
