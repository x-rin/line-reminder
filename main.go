package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kataras/iris"
	"github.com/kutsuzawa/line-reminder/handler"
	"github.com/kutsuzawa/line-reminder/reminder"
	"github.com/kutsuzawa/line-reminder/service"

	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

var (
	ReminderMessage = os.Getenv("REMINDER_MESSAGE")
	ReportMessage   = os.Getenv("REPORT_MESSAGE")
	ReplyMessage    = os.Getenv("REPLY_MESSAGE")
	CheckedMessage  = os.Getenv("CHECKED_MESSAGE")
)

type Handler struct {
	channelID     string
	channelSecret string
	groupID       string
	logger        *zap.Logger
}

// SetupRouter - ルーターの初期化を行う
func (h *Handler) SetupRouter() *iris.Application {
	router := iris.Default()
	v1 := router.Party("/api/v1")
	{
		v1.Post("/reminder", h.Remind)
		v1.Post("/report", h.Report)
		v1.Post("/check", h.Check)
		v1.Post("/webhook", h.Reply)
	}
	return router
}

func (h *Handler) createNewController() (*handler.LineController, error) {
	//channelToken, err := reminder.GetChannelToken(h.channelID, h.channelSecret)
	//if err != nil {
	//	h.logger.Error("failed to get channel token")
	//	return nil, err
	//}
	//client, err := linebot.New(h.channelSecret, *channelToken)
	//if err != nil {
	//	h.logger.Error("failed to create line client")
	//	return nil, err
	//}
	//service := service.NewLineService(client)
	//controller := handler.NewLineController(h.groupID, service)
	//return controller, nil
	return nil, nil
}

func (h *Handler) do(action string, ctx iris.Context) {
	controller, err := h.createNewController()
	if err != nil {
		return
	}
	id := ctx.FormValue("id")
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
		status, statusErr = controller.ReplyByWord(ctx.Request(), ReplyMessage, ReportMessage)
	}
	if statusErr != nil {
		h.logger.Error("failed to "+action,
			zap.String("message", statusErr.Error()))
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	h.logger.Info("response returned",
		zap.String("status", status))
	ctx.StatusCode(http.StatusOK)
	return
}

// Check - ステータスチェックのリクエストを受け取った際のハンドラ
func (h *Handler) Check(ctx iris.Context) {
	h.do("check", ctx)
	return
}

// Remind - リマインダーのリクエストを受け取った際のハンドラ
func (h *Handler) Remind(ctx iris.Context) {
	h.do("remind", ctx)
	return
}

// Report - レポートのリクエストを受け取った際のハンドラ
func (h *Handler) Report(ctx iris.Context) {
	h.do("report", ctx)
	return
}

// Reply - Webhookを受け取った際のハンドラ
func (h *Handler) Reply(ctx iris.Context) {
	h.do("reply", ctx)
	return
}

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	//handler := &Handler{
	//	logger:        logger,
	//	channelSecret: os.Getenv("CHANNEL_SECRET"),
	//	channelID:     os.Getenv("CHANNEL_ID"),
	//	groupID:       os.Getenv("GROUP_ID"),
	//}
	fmt.Println(os.Getenv("CHANNEL_ID"))
	fmt.Println(os.Getenv("GROUP_ID"))
	fmt.Println(os.Getenv("CHANNEL_SECRET"))

	channelToken, err := reminder.GetChannelToken(os.Getenv("CHANNEL_ID"), os.Getenv("CHANNEL_SECRET"))
	if err != nil {
		logger.Error("failed to get channel token")
		log.Fatal(err)
	}
	client, err := linebot.New(os.Getenv("CHANNEL_SECRET"), *channelToken)
	if err != nil {
		logger.Error("failed to create line client")
		log.Fatal(err)
	}
	service := service.NewLineService(client)
	handler := handler.NewLineController(
		os.Getenv("GROUP_ID"),
		service,
		logger,
		os.Getenv("REMINDER_MESSAGE"),
		os.Getenv("REPORT_MESSAGE"),
		os.Getenv("REPLY_MESSAGE"),
		os.Getenv("CHECKED_MESSAGE"),
	)

	//router := handler.SetupRouter()

	port := os.Getenv("PORT")
	handler.Run(port)
	//router.Run(iris.Addr(fmt.Sprintf(":%s", port)))
}
