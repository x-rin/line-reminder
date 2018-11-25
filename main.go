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

func (h *handler) createNewController() (*reminder.LineController, error) {
	channelToken, err := reminder.GetChannelToken(h.channelID, h.channelSecret)
	if err != nil {
		h.logger.Error("failed to get channel token")
		return nil, err
	}
	client, err := linebot.New(h.channelSecret, *channelToken)
	if err != nil {
		h.logger.Error("failed to create line client")
		return nil, err
	}
	service := reminder.NewLineService(client)
	controller := reminder.NewLineController(h.groupID, service)
	return controller, nil
}

func (h *handler) do(action string, c *gin.Context) {
	controller, err := h.createNewController()
	if err != nil {
		return
	}
	id := c.PostForm("id")
	var status string
	var statusErr error
	switch action {
	case "check":
		status, statusErr = controller.Check(id, CheckedMessage)
	case "remind":
		status, statusErr = controller.Remind(id, ReminderMessage)
	case "report":
		status, statusErr = controller.Report(id, ReportMessage)
	case "reply":
		status, statusErr = controller.ReplyByWord(c.Request, ReplyMessage, ReportMessage)
	}
	if statusErr != nil {
		h.logger.Error("failed to "+action,
			zap.String("message", statusErr.Error()))
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	h.logger.Info("response returned",
		zap.String("status", status))
	c.JSON(http.StatusOK, nil)
	return
}

// Check - ステータスチェックのリクエストを受け取った際のハンドラ
func (h *handler) Check(c *gin.Context) {
	h.do("check", c)
	return
}

// Remind - リマインダーのリクエストを受け取った際のハンドラ
func (h *handler) Remind(c *gin.Context) {
	h.do("remind", c)
	return
}

// Report - レポートのリクエストを受け取った際のハンドラ
func (h *handler) Report(c *gin.Context) {
	h.do("report", c)
	return
}

// Reply - Webhookを受け取った際のハンドラ
func (h *handler) Reply(c *gin.Context) {
	h.do("reply", c)
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
