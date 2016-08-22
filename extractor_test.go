package nuvi_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/svett/nuvi"
	"github.com/svett/nuvi/fakes"
)

var _ = Describe("LinkExtractor", func() {
	var (
		extractor nuvi.LinkExtractor
	)

	BeforeEach(func() {
		extractor = nuvi.LinkExtractor{}
	})

	It("extracts only the zip links", func() {
		reader, err := os.Open("assets/page.html")
		Expect(err).To(BeNil())

		links, err := extractor.Extract(reader)
		Expect(err).To(BeNil())
		Expect(links).To(HaveLen(3))
		Expect(links).To(ContainElement("file1.zip"))
		Expect(links).To(ContainElement("file2.zip"))
		Expect(links).To(ContainElement("file3.zip"))
	})

	Context("when the reader returns an error", func() {
		var reader *fakes.FakeReadCloser

		BeforeEach(func() {
			reader = &fakes.FakeReadCloser{}
			reader.ReadReturns(0, fmt.Errorf("Oh no!"))
		})

		It("returns the error", func() {
			links, err := extractor.Extract(reader)
			Expect(links).To(HaveLen(0))
			Expect(err).To(MatchError("Oh no!"))
		})
	})
})
