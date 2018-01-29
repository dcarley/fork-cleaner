package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

var cmdPath string

var _ = BeforeSuite(func() {
	var err error
	cmdPath, err = gexec.Build("github.com/dcarley/fork-cleaner")
	Expect(err).ToNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})

func TestForkCleaner(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "ForkCleaner Suite")
}
