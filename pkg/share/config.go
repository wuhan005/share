// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package share

import (
	"io"
	"os"
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

var servers = struct {
	sync.RWMutex
	m map[string]*Server
}{
	m: make(map[string]*Server),
}

func LoadServerFromFile(path string) (map[string]*Server, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	return LoadServersFromReader(file)
}

func LoadServersFromReader(reader io.Reader) (map[string]*Server, error) {
	var s map[string]*Server
	if err := yaml.NewDecoder(reader).Decode(&s); err != nil {
		return nil, errors.Wrap(err, "decode yaml")
	}
	return s, nil
}

func SetServers(s map[string]*Server) {
	servers.Lock()
	defer func() { servers.Unlock() }()

	servers.m = s
}
