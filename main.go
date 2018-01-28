package main

import (
	"github.com/gin-gonic/gin"
	. "github.com/kutsuzawa/line-reminder/reminder"
	"log"
	"net/http"
	"os"
)

func SetupRouter() *gin.Engine {
	router := gin.New()
	v1 := router.Group("/api/v1/")
	{
		v1.POST("reminder", PostReminderCtr)
		v1.POST("report", PostReportCtr)
		v1.POST("check", CheckCtr)
		v1.POST("webhook", GetWebHookCtr)
	}
	return router
}

func CheckCtr(c *gin.Context) {
	id := c.PostForm("id")
	con := NewLineConfig()
	status, err := con.Check(id)
	if err != nil {
		Response(c, "", err)
	} else {
		Response(c, status, nil)
	}
}

func PostReminderCtr(c *gin.Context) {
	id := c.PostForm("id")
	con := NewLineConfig()
	status, err := con.PostReminder(id)
	if err != nil {
		Response(c, "", err)
	} else {
		Response(c, status, nil)
	}
}

func PostReportCtr(c *gin.Context) {
	id := c.PostForm("id")
	con := NewLineConfig()
	status, err := con.PostReport(id)
	if err != nil {
		Response(c, "", err)
	} else {
		Response(c, status, nil)
	}
}

func GetWebHookCtr(c *gin.Context) {
	con := NewLineConfig()
	status, err := con.GetWebHook(c.Request)
	if err != nil {
		Response(c, "", nil)
	} else {
		Response(c, status, nil)
	}
}

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

func main() {

	router := SetupRouter()
	log.Printf("Start Go HTTP Server")

	port := os.Getenv("PORT")
	router.Run(":" + port)
}
