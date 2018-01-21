package main

import (
	"log"
	"os"
	. "github.com/kutsuzawa/line-reminder/message"
	"github.com/gin-gonic/gin"
)

//リクエストに応じたfuncを定義
func main() {
	//routingはここ
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("messages", GetMessage)
		v1.POST("messages", PostMessage)
	}

	log.Printf("Start Go HTTP Server")

	port := os.Getenv("PORT")
	router.Run(":" + port)
}
