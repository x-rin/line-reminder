package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/kutsuzawa/line-reminder"
	"os"
	"net/http/httptest"
	"net/http"
	"github.com/gin-gonic/gin"
	"io"
)

var _ = Describe("LineReminder", func() {

	Describe("GET: /api/v1/status", func() {

		Context("status is true", func() {
			router, w, req := RouterInit("GET", "/api/v1/status/test", nil)
			os.Setenv("TEST_STATUS", "true")
			router.ServeHTTP(w, req)
			It("should return true status", func() {
				Expect(w.Code).To(Equal(200))
				Expect(w.Body.String()).To(Equal("{\"status\":\"true\"}"))
			})
		})

		Context("status is false", func() {
			router, w, req := RouterInit("GET", "/api/v1/status/test", nil)
			os.Setenv("TEST_STATUS", "false")
			router.ServeHTTP(w, req)
			It("should access POST /api/v1/report ", func() {
				Expect(w.Code).To(Equal(200))
				//TODO: PostReportの実装次第変更を加えること
				Expect(w.Body.String()).To(Equal("Post Reminder"))
			})
		})
	})
})

func RouterInit (method string, url string, body io.Reader) (*gin.Engine, *httptest.ResponseRecorder, *http.Request){
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body)
	return router, w, req
}
