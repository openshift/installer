package ocm

import (
	"fmt"
	"strings"
)

const (
	LowerCaseLetters     = "abcdefghijklmnopqrstuvwxyz"
	UpperCaseLetters     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits               = "0123456789"
	ValidIdentifierChars = LowerCaseLetters + UpperCaseLetters + Digits + "_"
)

// FeatureFilter is a wrapper for a string providing a filter for instance types
// with a specific feature. Leaving the fields unexported allows us to prevent construction of
// an invalid filter by package-external callers
type FeatureFilter struct {
	featureName string
}

type FeatureFilters []FeatureFilter

func validateIdentifier(input string) error {
	if input == "" {
		return fmt.Errorf("valid identifier may not be empty")
	}
	if len(input) > 32 {
		return fmt.Errorf("identifier must not exceed 32 characters in length: this input is %d characters", len(input))
	}
	for _, c := range input {
		if !strings.ContainsRune(ValidIdentifierChars, c) {
			return fmt.Errorf("input contains invalid character '%c'", c)
		}
	}
	return nil
}

func (fs FeatureFilters) String() string {
	strs := make([]string, 0, len(fs))
	for _, f := range fs {
		strs = append(strs, f.String())
	}
	return strings.Join(strs, " AND ")
}

func (f FeatureFilter) String() string {
	return fmt.Sprintf("features.%s = 'true'", f.featureName)
}

func ParseFeatureFilter(name string) (FeatureFilter, error) {
	if err := validateIdentifier(name); err != nil {
		return FeatureFilter{}, fmt.Errorf("feature name '%s' is not a valid identifier: %w", name, err)
	}
	return FeatureFilter{
		featureName: name,
	}, nil
}
