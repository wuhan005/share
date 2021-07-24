// Copyright 2021 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"os"

	"github.com/urfave/cli/v2"
	log "unknwon.dev/clog/v2"

	"github.com/wuhan005/share/internal/cmd"
)

func main() {
	defer log.Stop()
	err := log.NewConsole()
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "share"
	app.Commands = []*cli.Command{
		cmd.Login,
	}

	if err := app.Run(os.Args); err != nil {
		log.Error("%v", err)
	}
}
