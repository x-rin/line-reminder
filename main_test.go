package main

import (
	"net/http"
	"os"
	"testing"

	"github.com/kataras/iris/httptest"
	"go.uber.org/zap"
)

func TestHandler_WithoutReply(t *testing.T) {
	t.Helper()
	cases := []struct {
		name     string
		endpoint string
		expect   int
	}{
		{name: "remind", endpoint: "/api/v1/reminder", expect: http.StatusOK},
		{name: "check", endpoint: "/api/v1/check", expect: http.StatusOK},
		{name: "report", endpoint: "/api/v1/report", expect: http.StatusOK},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			logger, _ := zap.NewProduction()
			handler := &handler{
				logger:        logger,
				channelSecret: os.Getenv("CHANNEL_SECRET"),
				channelID:     os.Getenv("CHANNEL_ID"),
				groupID:       os.Getenv("GROUP_ID"),
			}
			router := handler.SetupRouter()

			testServer := httptest.New(t, router)
			testServer.POST(c.endpoint).WithHeader("Content-Type", "application/x-www-form-urlencoded").
				WithFormField("id", os.Getenv("TEST_USER_ID")).Expect().Status(http.StatusOK)
		})
	}
}
