/*
 * Copyright (c) Microsoft Corporation.
 * Licensed under the MIT license.
 */

package genruntime

import "strings"

// ToEnum does a case-insensitive conversion of a string to an enum using a provided conversion map.
// If the required value is not found, a literal cast will be used to return the enum.
func ToEnum[T ~string](str string, enumMap map[string]T) T {
	if val, ok := enumMap[strings.ToLower(str)]; ok {
		return val
	}

	return T(str)
}
