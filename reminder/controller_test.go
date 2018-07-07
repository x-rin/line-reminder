package reminder_test

import (
	"net/http"
	"testing"

	"github.com/kutsuzawa/line-reminder/reminder"
	"github.com/line/line-bot-sdk-go/linebot"
	"os"
	"strings"
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
	run := m.Run()
	os.Exit(run)
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

func TestLineController_Check(t *testing.T) {
	t.Helper()
	cases := []struct {
		name    string
		inputID string
		expect  string
	}{
		{name: "checkForExistTrueUser", inputID: "existTrueUserID", expect: "true"},
		{name: "checkForExistFalseUser", inputID: "existFalseUserID", expect: "false"},
		{name: "checkForDoesNotExistUser", inputID: "doesNotExistUserID", expect: ""},
	}

	reminder.SetStatus("existTrueUserID", "true")
	reminder.SetStatus("existFalseUserID", "false")
	os.Setenv("EXISTFALSEUSERID_STATUS", "false")
	for _, c := range cases {
		mockService.SendCount = 0
		actual, _ := controller.Check(c.inputID, "testMsg")
		if actual != c.expect {
			t.Errorf("status should be %s, actual is %s", c.expect, actual)
		}
		if c.name == "checkForExistFalseUser" {
			if mockService.SendCount != 1 {
				t.Errorf("send Method should be called once, call is %d", mockService.SendCount)
			}
		} else {
			if mockService.SendCount != 0 {
				t.Errorf("send Method should not be called, call is %d", mockService.SendCount)
			}
		}
	}
}

type MockService struct {
	SendCount int
}

func (ms *MockService) GetNameByID(id string) (string, error) {
	name := strings.Replace(id, "ID", "Name", 1)
	return name, nil
}

func (ms *MockService) Send(groupID, message string) error {
	ms.SendCount++
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
