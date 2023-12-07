// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package share

import (
	"io"
	"net/http"
	"os"

	"github.com/pkg/errors"

	"github.com/wuhan005/share/internal/parser"
)

type Server struct {
	Method         string            `yaml:"method"`
	URL            string            `yaml:"url"`
	Headers        map[string]string `yaml:"headers"`
	Body           map[string]string `yaml:"body"`
	ResponseParser parser.Chain      `yaml:"response"`
}

const bodyFileFieldKey = "${file}"

func Share(server string, filePath string) (string, error) {
	servers.RLock()
	defer func() { servers.RUnlock() }()

	s, ok := servers.m[server]
	if !ok {
		return "", errors.Errorf("server %q not found", server)
	}
	requestHeader := parseHeader(s.Headers)

	// Open file.
	uploadedFile, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return "", errors.Wrap(err, "open uploaded file")
	}
	defer func() { _ = uploadedFile.Close() }()

	var body io.Reader
	switch requestHeader.Get("Content-Type") {
	case "multipart/form-data":
		requestBody, contentType, err := makeMultipartFormData(makeMultipartFormDataOptions{
			UploadedFile: uploadedFile,
			Body:         s.Body,
		})
		if err != nil {
			return "", errors.Wrap(err, "make multipart form data")
		}
		requestHeader.Set("Content-Type", contentType)
		body = requestBody

	case "application/x-www-form-urlencoded":
	case "application/json":
	}

	req, err := http.NewRequest(s.Method, s.URL, body)
	if err != nil {
		return "", errors.Wrap(err, "new request")
	}
	req.Header = requestHeader

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "do request")
	}
	defer func() { _ = resp.Body.Close() }()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "read response body")
	}

	output := string(responseBytes)
	for _, parser := range s.ResponseParser {
		output, err = parser.Do(output)
		if err != nil {
			return "", errors.Wrapf(err, "parser: %q", parser.Type)
		}
	}

	return output, nil
}

func parseHeader(h map[string]string) http.Header {
	header := make(http.Header)
	for k, v := range h {
		header.Set(k, v)
	}
	return header
}
