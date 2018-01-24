package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"os"
	"strconv"
	"log"
)

func GetStatus(c *gin.Context) {
	target := c.PostForm("id")
	statusKey := strings.ToUpper(target) + "_STATUS"
	log.Println("status: env key: " + statusKey)
	status := os.Getenv(statusKey)
	log.Printf("status: " + status)
	statusFlag, _ := strconv.ParseBool(status)

	if statusFlag {
		log.Println("STATUS TRUE")
		c.JSON(http.StatusOK, gin.H{
			"status": status,
		})
	} else {
		log.Println("STATS FALSE")
		PostReminder(c)
	}
}
