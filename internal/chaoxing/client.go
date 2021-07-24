// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chaoxing

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/pkg/errors"
)

const userAgent = `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36`

// Client represents a Chaoxing client.
type Client struct {
	phone, password string

	client *http.Client
}

func NewClient(phone, password string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "new cookie jar")
	}
	client := &http.Client{Jar: jar}

	return &Client{
		phone:    phone,
		password: password,
		client:   client,
	}, nil
}
