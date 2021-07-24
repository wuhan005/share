// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chaoxing

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
	"github.com/wuhan005/gadget"
)

type LoginResponse struct {
	Name   string `json:"name"`
	Status bool   `json:"status"`
	// Message is the error message when login failed.
	Message string `json:"msg2"`
}

// Login logs in the Chaoxing with the given phone and password.
func (c *Client) Login() error {
	password := gadget.Base64Encode(c.password)
	body := url.Values{
		"fid":              []string{"-1"},
		"uname":            []string{c.phone},
		"password":         []string{password},
		"refer":            []string{"http%3A%2F%2Fpan-yz.chaoxing.com%2F"},
		"t":                []string{"true"},
		"forbidotherlogin": []string{"0"},
	}

	req, err := http.NewRequest(http.MethodPost, "https://passport2.chaoxing.com/fanyalogin", strings.NewReader(body.Encode()))
	if err != nil {
		return errors.Wrap(err, "new request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("User-Agent", userAgent)

	resp, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "do request")
	}
	defer func() { _ = resp.Body.Close() }()

	// It always returns 200 OK whenever login succeeded or failed.
	if resp.StatusCode/100 != 2 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "read response body")
		}
		return errors.Errorf("unexpected status code: %d, response body: %q", resp.StatusCode, string(respBody))
	}

	var loginResponse LoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
		return errors.Wrap(err, "decode JSON")
	}
	if !loginResponse.Status {
		return errors.Errorf("login failed: %q", loginResponse.Message)
	}
	return nil
}
