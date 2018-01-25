package main

import (
	"github.com/gin-gonic/gin"
	. "github.com/kutsuzawa/line-reminder/api"
	"log"
	"os"
)

func SetupRouter() *gin.Engine {
	router := gin.New()
	v1 := router.Group("/api/v1/")
	{
		v1.POST("reminder", PostReminder)
		v1.POST("report", PostReport)
		v1.POST("status", GetStatus)
		v1.POST("check", Check)
	}
	return router
}

func main() {

	router := SetupRouter()
	log.Printf("Start Go HTTP Server")

	port := os.Getenv("PORT")
	router.Run(":" + port)
}
