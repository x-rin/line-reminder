package reminder_test

import (
	. "github.com/kutsuzawa/line-reminder/reminder"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/kutsuzawa/line-reminder/mock_reminder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Report", func() {
	var (
		c        *gomock.Controller
		config   *LineConfig
		mockLine *mock_reminder.MockLineService
	)

	BeforeEach(func() {
		c = gomock.NewController(GinkgoT())
		mockLine = mock_reminder.NewMockLineService(c)
		config = &LineConfig{
			AccessToken:   "hogehoge",
			ExpiresIn:     1234,
			TokenType:     "fugafuga",
			ChannelSecret: "foofoo",
		}
	})

	AfterEach(func() {
		c.Finish()
	})

	Describe("PostReport", func() {
		BeforeEach(func() {
			mockLine.EXPECT().PostMessage("test report").Return(nil).Times(1)
		})
		Context("hoge", func() {
			os.Setenv("REPORT_MESSAGE", "test report")
			mockParam := gin.Param{
				Key:   "id",
				Value: "Test12345",
			}
			var mockParams []gin.Param
			mockParams = append(mockParams, mockParam)
			mockContext := gin.Context{
				Params: mockParams,
			}

			It("fugafuga", func() {
				PostReport(&mockContext)
			})
		})
	})
})
