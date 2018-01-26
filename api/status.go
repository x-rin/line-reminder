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

	if !statusFlag {
		err := PostMessage(os.Getenv("STATUS_MESSAGE"))
		if err != nil {
			log.Println(err.Error())
		}
	}

	Response(c, status)
}
