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

// CIDRBlockAllowAccess represents the values of the 'CIDR_block_allow_access' type.
type CIDRBlockAllowAccess struct {
	fieldSet_ []bool
	mode      string
	values    []string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *CIDRBlockAllowAccess) Empty() bool {
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

// Mode returns the value of the 'mode' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// There are two modes: "allow_all" and "allow_list"; if "allow_list" is provided than a non-empty 'values' list must be provided.
// Otherwise, if "allow_all" is provided then 'values' list should be omitted.
func (o *CIDRBlockAllowAccess) Mode() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.mode
	}
	return ""
}

// GetMode returns the value of the 'mode' attribute and
// a flag indicating if the attribute has a value.
//
// There are two modes: "allow_all" and "allow_list"; if "allow_list" is provided than a non-empty 'values' list must be provided.
// Otherwise, if "allow_all" is provided then 'values' list should be omitted.
func (o *CIDRBlockAllowAccess) GetMode() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.mode
	}
	return
}

// Values returns the value of the 'values' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The 'values' list should contain a CIDR block list (An IPV4 address range in the format `<ipv4_address>/<network_mask>`).
// The maximum number of CIDR blocks supported is 500. The CIDR blocks should be non-overlapping and valid.
// The value "0.0.0.0/0" is not considered a valid value, as the user can use "allow_all" mode to indicate this behavior.
// The values should not contain the set of Private IP address ranges.
func (o *CIDRBlockAllowAccess) Values() []string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.values
	}
	return nil
}

// GetValues returns the value of the 'values' attribute and
// a flag indicating if the attribute has a value.
//
// The 'values' list should contain a CIDR block list (An IPV4 address range in the format `<ipv4_address>/<network_mask>`).
// The maximum number of CIDR blocks supported is 500. The CIDR blocks should be non-overlapping and valid.
// The value "0.0.0.0/0" is not considered a valid value, as the user can use "allow_all" mode to indicate this behavior.
// The values should not contain the set of Private IP address ranges.
func (o *CIDRBlockAllowAccess) GetValues() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.values
	}
	return
}

// CIDRBlockAllowAccessListKind is the name of the type used to represent list of objects of
// type 'CIDR_block_allow_access'.
const CIDRBlockAllowAccessListKind = "CIDRBlockAllowAccessList"

// CIDRBlockAllowAccessListLinkKind is the name of the type used to represent links to list
// of objects of type 'CIDR_block_allow_access'.
const CIDRBlockAllowAccessListLinkKind = "CIDRBlockAllowAccessListLink"

// CIDRBlockAllowAccessNilKind is the name of the type used to nil lists of objects of
// type 'CIDR_block_allow_access'.
const CIDRBlockAllowAccessListNilKind = "CIDRBlockAllowAccessListNil"

// CIDRBlockAllowAccessList is a list of values of the 'CIDR_block_allow_access' type.
type CIDRBlockAllowAccessList struct {
	href  string
	link  bool
	items []*CIDRBlockAllowAccess
}

// Len returns the length of the list.
func (l *CIDRBlockAllowAccessList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *CIDRBlockAllowAccessList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *CIDRBlockAllowAccessList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *CIDRBlockAllowAccessList) SetItems(items []*CIDRBlockAllowAccess) {
	l.items = items
}

// Items returns the items of the list.
func (l *CIDRBlockAllowAccessList) Items() []*CIDRBlockAllowAccess {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *CIDRBlockAllowAccessList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *CIDRBlockAllowAccessList) Get(i int) *CIDRBlockAllowAccess {
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
func (l *CIDRBlockAllowAccessList) Slice() []*CIDRBlockAllowAccess {
	var slice []*CIDRBlockAllowAccess
	if l == nil {
		slice = make([]*CIDRBlockAllowAccess, 0)
	} else {
		slice = make([]*CIDRBlockAllowAccess, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *CIDRBlockAllowAccessList) Each(f func(item *CIDRBlockAllowAccess) bool) {
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
func (l *CIDRBlockAllowAccessList) Range(f func(index int, item *CIDRBlockAllowAccess) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
