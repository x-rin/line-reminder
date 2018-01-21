package message

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func (m *message) PostReminder (c *gin.Context) {
	message := "Post Reminder"
	c.String(http.StatusOK, "%s", message)
}
