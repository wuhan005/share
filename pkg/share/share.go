// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package share

import (
	"io"
	"math/rand"
	"net/http"
	"os"

	"github.com/imroc/req/v3"
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
	contentType := requestHeader.Get("Content-Type")

	// Open file.
	uploadedFile, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return "", errors.Wrap(err, "open uploaded file")
	}
	defer func() { _ = uploadedFile.Close() }()

	var body io.Reader
	switch contentType {
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
		requestBody, err := makeJson(makeJsonOptions{
			UploadedFile: uploadedFile,
			Body:         s.Body,
		})
		if err != nil {
			return "", errors.Wrap(err, "make json")
		}
		body = requestBody
	default:
		return "", errors.Errorf("unknown content type: %q", contentType)
	}

	client := req.C()
	request := client.R()
	request.Headers = requestHeader
	request = request.
		SetHeader("X-Forwarded-For", s.URL).
		SetBody(body)

	resp, err := request.Send(s.Method, s.URL)
	if err != nil {
		return "", errors.Wrap(err, "do request")
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("response status code: %d", resp.StatusCode)
	}

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

func RandomServer() string {
	servers.RLock()
	defer func() { servers.RUnlock() }()

	names := make([]string, 0, len(servers.m))
	for name := range servers.m {
		names = append(names, name)
	}
	return names[rand.Intn(len(names))]
}

func parseHeader(h map[string]string) http.Header {
	header := make(http.Header)
	for k, v := range h {
		header.Set(k, v)
	}
	return header
}
