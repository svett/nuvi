package nuvi

import (
	"fmt"
	"io"
	"sync"
)

//go:generate counterfeiter . Downloader
//go:generate counterfeiter . Extractor
//go:generate counterfeiter . ArchiveWalker
//go:generate counterfeiter . Cacher
//go:generate counterfeiter . Logger

// Logger logs messages
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

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
	MaxConn       int
	Logger        Logger
}

// Scrape scrapes a web page
func (scraper *Scraper) Scrape(url string) error {
	scraper.Logger.Printf("Downloading %s", url)
	reader, err := scraper.Downloader.Download(url)
	if err != nil {
		return err
	}
	defer reader.Close()

	scraper.Logger.Printf("Extracting %s", url)
	archives, err := scraper.Extractor.Extract(reader)
	if err != nil {
		return err
	}

	if scraper.MaxConn == 0 {
		scraper.MaxConn = 4
	}

	count := 1
	wg := &sync.WaitGroup{}

	for _, archive := range archives {
		count++
		if count > scraper.MaxConn {
			scraper.Logger.Printf("Waiting %d files to be downloaded", count)
			wg.Wait()
			count = 1
		}

		archiveURL := fmt.Sprintf("%s/%s", url, archive)
		wg.Add(1)
		go scraper.downloadAndCache(archiveURL, wg)
	}

	wg.Wait()
	return nil
}

func (scraper *Scraper) downloadAndCache(archiveURL string, wg *sync.WaitGroup) {
	defer wg.Done()

	scraper.Logger.Printf("Downloading %s", archiveURL)
	zipfile, err := scraper.Downloader.Download(archiveURL)
	if err != nil {
		scraper.Logger.Printf("Downloading %s failed with error %v", archiveURL, err)
		return
	}

	scraper.Logger.Printf("Browsing %s", archiveURL)
	scraper.ArchiveWalker.Walk(zipfile, func(file io.Reader) {
		scraper.Cacher.Cache(file)
	})

	zipfile.Close()
}
