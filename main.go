package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"os"
	. "github.com/kutsuzawa/line-reminder/message"
)

//リクエストに応じたfuncを定義
func main() {
	//routingはここ
	//reference: https://godoc.org/github.com/julienschmidt/httprouter
	router := httprouter.New()
	router.GET("/api/v1/messages", GetMessage)
	router.POST("/api/v1/messages", PostMessage)

	log.Printf("Start Go HTTP Server")

	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(":" + port, router))
}
