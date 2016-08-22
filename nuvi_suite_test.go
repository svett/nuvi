package nuvi_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNuvi(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nuvi Suite")
}
