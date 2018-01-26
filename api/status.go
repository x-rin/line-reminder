package api

import (
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"strconv"
	"strings"
)

func GetStatus(c *gin.Context) {
	envKey := c.PostForm("id")
	target, err := GetProfile(envKey)
	if err != nil {
		log.Fatal(err.Error())
	}

	statusKey := strings.ToUpper(envKey) + "_STATUS"
	status := os.Getenv(statusKey)
	statusFlag, _ := strconv.ParseBool(status)

	if !statusFlag {
		err := PostMessage(target + ": " + os.Getenv("STATUS_MESSAGE"))
		if err != nil {
			log.Println(err.Error())
		}
	}

	Response(c, status)
}
