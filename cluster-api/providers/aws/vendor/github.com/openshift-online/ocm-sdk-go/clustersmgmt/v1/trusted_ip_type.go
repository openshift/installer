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

// TrustedIpKind is the name of the type used to represent objects
// of type 'trusted_ip'.
const TrustedIpKind = "TrustedIp"

// TrustedIpLinkKind is the name of the type used to represent links
// to objects of type 'trusted_ip'.
const TrustedIpLinkKind = "TrustedIpLink"

// TrustedIpNilKind is the name of the type used to nil references
// to objects of type 'trusted_ip'.
const TrustedIpNilKind = "TrustedIpNil"

// TrustedIp represents the values of the 'trusted_ip' type.
//
// Representation of a trusted ip address in clusterdeployment.
type TrustedIp struct {
	bitmap_ uint32
	id      string
	href    string
	enabled bool
}

// Kind returns the name of the type of the object.
func (o *TrustedIp) Kind() string {
	if o == nil {
		return TrustedIpNilKind
	}
	if o.bitmap_&1 != 0 {
		return TrustedIpLinkKind
	}
	return TrustedIpKind
}

// Link returns true iif this is a link.
func (o *TrustedIp) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *TrustedIp) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *TrustedIp) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *TrustedIp) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *TrustedIp) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *TrustedIp) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The boolean set to show if the ip is enabled.
func (o *TrustedIp) Enabled() bool {
	if o != nil && o.bitmap_&8 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// The boolean set to show if the ip is enabled.
func (o *TrustedIp) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.enabled
	}
	return
}

// TrustedIpListKind is the name of the type used to represent list of objects of
// type 'trusted_ip'.
const TrustedIpListKind = "TrustedIpList"

// TrustedIpListLinkKind is the name of the type used to represent links to list
// of objects of type 'trusted_ip'.
const TrustedIpListLinkKind = "TrustedIpListLink"

// TrustedIpNilKind is the name of the type used to nil lists of objects of
// type 'trusted_ip'.
const TrustedIpListNilKind = "TrustedIpListNil"

// TrustedIpList is a list of values of the 'trusted_ip' type.
type TrustedIpList struct {
	href  string
	link  bool
	items []*TrustedIp
}

// Kind returns the name of the type of the object.
func (l *TrustedIpList) Kind() string {
	if l == nil {
		return TrustedIpListNilKind
	}
	if l.link {
		return TrustedIpListLinkKind
	}
	return TrustedIpListKind
}

// Link returns true iif this is a link.
func (l *TrustedIpList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *TrustedIpList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *TrustedIpList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *TrustedIpList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *TrustedIpList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *TrustedIpList) Get(i int) *TrustedIp {
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
func (l *TrustedIpList) Slice() []*TrustedIp {
	var slice []*TrustedIp
	if l == nil {
		slice = make([]*TrustedIp, 0)
	} else {
		slice = make([]*TrustedIp, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *TrustedIpList) Each(f func(item *TrustedIp) bool) {
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
func (l *TrustedIpList) Range(f func(index int, item *TrustedIp) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
