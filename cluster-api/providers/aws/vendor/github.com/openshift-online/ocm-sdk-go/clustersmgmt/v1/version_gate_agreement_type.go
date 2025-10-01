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

import (
	time "time"
)

// VersionGateAgreementKind is the name of the type used to represent objects
// of type 'version_gate_agreement'.
const VersionGateAgreementKind = "VersionGateAgreement"

// VersionGateAgreementLinkKind is the name of the type used to represent links
// to objects of type 'version_gate_agreement'.
const VersionGateAgreementLinkKind = "VersionGateAgreementLink"

// VersionGateAgreementNilKind is the name of the type used to nil references
// to objects of type 'version_gate_agreement'.
const VersionGateAgreementNilKind = "VersionGateAgreementNil"

// VersionGateAgreement represents the values of the 'version_gate_agreement' type.
//
// VersionGateAgreement represents a version gate that the user agreed to for a specific cluster.
type VersionGateAgreement struct {
	bitmap_         uint32
	id              string
	href            string
	agreedTimestamp time.Time
	versionGate     *VersionGate
}

// Kind returns the name of the type of the object.
func (o *VersionGateAgreement) Kind() string {
	if o == nil {
		return VersionGateAgreementNilKind
	}
	if o.bitmap_&1 != 0 {
		return VersionGateAgreementLinkKind
	}
	return VersionGateAgreementKind
}

// Link returns true if this is a link.
func (o *VersionGateAgreement) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *VersionGateAgreement) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *VersionGateAgreement) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *VersionGateAgreement) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *VersionGateAgreement) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *VersionGateAgreement) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// AgreedTimestamp returns the value of the 'agreed_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The time the user agreed to the version gate
func (o *VersionGateAgreement) AgreedTimestamp() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.agreedTimestamp
	}
	return time.Time{}
}

// GetAgreedTimestamp returns the value of the 'agreed_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// The time the user agreed to the version gate
func (o *VersionGateAgreement) GetAgreedTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.agreedTimestamp
	}
	return
}

// VersionGate returns the value of the 'version_gate' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// link to the version gate that the user agreed to
func (o *VersionGateAgreement) VersionGate() *VersionGate {
	if o != nil && o.bitmap_&16 != 0 {
		return o.versionGate
	}
	return nil
}

// GetVersionGate returns the value of the 'version_gate' attribute and
// a flag indicating if the attribute has a value.
//
// link to the version gate that the user agreed to
func (o *VersionGateAgreement) GetVersionGate() (value *VersionGate, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.versionGate
	}
	return
}

// VersionGateAgreementListKind is the name of the type used to represent list of objects of
// type 'version_gate_agreement'.
const VersionGateAgreementListKind = "VersionGateAgreementList"

// VersionGateAgreementListLinkKind is the name of the type used to represent links to list
// of objects of type 'version_gate_agreement'.
const VersionGateAgreementListLinkKind = "VersionGateAgreementListLink"

// VersionGateAgreementNilKind is the name of the type used to nil lists of objects of
// type 'version_gate_agreement'.
const VersionGateAgreementListNilKind = "VersionGateAgreementListNil"

// VersionGateAgreementList is a list of values of the 'version_gate_agreement' type.
type VersionGateAgreementList struct {
	href  string
	link  bool
	items []*VersionGateAgreement
}

// Kind returns the name of the type of the object.
func (l *VersionGateAgreementList) Kind() string {
	if l == nil {
		return VersionGateAgreementListNilKind
	}
	if l.link {
		return VersionGateAgreementListLinkKind
	}
	return VersionGateAgreementListKind
}

// Link returns true iif this is a link.
func (l *VersionGateAgreementList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *VersionGateAgreementList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *VersionGateAgreementList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *VersionGateAgreementList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *VersionGateAgreementList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *VersionGateAgreementList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *VersionGateAgreementList) SetItems(items []*VersionGateAgreement) {
	l.items = items
}

// Items returns the items of the list.
func (l *VersionGateAgreementList) Items() []*VersionGateAgreement {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *VersionGateAgreementList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *VersionGateAgreementList) Get(i int) *VersionGateAgreement {
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
func (l *VersionGateAgreementList) Slice() []*VersionGateAgreement {
	var slice []*VersionGateAgreement
	if l == nil {
		slice = make([]*VersionGateAgreement, 0)
	} else {
		slice = make([]*VersionGateAgreement, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *VersionGateAgreementList) Each(f func(item *VersionGateAgreement) bool) {
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
func (l *VersionGateAgreementList) Range(f func(index int, item *VersionGateAgreement) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
