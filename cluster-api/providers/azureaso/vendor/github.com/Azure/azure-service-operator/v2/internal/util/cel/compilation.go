/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package cel

import (
	"github.com/google/cel-go/cel"
	"github.com/pkg/errors"
)

var (
	StringType = cel.StringType
	MapType    = cel.MapType(cel.StringType, cel.StringType)
)

// AllowedOutputTypes defines the set of allowed CEL expression output types supported by ASO.
// Any CEL expression whose result is a type other than one of these will be rejected.
func AllowedOutputTypes() []*cel.Type {
	return []*cel.Type{
		StringType,
		MapType,
	}
}

type CompilationResult struct {
	AST     *cel.Ast
	Program cel.Program
}

// Compile builds the specified expression and returns a CompilationResult containing
// the cel.AST and cel.Program.
func Compile(env *cel.Env, expression string) (*CompilationResult, error) {
	ast, iss := env.Compile(expression)
	if iss.Err() != nil {
		return nil, errors.Wrapf(iss.Err(), "failed to compile CEL expression: %q", expression)
	}

	err := CheckOutputTypeAllowed(ast, AllowedOutputTypes()...)
	if err != nil {
		return nil, err
	}

	program, err := env.Program(ast)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate program from CEL AST: %q", expression)
	}

	return &CompilationResult{
		AST:     ast,
		Program: program,
	}, nil
}
