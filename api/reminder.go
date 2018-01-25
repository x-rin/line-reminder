package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

func PostReminder(c *gin.Context) {
	target := c.PostForm("id")
	err := PostMessage(target + ": " + os.Getenv("REMINDER_MESSAGE"))
	if err != nil {
		log.Fatal(err.Error())
	}
	statusKey := strings.ToUpper(target) + "_STATUS"
	os.Setenv(statusKey, "false")
	c.JSON(http.StatusOK, gin.H{
		"status": os.Getenv(statusKey),
	})
}
