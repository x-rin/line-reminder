package message

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"os"
	"log"
)

func GetStatus(c *gin.Context) {
	target := c.Param("id")
	status := os.Getenv(strings.ToUpper(target) + "_STATUS")

	//if status == "failure" {
	//	//TODO: PostReportが実装され次第パラメータをセットしてあげる。
	//	messageType := "hogehoge"
	//	text := "fugafuga"
	//} else {
	//	c.JSON(http.StatusOK, gin.H{
	//		"status": status,
	//	})
	//}

	config, _ := GetAccessToken()
	log.Println(config)


	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}
