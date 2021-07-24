// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/wuhan005/share/internal/chaoxing"
)

var Login = &cli.Command{
	Name:   "login",
	Usage:  "Login to a Chaoxing account.",
	Action: login,
}

func login(c *cli.Context) error {
	phone := c.String("phone")
	password := c.String("password")

	client, err := chaoxing.NewClient(phone, password)
	if err != nil {
		return errors.Wrap(err, "new client")
	}
	return client.Login()
}
