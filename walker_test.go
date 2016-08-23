package nuvi_test

import (
	"fmt"
	"io"
	"os"

	"github.com/svett/nuvi"
	"github.com/svett/nuvi/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ZIPWalker", func() {
	var (
		fileExt string
		walker  *nuvi.ZIPWalker
	)

	BeforeEach(func() {
		fileExt = ".xml"
	})

	JustBeforeEach(func() {
		walker = &nuvi.ZIPWalker{
			FileExt: fileExt,
			Logger:  &fakes.FakeLogger{},
		}
	})

	It("walks through an XML files", func() {
		reader, err := os.Open("assets/info.zip")
		Expect(err).NotTo(HaveOccurred())

		count := 0
		walker.Walk(reader, func(file io.Reader) {
			count++
		})

		Expect(count).To(Equal(3))
	})

	Context("when txt file is extracted", func() {
		BeforeEach(func() {
			fileExt = ".txt"
		})
		It("walks through an TXT files", func() {
			reader, err := os.Open("assets/info.zip")
			Expect(err).NotTo(HaveOccurred())

			count := 0
			walker.Walk(reader, func(file io.Reader) {
				count++
			})

			Expect(count).To(Equal(1))
		})

	})

	Context("when the zip archive is corrupted", func() {
		It("returns the error", func() {
			count := 0
			reader := &fakes.FakeReadCloser{}
			reader.ReadReturns(0, fmt.Errorf("Oh no!"))
			walker.Walk(reader, func(file io.Reader) {
				count++
			})

			Expect(count).To(Equal(0))
		})
	})

})
