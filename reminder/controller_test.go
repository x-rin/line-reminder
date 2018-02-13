package reminder_test

import (
	. "github.com/kutsuzawa/line-reminder/reminder"

	"github.com/golang/mock/gomock"
	"github.com/kutsuzawa/line-reminder/reminder/mock_reminder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"os"
)

var _ = Describe("Controller", func() {
	var (
		c               *gomock.Controller
		mockLineService *mock_reminder.MockLineService
		lineController  LineController
	)

	BeforeEach(func() {
		c = gomock.NewController(GinkgoT())
		mockLineService = mock_reminder.NewMockLineService(c)
		lineController = NewLineController(mockLineService)
	})

	AfterEach(func() {
		c.Finish()
	})

	Describe("Check()", func() {
		Context("when TEST_USER_STATUS env is true", func() {
			BeforeEach(func() {
				os.Setenv("TEST_USER_STATUS", "true")
			})
			AfterEach(func() {
				os.Unsetenv("TEST_USER_STATUS")
			})
			It("status code is returned, err is nil", func() {
				status, err := lineController.Check("TEST_USER")
				Expect(err).To(BeNil())
				Expect(status).To(Equal("true"))
			})
		})
		Context("when TEST_USER_STATUS env is false and problem is nothing", func() {
			BeforeEach(func() {
				os.Setenv("TEST_USER_STATUS", "false")
				mockLineService.EXPECT().GetTargetName("test_user").Return("TEST_USER_NAME", nil).Times(1)
				os.Setenv("CHECKED_MESSAGE", "hoge")
				message := "To " + "TEST_USER_NAME" + "\n" + os.Getenv("CHECKED_MESSAGE")
				mockLineService.EXPECT().Send(message).Return(nil).Times(1)
			})
			AfterEach(func() {
				os.Unsetenv("TEST_USER_STATUS")
			})
			It("status code is returned, err is nil", func() {
				status, err := lineController.Check("test_user")
				Expect(err).To(BeNil())
				Expect(status).To(Equal("false"))
			})
		})
		Context("when failing to get profile from line api", func() {
			BeforeEach(func() {
				os.Setenv("TEST_USER_STATUS", "false")
				mockLineService.EXPECT().GetTargetName("test_user").Return("", errors.New("some error")).Times(1)
			})
			AfterEach(func() {
				os.Unsetenv("TEST_USER_STATUS")
			})
			It("status code is empty, err is returned", func() {
				status, err := lineController.Check("test_user")
				Expect(err).NotTo(BeNil())
				Expect(status).To(BeEmpty())
			})
		})
		Context("when failing to send message", func() {
			BeforeEach(func() {
				os.Setenv("TEST_USER_STATUS", "false")
				mockLineService.EXPECT().GetTargetName("test_user").Return("TEST_USER_NAME", nil).Times(1)
				os.Setenv("CHECKED_MESSAGE", "hoge")
				message := "To " + "TEST_USER_NAME" + "\n" + os.Getenv("CHECKED_MESSAGE")
				mockLineService.EXPECT().Send(message).Return(errors.New("some error")).Times(1)
			})
			AfterEach(func() {
				os.Unsetenv("TEST_USER_STATUS")
			})
			It("status code is empty, err is returned", func() {
				status, err := lineController.Check("test_user")
				Expect(err).NotTo(BeNil())
				Expect(status).To(BeEmpty())
			})
		})
	})
	//TODO: あとで
	Describe("Remind()", func() {

	})
	//TODO: あとで
	Describe("Report()", func() {

	})
	//TODO: あとで
	Describe("Reply()", func() {

	})
})
