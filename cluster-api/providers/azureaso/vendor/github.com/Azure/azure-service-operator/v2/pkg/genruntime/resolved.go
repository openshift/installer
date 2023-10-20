/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"fmt"

	"github.com/pkg/errors"
)

type reference interface {
	fmt.Stringer
	comparable
}

// Resolved is a set of references which have been resolved for a particular resource.
type Resolved[T reference] struct {
	// resolved is a map of T to value.
	resolved map[T]string
}

// MakeResolved creates a Resolved
func MakeResolved[T reference](resolvedMap map[T]string) Resolved[T] {
	return Resolved[T]{
		resolved: resolvedMap,
	}
}

// Lookup looks up the value for the given reference. If it cannot be found, an error is returned.
func (r Resolved[T]) Lookup(ref T) (string, error) {
	result, ok := r.resolved[ref]
	if !ok {
		return "", errors.Errorf("couldn't find resolved %T %s", ref, ref.String())
	}
	return result, nil
}

// LookupFromPtr looks up the value for the given reference. If the reference is nil, an error is returned.
// If the value cannot be found, an error is returned
func (r Resolved[T]) LookupFromPtr(ref *T) (string, error) {
	if ref == nil {
		return "", errors.Errorf("cannot look up secret from nil SecretReference")
	}

	return r.Lookup(*ref)
}
