package reminder_test

import (
	"net/http"
	"testing"

	"github.com/kutsuzawa/line-reminder/reminder"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
)

var (
	mockService *MockService
	groupID     string
	controller  *reminder.LineController
)

func TestMain(m *testing.M) {
	mockService = &MockService{}
	groupID = "groupID"
	controller = reminder.NewLineController(groupID, mockService)
}

func TestLineController_ReplyByWord(t *testing.T) {
	t.Helper()
	cases := []struct {
		name      string
		inputWord string
		expect    string
	}{
		{name: "matchedWord", inputWord: "matched", expect: "true"},
		{name: "unmatchedWord", inputWord: "unmatched", expect: "false"},
	}
	for _, c := range cases {
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

func TestLineController_Report(t *testing.T) {
	t.Helper()
	cases := []struct {
		name    string
		inputID string
		expect  string
	}{
		{name: "reportByExistUser", inputID: "existUserID", expect: "true"},
		{name: "reportByNewUser", inputID: "newUserID", expect: "true"},
	}
	for _, c := range cases {
		actual, err := controller.Report(c.inputID, "testMsg")
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

func (ms *MockService) GetNameByID(id string) (string, error) {
	if id == "existUserID" {
		return "existUserName", nil
	} else if id == "newUserID" {
		return "newUserName", nil
	} else {
		return "", errors.New("error!")
	}
}

func (ms *MockService) Send(groupID, message string) error {
	return nil
}

func (ms *MockService) Hear(request *http.Request) (linebot.Event, error) {
	return linebot.Event{
		ReplyToken: "replyToken",
		Type:       "type",
		Message:    linebot.NewTextMessage("matched"),
		Source: &linebot.EventSource{
			UserID: "userID",
		},
	}, nil
}

func (ms *MockService) Reply(replyToken string, message string) error {
	return nil
}
