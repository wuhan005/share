// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/urfave/cli/v2"

	"github.com/wuhan005/share/internal/chaoxing"
)

var Login = &cli.Command{
	Name:   "login",
	Usage:  "Login to a Chaoxing account.",
	Action: login,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "phone",
			Usage:    "Phone number of the Chaoxing account.",
			EnvVars:  []string{"CHAOXING_ACCOUNT_PHONE"},
			Required: true,
		},
		&cli.StringFlag{
			Name:     "password",
			Usage:    "Password number of the Chaoxing account.",
			EnvVars:  []string{"CHAOXING_ACCOUNT_PASSWORD"},
			Required: true,
		},
	},
}

func login(c *cli.Context) error {
	phone := c.String("phone")
	password := c.String("password")

	client := chaoxing.NewClient(phone, password)
	return client.Login()
}
