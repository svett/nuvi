package nuvi_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/svett/nuvi"
	"github.com/svett/nuvi/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scraper", func() {
	var (
		page             *fakes.FakeReadCloser
		reportArchive    *fakes.FakeReadCloser
		photosArchive    *fakes.FakeReadCloser
		xmlReportArchive *fakes.FakeReadCloser
		xmlPhotosArchive *fakes.FakeReadCloser
		downloader       *fakes.FakeDownloader
		extractor        *fakes.FakeExtractor
		archiver         *fakes.FakeArchiver
		cacher           *fakes.FakeCacher

		scraper *nuvi.Scraper
	)

	BeforeEach(func() {
		page = &fakes.FakeReadCloser{}

		reportArchive = &fakes.FakeReadCloser{}
		photosArchive = &fakes.FakeReadCloser{}

		xmlReportArchive = &fakes.FakeReadCloser{}
		xmlPhotosArchive = &fakes.FakeReadCloser{}

		downloader = &fakes.FakeDownloader{}
		downloader.DownloadStub = func(url string) (io.ReadCloser, error) {
			if url == "www.example.com" {
				return page, nil
			} else if url == "www.example.com/error.zip" {
				return nil, fmt.Errorf("error.zip cannot be downloaded")
			} else if url == "www.example.com/report.zip" {
				return reportArchive, nil
			} else if url == "www.example.com/photos.zip" {
				return photosArchive, nil
			}

			return ioutil.NopCloser(strings.NewReader(url)), nil
		}

		extractor = &fakes.FakeExtractor{}
		extractor.ExtractReturns([]string{"report.zip", "photos.zip"})

		archiver = &fakes.FakeArchiver{}
		archiver.UnzipStub = func(reader io.Reader) (io.ReadCloser, error) {
			if reader == reportArchive {
				Expect(reportArchive.CloseCallCount()).To(Equal(0))
				return xmlReportArchive, nil
			} else if reader == photosArchive {
				Expect(photosArchive.CloseCallCount()).To(Equal(0))
				return xmlPhotosArchive, nil
			}
			return nil, fmt.Errorf("Oh no unzip error!")
		}

		cacher = &fakes.FakeCacher{}
		cacher.CacheStub = func(reader io.Reader) {
			if fakeReader, ok := reader.(*fakes.FakeReadCloser); ok {
				Expect(fakeReader.CloseCallCount()).To(Equal(0))
			}
		}

		scraper = &nuvi.Scraper{
			Downloader: downloader,
			Extractor:  extractor,
			Archiver:   archiver,
			Cacher:     cacher,
		}
	})

	It("downloads a web page for provided url", func() {
		Expect(scraper.Scrape("www.example.com")).To(Succeed())
		Expect(downloader.DownloadCallCount()).To(Equal(3))
		Expect(downloader.DownloadArgsForCall(0)).To(Equal("www.example.com"))
		Expect(page.CloseCallCount()).To(Equal(1))
	})

	Context("when Downloader cannot download a file", func() {
		BeforeEach(func() {
			downloader.DownloadReturns(nil, fmt.Errorf("Oh no!"))
		})

		It("returns the error", func() {
			Expect(scraper.Scrape("www.example.com")).To(MatchError("Oh no!"))
			Expect(page.CloseCallCount()).To(Equal(0))
		})
	})

	It("extracts the links from the downloaded page", func() {
		Expect(scraper.Scrape("www.example.com")).To(Succeed())
		Expect(extractor.ExtractCallCount()).To(Equal(1))
		Expect(extractor.ExtractArgsForCall(0)).To(Equal(page))
	})

	It("downloads the extracted files", func() {
		Expect(scraper.Scrape("www.example.com")).To(Succeed())
		Expect(downloader.DownloadCallCount()).To(Equal(3))
		Expect(downloader.DownloadArgsForCall(1)).To(Equal("www.example.com/report.zip"))
		Expect(downloader.DownloadArgsForCall(2)).To(Equal("www.example.com/photos.zip"))
	})

	Context("when download of extracted files fails", func() {
		BeforeEach(func() {
			extractor.ExtractReturns([]string{"report.zip", "error.zip", "photos.zip"})
		})

		It("continues to download the rest", func() {
			Expect(scraper.Scrape("www.example.com")).To(Succeed())
			Expect(downloader.DownloadCallCount()).To(Equal(4))
			Expect(downloader.DownloadArgsForCall(3)).To(Equal("www.example.com/photos.zip"))

			Expect(archiver.UnzipCallCount()).To(Equal(2))
			Expect(archiver.UnzipArgsForCall(0)).To(Equal(reportArchive))
			Expect(archiver.UnzipArgsForCall(1)).To(Equal(photosArchive))

			Expect(reportArchive.CloseCallCount()).To(Equal(1))
			Expect(photosArchive.CloseCallCount()).To(Equal(1))
		})
	})

	It("unzip the downloaded files", func() {
		Expect(scraper.Scrape("www.example.com")).To(Succeed())
		Expect(archiver.UnzipCallCount()).To(Equal(2))
		Expect(archiver.UnzipArgsForCall(0)).To(Equal(reportArchive))
		Expect(archiver.UnzipArgsForCall(1)).To(Equal(photosArchive))
	})

	Context("when the unzip fails", func() {
		BeforeEach(func() {
			archiver.UnzipStub = func(reader io.Reader) (io.ReadCloser, error) {
				if reader == photosArchive {
					return xmlPhotosArchive, nil
				}
				return nil, fmt.Errorf("Oh no unzip error!")
			}
			Expect(scraper.Scrape("www.example.com")).To(Succeed())
		})

		It("does not cache all files", func() {
			Expect(cacher.CacheCallCount()).To(Equal(1))
			Expect(cacher.CacheArgsForCall(0)).To(Equal(xmlPhotosArchive))
			Expect(xmlPhotosArchive.CloseCallCount()).To(Equal(1))
		})

		It("continues to unzip the other files", func() {
			Expect(archiver.UnzipCallCount()).To(Equal(2))
			Expect(archiver.UnzipArgsForCall(0)).To(Equal(reportArchive))
			Expect(archiver.UnzipArgsForCall(1)).To(Equal(photosArchive))
		})

		It("does close the file", func() {
			Expect(reportArchive.CloseCallCount()).To(Equal(1))
			Expect(photosArchive.CloseCallCount()).To(Equal(1))
		})
	})
})
