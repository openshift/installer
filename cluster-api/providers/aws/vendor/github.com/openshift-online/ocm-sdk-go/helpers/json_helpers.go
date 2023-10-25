/*
Copyright (c) 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package helpers // github.com/openshift-online/ocm-sdk-go/helpers

import (
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// NewIterator creates a new JSON iterator that will read to the given source, which
// can be a slice of bytes, a string, a reader or an existing iterator.
func NewIterator(source interface{}) (iterator *jsoniter.Iterator, err error) {
	config := jsoniter.Config{}
	api := config.Froze()
	switch typed := source.(type) {
	case []byte:
		iterator = jsoniter.ParseBytes(api, typed)
	case string:
		iterator = jsoniter.ParseString(api, typed)
	case io.Reader:
		iterator = jsoniter.Parse(api, typed, 4096)
	case *jsoniter.Iterator:
		iterator = typed
	default:
		err = fmt.Errorf(
			"expected slice of bytes, string, reader or iterator but got '%T'",
			source,
		)
	}
	return
}

// NewStream creates a new JSON stream that will write to the given writer.
func NewStream(writer io.Writer) *jsoniter.Stream {
	config := jsoniter.Config{
		IndentionStep: 2,
	}
	api := config.Froze()
	return jsoniter.NewStream(api, writer, 0)
}

// NewBoolean allocates a new bool in the heap and returns a pointer to it.
func NewBoolean(value bool) *bool {
	return &value
}

// NewInteger allocates a new integer in the heap and returns a pointer to it.
func NewInteger(value int) *int {
	return &value
}

// NewFloat allocates a new floating point value in the heap and returns an pointer
// to it.
func NewFloat(value float64) *float64 {
	return &value
}

// NewString allocates a new string in the heap and returns a pointer to it.
func NewString(value string) *string {
	return &value
}

// NewDate allocates a new date in the heap and returns a pointer to it.
func NewDate(value time.Time) *time.Time {
	return &value
}

// ParseInteger reads a string and parses it to integer,
// if an error occurred it returns a non-nil error.
func ParseInteger(query url.Values, parameterName string) (*int, error) {
	values := query[parameterName]
	count := len(values)
	if count == 0 {
		return nil, nil
	}
	if count > 1 {
		err := fmt.Errorf(
			"expected at most one value for parameter '%s' but got %d",
			parameterName, count,
		)
		return nil, err
	}
	value := values[0]
	parsedInt64, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf(
			"value '%s' isn't valid for the '%s' parameter because it isn't an integer: %v",
			value, parameterName, err,
		)
	}
	parsedInt := int(parsedInt64)
	return &parsedInt, nil
}

// ParseFloat reads a string and parses it to float,
// if an error occurred it returns a non-nil error.
func ParseFloat(query url.Values, parameterName string) (*float64, error) {
	values := query[parameterName]
	count := len(values)
	if count == 0 {
		return nil, nil
	}
	if count > 1 {
		err := fmt.Errorf(
			"expected at most one value for parameter '%s' but got %d",
			parameterName, count,
		)
		return nil, err
	}
	value := values[0]
	parsedFloat, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, fmt.Errorf(
			"value '%s' isn't valid for the '%s' parameter because it isn't a float: %v",
			value, parameterName, err,
		)
	}
	return &parsedFloat, nil
}

// ParseString returns a pointer to the string and nil error.
func ParseString(query url.Values, parameterName string) (*string, error) {
	values := query[parameterName]
	count := len(values)
	if count == 0 {
		return nil, nil
	}
	if count > 1 {
		err := fmt.Errorf(
			"expected at most one value for parameter '%s' but got %d",
			parameterName, count,
		)
		return nil, err
	}
	return &values[0], nil
}

// ParseBoolean reads a string and parses it to boolean,
// if an error occurred it returns a non-nil error.
func ParseBoolean(query url.Values, parameterName string) (*bool, error) {
	values := query[parameterName]
	count := len(values)
	if count == 0 {
		return nil, nil
	}
	if count > 1 {
		err := fmt.Errorf(
			"expected at most one value for parameter '%s' but got %d",
			parameterName, count,
		)
		return nil, err
	}
	value := values[0]
	parsedBool, err := strconv.ParseBool(value)
	if err != nil {
		return nil, fmt.Errorf(
			"value '%s' isn't valid for the '%s' parameter because it isn't a boolean: %v",
			value, parameterName, err,
		)
	}
	return &parsedBool, nil
}

// ParseDate reads a string and parses it to a time.Time,
// if an error occurred it returns a non-nil error.
func ParseDate(query url.Values, parameterName string) (*time.Time, error) {
	values := query[parameterName]
	count := len(values)
	if count == 0 {
		return nil, nil
	}
	if count > 1 {
		err := fmt.Errorf(
			"expected at most one value for parameter '%s' but got %d",
			parameterName, count,
		)
		return nil, err
	}
	value := values[0]
	parsedTime, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, fmt.Errorf(
			"value '%s' isn't valid for the '%s' parameter because it isn't a date: %v",
			value, parameterName, err,
		)
	}
	return &parsedTime, nil
}
