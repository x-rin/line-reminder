package main

import (
	"log"
	"os"
	. "github.com/kutsuzawa/line-reminder/message"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.POST("reminder", PostReminder)
		v1.POST("report", PostReport)
		v1.GET("status/:id", GetStatus)
	}
	return router
}

func main() {

	router := SetupRouter()
	log.Printf("Start Go HTTP Server")

	port := os.Getenv("PORT")
	router.Run(":" + port)
}
