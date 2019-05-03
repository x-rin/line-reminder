package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/kutsuzawa/line-reminder/handler"
	"github.com/kutsuzawa/line-reminder/reminder"
	"github.com/kutsuzawa/line-reminder/scheduler"
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

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	channelToken, err := reminder.GetChannelToken(os.Getenv("CHANNEL_ID"), os.Getenv("CHANNEL_SECRET"))
	if err != nil {
		log.Fatal(err)
	}
	client, err := linebot.New(os.Getenv("CHANNEL_SECRET"), *channelToken)
	if err != nil {
		log.Fatal(err)
	}
	service := service.NewLineService(client)
	handler := handler.NewLineHandler(
		os.Getenv("GROUP_ID"),
		service,
		logger,
		os.Getenv("REMINDER_MESSAGE"),
		os.Getenv("REPORT_MESSAGE"),
		os.Getenv("REPLY_MESSAGE"),
		os.Getenv("CHECKED_MESSAGE"),
	)

	port := os.Getenv("PORT")
	go func() {
		if err := handler.Run(port); err != nil {
			log.Fatal(err)
		}
	}()

	checker := &scheduler.Checker{
		Message:  os.Getenv("CHECKED_MESSAGE"),
		GroupID:  os.Getenv("GROUP_ID"),
		Line:     service,
		Duration: 60 * time.Minute,
	}

	targets := strings.Split(os.Getenv("TARGET_IDS"), ",")
	if err := checker.Schedule(targets); err != nil {
		log.Fatal(err)
	}
}
