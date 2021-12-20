package core

// (C) Copyright IBM Corp. 2019.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	validator "gopkg.in/go-playground/validator.v9"
)

// Validate is a shared validator instance used to perform validation of structs.
var Validate *validator.Validate

func init() {
	Validate = validator.New()
}

const (
	jsonMimePattern      = "(?i)^application\\/((json)|(merge\\-patch\\+json))(;.*)?$"
	jsonPatchMimePattern = "(?i)^application\\/json\\-patch\\+json(;.*)?$"
)

// isNil checks if the specified object is nil or not.
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false
}

// ValidateNotNil returns the specified error if 'object' is nil, nil otherwise.
func ValidateNotNil(object interface{}, errorMsg string) error {
	if isNil(object) {
		return errors.New(errorMsg)
	}
	return nil
}

// ValidateStruct validates 'param' (assumed to be a struct) according to the
// annotations attached to its fields.
func ValidateStruct(param interface{}, paramName string) error {
	if param != nil {
		if err := Validate.Struct(param); err != nil {
			// If there were validation errors then return an error containing the field errors
			if fieldErrors, ok := err.(validator.ValidationErrors); ok {
				return fmt.Errorf("%s failed validation:\n%s", paramName, fieldErrors.Error())
			}
			return fmt.Errorf("An unexpected system error occurred while validating %s\n%s", paramName, err.Error())
		}
	}
	return nil
}

// StringPtr returns a pointer to string literal.
func StringPtr(literal string) *string {
	return &literal
}

// BoolPtr returns a pointer to boolean literal.
func BoolPtr(literal bool) *bool {
	return &literal
}

// Int64Ptr returns a pointer to int64 literal.
func Int64Ptr(literal int64) *int64 {
	return &literal
}

// Float32Ptr returns a pointer to float32 literal.
func Float32Ptr(literal float32) *float32 {
	return &literal
}

// Float64Ptr returns a pointer to float64 literal.
func Float64Ptr(literal float64) *float64 {
	return &literal
}

// IsJSONMimeType Returns true iff the specified mimeType value represents a
// "JSON" mimetype.
func IsJSONMimeType(mimeType string) bool {
	if mimeType == "" {
		return false
	}
	matched, err := regexp.MatchString(jsonMimePattern, mimeType)
	if err != nil {
		return false
	}
	return matched
}

// IsJSONPatchMimeType returns true iff the specified mimeType value represents
// a "JSON Patch" mimetype.
func IsJSONPatchMimeType(mimeType string) bool {
	if mimeType == "" {
		return false
	}
	matched, err := regexp.MatchString(jsonPatchMimePattern, mimeType)
	if err != nil {
		return false
	}
	return matched
}

// StringNilMapper de-references the parameter 's' and returns the result, or ""
// if 's' is nil.
func StringNilMapper(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// HasBadFirstOrLastChar checks if the string starts with `{` or `"`
// or ends with `}` or `"`.
func HasBadFirstOrLastChar(str string) bool {
	return strings.HasPrefix(str, "{") || strings.HasPrefix(str, "\"") ||
		strings.HasSuffix(str, "}") || strings.HasSuffix(str, "\"")
}

// UserHomeDir returns the user home directory.
func UserHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// SystemInfo returns the system information.
func SystemInfo() string {
	return fmt.Sprintf("(arch=%s; os=%s; go.version=%s)", runtime.GOARCH, runtime.GOOS, runtime.Version())
}

// PrettyPrint print pretty.
func PrettyPrint(result interface{}, resultName string) {
	output, err := json.MarshalIndent(result, "", "    ")

	if err == nil {
		fmt.Printf("%v:\n%+v\n\n", resultName, string(output))
	}
}

// GetCurrentTime returns the current Unix time.
func GetCurrentTime() int64 {
	return time.Now().Unix()
}
