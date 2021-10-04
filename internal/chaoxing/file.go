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
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/schollz/progressbar/v3"
	log "unknwon.dev/clog/v2"
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

	// Progress bar
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "get file info")
	}
	bar := progressbar.NewOptions64(
		fileInfo.Size(),
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			_, _ = fmt.Fprint(os.Stderr, "\n")
		}),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionFullWidth(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerPadding: " ",
			BarStart:      "|",
			BarEnd:        "|",
			SaucerHead:    ">",
		}),
	)

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

	log.Trace("Upload file %q...", file.Name())
	reader := progressbar.NewReader(buf, bar)
	resp, err := c.client.Post("https://notice.chaoxing.com/pc/files/uploadNoticeFile", bufferWriter.FormDataContentType(), &reader)
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

func GetDownloadURL(objectID string) string {
	return "http://cloud.ananas.chaoxing.com/view/fileviewDownload?objectId=" + objectID
}
