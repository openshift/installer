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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// CCSKind is the name of the type used to represent objects
// of type 'CCS'.
const CCSKind = "CCS"

// CCSLinkKind is the name of the type used to represent links
// to objects of type 'CCS'.
const CCSLinkKind = "CCSLink"

// CCSNilKind is the name of the type used to nil references
// to objects of type 'CCS'.
const CCSNilKind = "CCSNil"

// CCS represents the values of the 'CCS' type.
type CCS struct {
	fieldSet_        []bool
	id               string
	href             string
	disableSCPChecks bool
	enabled          bool
}

// Kind returns the name of the type of the object.
func (o *CCS) Kind() string {
	if o == nil {
		return CCSNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return CCSLinkKind
	}
	return CCSKind
}

// Link returns true if this is a link.
func (o *CCS) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *CCS) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *CCS) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *CCS) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *CCS) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *CCS) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}

	// Check all fields except the link flag (index 0)
	for i := 1; i < len(o.fieldSet_); i++ {
		if o.fieldSet_[i] {
			return false
		}
	}
	return true
}

// DisableSCPChecks returns the value of the 'disable_SCP_checks' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if cloud permissions checks are disabled,
// when attempting installation of the cluster.
func (o *CCS) DisableSCPChecks() bool {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.disableSCPChecks
	}
	return false
}

// GetDisableSCPChecks returns the value of the 'disable_SCP_checks' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if cloud permissions checks are disabled,
// when attempting installation of the cluster.
func (o *CCS) GetDisableSCPChecks() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.disableSCPChecks
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if Customer Cloud Subscription is enabled on the cluster.
func (o *CCS) Enabled() bool {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if Customer Cloud Subscription is enabled on the cluster.
func (o *CCS) GetEnabled() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.enabled
	}
	return
}

// CCSListKind is the name of the type used to represent list of objects of
// type 'CCS'.
const CCSListKind = "CCSList"

// CCSListLinkKind is the name of the type used to represent links to list
// of objects of type 'CCS'.
const CCSListLinkKind = "CCSListLink"

// CCSNilKind is the name of the type used to nil lists of objects of
// type 'CCS'.
const CCSListNilKind = "CCSListNil"

// CCSList is a list of values of the 'CCS' type.
type CCSList struct {
	href  string
	link  bool
	items []*CCS
}

// Kind returns the name of the type of the object.
func (l *CCSList) Kind() string {
	if l == nil {
		return CCSListNilKind
	}
	if l.link {
		return CCSListLinkKind
	}
	return CCSListKind
}

// Link returns true iif this is a link.
func (l *CCSList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *CCSList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *CCSList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *CCSList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *CCSList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *CCSList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *CCSList) SetItems(items []*CCS) {
	l.items = items
}

// Items returns the items of the list.
func (l *CCSList) Items() []*CCS {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *CCSList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *CCSList) Get(i int) *CCS {
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
func (l *CCSList) Slice() []*CCS {
	var slice []*CCS
	if l == nil {
		slice = make([]*CCS, 0)
	} else {
		slice = make([]*CCS, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *CCSList) Each(f func(item *CCS) bool) {
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
func (l *CCSList) Range(f func(index int, item *CCS) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
