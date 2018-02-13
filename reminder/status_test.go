package reminder_test

import (
	. "github.com/kutsuzawa/line-reminder/reminder"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"strings"
)

var _ = Describe("Status", func() {
	Describe("GetStatus", func() {
		Describe("normal pattern", func() {
			Context("status env is true", func() {
				testId := "Test1234"
				testStatusKey := strings.ToUpper(testId) + "_STATUS"
				os.Setenv(testStatusKey, "true")
				It("flag is true, status is \"true\", err is nil", func() {
					acStFlag, acSt, acErr := GetStatus(testId)
					Expect(acStFlag).To(BeTrue())
					Expect(acSt).To(Equal("true"))
					Expect(acErr).To(BeNil())
				})
			})

			Context("status env is false", func() {
				testId := "Test5678"
				testStatusKey := strings.ToUpper(testId) + "_STATUS"
				os.Setenv(testStatusKey, "false")
				It("flag is false, status is \"false\", err is nil", func() {
					acStFlag, acSt, acErr := GetStatus(testId)
					Expect(acStFlag).To(BeFalse())
					Expect(acSt).To(Equal("false"))
					Expect(acErr).To(BeNil())
				})
			})
		})
		Describe("exception pattern", func() {
			Context("id is missing", func() {
				testId := "Test2345"
				It("flag is false, status is empty, err is not nil", func() {
					acStFlag, acSt, acErr := GetStatus(testId)
					Expect(acStFlag).To(BeFalse())
					Expect(acSt).To(BeEmpty())
					Expect(acErr).NotTo(BeNil())
				})
			})
		})
	})

	Describe("SetStatus", func() {
		Describe("normal pattern", func() {
			Context("set true", func() {
				testId := "1234Test"
				It("status is \"true\", err is nil", func() {
					acSt := SetStatus(testId, "true")
					Expect(acSt).To(Equal("true"))
				})
			})
			Context("set flase", func() {
				testId := "5678Test"
				It("status is \"false\", err is nil", func() {
					acSt := SetStatus(testId, "false")
					Expect(acSt).To(Equal("false"))
				})
			})
		})
	})
})
