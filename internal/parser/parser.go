// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package parser

import (
	"github.com/google/cel-go/cel"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"

	"github.com/wuhan005/share/internal/expression"
	"github.com/wuhan005/share/internal/yaml"
)

type Chain []*Parser

type Parser struct {
	Type   string          `yaml:"type"`
	Action yaml.RawMessage `yaml:"action"`
}

func (p *Parser) Do(input string) (string, error) {
	output := input
	var err error

	switch p.Type {
	case "json":
		var jsonPath string
		if err := p.Action.Unmarshal(&jsonPath); err != nil {
			return "", errors.Wrap(err, "unmarshal json path")
		}
		output = gjson.Get(output, jsonPath).String()

	case "expression":
		var expressionStr string
		if err := p.Action.Unmarshal(&expressionStr); err != nil {
			return "", errors.Wrap(err, "unmarshal expression")
		}

		output, err = expression.Parse(expression.ParseOptions{
			Expression: expressionStr,
			VariableDefs: map[string]*cel.Type{
				"input": cel.StringType,
			},
			Variables: map[string]interface{}{
				"input": output,
			},
		})
		if err != nil {
			return "", errors.Wrap(err, "parse expression")
		}

	default:
		return "", errors.Errorf("unknown parser type: %q", p.Type)
	}

	return output, nil
}
