package message

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

func GetStatus(c *gin.Context) {
	target := c.Param("id")
	status := os.Getenv(strings.ToUpper(target) + "_STATUS")

	if status == "failure" {
		//TODO: PostReportが実装され次第パラメータをセットしてあげる。
		PostReport(c)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": status,
		})
	}
}
