package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strings"
)

func PostReminder(c *gin.Context) {
	envKey := c.PostForm("id")
	target, err := GetProfile(envKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	rmdErr := PostMessage(target + ": " + os.Getenv("REMINDER_MESSAGE"))
	if rmdErr != nil {
		log.Fatal(rmdErr.Error())
	}

	statusKey := strings.ToUpper(envKey) + "_STATUS"
	os.Setenv(statusKey, "false")
	Response(c, os.Getenv(statusKey))
}
