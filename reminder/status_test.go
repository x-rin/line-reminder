package reminder_test

import (
	"os"
	"testing"

	"github.com/kutsuzawa/line-reminder/reminder"
)

func TestGetStatus(t *testing.T) {
	t.Helper()
	cases := []struct {
		name   string
		input  string
		expect bool
	}{
		{name: "getExistIDStatus", input: "existID", expect: true},
		{name: "getDoesNotExistIDStatus", input: "doesNotExistID", expect: false},
	}
	os.Setenv("EXISTID_STATUS", "true")
	for i, c := range cases {
		actual, err := reminder.GetStatus(c.input)
		if i == 0 {
			if err != nil {
				t.Error("err should not occur")
			}
		} else {
			if err == nil {
				t.Error("err should occur")
			}
		}
		if actual != c.expect {
			t.Errorf("actual should be %v, actual %v", c.expect, actual)
		}
	}
}

func TestSetStatus(t *testing.T) {
	t.Helper()
	cases := []struct {
		name   string
		input  string
		expect bool
	}{
		{name: "setTrue", input: "true", expect: true},
		{name: "setFalse", input: "false", expect: false},
	}
	for _, c := range cases {
		actual, err := reminder.SetStatus("hoge", c.input)
		if err != nil {
			t.Error("err should not occur")
		}
		if actual != c.expect {
			t.Errorf("actual should be %v", c.expect)
		}
	}
	os.Unsetenv("HOGE_STATUS")
}
