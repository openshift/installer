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
	"reflect"
	re "regexp"
	"strconv"
	"strings"

	glog "github.com/golang/glog"
)

// MoveMapEntry moves the entry at `from` to `to`.  `from` and `to` are slices
// of string keys.  Each key except the last must refer to a map[string]interface{}
// in m - we will descend into m following those keys.  If the maps at the levels
// above the target are empty after the move, they will be deleted.  If there
// are no maps along the path to `to`, they will be created.  If a map above
// the level of the target is missing, nothing will be done.  If the map exists
// but `target` is not present, `nil` will be inserted at `to`.
func MoveMapEntry(m map[string]interface{}, from, to []string) error {
	fetch := m
	// All elements before the last must point to a map[string]interface{} -
	// this ranges over all those elements, so at the end of this loop, we have
	// the map which contains the actual final element to move.
	for _, idx := range from[:len(from)-1] {
		f, ok := fetch[idx]
		if !ok {
			// Nothing to move, so it's not an error not to move it.
			return nil
		}
		fetch, ok = f.(map[string]interface{})
		if !ok {
			return fmt.Errorf("could not fetch %q from %v", idx, fetch)
		}
	}
	value, ok := fetch[from[len(from)-1]]
	if !ok {
		value = nil
	}
	delete(fetch, from[len(from)-1])
	if len(to) > 0 {
		fetch = m
		for _, idx := range to[:len(to)-1] {
			f, ok := fetch[idx]
			if !ok {
				fetch[idx] = make(map[string]interface{})
				f = fetch[idx]
			}
			fetch, ok = f.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%v is not map[string]interface{}", f)
			}
		}
		fetch[to[len(to)-1]] = value
	}
	return deleteIfEmpty(m, from)
}

// GetMapEntry returns the value at `path` from `m`, following the same rules as
// `MoveMapEntry` except that a missing map or value is an error.
func GetMapEntry(m map[string]interface{}, path []string) (interface{}, error) {
	if len(path) == 0 {
		return m, nil
	}
	fetch := m
	// All elements before the last must point to a map[string]interface{} -
	// this ranges over all those elements, so at the end of this loop, we have
	// the map which contains the element to fetch.
	for _, idx := range path[:len(path)-1] {
		f, err := mapEntry(fetch, idx)
		if err != nil {
			return nil, err
		}
		var ok bool
		fetch, ok = f.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("could not fetch %q from %v", idx, fetch)
		}
	}

	value, err := mapEntry(fetch, path[len(path)-1])
	if err != nil {
		return nil, err
	}
	return value, nil
}

// mapEntry grabs item from fetch, and indexes into the array if the [num] notation is present.
func mapEntry(fetch map[string]interface{}, item string) (interface{}, error) {
	// Check if we're fetching from an array.
	arrayRegexp := re.MustCompile(`\[([0-9]*)\]`)
	if arrayRegexp.MatchString(item) {
		field := strings.Split(item, "[")[0]
		items := arrayRegexp.FindAllStringSubmatch(item, 1)
		index, err := strconv.Atoi(items[0][1])
		if err != nil {
			return nil, err
		}

		f, ok := fetch[field]
		if !ok {
			return nil, fmt.Errorf("could not find %q in %v", item, fetch)
		}

		if f == nil {
			return nil, &AttemptToIndexNilArray{FieldName: field}
		}

		fetch, ok := f.([]interface{})
		if !ok {
			return nil, fmt.Errorf("field %s is a %T, not an array", field, f)
		}

		if len(fetch) < index {
			return nil, fmt.Errorf("field %s only has %v elements, needs %v", field, len(fetch), index)
		}

		return fetch[index], nil
	}

	f, ok := fetch[item]
	if !ok {
		return nil, fmt.Errorf("could not find %q in %v", item, fetch)
	}
	return f, nil
}

func deleteIfEmpty(m map[string]interface{}, from []string) error {
	if len(from) > 1 {
		sub, ok := m[from[0]]
		if !ok {
			return fmt.Errorf("could not fetch %q from %v", from[0], m)
		}
		smap, ok := sub.(map[string]interface{})
		if !ok {
			glog.Warningf("In deleting empty map while marshalling, %v not map[string]interface{}", sub)
			return nil
		}
		deleteIfEmpty(smap, from[1:])
	}
	if len(from) >= 1 {
		if sub, ok := m[from[0]]; ok {
			if subm, ok := sub.(map[string]interface{}); ok && len(subm) == 0 {
				delete(m, from[0])
			}
		}
	}
	return nil
}

// PutMapEntry inserts `item` at `path` into `m` - the inverse of GetMapEntry.
func PutMapEntry(m map[string]interface{}, path []string, item interface{}) error {
	if len(path) == 0 {
		return fmt.Errorf("cannot insert value at empty path")
	}
	put := m
	// All elements before the last must point to a map[string]interface{} -
	// this ranges over all those elements, so at the end of this loop, we have
	// the map which contains the element to fetch.
	for _, idx := range path[:len(path)-1] {
		f, ok := put[idx]
		if !ok {
			f = make(map[string]interface{})
			put[idx] = f
		}
		put, ok = f.(map[string]interface{})
		if !ok {
			return fmt.Errorf("could not cast %q from %v as map[string]interface{}", idx, put)
		}
	}
	put[path[len(path)-1]] = item
	return nil
}

