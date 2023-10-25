/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package set

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// Set provides a standard way to have a set of distinct things
type Set[T comparable] map[T]struct{}

// Make creates a new set with the given values
func Make[T comparable](ts ...T) Set[T] {
	result := make(Set[T], len(ts))

	for _, x := range ts {
		result.Add(x)
	}

	return result
}

// Contains does a check to see if the provided value is in the set
func (set Set[T]) Contains(x T) bool {
	_, ok := set[x]
	return ok
}

// Add adds the provided value into the set
// Nothing happens if the value is already present
func (set Set[T]) Add(x T) {
	set[x] = struct{}{}
}

// AddAll adds the provided set into the existing set.
// Nothing happens for values which are already present
func (set Set[T]) AddAll(x Set[T]) {
	for _, val := range x.Values() {
		set.Add(val)
	}
}

// Remove deletes the provided value from the set
// Nothing happens if the value is not present
func (set Set[T]) Remove(x T) {
	delete(set, x)
}

// Copy returns an independent copy of the set
func (set Set[T]) Copy() Set[T] {
	return maps.Clone(set)
}

// Clear removes all the items from this set.
func (set Set[T]) Clear() {
	maps.Clear(set)
}

// Equals checks to see if the two sets are equivalent
func (set Set[T]) Equals(other Set[T]) bool {
	return maps.Equal(set, other)
}

func AreEqual[T comparable](left, right Set[T]) bool {
	return maps.Equal(left, right)
}

// Values returns a slice of all the values in the set
func (set Set[T]) Values() []T {
	return maps.Keys(set)
}

// Where returns a new set with only the set of values which match the predicate
func (set Set[T]) Where(predicate func(T) bool) Set[T] {
	result := make(Set[T], len(set))

	for val := range set {
		if predicate(val) {
			result.Add(val)
		}
	}

	return result
}

// Except returns a new set with only the set of values which are not in the other set
func (set Set[T]) Except(other Set[T]) Set[T] {
	result := make(Set[T], len(set))

	for val := range set {
		if !other.Contains(val) {
			result.Add(val)
		}
	}

	return result
}

// AsSortedSlice returns a sorted slice of values from this set
func AsSortedSlice[T constraints.Ordered](set Set[T]) []T {
	result := maps.Keys(set)
	slices.Sort(result)
	return result
}
