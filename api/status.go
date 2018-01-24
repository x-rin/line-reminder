package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func GetStatus(c *gin.Context) {
	target := c.PostForm("id")
	statusKey := strings.ToUpper(target) + "_STATUS"
	status := os.Getenv(statusKey)
	statusFlag, _ := strconv.ParseBool(status)

	if statusFlag {
		c.JSON(http.StatusOK, gin.H{
			"status": status,
		})
	} else {
		PostReminder(c)
	}
}
