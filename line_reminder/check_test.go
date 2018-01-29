package line_reminder_test

import (
	. "github.com/kutsuzawa/line-reminder/line_reminder"

	"github.com/golang/mock/gomock"
	"github.com/kutsuzawa/line-reminder/line_reminder/mock_line_reminder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"os"
)

var _ = Describe("Check", func() {
	var (
		c              *gomock.Controller
		mockLineClient *mock_line_reminder.MockLineClient
		lineReminder   LineReminder
	)

	BeforeEach(func() {
		c = gomock.NewController(GinkgoT())
		mockLineClient = mock_line_reminder.NewMockLineClient(c)
		lineReminder = NewLineReminder(mockLineClient)
	})

	AfterEach(func() {
		c.Finish()
	})

	Describe("Check", func() {
		Describe("normal pattern", func() {
			Context("done report", func() {
				os.Setenv("CHECKUSER12_STATUS", "true")
				It("status is \"true\", err is nil", func() {
					status, err := lineReminder.Check("checkUser12")
					Expect(status).To(Equal("true"))
					Expect(err).To(BeNil())
				})
			})
			Context("Not done report", func() {
				BeforeEach(func() {
					mockLineClient.EXPECT().GetProfile("checkUser34").Return("test", nil).Times(1)
					os.Setenv("CHECKED_MESSAGE", "checkedMessage")
					mockLineClient.EXPECT().PostMessage("To test\n" + os.Getenv("CHECKED_MESSAGE")).Return(nil).Times(1)
				})
				os.Setenv("CHECKUSER34_STATUS", "false")
				It("status is \"false\", err is nil", func() {
					status, err := lineReminder.Check("checkUser34")
					Expect(status).To(Equal("false"))
					Expect(err).To(BeNil())
				})
			})
		})
		Describe("exception pattern", func() {
			Context("user id is missing", func() {
				It("status is empty, err is not nil", func() {
					status, err := lineReminder.Check("checkUser56")
					Expect(status).To(BeEmpty())
					Expect(err).ToNot(BeNil())
				})
			})
			Context("failed to postMessage", func() {
				BeforeEach(func() {
					mockLineClient.EXPECT().GetProfile("checkUser78").Return("test", nil).Times(1)
					os.Setenv("CHECKED_MESSAGE", "checkedMessage")
					mockLineClient.EXPECT().PostMessage("To test\n" + os.Getenv("CHECKED_MESSAGE")).Return(errors.New("some error")).Times(1)
				})
				os.Setenv("CHECKUSER78_STATUS", "false")
				It("status is empty, err is not nil", func() {
					status, err := lineReminder.Check("checkUser78")
					Expect(status).To(BeEmpty())
					Expect(err).ToNot(BeNil())
				})
			})
			Context("failed to getProfile", func() {
				BeforeEach(func() {
					mockLineClient.EXPECT().GetProfile("checkUser91").Return("", errors.New("some error")).Times(1)
					os.Setenv("CHECKED_MESSAGE", "checkedMessage")
				})
				os.Setenv("CHECKUSER91_STATUS", "false")
				It("status is empty, err is not nil", func() {
					status, err := lineReminder.Check("checkUser91")
					Expect(status).To(BeEmpty())
					Expect(err).ToNot(BeNil())
				})
			})
		})
	})
})
