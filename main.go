package main

import (
	"github.com/gin-gonic/gin"
	. "github.com/kutsuzawa/line-reminder/reminder"
	"log"
	"os"
)

func SetupRouter() *gin.Engine {
	router := gin.New()
	v1 := router.Group("/api/v1/")
	{
		v1.POST("reminder", PostReminder)
		v1.POST("report", PostReport)
		v1.POST("check", Check)
		v1.POST("webhook", GetWebHook)
	}
	return router
}

func main() {

	router := SetupRouter()
	log.Printf("Start Go HTTP Server")

	port := os.Getenv("PORT")
	router.Run(":" + port)
}
