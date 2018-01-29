package line_reminder_test

import (
	. "github.com/kutsuzawa/line-reminder/line_reminder"

	"errors"
	"github.com/golang/mock/gomock"
	"github.com/kutsuzawa/line-reminder/line_reminder/mock_line_reminder"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
)

var _ = Describe("Reminder", func() {
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

	Describe("PostReminder", func() {
		Describe("normal pattern", func() {
			Context("success", func() {
				BeforeEach(func() {
					mockLineClient.EXPECT().GetProfile("reminderUser12").Return("test", nil).Times(1)
					os.Setenv("REMINDER_MESSAGE", "reminderMessage")
					mockLineClient.EXPECT().PostMessage("To test" + "\n" + os.Getenv("REMINDER_MESSAGE")).Return(nil).Times(1)
				})

				It("status is \"false\", err is nil", func() {
					status, err := lineReminder.PostReminder("reminderUser12")
					Expect(status).To(Equal("false"))
					Expect(err).To(BeNil())
				})
			})
		})
		Describe("exception pattern", func() {
			Context("failed to getProfile", func() {
				BeforeEach(func() {
					mockLineClient.EXPECT().GetProfile("reminderUser34").Return("", errors.New("some error")).Times(1)
					os.Setenv("CHECKED_MESSAGE", "reminderMessage")
				})
				It("status is empty, err is not nil", func() {
					status, err := lineReminder.PostReminder("reminderUser34")
					Expect(status).To(BeEmpty())
					Expect(err).ToNot(BeNil())
				})
			})
			Context("failed to postMessage", func() {
				BeforeEach(func() {
					mockLineClient.EXPECT().GetProfile("reminderUser56").Return("test", nil).Times(1)
					os.Setenv("CHECKED_MESSAGE", "reminderMessage")
					mockLineClient.EXPECT().PostMessage("To test" + "\n" + os.Getenv("REMINDER_MESSAGE")).Return(errors.New("some error")).Times(1)
				})
				It("status is empty, err is not nil", func() {
					status, err := lineReminder.PostReminder("reminderUser56")
					Expect(status).To(BeEmpty())
					Expect(err).ToNot(BeNil())
				})
			})
		})
	})
})
