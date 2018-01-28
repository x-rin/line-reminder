package reminder_test

//import (
//	"github.com/golang/mock/gomock"
//	"github.com/kutsuzawa/line-reminder/mock_reminder"
//	. "github.com/kutsuzawa/line-reminder/reminder"
//	. "github.com/onsi/ginkgo"
//	//. "github.com/onsi/gomega"
//	"os"
//	"log"
//)
//
//var _ = Describe("Report", func() {
//	var (
//		c        *gomock.Controller
//		config   *LineConfig
//		mockLine *mock_reminder.MockLineService
//	)
//
//	BeforeEach(func() {
//		c = gomock.NewController(GinkgoT())
//		mockLine = mock_reminder.NewMockLineService(c)
//		config = &LineConfig{
//			AccessToken:   "hogehoge",
//			ExpiresIn:     1234,
//			TokenType:     "fugafuga",
//			ChannelSecret: "foofoo",
//		}
//	})
//
//	AfterEach(func() {
//		c.Finish()
//	})
//
//	Describe("PostReport", func() {
//		BeforeEach(func() {
//			mockLine.EXPECT().GetProfile("reportUser12").Return("test", nil).Times(1)
//			mockLine.EXPECT().PostMessage("test: test report").Return(nil).Times(1)
//		})
//		Context("hoge", func() {
//			os.Setenv("REPORT_MESSAGE", "test report")
//
//			It("fugafuga", func() {
//				status, err := config.PostReport("reportUser12")
//				log.Println(err.Error())
//				log.Println(status)
//			})
//		})
//	})
//})
