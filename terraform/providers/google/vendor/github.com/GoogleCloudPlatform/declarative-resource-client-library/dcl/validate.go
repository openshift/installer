// Copyright 2024 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package dcl

import (
	"fmt"
	"strings"
)

// Required takes in a DCL resource (represented as an interface) and a dot-notated path (with JSON names). If the path is not set, an error will be returned.
func Required(r interface{}, path string) error {
	var m map[string]interface{}
	m, err := ConvertToMap(r)
	if err != nil {
		return err
	}
	val, err := GetMapEntry(m, strings.Split(path, "."))
	if err != nil {
		return err
	} else if val == nil {
		return fmt.Errorf("required value %q could not be found", path)
	}
	return nil
}

// RequiredParameter takes in a value (typically one that's not exported in JSON) and returns an error if it is not set.
func RequiredParameter(v interface{}, name string) error {
	if IsZeroValue(v) {
		return fmt.Errorf("%q must be set", name)
	}
	return nil
}

func countOfNonempty(v ...interface{}) int {
	i := 0
	for _, val := range v {
		if !IsEmptyValueIndirect(val) {
			i++
		}
	}
	return i
}

// ValidateAtMostOneOfFieldsSet returns an error if more than one of the provided values is nonempty.
func ValidateAtMostOneOfFieldsSet(fieldNames []string, v ...interface{}) error {
	if countOfNonempty(v...) > 1 {
		return fmt.Errorf("more than one value set: %v", fieldNames)
	}
	return nil
}

// ValidateAtLeastOneOfFieldsSet returns an error if none of the provided values is nonempty.
func ValidateAtLeastOneOfFieldsSet(fieldNames []string, v ...interface{}) error {
	if countOfNonempty(v...) == 0 {
		return fmt.Errorf("zero values set: %v", fieldNames)
	}
	return nil
}

// ValidateExactlyOneOfFieldsSet returns an error if 0 or 2+ of the provided values is nonempty.
func ValidateExactlyOneOfFieldsSet(fieldNames []string, v ...interface{}) error {
	if countOfNonempty(v...) != 1 {
		return fmt.Errorf("not exactly one value set: %v", fieldNames)
	}
	return nil
}

// AnySet returns true if any of the values provided is nonempty.
func AnySet(v ...interface{}) bool {
	return countOfNonempty(v...) > 0
}
