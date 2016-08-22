package nuvi

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"strings"
)

// ZIPWalker unzip *.zip files
type ZIPWalker struct{}

// Unzip unzips a *.zip files
//
// The original implementation can be found in my blog
// http://blog.ralch.com/tutorial/golang-working-with-zip/
//
// ZIP algorithm is using random access so unforthunately
// we need to read the whole file before we unzip it
func (walker ZIPWalker) Walk(reader io.Reader, walk ArchiveWalkerFunc) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return
	}

	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return
	}

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		if !strings.HasSuffix(file.Name, ".xml") {
			continue
		}

		unzippedFile, err := file.Open()
		if err == nil {
			walk(unzippedFile)
			unzippedFile.Close()
		}
	}
}
