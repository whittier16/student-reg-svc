package config_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/whittier16/student-reg-svc/internal/app/config"
)

var _ = Describe("Config", func() {
	It("LoadConfig should return error when file path is not present", func() {
		testCfgPath := "this-is-a-test-path"
		v, err := config.LoadConfig(testCfgPath, "yml")
		Expect(err).ShouldNot(BeNil())
		Expect(v).Should(BeNil())
	})
	It("LoadConfig should return success when file path is correct", func() {
		testCfgPath := "../../../configs/config-test"
		v, err := config.LoadConfig(testCfgPath, "yml")
		Expect(err).Should(BeNil())
		Expect(v).ShouldNot(BeNil())
	})
})
