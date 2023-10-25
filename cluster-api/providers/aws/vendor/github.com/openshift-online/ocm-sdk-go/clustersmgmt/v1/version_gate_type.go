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

// VersionGateKind is the name of the type used to represent objects
// of type 'version_gate'.
const VersionGateKind = "VersionGate"

// VersionGateLinkKind is the name of the type used to represent links
// to objects of type 'version_gate'.
const VersionGateLinkKind = "VersionGateLink"

// VersionGateNilKind is the name of the type used to nil references
// to objects of type 'version_gate'.
const VersionGateNilKind = "VersionGateNil"

// VersionGate represents the values of the 'version_gate' type.
//
// Representation of an _OpenShift_ version gate.
type VersionGate struct {
	bitmap_            uint32
	id                 string
	href               string
	creationTimestamp  time.Time
	description        string
	documentationURL   string
	label              string
	value              string
	versionRawIDPrefix string
	warningMessage     string
	stsOnly            bool
}

// Kind returns the name of the type of the object.
func (o *VersionGate) Kind() string {
	if o == nil {
		return VersionGateNilKind
	}
	if o.bitmap_&1 != 0 {
		return VersionGateLinkKind
	}
	return VersionGateKind
}

// Link returns true iif this is a link.
func (o *VersionGate) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *VersionGate) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *VersionGate) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *VersionGate) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *VersionGate) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *VersionGate) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// STSOnly returns the value of the 'STS_only' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// STSOnly indicates if this version gate is for STS clusters only
func (o *VersionGate) STSOnly() bool {
	if o != nil && o.bitmap_&8 != 0 {
		return o.stsOnly
	}
	return false
}

// GetSTSOnly returns the value of the 'STS_only' attribute and
// a flag indicating if the attribute has a value.
//
// STSOnly indicates if this version gate is for STS clusters only
func (o *VersionGate) GetSTSOnly() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.stsOnly
	}
	return
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// CreationTimestamp is the date and time when the version gate was created,
// format defined in https://www.ietf.org/rfc/rfc3339.txt[RC3339].
func (o *VersionGate) CreationTimestamp() time.Time {
	if o != nil && o.bitmap_&16 != 0 {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// CreationTimestamp is the date and time when the version gate was created,
// format defined in https://www.ietf.org/rfc/rfc3339.txt[RC3339].
func (o *VersionGate) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.creationTimestamp
	}
	return
}

// Description returns the value of the 'description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Description of the version gate.
func (o *VersionGate) Description() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.description
	}
	return ""
}

// GetDescription returns the value of the 'description' attribute and
// a flag indicating if the attribute has a value.
//
// Description of the version gate.
func (o *VersionGate) GetDescription() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.description
	}
	return
}

// DocumentationURL returns the value of the 'documentation_URL' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// DocumentationURL is the URL for the documentation of the version gate.
func (o *VersionGate) DocumentationURL() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.documentationURL
	}
	return ""
}

// GetDocumentationURL returns the value of the 'documentation_URL' attribute and
// a flag indicating if the attribute has a value.
//
// DocumentationURL is the URL for the documentation of the version gate.
func (o *VersionGate) GetDocumentationURL() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.documentationURL
	}
	return
}

// Label returns the value of the 'label' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Label representing the version gate in OpenShift.
func (o *VersionGate) Label() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.label
	}
	return ""
}

// GetLabel returns the value of the 'label' attribute and
// a flag indicating if the attribute has a value.
//
// Label representing the version gate in OpenShift.
func (o *VersionGate) GetLabel() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.label
	}
	return
}

// Value returns the value of the 'value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Value represents the required value of the label.
func (o *VersionGate) Value() string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.value
	}
	return ""
}

// GetValue returns the value of the 'value' attribute and
// a flag indicating if the attribute has a value.
//
// Value represents the required value of the label.
func (o *VersionGate) GetValue() (value string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.value
	}
	return
}

// VersionRawIDPrefix returns the value of the 'version_raw_ID_prefix' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// VersionRawIDPrefix represents the versions prefix that the gate applies to.
func (o *VersionGate) VersionRawIDPrefix() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.versionRawIDPrefix
	}
	return ""
}

// GetVersionRawIDPrefix returns the value of the 'version_raw_ID_prefix' attribute and
// a flag indicating if the attribute has a value.
//
// VersionRawIDPrefix represents the versions prefix that the gate applies to.
func (o *VersionGate) GetVersionRawIDPrefix() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.versionRawIDPrefix
	}
	return
}

// WarningMessage returns the value of the 'warning_message' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// WarningMessage is a warning that will be displayed to the user before they acknowledge the gate
func (o *VersionGate) WarningMessage() string {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.warningMessage
	}
	return ""
}

// GetWarningMessage returns the value of the 'warning_message' attribute and
// a flag indicating if the attribute has a value.
//
// WarningMessage is a warning that will be displayed to the user before they acknowledge the gate
func (o *VersionGate) GetWarningMessage() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.warningMessage
	}
	return
}

// VersionGateListKind is the name of the type used to represent list of objects of
// type 'version_gate'.
const VersionGateListKind = "VersionGateList"

// VersionGateListLinkKind is the name of the type used to represent links to list
// of objects of type 'version_gate'.
const VersionGateListLinkKind = "VersionGateListLink"

// VersionGateNilKind is the name of the type used to nil lists of objects of
// type 'version_gate'.
const VersionGateListNilKind = "VersionGateListNil"

// VersionGateList is a list of values of the 'version_gate' type.
type VersionGateList struct {
	href  string
	link  bool
	items []*VersionGate
}

// Kind returns the name of the type of the object.
func (l *VersionGateList) Kind() string {
	if l == nil {
		return VersionGateListNilKind
	}
	if l.link {
		return VersionGateListLinkKind
	}
	return VersionGateListKind
}

// Link returns true iif this is a link.
func (l *VersionGateList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *VersionGateList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *VersionGateList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *VersionGateList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *VersionGateList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *VersionGateList) Get(i int) *VersionGate {
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
func (l *VersionGateList) Slice() []*VersionGate {
	var slice []*VersionGate
	if l == nil {
		slice = make([]*VersionGate, 0)
	} else {
		slice = make([]*VersionGate, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *VersionGateList) Each(f func(item *VersionGate) bool) {
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
func (l *VersionGateList) Range(f func(index int, item *VersionGate) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
