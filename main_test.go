package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/kutsuzawa/line-reminder/handler"
	"github.com/kutsuzawa/line-reminder/reminder"
	"github.com/kutsuzawa/line-reminder/service"
	"github.com/line/line-bot-sdk-go/linebot"
	httpdoc "go.mercari.io/go-httpdoc"
	"go.uber.org/zap"
)

func TestJoin(t *testing.T) {
	t.Helper()
	base := "/api/v1"
	logger, _ := zap.NewProduction()
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
	mux := http.NewServeMux()
	api.registHandler(mux)

	cases := []struct {
		name     string
		method   string
		endpoint string
		expect   int
		handler  http.Handler
	}{
		{name: "report", method: http.MethodPost, endpoint: base + "/report", expect: http.StatusOK, handler: mux},
		{name: "health", method: http.MethodGet, endpoint: base + "/health", expect: http.StatusOK, handler: mux},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			name := fmt.Sprintf("line-reminder %s endpoint", c.name)
			document := &httpdoc.Document{
				Name: name,
				ExcludeHeaders: []string{
					"Content-Length",
					"Accept-Encoding",
					"User-Agent",
				},
			}
			defer func() {
				docName := fmt.Sprintf("doc/%s.md", c.name)
				if err := document.Generate(docName); err != nil {
					t.Fatalf("err :%s", err)
				}
			}()
			mux := http.NewServeMux()
			description := fmt.Sprintf("%s for a target user", c.name)
			mux.Handle(c.endpoint, httpdoc.Record(c.handler, document, &httpdoc.RecordOption{
				Description: description,
			}))

			testServer := httptest.NewServer(mux)
			defer testServer.Close()

			req := testNewRequest(t, c.method, testServer.URL+c.endpoint)
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("err: %s", err)
			}
			if res.StatusCode != c.expect {
				t.Fatalf("statusCode should be %d, actual is %d", c.expect, res.StatusCode)
			}

		})
	}
}

func testNewRequest(t *testing.T, method, urlStr string) *http.Request {
	var req *http.Request
	var err error
	switch method {
	case http.MethodPost:
		values := url.Values{}
		values.Set("id", os.Getenv("TEST_USER_ID"))
		req, err = http.NewRequest(method, urlStr, strings.NewReader(values.Encode()))
	case http.MethodGet:
		req, err = http.NewRequest(method, urlStr, nil)
	}
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}
