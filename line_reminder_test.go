package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/kutsuzawa/line-reminder"
	"os"
	"net/http/httptest"
	"net/http"
	"github.com/gin-gonic/gin"
)

var _ = Describe("LineReminder", func() {

	Describe("GET: /api/v1/status", func() {

		Context("status is success", func() {
			router, w, req := StatusRouterInit()
			os.Setenv("TEST_STATUS", "success")
			router.ServeHTTP(w, req)
			It("should return success status", func() {
				Expect(w.Code).To(Equal(200))
				Expect(w.Body.String()).To(Equal("{\"status\":\"success\"}"))
			})
		})

		Context("status is failure", func() {
			router, w, req := StatusRouterInit()
			os.Setenv("TEST_STATUS", "failure")
			router.ServeHTTP(w, req)
			It("should access POST /api/v1/report ", func() {
				Expect(w.Code).To(Equal(200))
				//TODO: PostReportの実装次第変更を加えること
				Expect(w.Body.String()).To(Equal("Post Reminder"))
			})
		})
	})
})

func StatusRouterInit () (*gin.Engine, *httptest.ResponseRecorder, *http.Request){
	router := SetupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/status/test", nil)
	return router, w, req
}
