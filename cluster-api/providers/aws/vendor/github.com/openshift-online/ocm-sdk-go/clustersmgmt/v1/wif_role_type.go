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

// WifRole represents the values of the 'wif_role' type.
type WifRole struct {
	bitmap_     uint32
	permissions []string
	roleId      string
	predefined  bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *WifRole) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Permissions returns the value of the 'permissions' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifRole) Permissions() []string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.permissions
	}
	return nil
}

// GetPermissions returns the value of the 'permissions' attribute and
// a flag indicating if the attribute has a value.
func (o *WifRole) GetPermissions() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.permissions
	}
	return
}

// Predefined returns the value of the 'predefined' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifRole) Predefined() bool {
	if o != nil && o.bitmap_&2 != 0 {
		return o.predefined
	}
	return false
}

// GetPredefined returns the value of the 'predefined' attribute and
// a flag indicating if the attribute has a value.
func (o *WifRole) GetPredefined() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.predefined
	}
	return
}

// RoleId returns the value of the 'role_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifRole) RoleId() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.roleId
	}
	return ""
}

// GetRoleId returns the value of the 'role_id' attribute and
// a flag indicating if the attribute has a value.
func (o *WifRole) GetRoleId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.roleId
	}
	return
}

// WifRoleListKind is the name of the type used to represent list of objects of
// type 'wif_role'.
const WifRoleListKind = "WifRoleList"

// WifRoleListLinkKind is the name of the type used to represent links to list
// of objects of type 'wif_role'.
const WifRoleListLinkKind = "WifRoleListLink"

// WifRoleNilKind is the name of the type used to nil lists of objects of
// type 'wif_role'.
const WifRoleListNilKind = "WifRoleListNil"

// WifRoleList is a list of values of the 'wif_role' type.
type WifRoleList struct {
	href  string
	link  bool
	items []*WifRole
}

// Len returns the length of the list.
func (l *WifRoleList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *WifRoleList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *WifRoleList) Get(i int) *WifRole {
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
func (l *WifRoleList) Slice() []*WifRole {
	var slice []*WifRole
	if l == nil {
		slice = make([]*WifRole, 0)
	} else {
		slice = make([]*WifRole, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *WifRoleList) Each(f func(item *WifRole) bool) {
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
func (l *WifRoleList) Range(f func(index int, item *WifRole) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
