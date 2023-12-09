// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package share

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"

	"github.com/google/cel-go/cel"
	"github.com/pkg/errors"

	"github.com/wuhan005/share/internal/expression"
)

type makeMultipartFormDataOptions struct {
	UploadedFile io.Reader
	Body         map[string]string
}

func makeMultipartFormData(options makeMultipartFormDataOptions) (io.Reader, string, error) {
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	for k, v := range options.Body {
		if v == bodyFileFieldKey {
			fileWriter, err := writer.CreateFormFile(k, "image.png")
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

type makeJsonOptions struct {
	UploadedFile io.Reader
	Body         map[string]string
}

func makeJson(options makeJsonOptions) (io.Reader, error) {
	uploadedFileBytes, err := io.ReadAll(options.UploadedFile)
	if err != nil {
		return nil, errors.Wrap(err, "read uploaded file")
	}

	for k, v := range options.Body {
		if v == bodyFileFieldKey {
			options.Body[k] = string(uploadedFileBytes)
		} else if len(v) > 3 && v[0:2] == "${" && v[len(v)-1:] == "}" {
			expressionStr := v[2 : len(v)-1]
			result, err := expression.Parse(expression.ParseOptions{
				VariableDefs: map[string]*cel.Type{"file": cel.BytesType},
				Variables:    map[string]interface{}{"file": uploadedFileBytes},
				Expression:   expressionStr,
			})
			if err != nil {
				return nil, errors.Wrapf(err, "parse expression of %q: %q", k, expressionStr)
			}
			options.Body[k] = result
		}
	}

	var requestBody bytes.Buffer
	if err := json.NewEncoder(&requestBody).Encode(options.Body); err != nil {
		return nil, errors.Wrap(err, "encode json")
	}
	return &requestBody, nil
}
