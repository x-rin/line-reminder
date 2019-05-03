package util_test

import (
	"os"
	"testing"

	"github.com/kutsuzawa/line-reminder/util"
)

func TestGetStatus(t *testing.T) {
	t.Helper()
	cases := []struct {
		name   string
		input  string
		expect bool
	}{
		{name: "getExistIDStatus", input: "existID", expect: true},
		{name: "getDoesNotExistIDStatus", input: "doesNotExistID", expect: true},
	}
	os.Setenv("EXISTID_STATUS", "true")
	for _, c := range cases {
		actual, err := util.GetStatus(c.input)
		if err != nil {
			t.Error("err should not occur")
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
		actual, err := util.SetStatus("hoge", c.input)
		if err != nil {
			t.Error("err should not occur")
		}
		if actual != c.expect {
			t.Errorf("actual should be %v", c.expect)
		}
	}
	os.Unsetenv("HOGE_STATUS")
}
