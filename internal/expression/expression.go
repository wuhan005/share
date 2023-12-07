// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package expression

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/ext"
	"github.com/pkg/errors"
)

type ParseOptions struct {
	VariableDefs map[string]*cel.Type
	Variables    map[string]interface{}
	Expression   string
}

func Parse(options ParseOptions) (string, error) {
	envOptions := []cel.EnvOption{
		ext.Encoders(),
		ext.Strings(),
		Randoms(),
	}
	for k, v := range options.VariableDefs {
		v := v
		envOptions = append(envOptions, cel.Variable(k, v))
	}

	env, err := cel.NewEnv(envOptions...)
	if err != nil {
		return "", errors.Wrap(err, "new env")
	}

	ast, iss := env.Compile(options.Expression)
	if iss.Err() != nil {
		return "", errors.Wrap(iss.Err(), "compile expression")
	}
	prg, err := env.Program(ast)
	if err != nil {
		return "", errors.Wrap(err, "program")
	}

	out, _, err := prg.Eval(options.Variables)
	if err != nil {
		return "", errors.Wrap(err, "eval")
	}
	return out.Value().(string), nil
}
