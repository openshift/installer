/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package match

import (
	"strings"

	"github.com/Azure/azure-service-operator/v2/internal/util/typo"
)

// literalMatcher is a StringMatcher that provides a case-insensitive match against a given string
type literalMatcher struct {
	literal string
	matched bool
	advisor *typo.Advisor
}

var _ StringMatcher = &literalMatcher{}

// newLiteralMatcher returns a new literalMatcher for the given string
func newLiteralMatcher(literal string) StringMatcher {
	return &literalMatcher{
		literal: strings.TrimSpace(literal),
		advisor: typo.NewAdvisor(),
	}
}

func (lm *literalMatcher) Matches(value string) Result {
	if strings.EqualFold(lm.literal, strings.TrimSpace(value)) {
		if !lm.matched {
			// First time we match, clear out our advisory as we won't be using it
			lm.matched = true
			lm.advisor.ClearTerms()
		}

		return matchFound(lm.literal)
	}

	if !lm.matched {
		// Still collecting potential suggestions
		lm.advisor.AddTerm(value)
	}

	return matchNotFound()
}

// WasMatched returns an error if we didn't match anything, nil otherwise
func (lm *literalMatcher) WasMatched() error {
	if lm.matched {
		return nil
	}

	return lm.advisor.Errorf(lm.literal, "no match for %q", lm.literal)
}

// IsRestrictive returns true if we have a non-empty literal to match
func (lm *literalMatcher) IsRestrictive() bool {
	return lm.literal != ""
}

// String returns the literal we match
func (lm *literalMatcher) String() string {
	return lm.literal
}
