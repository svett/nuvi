package nuvi_test

import (
	"fmt"
	"net/http"

	"github.com/svett/nuvi"
	"github.com/svett/nuvi/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("HTTPDownloader", func() {
	It("returns the response body", func() {
		response := &http.Response{Body: &fakes.FakeReadCloser{}}
		downloader := nuvi.HTTPDownloader(func(url string) (*http.Response, error) {
			return response, nil
		})

		reader, err := downloader.Download("www.example.com")
		Expect(reader).To(Equal(response.Body))
		Expect(err).To(BeNil())
	})

	Context("when the dowloader fails", func() {
		It("returns the error", func() {
			downloader := nuvi.HTTPDownloader(func(url string) (*http.Response, error) {
				return nil, fmt.Errorf("Oh no!")
			})

			reader, err := downloader.Download("www.example.com")
			Expect(reader).To(BeNil())
			Expect(err).To(MatchError("Oh no!"))
		})
	})

})
