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
	router.GET("/api/v1/messages", GetMessage)
	router.POST("/api/v1/messages", PostMessage)

	log.Printf("Start Go HTTP Server")

	port := os.Getenv("PORT")
	router.Run(":" + port)
}
