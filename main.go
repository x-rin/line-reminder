package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kutsuzawa/line-reminder/reminder"
	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

var (
	ReminderMessage = os.Getenv("REMINDER_MESSAGE")
	ReportMessage   = os.Getenv("REPORT_MESSAGE")
	ReplyMessage    = os.Getenv("REPLY_MESSAGE")
	CheckedMessage  = os.Getenv("CHECKED_MESSAGE")
)

type handler struct {
	channelID     string
	channelSecret string
	groupID       string
	logger        *zap.Logger
}

// SetupRouter - ルーターの初期化を行う
func (h *handler) SetupRouter() *gin.Engine {
	router := gin.New()
	v1 := router.Group("/api/v1/")
	{
		v1.POST("reminder", h.Remind)
		v1.POST("report", h.Report)
		v1.POST("check", h.Check)
		v1.POST("webhook", h.Reply)
	}
	return router
}

// Check - ステータスチェックのリクエストを受け取った際のハンドラ
func (h *handler) Check(c *gin.Context) {
	channelToken, err := reminder.GetChannelToken(h.channelID, h.channelSecret)
	if err != nil {
		h.logger.Error("failed to get channel token")
		return
	}
	client, err := linebot.New(h.channelSecret, *channelToken)
	if err != nil {
		h.logger.Error("failed to create line client")
	}
	service := reminder.NewLineService(client)
	controller := reminder.NewLineController(h.groupID, service)
	id := c.PostForm("id")
	status, err := controller.Check(id, CheckedMessage)
	if err != nil {
		h.logger.Error("failed to check",
			zap.String("message", err.Error()))
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	h.logger.Info("response returned",
		zap.String("status", status))
	c.JSON(http.StatusOK, nil)
	return
}

// Remind - リマインダーのリクエストを受け取った際のハンドラ
func (h *handler) Remind(c *gin.Context) {
	channelToken, err := reminder.GetChannelToken(h.channelID, h.channelSecret)
	if err != nil {
		h.logger.Error("failed to get channel token")
		return
	}
	client, err := linebot.New(h.channelSecret, *channelToken)
	if err != nil {
		h.logger.Error("failed to create line client")
		return
	}
	service := reminder.NewLineService(client)
	controller := reminder.NewLineController(h.groupID, service)
	id := c.PostForm("id")
	status, err := controller.Remind(id, ReminderMessage)
	if err != nil {
		h.logger.Error("failed to remind",
			zap.String("message", err.Error()))
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	h.logger.Info("response returned",
		zap.String("status", status))
	c.JSON(http.StatusOK, nil)
	return
}

// Report - レポートのリクエストを受け取った際のハンドラ
func (h *handler) Report(c *gin.Context) {
	channelToken, err := reminder.GetChannelToken(h.channelID, h.channelSecret)
	if err != nil {
		h.logger.Error("failed to get channel token")
		return
	}
	client, err := linebot.New(h.channelSecret, *channelToken)
	if err != nil {
		h.logger.Error("failed to create line client")
		return
	}
	service := reminder.NewLineService(client)
	controller := reminder.NewLineController(h.groupID, service)
	id := c.PostForm("id")
	status, err := controller.Report(id, ReportMessage)
	if err != nil {
		h.logger.Error("failed to report",
			zap.String("message", err.Error()))
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	h.logger.Info("response returned",
		zap.String("status", status))
	c.JSON(http.StatusOK, nil)
	return
}

// Reply - Webhookを受け取った際のハンドラ
func (h *handler) Reply(c *gin.Context) {
	channelToken, err := reminder.GetChannelToken(h.channelID, h.channelSecret)
	if err != nil {
		h.logger.Error("failed to get channel token")
		return
	}
	client, err := linebot.New(h.channelSecret, *channelToken)
	h.logger.Info("get token",
		zap.String("token", *channelToken))
	if err != nil {
		h.logger.Error("failed to create line client")
		return
	}
	service := reminder.NewLineService(client)
	controller := reminder.NewLineController(h.groupID, service)
	status, err := controller.ReplyByWord(c.Request, ReplyMessage, ReportMessage)
	if err != nil {
		h.logger.Error("failed to reply",
			zap.String("message", err.Error()))
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	h.logger.Info("response returned",
		zap.String("status", status))
	c.JSON(http.StatusOK, nil)
	return
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	handler := &handler{
		logger:        logger,
		channelSecret: os.Getenv("CHANNEL_SECRET"),
		channelID:     os.Getenv("CHANNEL_ID"),
		groupID:       os.Getenv("GROUP_ID"),
	}
	router := handler.SetupRouter()

	port := os.Getenv("PORT")
	router.Run(":" + port)
}
