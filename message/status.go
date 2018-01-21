package message

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func (m *message) GetStatus(c *gin.Context) {
	target := c.Param("id")
	status := os.Getenv(strings.ToUpper(target) + "_STATUS")

	if status == "failure" {
		//TODO: PostReportが実装され次第パラメータをセットしてあげる。
		messageType := "hogehoge"
		text := "fugafuga"
		NewMessage(messageType, text).PostReminder(c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": status,
		})
	}
}
