// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/share/internal/chaoxing"
)

var Share = func(c *cli.Context) error {
	phone := c.String("phone")
	password := c.String("password")

	client, err := chaoxing.NewClient(phone, password)
	if err != nil {
		return errors.Wrap(err, "new client")
	}
	if err := client.Login(); err != nil {
		return errors.Wrap(err, "login")
	}

	if c.NArg() == 0 {
		return errors.New("empty argument")
	}
	filePath := c.Args().Get(0)
	resp, err := client.Upload(filePath)
	if err != nil {
		return errors.Wrap(err, "upload")
	}

	fileID := resp.AttFile.AttClouddisk.FileId
	downloadURL, err := chaoxing.GetDownloadURL(fileID)
	if err != nil {
		return errors.Wrap(err, "get download URL")
	}
	log.Info("%s", downloadURL)
	return nil
}
