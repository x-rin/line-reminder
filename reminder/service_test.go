package reminder_test

import (
	. "github.com/kutsuzawa/line-reminder/reminder"

	"github.com/golang/mock/gomock"
	"github.com/kutsuzawa/line-reminder/reminder/mock_reminder"
	"github.com/line/line-bot-sdk-go/linebot"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"os"
)

var _ = Describe("Service", func() {
	var (
		c           *gomock.Controller
		mockLineAPI *mock_reminder.MockLineAPI
		lineService LineService
	)

	BeforeEach(func() {
		c = gomock.NewController(GinkgoT())
		mockLineAPI = mock_reminder.NewMockLineAPI(c)
		lineService = NewLineService(mockLineAPI)
	})

	AfterEach(func() {
		c.Finish()
	})

	Describe("GetTargetName()", func() {
		Context("when correct id is passed", func() {
			response := &linebot.UserProfileResponse{
				UserID:      "correctID",
				DisplayName: "correctName",
			}
			BeforeEach(func() {
				mockLineAPI.EXPECT().GetProfile("correctID").Return(response, nil).Times(1)
			})
			It("targetDisplayName is returned, err is nil", func() {
				name, err := lineService.GetTargetName("correctID")
				Expect(err).To(BeNil())
				Expect(name).To(Equal("correctName"))
			})
		})
		Context("when failing to get profile from line api", func() {
			BeforeEach(func() {
				mockLineAPI.EXPECT().GetProfile("correctID").Return(nil, errors.New("some error")).Times(1)
			})
			It("targetDisplayName is empty, err is returned", func() {
				name, err := lineService.GetTargetName("correctID")
				Expect(err).NotTo(BeNil())
				Expect(name).To(BeEmpty())
			})
		})
	})
	Describe("Send()", func() {
		response := &linebot.BasicResponse{}
		Context("when correct message is passed", func() {
			BeforeEach(func() {
				os.Setenv("GROUP_ID", "hoge")
				mockLineAPI.EXPECT().PushMessage(os.Getenv("GROUP_ID"), linebot.NewTextMessage("hogehoge")).Return(response, nil).Times(1)
			})
			AfterEach(func() {
				os.Unsetenv("GROUP_ID")
			})
			It("err is nil", func() {
				err := lineService.Send("hogehoge")
				Expect(err).To(BeNil())
			})
		})
		Context("when failing to push message via line api", func() {
			BeforeEach(func() {
				os.Setenv("GROUP_ID", "hoge")
				mockLineAPI.EXPECT().PushMessage(os.Getenv("GROUP_ID"), linebot.NewTextMessage("hogehoge")).Return(nil, errors.New("some error")).Times(1)
			})
			It("err is returned", func() {
				err := lineService.Send("hogehoge")
				Expect(err).NotTo(BeNil())
			})
		})
	})
	//TODO: あとで
	Describe("Hear()", func() {

	})
	//TODO: あとで
	Describe("Reply()", func() {

	})
})
