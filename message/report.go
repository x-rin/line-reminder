package message

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostReport(c *gin.Context) {
	message := "Post Report!"
	c.String(http.StatusOK, "%s", message)
}
