// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/wuhan005/share/pkg/share"
)

var Share = func(c *cli.Context) error {
	server := c.String("server")
	filePath := c.Args().Get(0)

	url, err := share.Share(server, filePath)
	if err != nil {
		return errors.Wrap(err, "share")
	}

	fmt.Println(url)
	return nil
}
