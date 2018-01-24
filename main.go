package main

import (
	"github.com/gin-gonic/gin"
	. "github.com/kutsuzawa/line-reminder/message"
	"log"
	"os"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	messageType, text := "",""
	{
		v1.POST("reminder", NewMessage(messageType, text).PostReminder)
		v1.POST("report", NewMessage(messageType, text).PostReport)
		v1.GET("status/:id", NewMessage(messageType, text).GetStatus)
		v1.POST("check", NewMessage(messageType, text).Check)
	}
	return router
}

func main() {

	router := SetupRouter()
	log.Printf("Start Go HTTP Server")

	port := os.Getenv("PORT")
	router.Run(":" + port)
}
