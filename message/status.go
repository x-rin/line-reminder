package message

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"os"
)

func GetStatus(c *gin.Context) {
	target := c.Param("id")
	status := os.Getenv(strings.ToUpper(target) + "_STATUS")

	if status == "failure" {
		PostReminder(c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": status,
		})
	}
}
