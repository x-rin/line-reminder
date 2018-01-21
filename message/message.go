package message

import (
	"github.com/gin-gonic/gin"
)

type Message interface {
	PostReminder(c *gin.Context)
	PostReport(c *gin.Context)
	GetStatus(c *gin.Context)
	Check(c *gin.Context)
}

// Reference: https://developers.line.me/ja/docs/messaging-api/reference/#message-objects
type message struct {
	Type string
	Text string
}

func NewMessage(messageType string, text string) Message {
	return &message{
		Type: messageType,
		Text: text,
	}
}
