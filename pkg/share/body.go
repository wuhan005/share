// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package share

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"

	"github.com/pkg/errors"
)

type makeMultipartFormDataOptions struct {
	UploadedFile *os.File
	Body         map[string]string
}

func makeMultipartFormData(options makeMultipartFormDataOptions) (io.Reader, string, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	for k, v := range options.Body {
		if v == bodyFileFieldKey {
			fileWriter, err := writer.CreateFormFile(k, options.UploadedFile.Name())
			if err != nil {
				return nil, "", errors.Wrap(err, "create form file field")
			}
			if _, err := io.Copy(fileWriter, options.UploadedFile); err != nil {
				return nil, "", errors.Wrap(err, "copy file")
			}
		} else {
			if err := writer.WriteField(k, v); err != nil {
				return nil, "", errors.Wrap(err, "write field")
			}
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", errors.Wrap(err, "close writer")
	}

	return &requestBody, writer.FormDataContentType(), nil
}
