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

// DeleteProtection represents the values of the 'delete_protection' type.
//
// DeleteProtection configuration.
type DeleteProtection struct {
	bitmap_ uint32
	enabled bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *DeleteProtection) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Boolean flag indicating if the cluster should be be using _DeleteProtection_.
//
// By default this is `false`.
//
// To enable it a SREP needs to patch the value through OCM API
func (o *DeleteProtection) Enabled() bool {
	if o != nil && o.bitmap_&1 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Boolean flag indicating if the cluster should be be using _DeleteProtection_.
//
// By default this is `false`.
//
// To enable it a SREP needs to patch the value through OCM API
func (o *DeleteProtection) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.enabled
	}
	return
}

// DeleteProtectionListKind is the name of the type used to represent list of objects of
// type 'delete_protection'.
const DeleteProtectionListKind = "DeleteProtectionList"

// DeleteProtectionListLinkKind is the name of the type used to represent links to list
// of objects of type 'delete_protection'.
const DeleteProtectionListLinkKind = "DeleteProtectionListLink"

// DeleteProtectionNilKind is the name of the type used to nil lists of objects of
// type 'delete_protection'.
const DeleteProtectionListNilKind = "DeleteProtectionListNil"

// DeleteProtectionList is a list of values of the 'delete_protection' type.
type DeleteProtectionList struct {
	href  string
	link  bool
	items []*DeleteProtection
}

// Len returns the length of the list.
func (l *DeleteProtectionList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *DeleteProtectionList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *DeleteProtectionList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *DeleteProtectionList) SetItems(items []*DeleteProtection) {
	l.items = items
}

// Items returns the items of the list.
func (l *DeleteProtectionList) Items() []*DeleteProtection {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *DeleteProtectionList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *DeleteProtectionList) Get(i int) *DeleteProtection {
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
func (l *DeleteProtectionList) Slice() []*DeleteProtection {
	var slice []*DeleteProtection
	if l == nil {
		slice = make([]*DeleteProtection, 0)
	} else {
		slice = make([]*DeleteProtection, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *DeleteProtectionList) Each(f func(item *DeleteProtection) bool) {
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
func (l *DeleteProtectionList) Range(f func(index int, item *DeleteProtection) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
