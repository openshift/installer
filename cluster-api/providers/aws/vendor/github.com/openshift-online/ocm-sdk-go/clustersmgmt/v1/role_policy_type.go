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

// RolePolicy represents the values of the 'role_policy' type.
type RolePolicy struct {
	bitmap_ uint32
	arn     string
	name    string
	type_   string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *RolePolicy) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Arn returns the value of the 'arn' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RolePolicy) Arn() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.arn
	}
	return ""
}

// GetArn returns the value of the 'arn' attribute and
// a flag indicating if the attribute has a value.
func (o *RolePolicy) GetArn() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.arn
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RolePolicy) Name() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
func (o *RolePolicy) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.name
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *RolePolicy) Type() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.type_
	}
	return ""
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
func (o *RolePolicy) GetType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.type_
	}
	return
}

// RolePolicyListKind is the name of the type used to represent list of objects of
// type 'role_policy'.
const RolePolicyListKind = "RolePolicyList"

// RolePolicyListLinkKind is the name of the type used to represent links to list
// of objects of type 'role_policy'.
const RolePolicyListLinkKind = "RolePolicyListLink"

// RolePolicyNilKind is the name of the type used to nil lists of objects of
// type 'role_policy'.
const RolePolicyListNilKind = "RolePolicyListNil"

// RolePolicyList is a list of values of the 'role_policy' type.
type RolePolicyList struct {
	href  string
	link  bool
	items []*RolePolicy
}

// Len returns the length of the list.
func (l *RolePolicyList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *RolePolicyList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *RolePolicyList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *RolePolicyList) SetItems(items []*RolePolicy) {
	l.items = items
}

// Items returns the items of the list.
func (l *RolePolicyList) Items() []*RolePolicy {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *RolePolicyList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *RolePolicyList) Get(i int) *RolePolicy {
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
func (l *RolePolicyList) Slice() []*RolePolicy {
	var slice []*RolePolicy
	if l == nil {
		slice = make([]*RolePolicy, 0)
	} else {
		slice = make([]*RolePolicy, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *RolePolicyList) Each(f func(item *RolePolicy) bool) {
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
func (l *RolePolicyList) Range(f func(index int, item *RolePolicy) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
