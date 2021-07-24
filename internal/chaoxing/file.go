// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package chaoxing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/pkg/errors"
)

type UploadResponse struct {
	AttFile struct {
		AttClouddisk struct {
			Size       int    `json:"size"`
			FileSize   string `json:"fileSize"`
			Isfile     bool   `json:"isfile"`
			Modtime    int64  `json:"modtime"`
			ParentPath string `json:"parentPath"`
			Icon       string `json:"icon"`
			Name       string `json:"name"`
			DownPath   string `json:"downPath"`
			ShareUrl   string `json:"shareUrl"`
			Suffix     string `json:"suffix"`
			FileId     string `json:"fileId"`
		} `json:"att_clouddisk"`
		AttachmentType int `json:"attachmentType"`
	} `json:"att_file"`
	Type   string `json:"type"`
	Url    string `json:"url"`
	Status bool   `json:"status"`
}

func (c *Client) Upload(path string) (*UploadResponse, error) {
	buf := new(bytes.Buffer)
	bufferWriter := multipart.NewWriter(buf)

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	defer func() { _ = file.Close() }()
	fileWriter, err := bufferWriter.CreateFormFile("attrFile", file.Name())
	if err != nil {
		return nil, errors.Wrap(err, "create file writer")
	}
	if _, err := io.Copy(fileWriter, file); err != nil {
		return nil, errors.Wrap(err, "copy file")
	}
	if err := bufferWriter.Close(); err != nil {
		return nil, errors.Wrap(err, "close buffer writer")
	}

	resp, err := c.client.Post("https://notice.chaoxing.com/pc/files/uploadNoticeFile", bufferWriter.FormDataContentType(), buf)
	if err != nil {
		return nil, errors.Wrap(err, "request")
	}
	defer func() { _ = resp.Body.Close() }()

	var uploadResponse UploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResponse); err != nil {
		return nil, errors.Wrap(err, "decode json")
	}

	return &uploadResponse, nil
}

func (c *Client) UploadToSpace() {
	panic("implement me")
}

func GetShareURL(objectURL string) string {
	return fmt.Sprintf("https://pan-yz.chaoxing.com/external/m/file/" + objectURL)
}

func GetDownloadURL(objectURL string) (string, error) {
	var respJSON struct {
		Message     string `json:"msg"`
		DownloadURL string `json:"download"`
		PreviewURL  string `json:"url"`
		Status      bool   `json:"status"`
	}

	resp, err := http.Get("https://noteyd.chaoxing.com/screen/note_note/files/status/" + objectURL)
	if err != nil {
		return "", errors.Wrap(err, "request")
	}
	defer func() { _ = resp.Body.Close() }()

	// It always returns 200 OK whenever file is found.
	if resp.StatusCode/100 != 2 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", errors.Wrap(err, "read response body")
		}
		return "", errors.Errorf("unexpected status code: %d, response body: %q", resp.StatusCode, string(respBody))
	}

	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return "", errors.Wrap(err, "json decode")
	}

	if !respJSON.Status {
		return "", errors.Errorf("get download url: %q", respJSON.Message)
	}
	return respJSON.DownloadURL, nil
}
