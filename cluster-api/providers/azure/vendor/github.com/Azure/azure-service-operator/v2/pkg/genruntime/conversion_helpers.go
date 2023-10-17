/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import (
	"math"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/conditions"
)

/*
This file contains manual implementations to reduce code bloat in generated code.
*/

// CloneSliceOfString clones the given []string. It is used (indirectly)
// by the generator when generating property conversions.
func CloneSliceOfString(slice []string) []string {
	if slice == nil {
		return nil // preserve nils
	}

	result := make([]string, len(slice))
	copy(result, slice)
	return result
}

// CloneMapOfStringToString clones the given map[string]string. It is used
// (indirectly) by the generator when generating property conversions.
func CloneMapOfStringToString(input map[string]string) map[string]string {
	if input == nil {
		return nil // preserve nils
	}

	result := make(map[string]string, len(input))
	for k, v := range input {
		result[k] = v
	}

	return result
}

// CloneSliceOfCondition clones the given []conditions.Condition. It is used
// (indirectly) by the generator when generating property conversions.
func CloneSliceOfCondition(slice []conditions.Condition) []conditions.Condition {
	if slice == nil {
		return nil // preserve nils
	}

	result := make([]conditions.Condition, len(slice))
	copy(result, slice)
	// we don’t need to invoke .Copy() manually as Conditions are “value objects”
	// see implementation of Condition.Copy
	return result
}

// ClonePointerToInt clones the given *int. It is used
// (indirectly) by the generator when generating property conversions.
func ClonePointerToInt(ptr *int) *int {
	if ptr == nil {
		return nil
	}

	result := *ptr
	return &result
}

// ClonePointerToString clones the given *string. It is used
// (indirectly) by the generator when generating property conversions.
func ClonePointerToString(ptr *string) *string {
	if ptr == nil {
		return nil
	}

	result := *ptr
	return &result
}

// GetOptionalStringValue clones the given *string, or
// returns empty string if the pointer is nil. It is used
// (indirectly) by the generator when generating property conversions.
func GetOptionalStringValue(ptr *string) string {
	if ptr == nil {
		return ""
	}

	return *ptr
}

// GetOptionalIntValue clones the given *int, or
// returns 0 if the pointer is nil. It is used (indirectly)
// by the generator when generating property conversions.
func GetOptionalIntValue(ptr *int) int {
	if ptr == nil {
		return 0
	}

	return *ptr
}

// GetIntFromFloat returns the int value of the given float64.
// Primarily used when initializing Spec properties from Status properties.
func GetIntFromFloat(f float64) int {
	return int(math.Round(f))
}
