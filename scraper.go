package nuvi

import (
	"fmt"
	"io"
)

//go:generate counterfeiter . Downloader
//go:generate counterfeiter . Extractor
//go:generate counterfeiter . ArchiveWalker
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

// ArchiveWalkerFunc callback function
type ArchiveWalkerFunc func(io.Reader)

// ArchiveWalker unarchive zip archives
type ArchiveWalker interface {
	// Walk walks throu the content of io.Reader
	Walk(reader io.Reader, walker ArchiveWalkerFunc)
}

// Cacher caches any content
type Cacher interface {
	// Cache caches the content provided by the reader
	Cache(reader io.Reader)
}

// Scraper scrapes a web content
type Scraper struct {
	Downloader    Downloader
	Extractor     Extractor
	ArchiveWalker ArchiveWalker
	Cacher        Cacher
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

		scraper.ArchiveWalker.Walk(zipfile, func(file io.Reader) {
			scraper.Cacher.Cache(file)
		})

		zipfile.Close()
	}
	return nil
}
