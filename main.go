package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kutsuzawa/line-reminder/handler"
	"github.com/kutsuzawa/line-reminder/handler/middleware"
	"github.com/kutsuzawa/line-reminder/reminder"
	"github.com/kutsuzawa/line-reminder/scheduler"
	"github.com/kutsuzawa/line-reminder/service"

	"github.com/line/line-bot-sdk-go/linebot"
	"go.uber.org/zap"
)

type API struct {
	handler *handler.LineHandler
}

func (a *API) registHandler(mux *http.ServeMux) {
	endpointPrefix := "/api/v1"
	mux.HandleFunc(endpointPrefix+"/report", middleware.GetID(a.handler.Report))
	mux.HandleFunc(endpointPrefix+"/webhook", middleware.GetID(a.handler.Reply))
	// for waking up heroku app
	mux.HandleFunc(endpointPrefix+"/health", a.handler.Health)
}

func (a *API) run(port string) error {
	mux := http.NewServeMux()
	a.registHandler(mux)
	return http.ListenAndServe(":"+port, mux)
}

type remindTimer struct {
	hours []string
}

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
		os.Getenv("REPORT_MESSAGE"),
		os.Getenv("REPLY_MESSAGE"),
	)
	api := &API{handler}

	port := os.Getenv("PORT")
	go func() {
		if err := api.run(port); err != nil {
			log.Fatal(err)
		}
	}()

	targets := strings.Split(os.Getenv("TARGET_IDS"), ",")

	reminder := &scheduler.Reminder{
		Message: os.Getenv("REMINDER_MESSAGE"),
		GroupID: os.Getenv("GROUP_ID"),
		Line:    service,
		Hours:   []string{"08:00", "19:00"},
	}

	go func() {
		if err := reminder.Schedule(targets); err != nil {
			log.Println(err)
		}
	}()

	checker := &scheduler.Checker{
		Message:  os.Getenv("CHECKED_MESSAGE"),
		GroupID:  os.Getenv("GROUP_ID"),
		Line:     service,
		Duration: 60 * time.Minute,
	}

	if err := checker.Schedule(targets); err != nil {
		log.Println(err)
	}
}
