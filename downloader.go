package nuvi

import (
	"io"
	"net/http"
)

// HTTPDownloader downloads
type HTTPDownloader func(string) (*http.Response, error)

func (downloader HTTPDownloader) Download(url string) (io.ReadCloser, error) {
	response, err := downloader(url)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}
