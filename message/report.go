package message

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func (m *message) PostReport (c *gin.Context) {
	message := "Post Report!"
	c.String(http.StatusOK, "%s", message)
}
