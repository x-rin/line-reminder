package reminder

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func PostReminder(c *gin.Context) {
	id := c.PostForm("id")
	target, err := GetProfile(id)
	if err != nil {
		log.Fatal(err.Error())
	}

	rmdErr := PostMessage(target + ": " + os.Getenv("REMINDER_MESSAGE"))
	if rmdErr != nil {
		log.Fatal(rmdErr.Error())
	}

	status := setStatus(id, "false")

	Response(c, status)
}
