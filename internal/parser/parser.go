// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package parser

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/wuhan005/share/internal/yaml"
)

type Chain []*Parser

type Parser struct {
	Type   string          `yaml:"type"`
	Action yaml.RawMessage `yaml:"action"`
}

func (p *Parser) Do(input string) (string, error) {
	output := input

	switch p.Type {
	case "replace":
		var group []string
		if err := p.Action.Unmarshal(&group); err != nil {
			return "", errors.Wrap(err, "unmarshal replace group")
		}
		if len(group)%2 != 0 {
			return "", errors.New("replace group must be even")
		}
		output = strings.NewReplacer(group...).Replace(output)

	case "json":
		var jsonPath string
		if err := p.Action.Unmarshal(&jsonPath); err != nil {
			return "", errors.Wrap(err, "unmarshal json path")
		}
		output = gjson.Get(output, jsonPath).String()

	default:
		return "", errors.Errorf("unknown parser type: %q", p.Type)
	}

	return output, nil
}
