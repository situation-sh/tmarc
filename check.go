package main

import (
	"archive/zip"
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
)

// type for XML DMARC report
var dmarcType = filetype.NewType("dmarc", "application/dmarc")

func dmarcMatcher(data []byte) bool {
	buffer := bytes.NewBuffer(data)
	scanner := bufio.NewScanner(buffer)
	// skip first line (xml type)
	if !scanner.Scan() {
		return false
	}
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "<feedback>") {
			return true
		}
	}
	return false
}

func init() {
	// Register the new matcher and its type
	filetype.AddMatcher(dmarcType, dmarcMatcher)
}

func checkFile(path string) ([]byte, error) {
	header := make([]byte, 128)
	var reader io.Reader
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	t, err := filetype.MatchFile(path)
	if err != nil {
		return nil, err
	}

	switch t {
	case matchers.TypeGz:
		// gzip archive case
		gzipReader, err := gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
		if _, err := gzipReader.Read(header); err != nil {
			return nil, err
		}
		gzipReader.Close()
		// reopen that file
		file.Seek(0, 0)
		reader, _ = gzip.NewReader(file)
	case matchers.TypeZip:
		// zip archive case
		zipReader, err := zip.OpenReader(path)
		if err != nil {
			return nil, err
		}
		defer zipReader.Close()
		n := len(zipReader.File)
		if n != 1 {
			return nil, fmt.Errorf("the zip archive has not a single file (%d)", n)
		}

		fileReader, err := zipReader.File[0].Open()
		if err != nil {
			return nil, err
		}

		// just read that header
		if _, err := fileReader.Read(header); err != nil {
			return nil, err
		}
		fileReader.Close()
		// reopen the reader
		reader, _ = zipReader.File[0].Open()
	default:
		// assume raw case
		file.Seek(0, 0)
		if _, err := file.Read(header); err != nil {
			return nil, err
		}
		// reset the reader
		file.Seek(0, 0)
		reader = file

	}

	if filetype.IsType(header, dmarcType) {
		var buffer bytes.Buffer
		if _, err := io.Copy(&buffer, reader); err != nil {
			return nil, err
		} else {
			return buffer.Bytes(), nil
		}
		// if bytes, err := ioutil.ReadAll(reader); err != nil {
		// 	return nil, err
		// } else {
		// 	return bytes, nil
		// }
	}
	return nil, fmt.Errorf("the file is not a DMARC report")
}
