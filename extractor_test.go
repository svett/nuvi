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
		fileExt   string
		extractor *nuvi.LinkExtractor
	)

	BeforeEach(func() {
		fileExt = ".zip"
	})

	JustBeforeEach(func() {
		extractor = &nuvi.LinkExtractor{
			FileExt: fileExt,
		}
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

	Context("when the file extension is .sh", func() {
		BeforeEach(func() {
			fileExt = ".sh"
		})

		It("extractor only the script files", func() {
			reader, err := os.Open("assets/page.html")
			Expect(err).To(BeNil())

			links, err := extractor.Extract(reader)
			Expect(err).To(BeNil())
			Expect(links).To(HaveLen(1))
			Expect(links).To(ContainElement("script.sh"))
		})
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
