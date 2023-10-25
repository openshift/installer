/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package typo

import (
	"fmt"
	"strings"
	"sync"

	"github.com/hbollon/go-edlib"
	"github.com/pkg/errors"

	"github.com/Azure/azure-service-operator/v2/internal/set"
)

// Advisor is a utility that helps augment errors with guidance when mistakes are made in configuration.
type Advisor struct {
	lock  sync.RWMutex
	terms set.Set[string] // set of terms we know to exist
}

func NewAdvisor() *Advisor {
	return &Advisor{
		terms: set.Make[string](),
	}
}

// AddTerm records that we saw the specified item
func (advisor *Advisor) AddTerm(item string) {
	advisor.lock.Lock()
	defer advisor.lock.Unlock()
	advisor.terms.Add(item)
}

// HasTerms returns true if this advisor has any terms available
func (advisor *Advisor) HasTerms() bool {
	return len(advisor.terms) > 0
}

// HasTerm returns true if the specified term is known to the advisor
func (advisor *Advisor) HasTerm(term string) bool {
	advisor.lock.RLock()
	defer advisor.lock.RUnlock()
	return advisor.terms.Contains(term)
}

// ClearTerms removes all the terms from ths advisor
func (advisor *Advisor) ClearTerms() {
	advisor.terms.Clear()
}

// Errorf creates a new error with advice, or a simple error if no advice possible
func (advisor *Advisor) Errorf(typo string, format string, args ...interface{}) error {
	advisor.lock.RLock()
	defer advisor.lock.RUnlock()

	if !advisor.HasTerms() || advisor.terms.Contains(typo) {
		// Can't make any suggestions,
		return errors.Errorf(format, args...)
	}

	msg := fmt.Sprintf(format, args...)

	suggestion, err := edlib.FuzzySearch(
		strings.ToLower(typo),
		set.AsSortedSlice(advisor.terms),
		edlib.Levenshtein)
	if err != nil {
		// Can't offer a suggestion
		return errors.Errorf(
			"%s (unable to provide suggestion: %s)",
			msg,
			err)
	}

	return errors.Errorf(
		"%s (did you mean %s?)",
		msg,
		suggestion)
}

// Wrapf adds any guidance to the provided error, if possible
func (advisor *Advisor) Wrapf(originalError error, typo string, format string, args ...interface{}) error {
	advisor.lock.RLock()
	defer advisor.lock.RUnlock()

	if !advisor.HasTerms() || advisor.terms.Contains(typo) || originalError == nil {
		// Can't make any suggestions, or don't need one
		return originalError
	}

	msg := fmt.Sprintf(format, args...)

	suggestion, err := edlib.FuzzySearch(
		strings.ToLower(typo),
		set.AsSortedSlice(advisor.terms),
		edlib.Levenshtein)
	if err != nil {
		// Can't offer a suggestion
		return errors.Wrapf(
			originalError,
			"%s (unable to provide suggestion: %s)",
			msg,
			err)
	}

	return errors.Wrapf(
		originalError,
		"%s (did you mean %s?)",
		msg,
		suggestion)
}
