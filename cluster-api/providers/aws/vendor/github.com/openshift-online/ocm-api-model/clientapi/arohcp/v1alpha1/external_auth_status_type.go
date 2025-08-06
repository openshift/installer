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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

// ExternalAuthStatus represents the values of the 'external_auth_status' type.
//
// Representation of the status of an external authentication provider.
type ExternalAuthStatus struct {
	fieldSet_ []bool
	message   string
	state     *ExternalAuthState
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ExternalAuthStatus) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}
	for _, set := range o.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// Message returns the value of the 'message' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// A descriptive message providing additional context about the current
// state of the external authentication provider.
func (o *ExternalAuthStatus) Message() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.message
	}
	return ""
}

// GetMessage returns the value of the 'message' attribute and
// a flag indicating if the attribute has a value.
//
// A descriptive message providing additional context about the current
// state of the external authentication provider.
func (o *ExternalAuthStatus) GetMessage() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.message
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The current state of the external authentication provider.
func (o *ExternalAuthStatus) State() *ExternalAuthState {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.state
	}
	return nil
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
//
// The current state of the external authentication provider.
func (o *ExternalAuthStatus) GetState() (value *ExternalAuthState, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.state
	}
	return
}

// ExternalAuthStatusListKind is the name of the type used to represent list of objects of
// type 'external_auth_status'.
const ExternalAuthStatusListKind = "ExternalAuthStatusList"

// ExternalAuthStatusListLinkKind is the name of the type used to represent links to list
// of objects of type 'external_auth_status'.
const ExternalAuthStatusListLinkKind = "ExternalAuthStatusListLink"

// ExternalAuthStatusNilKind is the name of the type used to nil lists of objects of
// type 'external_auth_status'.
const ExternalAuthStatusListNilKind = "ExternalAuthStatusListNil"

// ExternalAuthStatusList is a list of values of the 'external_auth_status' type.
type ExternalAuthStatusList struct {
	href  string
	link  bool
	items []*ExternalAuthStatus
}

// Len returns the length of the list.
func (l *ExternalAuthStatusList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ExternalAuthStatusList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ExternalAuthStatusList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ExternalAuthStatusList) SetItems(items []*ExternalAuthStatus) {
	l.items = items
}

// Items returns the items of the list.
func (l *ExternalAuthStatusList) Items() []*ExternalAuthStatus {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ExternalAuthStatusList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ExternalAuthStatusList) Get(i int) *ExternalAuthStatus {
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
func (l *ExternalAuthStatusList) Slice() []*ExternalAuthStatus {
	var slice []*ExternalAuthStatus
	if l == nil {
		slice = make([]*ExternalAuthStatus, 0)
	} else {
		slice = make([]*ExternalAuthStatus, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ExternalAuthStatusList) Each(f func(item *ExternalAuthStatus) bool) {
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
func (l *ExternalAuthStatusList) Range(f func(index int, item *ExternalAuthStatus) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
