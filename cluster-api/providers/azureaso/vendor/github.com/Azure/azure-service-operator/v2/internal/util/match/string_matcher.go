/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package match

import (
	"fmt"
)

type Result struct {
	Matched         bool
	MatchingPattern string
}

func matchFound(pattern string) Result {
	return Result{
		Matched:         true,
		MatchingPattern: pattern,
	}
}

func matchNotFound() Result {
	return Result{
		Matched:         false,
		MatchingPattern: "",
	}
}

// StringMatcher is an interface implemented by predicates used to test string values
type StringMatcher interface {
	fmt.Stringer

	// Matches returns true if the provided value is matched by the matcher
	Matches(value string) Result
	// WasMatched returns nil if the matcher had a match, otherwise returning a diagnostic error
	WasMatched() error
	// IsRestrictive returns true if the matcher is populated and will restrict matches
	IsRestrictive() bool
}

// NewStringMatcher returns a matcher for the specified string
// Different strings may return different implementations:
// o If the string contains ';', a multi-matcher of sub-matches, one for each item in the string separated by ';'
// o If the string contains '*' or '?' a globbing wildcard matcher
// o Otherwise a case-insensitive literal string matcher
func NewStringMatcher(matcher string) StringMatcher {
	if HasMultipleMatchers(matcher) {
		return newMultiMatcher(matcher)
	}

	if HasWildCards(matcher) {
		return newGlobMatcher(matcher)
	}

	return newLiteralMatcher(matcher)
}
