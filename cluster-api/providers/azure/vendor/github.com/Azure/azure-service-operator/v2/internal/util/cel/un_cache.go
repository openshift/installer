/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package cel

import (
	"reflect"

	"github.com/google/cel-go/cel"
	"github.com/pkg/errors"
)

// UnCache is a cache that doesn't cache, useful mostly for testing purposes
type UnCache struct {
	newEnv  func(resource reflect.Type) (*cel.Env, error)
	compile func(env *cel.Env, expression string) (*CompilationResult, error)
}

var _ ProgramCacher = &UnCache{}

func NewUnCache(
	newEnv func(resource reflect.Type) (*cel.Env, error),
	compile func(env *cel.Env, expression string) (*CompilationResult, error),
) *UnCache {
	return &UnCache{
		newEnv:  newEnv,
		compile: compile,
	}
}

func (c *UnCache) Start() {}

func (c *UnCache) Stop() {}

func (c *UnCache) Get(resource reflect.Type, expression string) (*CompilationResult, error) {
	env, err := c.newEnv(resource)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get CEL env")
	}

	return c.compile(env, expression)
}
