package common

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

type testElement struct {
	Name    string `validate:"required,uniqueField"`
	Address string `validate:"required,uniqueField"`
}

type testElementOptionalFields struct {
	Name     string `validate:"uniqueField"`
	Optional string `validate:"uniqueField"`
	Address  string `validate:"required,uniqueField"`
}

func noFilter([]byte) bool { return false }

func TestValidateUniqueAndRequiredFields(t *testing.T) {
	fldPath := field.NewPath("test")

	cases := []struct {
		name           string
		elements       []*testElement
		expectedErrors int
		expectedMsgs   []string
	}{
		{
			name: "valid unique elements",
			elements: []*testElement{
				{Name: "a", Address: "addr1"},
				{Name: "b", Address: "addr2"},
			},
			expectedErrors: 0,
		},
		{
			name: "duplicate name",
			elements: []*testElement{
				{Name: "a", Address: "addr1"},
				{Name: "a", Address: "addr2"},
			},
			expectedErrors: 1,
			expectedMsgs:   []string{"Duplicate value"},
		},
		{
			name: "duplicate address",
			elements: []*testElement{
				{Name: "a", Address: "addr1"},
				{Name: "b", Address: "addr1"},
			},
			expectedErrors: 1,
			expectedMsgs:   []string{"Duplicate value"},
		},
		{
			name: "missing required name",
			elements: []*testElement{
				{Name: "", Address: "addr1"},
			},
			expectedErrors: 1,
			expectedMsgs:   []string{"Required value"},
		},
		{
			name: "missing required address",
			elements: []*testElement{
				{Name: "a", Address: ""},
			},
			expectedErrors: 1,
			expectedMsgs:   []string{"Required value"},
		},
		{
			name: "multiple required fields missing",
			elements: []*testElement{
				{Name: "", Address: ""},
			},
			expectedErrors: 2,
		},
		{
			name: "duplicate and missing required",
			elements: []*testElement{
				{Name: "a", Address: "addr1"},
				{Name: "a", Address: ""},
			},
			expectedErrors: 2,
			expectedMsgs:   []string{"Duplicate value", "Required value"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateUniqueAndRequiredFields(tc.elements, fldPath, noFilter)
			assert.Len(t, errs, tc.expectedErrors)
			for _, msg := range tc.expectedMsgs {
				found := false
				for _, err := range errs {
					if assert.ObjectsAreEqual(err.Type.String(), msg) ||
						strings.Contains(err.Error(), msg) {
						found = true
						break
					}
				}
				assert.True(t, found, "expected error containing %q", msg)
			}
		})
	}
}

func TestValidateUniqueAndRequiredFields_ZeroValueSkipsUniqueness(t *testing.T) {
	fldPath := field.NewPath("test")

	cases := []struct {
		name           string
		elements       []*testElementOptionalFields
		expectedErrors int
		expectedMsgs   []string
	}{
		{
			name: "multiple empty optional fields are not duplicates",
			elements: []*testElementOptionalFields{
				{Name: "a", Optional: "", Address: "addr1"},
				{Name: "b", Optional: "", Address: "addr2"},
			},
			expectedErrors: 0,
		},
		{
			name: "multiple empty name fields are not duplicates",
			elements: []*testElementOptionalFields{
				{Name: "", Optional: "", Address: "addr1"},
				{Name: "", Optional: "", Address: "addr2"},
			},
			expectedErrors: 0,
		},
		{
			name: "non-empty optional fields still checked for uniqueness",
			elements: []*testElementOptionalFields{
				{Name: "a", Optional: "same", Address: "addr1"},
				{Name: "b", Optional: "same", Address: "addr2"},
			},
			expectedErrors: 1,
			expectedMsgs:   []string{"Duplicate value"},
		},
		{
			name: "mix of empty and non-empty optional does not conflict",
			elements: []*testElementOptionalFields{
				{Name: "a", Optional: "value", Address: "addr1"},
				{Name: "b", Optional: "", Address: "addr2"},
			},
			expectedErrors: 0,
		},
		{
			name: "required fields still enforced when optional is empty",
			elements: []*testElementOptionalFields{
				{Name: "a", Optional: "", Address: ""},
			},
			expectedErrors: 1,
			expectedMsgs:   []string{"Required value"},
		},
		{
			name: "empty required+uniqueField values are flagged as required not duplicate",
			elements: []*testElementOptionalFields{
				{Name: "a", Optional: "", Address: ""},
				{Name: "b", Optional: "", Address: ""},
			},
			expectedErrors: 2,
			expectedMsgs:   []string{"Required value"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			errs := ValidateUniqueAndRequiredFields(tc.elements, fldPath, noFilter)
			assert.Len(t, errs, tc.expectedErrors, "errors: %v", errs)
			for _, msg := range tc.expectedMsgs {
				found := false
				for _, err := range errs {
					if strings.Contains(err.Error(), msg) {
						found = true
						break
					}
				}
				assert.True(t, found, "expected error containing %q in %v", msg, errs)
			}
		})
	}
}
