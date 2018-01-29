package line_reminder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLineReminder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LineReminder Suite")
}
