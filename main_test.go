package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/kutsuzawa/line-reminder/factory"
	"github.com/kutsuzawa/line-reminder/handler"
	httpdoc "go.mercari.io/go-httpdoc"
	"go.uber.org/zap"
)

func TestJoin(t *testing.T) {
	t.Helper()
	base := "/api/v1"
	cases := []struct {
		name     string
		method   string
		endpoint string
		expect   int
	}{
		{name: "report", method: http.MethodPost, endpoint: base + "/report", expect: http.StatusOK},
		{name: "health", method: http.MethodGet, endpoint: base + "/health", expect: http.StatusOK},
	}

	for _, c := range cases {
		logger, _ := zap.NewProduction()
		serviceFactory := factory.NewServiceFactory(os.Getenv("CHANNEL_ID"), os.Getenv("CHANNEL_SECRET"))
		handler := handler.NewLineHandler(
			os.Getenv("GROUP_ID"),
			serviceFactory,
			logger,
			os.Getenv("REPORT_MESSAGE"),
			os.Getenv("REPLY_MESSAGE"),
		)
		api := &API{handler}
		router := http.NewServeMux()
		api.registHandler(router)

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
			mux.Handle(c.endpoint, httpdoc.Record(router, document, &httpdoc.RecordOption{
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
