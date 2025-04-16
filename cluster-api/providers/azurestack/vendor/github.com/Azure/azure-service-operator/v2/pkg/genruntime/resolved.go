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
type Resolved[T reference, V any] struct {
	// resolved is a map of T to value.
	resolved map[T]V
}

// MakeResolved creates a Resolved
func MakeResolved[T reference, V any](resolvedMap map[T]V) Resolved[T, V] {
	return Resolved[T, V]{
		resolved: resolvedMap,
	}
}

// Lookup looks up the value for the given reference. If it cannot be found, an error is returned.
func (r Resolved[T, V]) Lookup(ref T) (V, error) {
	result, ok := r.resolved[ref]
	if !ok {
		var ret V
		return ret, errors.Errorf("couldn't find resolved %T %s", ref, ref.String())
	}
	return result, nil
}

// LookupFromPtr looks up the value for the given reference. If the reference is nil, an error is returned.
// If the value cannot be found, an error is returned
func (r Resolved[T, V]) LookupFromPtr(ref *T) (V, error) {
	if ref == nil {
		var ret V
		return ret, errors.Errorf("cannot look up secret from nil SecretReference")
	}

	return r.Lookup(*ref)
}
