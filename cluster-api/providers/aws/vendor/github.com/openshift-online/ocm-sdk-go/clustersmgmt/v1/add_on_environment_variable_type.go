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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

// AddOnEnvironmentVariableKind is the name of the type used to represent objects
// of type 'add_on_environment_variable'.
const AddOnEnvironmentVariableKind = "AddOnEnvironmentVariable"

// AddOnEnvironmentVariableLinkKind is the name of the type used to represent links
// to objects of type 'add_on_environment_variable'.
const AddOnEnvironmentVariableLinkKind = "AddOnEnvironmentVariableLink"

// AddOnEnvironmentVariableNilKind is the name of the type used to nil references
// to objects of type 'add_on_environment_variable'.
const AddOnEnvironmentVariableNilKind = "AddOnEnvironmentVariableNil"

// AddOnEnvironmentVariable represents the values of the 'add_on_environment_variable' type.
//
// Representation of an add-on env object.
type AddOnEnvironmentVariable struct {
	bitmap_ uint32
	id      string
	href    string
	name    string
	value   string
}

// Kind returns the name of the type of the object.
func (o *AddOnEnvironmentVariable) Kind() string {
	if o == nil {
		return AddOnEnvironmentVariableNilKind
	}
	if o.bitmap_&1 != 0 {
		return AddOnEnvironmentVariableLinkKind
	}
	return AddOnEnvironmentVariableKind
}

// Link returns true iif this is a link.
func (o *AddOnEnvironmentVariable) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *AddOnEnvironmentVariable) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AddOnEnvironmentVariable) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AddOnEnvironmentVariable) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AddOnEnvironmentVariable) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AddOnEnvironmentVariable) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the env object.
func (o *AddOnEnvironmentVariable) Name() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the env object.
func (o *AddOnEnvironmentVariable) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.name
	}
	return
}

// Value returns the value of the 'value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Value of the env object.
func (o *AddOnEnvironmentVariable) Value() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.value
	}
	return ""
}

// GetValue returns the value of the 'value' attribute and
// a flag indicating if the attribute has a value.
//
// Value of the env object.
func (o *AddOnEnvironmentVariable) GetValue() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.value
	}
	return
}

// AddOnEnvironmentVariableListKind is the name of the type used to represent list of objects of
// type 'add_on_environment_variable'.
const AddOnEnvironmentVariableListKind = "AddOnEnvironmentVariableList"

// AddOnEnvironmentVariableListLinkKind is the name of the type used to represent links to list
// of objects of type 'add_on_environment_variable'.
const AddOnEnvironmentVariableListLinkKind = "AddOnEnvironmentVariableListLink"

// AddOnEnvironmentVariableNilKind is the name of the type used to nil lists of objects of
// type 'add_on_environment_variable'.
const AddOnEnvironmentVariableListNilKind = "AddOnEnvironmentVariableListNil"

// AddOnEnvironmentVariableList is a list of values of the 'add_on_environment_variable' type.
type AddOnEnvironmentVariableList struct {
	href  string
	link  bool
	items []*AddOnEnvironmentVariable
}

// Kind returns the name of the type of the object.
func (l *AddOnEnvironmentVariableList) Kind() string {
	if l == nil {
		return AddOnEnvironmentVariableListNilKind
	}
	if l.link {
		return AddOnEnvironmentVariableListLinkKind
	}
	return AddOnEnvironmentVariableListKind
}

// Link returns true iif this is a link.
func (l *AddOnEnvironmentVariableList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AddOnEnvironmentVariableList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AddOnEnvironmentVariableList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AddOnEnvironmentVariableList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AddOnEnvironmentVariableList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AddOnEnvironmentVariableList) Get(i int) *AddOnEnvironmentVariable {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *AddOnEnvironmentVariableList) Slice() []*AddOnEnvironmentVariable {
	var slice []*AddOnEnvironmentVariable
	if l == nil {
		slice = make([]*AddOnEnvironmentVariable, 0)
	} else {
		slice = make([]*AddOnEnvironmentVariable, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AddOnEnvironmentVariableList) Each(f func(item *AddOnEnvironmentVariable) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *AddOnEnvironmentVariableList) Range(f func(index int, item *AddOnEnvironmentVariable) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
