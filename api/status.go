package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"os"
	"strconv"
)

func GetStatus(c *gin.Context) {
	target := c.Param("id")
	status := os.Getenv(strings.ToUpper(target) + "_STATUS")
	statusFlag, _ := strconv.ParseBool(status)

	if statusFlag {
		c.JSON(http.StatusOK, gin.H{
			"status": status,
		})
	} else {
		PostReminder(c)
	}
}
