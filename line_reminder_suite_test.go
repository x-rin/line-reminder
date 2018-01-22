package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLineReminder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LineReminder Suite")
}
