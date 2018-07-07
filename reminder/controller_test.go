package reminder_test

import (
	"testing"
	"net/http"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/kutsuzawa/line-reminder/reminder"
)

func TestLineController_ReplyByWord(t *testing.T) {
	cases := []struct {
		name      string
		inputWord string
		expect    string
	}{
		{name: "matchedWord", inputWord: "matched", expect: "true"},
		{name: "unmatchedWord", inputWord: "unmatched", expect: "false"},
	}
	for _, c := range cases {
		mockService := &MockService{}
		groupID := "groupID"
		controller := reminder.NewLineController(groupID, mockService)
		req := &http.Request{}
		actual, err := controller.ReplyByWord(req, "sendMsg", c.inputWord)
		if err != nil {
			t.Error("err should not occur")
		}
		if actual != c.expect {
			t.Errorf("status should be %s, actual is %s", c.expect, actual)
		}
	}
}

type MockService struct {

}

func (ms *MockService) GetNameByID(id string) (string, error){
	return "hoge", nil
}

func (ms *MockService) Send(groupID, message string) error {
	return nil
}

func (ms *MockService) Hear(request *http.Request) (linebot.Event, error){
	return linebot.Event{
		ReplyToken: "replyToken",
		Type: "type",
		Message: linebot.NewTextMessage("matched"),
		Source: &linebot.EventSource{
			UserID: "userID",
		},
	}, nil
}

func (ms *MockService) Reply(replyToken string, message string) error {
	return nil
}
