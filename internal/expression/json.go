// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package expression

import (
	"encoding/json"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func JSONs() cel.EnvOption {
	return cel.Lib(jsonLib{})
}

type jsonLib struct{}

func (jsonLib) LibraryName() string {
	return "cel.lib.ext.jsons"
}

func (jsonLib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Function("json.loads",
			cel.Overload("json_loads", []*cel.Type{cel.StringType}, cel.StringType,
				cel.UnaryBinding(func(input ref.Val) ref.Val {
					str := string(input.(types.String))
					if err := json.Unmarshal([]byte(str), &str); err != nil {
						return types.String(err.Error())
					}
					return types.String(str)
				}))),
	}
}

func (jsonLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}
