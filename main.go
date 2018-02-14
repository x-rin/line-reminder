package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kutsuzawa/line-reminder/reminder"
	"log"
	"net/http"
	"os"
)

// SetupRouter - ルーターの初期化を行う
func SetupRouter() *gin.Engine {
	router := gin.New()
	v1 := router.Group("/api/v1/")
	{
		v1.POST("reminder", RemindCtr)
		v1.POST("report", ReportCtr)
		v1.POST("check", CheckCtr)
		v1.POST("webhook", ReplyCtr)
	}
	return router
}

// CheckCtr - ステータスチェックのリクエストを受け取った際のハンドラ
func CheckCtr(c *gin.Context) {
	controller := New()
	id := c.PostForm("id")
	status, err := controller.Check(id)
	if err != nil {
		Response(c, "", err)
	} else {
		Response(c, status, nil)
	}
}

// RemindCtr - リマインダーのリクエストを受け取った際のハンドラ
func RemindCtr(c *gin.Context) {
	controller := New()
	id := c.PostForm("id")
	status, err := controller.Remind(id)
	if err != nil {
		Response(c, "", err)
	} else {
		Response(c, status, nil)
	}
}

// ReportCtr - レポートのリクエストを受け取った際のハンドラ
func ReportCtr(c *gin.Context) {
	controller := New()
	id := c.PostForm("id")
	status, err := controller.Report(id)
	if err != nil {
		Response(c, "", err)
	} else {
		Response(c, status, nil)
	}
}

// ReplyCtr - Webhookを受け取った際のハンドラ
func ReplyCtr(c *gin.Context) {
	controller := New()
	status, err := controller.Reply(c.Request)
	if err != nil {
		Response(c, "", nil)
	} else {
		Response(c, status, nil)
	}
}

// Response - Responseを返す
func Response(c *gin.Context, status string, err error) {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": err.Error(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}

// New - Controllerを生成
func New() reminder.LineController {
	config, err := reminder.NewConfig()
	if err != nil {
		log.Fatalln("failing to create config: " + err.Error())
	}
	api, err := reminder.NewLineAPI(config)
	if err != nil {
		log.Fatalln("failing to create api: " + err.Error())
	}
	service := reminder.NewLineService(api)
	controller := reminder.NewLineController(service)
	return controller
}

func main() {
	router := SetupRouter()
	log.Printf("Start Go HTTP Server")

	port := os.Getenv("PORT")
	router.Run(":" + port)
}