// MapFromListOfKeyValues turns a [{"key": k, "value": v}, ...] format-map into a normal string-string map.
// This is useful for a handful of GCP APIs which have chosen to represent maps this way.  We
// expect relatively few of these in newer APIs - it is explicitly against https://aip.dev/apps/2717 -
// ("such a map is represented by a normal JSON object").
// That AIP didn't exist at the time of development of, for instance, Compute v1.
func MapFromListOfKeyValues(rawFetch map[string]interface{}, path []string, keyName, valueName string) (map[string]string, error) {
	i, err := GetMapEntry(rawFetch, path)
	if err != nil {
		// If there's nothing there, it's okay to ignore.
		glog.Warningf("In converting a map to [{\"key\": k, ...}, ...] format, no entry at %q in %v", path, rawFetch)
		return nil, nil
	}
	il, ok := i.([]interface{})
	if !ok {
		return nil, fmt.Errorf("could not cast %v to []interface{}", i)
	}
	var items []map[string]interface{}
	for _, it := range il {
		cast, ok := it.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("could not cast %v to map[string]interface{}", it)
		}
		items = append(items, cast)
	}

	m := make(map[string]string, len(items))
	for _, item := range items {
		key, ok := item[keyName].(string)
		if !ok {
			return nil, fmt.Errorf("could not find 'key' in %v", item)
		}
		value, ok := item[valueName].(string)
		if !ok {
			return nil, fmt.Errorf("could not find 'value' in %v", item)
		}
		m[key] = value
	}
	return m, nil
}

// ListOfKeyValuesFromMap is the opposite of MapFromListOfKeyValues, used in marshalling instead of unmarshalling.
func ListOfKeyValuesFromMap(m map[string]string, keyName, valueName string) ([]map[string]string, error) {
	var items []map[string]string
	for k, v := range m {
		items = append(items, map[string]string{
			keyName:   k,
			valueName: v,
		})
	}
	return items, nil
}

// ListOfKeyValuesFromMapInStruct returns the opposite of MapFromListOfKeyValues, except nested inside an struct under the subfield name.
func ListOfKeyValuesFromMapInStruct(m map[string]string, subfieldName, keyName, valueName string) (map[string][]map[string]string, error) {
	maps, err := ListOfKeyValuesFromMap(m, keyName, valueName)
	if err != nil {
		return nil, err
	}
	return map[string][]map[string]string{
		subfieldName: maps,
	}, nil
}

// ConvertToMap converts the specified object into the map[string]interface{} which can
// be serialized into the same json object as the input object.
func ConvertToMap(obj interface{}) (map[string]interface{}, error) {
	var m map[string]interface{}
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// ValueOrEmptyString takes a scalar or pointer to a scalar and returns either the empty string or its value.
func ValueOrEmptyString(i interface{}) string {
	if i == nil {
		return ""
	}
	v := reflect.ValueOf(i)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.IsValid() {
		switch v.Kind() {
		case reflect.Bool, reflect.Int, reflect.Int64, reflect.Float64, reflect.String:
			return fmt.Sprintf("%v", v.Interface())
		}
	}
	return ""
}

// ValueOrEmptyInt64 returns the value or the default value if the pointer is nil.
func ValueOrEmptyInt64(s *int64) int64 {
	if s == nil {
		return 0
	}
	return *s
}

// ValueOrEmptyBool returns the value or the default value if the pointer is nil.
func ValueOrEmptyBool(s *bool) bool {
	if s == nil {
		return false
	}
	return *s
}

// ValueOrEmptyDouble returns the value or the default value if the pointer is nil.
func ValueOrEmptyDouble(s *float64) float64 {
	if s == nil {
		return 0.0
	}
	return *s
}

// FindStringInArray returns true if value found in array of strings
func FindStringInArray(s string, items []string) bool {
	for _, v := range items {
		if v == s {
			return true
		}
	}
	return false
}

// ValueFromRegexOnField assigns val to the regex value on containerVal if val is unset
func ValueFromRegexOnField(fieldName string, val *string, containerVal *string, regex string) (*string, error) {
	containerGroupedVal := String("")
	// Fetch value from container if the container exists.
	if containerVal != nil && *containerVal != "" {
		r := re.MustCompile(regex)
		m := r.FindStringSubmatch(*containerVal)
		if m != nil && len(m) >= 2 {
			containerGroupedVal = String(m[1])
		} else if val == nil || *val == "" {
			// The regex didn't match and the value doesn't exist.
			return nil, fmt.Errorf("%s field parent has no matching values from regex %s in value %s", fieldName, regex, *containerVal)
		}
	}

	// If value exists + different from what's in container, error.
	if val != nil && *val != "" {
		if containerGroupedVal != nil && *containerGroupedVal != "" && *containerGroupedVal != *val {
			return nil, fmt.Errorf("%s field has conflicting values of %s (from parent) and %s (from self)", fieldName, *containerGroupedVal, *val)
		}
	}

	// If value does not exist, use the value in container.
	if val == nil || *val == "" {
		return containerGroupedVal, nil
	}

	return val, nil
}
