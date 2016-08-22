package nuvi

import (
	"fmt"
	"io"
)

//go:generate counterfeiter . Downloader
//go:generate counterfeiter . Extractor
//go:generate counterfeiter . Archiver
//go:generate counterfeiter . Cacher

// Downloader downloads a content from URL
type Downloader interface {
	// Download downloads the content provided by url
	Download(url string) (io.ReadCloser, error)
}

// Extractor extracts a content of the page
type Extractor interface {
	// Extract extracts a links/anchors from a io.Reader
	Extract(reader io.Reader) ([]string, error)
}

// Archiver unarchive zip archives
type Archiver interface {
	// Unzip unzips the content of io.Reader
	Unzip(reader io.Reader) ([]io.ReadCloser, error)
}

// Cacher caches any content
type Cacher interface {
	// Cache caches the content provided by the reader
	Cache(reader io.Reader)
}

// Scraper scrapes a web content
type Scraper struct {
	Downloader Downloader
	Extractor  Extractor
	Archiver   Archiver
	Cacher     Cacher
}

// Scrape scrapes a web page
func (scraper *Scraper) Scrape(url string) error {
	reader, err := scraper.Downloader.Download(url)
	if err != nil {
		return err
	}
	defer reader.Close()

	archives, err := scraper.Extractor.Extract(reader)
	if err != nil {
		return err
	}

	for _, archive := range archives {
		archiveURL := fmt.Sprintf("%s/%s", url, archive)
		zipfile, err := scraper.Downloader.Download(archiveURL)
		if err != nil {
			continue
		}
		files, err := scraper.Archiver.Unzip(zipfile)
		if err == nil {
			for _, file := range files {
				scraper.Cacher.Cache(file)
				file.Close()
			}
		}
		zipfile.Close()
	}
	return nil
}
