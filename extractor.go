package nuvi

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// LinkExtractor extract <a href="*.zip"> links
type LinkExtractor struct{}

// Extract extracts html anchor links
func (extractor LinkExtractor) Extract(reader io.Reader) ([]string, error) {
	var links []string
	tokenizer := html.NewTokenizer(reader)

	for {
		tag := tokenizer.Next()

		if tag == html.ErrorToken {
			if err := tokenizer.Err(); err != io.EOF {
				return []string{}, tokenizer.Err()
			}
			break
		}

		if tag != html.StartTagToken {
			continue
		}

		token := tokenizer.Token()

		if token.Data != "a" {
			continue
		}

		for _, attr := range token.Attr {
			if attr.Key == "href" && strings.HasSuffix(attr.Val, ".zip") {
				links = append(links, attr.Val)
			}
		}
	}

	return links, nil
}
