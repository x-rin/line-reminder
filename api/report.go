package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strings"
)

func PostReport(c *gin.Context) {
	source := c.PostForm("id")
	err := PostMessage(source + ": " + os.Getenv("REPORT_MESSAGE"))
	if err != nil {
		log.Fatal(err.Error())
	}
	statusKey := strings.ToUpper(source) + "_STATUS"
	os.Setenv(statusKey, "true")
	c.JSON(http.StatusOK, gin.H{
		"status": os.Getenv(statusKey),
	})
}
