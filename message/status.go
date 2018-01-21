package message

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetStatus (c *gin.Context) {
	message := "Get Status!"
	c.String(http.StatusOK, "%s", message)
}
