// Copyright 2016 CoreOS, Inc.
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

package merge

import (
	"reflect"

	"github.com/coreos/ignition/v2/config/util"
)

// Rules of Config Merging:
// 1) Parent and child configs must be the same version/type
// 2) Only valid configs can be merged
// 3) It is possible to merge two valid configs and get an invalid config
// 3) For structs:
//   a) Members that are structs get merged recursively (i.e. ignition.storage)
//   b) Members that are primitives get replaced by the child's member (e.g. ignition.storage.files[i].path)
//   c) Members that are pointers only get replaced by the child's value if the child's value is non-nil (e.g. ignition.config.replace.source)
//   d) List merging of a list with IgnoreDuplicates: append the lists (e.g. ignition.storage.files[i].append)
//   e) List merging of a list not merged with other lists: merge any entries with the same Key() and append the others (e.g. ignition.storage.filesystems by path)
//   f) List merging of a list merged with other lists: (e.g. ignition.storage.{files,links,directories} by path)
//      - merge entries with the same Key() that are in the same list
//      - remove entries from the parent with the same Key() that are not in the same list
//      - append entries that are unique to the child

// appendToSlice is a helper that appends to a slice without returning a new one.
// panics if len >= cap
func appendToSlice(s, v reflect.Value) {
	s.SetLen(s.Len() + 1)
	s.Index(s.Len() - 1).Set(v)
}

type handleKey struct {
	handle string
	key    string
}

// structInfo holds information about a struct being processed and has helper methods for querying that
// information in a way that is more clear what the intent is.
type structInfo struct {
	// set of field names to not do duplicate merging on
	ignoreDups map[string]struct{}

	// map from field names to a handle indicating all those with the same handle should have duplication
	// checking done across all fields that share that handle
	mergedKeys map[string]string

	// map from each handle + key() value to what list it came from
	keysToValues map[handleKey]reflect.Value

	// map from each handle + key() to the list it came from
	keysToLists map[handleKey]string
}

// returns if this field should not do duplicate checking/merging
func (s structInfo) ignoreField(name string) bool {
	_, ignore := s.ignoreDups[name]
	return ignore
}

// getChildEntryByKey takes the name of a field (not handle) in the parent and a key and looks that entry
// up in the child. It will look up across all slices that share the same handle. It return the value and
// name of the field in the child it was found in. The bool indicates whether it was found.
func (s structInfo) getChildEntryByKey(fieldName, key string) (reflect.Value, string, bool) {
	handle := fieldName
	if tmp, ok := s.mergedKeys[fieldName]; ok {
		handle = tmp
	}

	hkey := handleKey{
		handle: handle,
		key:    key,
	}
	if v, ok := s.keysToValues[hkey]; ok {
		return v, s.keysToLists[hkey], true
	}
	return reflect.Value{}, "", false
}

func newStructInfo(parent, child reflect.Value) structInfo {
	ignoreDups := map[string]struct{}{}
	if ignorer, ok := parent.Interface().(util.IgnoresDups); ok {
		ignoreDups = ignorer.IgnoreDuplicates()
	}

	mergedKeys := map[string]string{}
	if merger, ok := parent.Interface().(util.MergesKeys); ok {
		mergedKeys = merger.MergedKeys()
	}

	keysToValues := map[handleKey]reflect.Value{}
	keysToLists := map[handleKey]string{}
	for i := 0; i < child.NumField(); i++ {
		field := child.Field(i)
		if field.Kind() != reflect.Slice {
			continue
		}

		fieldName := child.Type().Field(i).Name
		if _, ok := ignoreDups[fieldName]; ok {
			continue
		}

		handle := fieldName
		if tmp, ok := mergedKeys[handle]; ok {
			handle = tmp
		}

		for j := 0; j < field.Len(); j++ {
			v := field.Index(j)
			hkey := handleKey{
				handle: handle,
				key:    util.CallKey(v),
			}
			keysToValues[hkey] = v
			keysToLists[hkey] = fieldName
		}
	}

	return structInfo{
		ignoreDups:   ignoreDups,
		mergedKeys:   mergedKeys,
		keysToValues: keysToValues,
		keysToLists:  keysToLists,
	}
}

// MergeStruct is intended for use by config/vX_Y/ packages only. They should expose their own Merge() that is properly
// typed. Use that one instead.
// parent and child MUST be the same type
func MergeStruct(parent, child reflect.Value) reflect.Value {
	// use New() so it's settable, addr-able, etc
	result := reflect.New(parent.Type()).Elem()
	info := newStructInfo(parent, child)

	for i := 0; i < parent.NumField(); i++ {
		fieldName := parent.Type().Field(i).Name
		parentField := parent.Field(i)
		childField := child.Field(i)
		resultField := result.Field(i)

		kind := parentField.Kind()
		switch {
		case util.IsPrimitive(kind):
			resultField.Set(childField)
		case kind == reflect.Ptr && childField.IsNil():
			resultField.Set(parentField)
		case kind == reflect.Ptr && !childField.IsNil():
			resultField.Set(childField)
		case kind == reflect.Struct:
			resultField.Set(MergeStruct(parentField, childField))
		case kind == reflect.Slice && info.ignoreField(fieldName):
			if parentField.Len()+childField.Len() == 0 {
				continue
			}
			resultField.Set(reflect.AppendSlice(parentField, childField))
		case kind == reflect.Slice && !info.ignoreField(fieldName):
			// ooph, this is a doosey
			maxlen := parentField.Len() + childField.Len()
			if maxlen == 0 {
				continue
			}
			resultField.Set(reflect.MakeSlice(parentField.Type(), 0, parentField.Len()+childField.Len()))
			parentKeys := getKeySet(parentField)

			for i := 0; i < parentField.Len(); i++ {
				parentItem := parentField.Index(i)
				key := util.CallKey(parentItem)

				if childItem, childList, ok := info.getChildEntryByKey(fieldName, key); ok {
					if childList == fieldName {
						// case 1: in child config in same list
						if childItem.Kind() == reflect.Struct {
							appendToSlice(resultField, MergeStruct(parentItem, childItem))
						} else if util.IsPrimitive(childItem.Kind()) {
							appendToSlice(resultField, childItem)
						} else {
							panic("List of pointers or slices or something else weird")
						}
					} else {
						// case 2: in child config in different list. Do nothing since it'll be handled iterating over that list
					}
				} else {
					// case 3: not in child config, append it
					appendToSlice(resultField, parentItem)
				}
			}
			for i := 0; i < childField.Len(); i++ {
				childItem := childField.Index(i)
				key := util.CallKey(childItem)
				if _, alreadyMerged := parentKeys[key]; !alreadyMerged {
					// We only check the parentMap for this field. If the parent had a matching entry in a differnt field
					// then it would be skipped as case 2 above
					appendToSlice(resultField, childItem)
				}
			}
		default:
			panic("unreachable code reached")
		}
	}

	return result
}

// getKeySet takes a value of a slice and returns the set of all the Key() values in that slice
func getKeySet(list reflect.Value) map[string]struct{} {
	m := map[string]struct{}{}
	for i := 0; i < list.Len(); i++ {
		m[util.CallKey(list.Index(i))] = struct{}{}
	}
	return m
}
