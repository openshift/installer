// Copyright 2023 Google LLC. All Rights Reserved.
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
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/kylelemons/godebug/pretty"
)

// LifecycleParam is used to specify what operations are acceptable.
type LifecycleParam int

const (
	// BlockDestruction blocks deleting the resource.  This will
	// cause some Apply operations to be impossible - some fields
	// cannot be modified and require a destroy/recreate.
	BlockDestruction LifecycleParam = iota
	// BlockAcquire will error if the resource already exists.
	BlockAcquire
	// BlockCreation will error if the resource does not exist.
	BlockCreation
	// BlockModification will error if the resource is not in the desired state.
	BlockModification
	// IgnoreIfMissing does not create (and does not error) if the resource
	// does not exist.
	IgnoreIfMissing
	// NoopOnDestroy does not destroy the resource, even if Delete() is
	// called.
	NoopOnDestroy
)

// HasLifecycleParam returns whether the given slice has the requested param.
func HasLifecycleParam(lps []LifecycleParam, p LifecycleParam) bool {
	for _, lp := range lps {
		if lp == p {
			return true
		}
	}
	return false
}

// SprintResourceCompact prints a struct into a compact single line string.
func SprintResourceCompact(v interface{}) string {
	prettyConfig := &pretty.Config{
		Compact:           true,
		IncludeUnexported: true,
	}
	return prettyConfig.Sprint(v)
}

// SprintResource prints a struct into a multiline string to display to readers.
func SprintResource(v interface{}) string {
	prettyConfig := &pretty.Config{
		Diffable:          true, // add line between braces and first/last val
		IncludeUnexported: true,
	}
	return prettyConfig.Sprint(v)
}

// EmptyValue returns an empty value to exclude PARAMETER-type values from
// being expanded
func EmptyValue() (map[string]interface{}, error) {
	return nil, nil
}

/*
 *
 * Default Flatten functions
 *
 * Flatten functions are expected to return the value stored in the interface{},
 * returning a zero value otherwise. For primitive types and arrays that's a nil,
 * but for types like time.Time it may be an empty object.
 *
 * Flattens need to consider the behaviour of https://developers.google.com/discovery/v1/type-format
 * and of https://golang.org/pkg/encoding/json/#Unmarshal. As well, values may
 * get inserted into the JSON map with their correct types already by decoders.
 */

// FlattenInteger turns an interface pointing to an arbitary type into *in64,
// taking into account that it may have been represented as various types per
// https://developers.google.com/discovery/v1/type-format.
// string, float64, int64, and int64 values will return a *int64.
// nil and unrecognised types will return 0.
func FlattenInteger(v interface{}) *int64 {
	if v == nil {
		return nil
	}
	// int64 values are represented as strings
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			return &intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int64(floatVal)
		return &intVal
	}

	// ints won't appear per https://golang.org/pkg/encoding/json/#Unmarshal
	// but may get inserted by decoders in unmarshalResource
	if int32Val, ok := v.(int32); ok {
		intVal := int64(int32Val)
		return &intVal
	}

	if intVal, ok := v.(int64); ok {
		return &intVal
	}

	if machineIntVal, ok := v.(int); ok {
		intVal := int64(machineIntVal)
		return &intVal
	}

	return Int64(0)
}

// FlattenDouble asserts that an interface is a float64 and returns a pointer to it,
// or to 0.0 if the value is invalid.
func FlattenDouble(v interface{}) *float64 {
	if v == nil {
		return nil
	}
	p, ok := v.(float64)
	if !ok {
		return Float64(0.0)
	}

	return &p
}

// FlattenKeyValuePairs asserts that an interface is a map[string]string and
// returns it, or an empty map if the value is invalid.
func FlattenKeyValuePairs(v interface{}) map[string]string {
	if v == nil {
		return nil
	}
	if ss, ok := v.(map[string]string); ok {
		return ss
	}
	p, ok := v.(map[string]interface{})
	if !ok {
		return map[string]string{}
	}

	return assertStringMap(p)
}

// FlattenKeyValueInterface returns a pointer to an interface.
// It can only be used for untyped maps.
func FlattenKeyValueInterface(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}

	if ss, ok := v.(map[string]interface{}); ok {
		return ss
	}

	return map[string]interface{}{}
}

// Returns a map[string]string from a map[string]interface{}
// Non-string values are skipped.
func assertStringMap(mi map[string]interface{}) map[string]string {
	ms := make(map[string]string)
	for k, v := range mi {
		if v == nil {
			continue
		}

		if s, ok := v.(string); ok {
			ms[k] = s
		}
	}
	return ms
}

