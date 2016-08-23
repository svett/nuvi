package integration_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"github.com/svett/nuvi/integration/utils"

	"testing"
)

var (
	server    *utils.RedisRunner
	scrapeBin string
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var _ = BeforeSuite(func() {
	var err error
	server = &utils.RedisRunner{}
	scrapeBin, err = gexec.Build("github.com/svett/nuvi/cmd/nuvi")
	Expect(err).NotTo(HaveOccurred())
	Expect(server.Start()).To(Succeed())
})

func runScraper(args ...string) (*gexec.Session, error) {
	cmd := exec.Command(scrapeBin, args...)
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter
	return gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
}

var _ = AfterSuite(func() {
	server.Stop()
	gexec.CleanupBuildArtifacts()
})
