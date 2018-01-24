package message

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PostReminder (c *gin.Context) {
	message := "Post Reminder"
	c.String(http.StatusOK, "%s", message)
}
