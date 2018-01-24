package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"log"
	"strings"
)

func PostReminder (c *gin.Context) {
	source := c.PostForm("id")
	err := PostMessage(source+": "+os.Getenv("REMINDER_MESSAGE"))
	if err != nil {
		log.Fatal(err.Error())
	}
	statusKey := strings.ToUpper(source) + "_STATUS"
	os.Setenv(statusKey, "false")
	c.JSON(http.StatusOK, gin.H{
		"status": os.Getenv(statusKey),
	})
}
