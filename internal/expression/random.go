// Copyright 2023 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package expression

import (
	"math/rand"
	"strconv"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func Randoms() cel.EnvOption {
	return cel.Lib(randomLib{})
}

type randomLib struct{}

func (randomLib) LibraryName() string {
	return "cel.lib.ext.randoms"
}

func (randomLib) CompileOptions() []cel.EnvOption {
	return []cel.EnvOption{
		cel.Function("random.randint",
			cel.Overload("random_randint", []*cel.Type{cel.IntType, cel.IntType}, cel.StringType,
				cel.BinaryBinding(func(from ref.Val, to ref.Val) ref.Val {
					fromInt := int(from.(types.Int))
					toInt := int(to.(types.Int))
					randomNum := rand.Intn(toInt-fromInt+1) + fromInt
					return types.String(strconv.Itoa(randomNum))
				}))),
	}
}

func (randomLib) ProgramOptions() []cel.ProgramOption {
	return []cel.ProgramOption{}
}
