// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	_ "embed"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/wuhan005/share/internal/cmd"
	"github.com/wuhan005/share/pkg/share"
)

//go:embed data/server.yaml
var configFile []byte

func init() {
	servers, err := share.LoadServersFromReader(bytes.NewReader(configFile))
	if err != nil {
		logrus.WithError(err).Fatal("Failed to load config file")
	}
	share.SetServers(servers)
}

func main() {
	app := cli.NewApp()
	app.Name = "share"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "server",
			Usage:    "Server to upload.",
			Required: true,
		},
	}
	app.Action = cmd.Share

	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("Failed to start application")
	}
}