// FlattenFloatSlice asserts that an interface is a []float64 and returns
// it.
func FlattenFloatSlice(v interface{}) []float64 {
	if v == nil {
		return nil
	}
	p, ok := v.([]interface{})
	if !ok {
		return []float64{}
	}

	return assertFloatSlice(p)
}

// Returns a []float64 from an []interface
// Non-float values are skipped.
func assertFloatSlice(id []interface{}) []float64 {
	dd := []float64{}
	for _, v := range id {
		if v == nil {
			continue
		}

		if d, ok := v.(float64); ok {
			dd = append(dd, d)
		}
	}

	return dd
}

// FlattenIntSlice asserts that an interface is a []int and returns
// it.
func FlattenIntSlice(v interface{}) []int64 {
	if v == nil {
		return nil
	}
	p, ok := v.([]interface{})
	if !ok {
		return []int64{}
	}

	return assertIntSlice(p)
}

// Returns a []int64 from an []interface
// Non-int values are skipped.
func assertIntSlice(id []interface{}) []int64 {
	dd := []int64{}
	for _, v := range id {
		if v == nil {
			continue
		}
		if f, ok := v.(float64); ok {
			dd = append(dd, int64(f))
			continue
		}
		if d, ok := v.(int64); ok {
			dd = append(dd, d)
		}
	}

	return dd
}

// FlattenStringSlice asserts that an interface is a []string and returns
// it.
func FlattenStringSlice(v interface{}) []string {
	if v == nil {
		return nil
	}
	p, ok := v.([]interface{})
	if !ok {
		return []string{}
	}

	return assertStringSlice(p)
}

// Returns a []string from an []interface
// Non-string values are skipped.
func assertStringSlice(is []interface{}) []string {
	ss := []string{}
	for _, v := range is {
		if v == nil {
			continue
		}

		if s, ok := v.(string); ok {
			ss = append(ss, s)
		}
	}

	return ss
}

// FlattenString asserts that an interface is a string and returns a pointer to
// it, or to the empty string if the value is invalid.
func FlattenString(v interface{}) *string {
	if v == nil {
		return nil
	}
	p, ok := v.(string)
	if !ok {
		return String("")
	}

	return &p
}

// FlattenBool asserts that an interface is a bool and returns a pointer to it, or
// a pointer to false if the value is invalid.
func FlattenBool(v interface{}) *bool {
	if v == nil {
		return nil
	}
	p, ok := v.(bool)
	if !ok {
		return Bool(false)
	}

	return &p
}

// FlattenTime asserts that an interface is a time.Time and returns it.
// Time values transmitted in JSON will be an RFC3339 time as per
// https://developers.google.com/discovery/v1/type-format
// Otherwise, it returns the empty time.
func FlattenTime(v interface{}) time.Time {
	if s, ok := v.(string); ok {
		t, err := time.Parse(time.RFC3339, s)
		if err == nil {
			return t
		}
	}

	// In case we inject a time.Time in custom code, convert it.
	if p, ok := v.(time.Time); ok {
		return p
	}

	return time.Time{}
}

// FlattenSecretValue behaves the same way as FlattenString, except that it
// returns nil if the value is not present.
func FlattenSecretValue(v interface{}) *string {
	p, ok := v.(string)
	if !ok {
		return nil
	}

	return &p

}

// ExtractElementFromList takes in bytes corresponding to a json object of the structure
// { "listKey": [{"foo": "bar", ...}, {"foo": "baz", ...}] }
// and returns the first element for which isElement returns true.
// isElement operates on the serialized json representation of each element -
// to the extent that json.Marshal(json.Unmarshal(X)) != X, this may differ from the
// actual elements in the input bytes - but this should be exclusively
// differences which are not semantically significant in json.
func ExtractElementFromList(b []byte, listKey string, isElement func([]byte) bool) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	l, ok := m[listKey]
	if !ok {
		return nil, NotFoundError{Cause: fmt.Errorf("could not find %q in %v, assuming list is empty and returning not found", listKey, m)}
	}
	list, ok := l.([]interface{})
	if !ok {
		return nil, fmt.Errorf("could not convert %v to list", l)
	}
	for _, v := range list {
		if subM, ok := v.(map[string]interface{}); ok {
			if subB, err := json.Marshal(subM); err != nil {
				continue
			} else if isElement(subB) {
				return subB, nil
			}
		}
	}
	// Return a 404-style error.
	return nil, NotFoundError{Cause: fmt.Errorf("could not find a match in %v", list)}
}
