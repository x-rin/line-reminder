package main

import (
	"fmt"
	"go.mercari.io/go-httpdoc"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

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

			logger, _ := zap.NewProduction()
			handler := &handler{
				logger:        logger,
				channelSecret: os.Getenv("CHANNEL_SECRET"),
				channelID:     os.Getenv("CHANNEL_ID"),
				groupID:       os.Getenv("GROUP_ID"),
			}
			router := handler.SetupRouter()
			mux := http.NewServeMux()
			description := fmt.Sprintf("%s for a target user", c.name)
			mux.Handle(c.endpoint, httpdoc.Record(router, document, &httpdoc.RecordOption{
				Description: description,
			}))

			testServer := httptest.NewServer(mux)
			defer testServer.Close()

			req := testNewRequest(t, testServer.URL+c.endpoint)
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

func testNewRequest(t *testing.T, urlStr string) *http.Request {
	values := url.Values{}
	values.Set("id", os.Getenv("TEST_USER_ID"))
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(values.Encode()))
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}
