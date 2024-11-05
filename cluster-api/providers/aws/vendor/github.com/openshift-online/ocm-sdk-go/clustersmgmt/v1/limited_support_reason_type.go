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

// LimitedSupportReasonKind is the name of the type used to represent objects
// of type 'limited_support_reason'.
const LimitedSupportReasonKind = "LimitedSupportReason"

// LimitedSupportReasonLinkKind is the name of the type used to represent links
// to objects of type 'limited_support_reason'.
const LimitedSupportReasonLinkKind = "LimitedSupportReasonLink"

// LimitedSupportReasonNilKind is the name of the type used to nil references
// to objects of type 'limited_support_reason'.
const LimitedSupportReasonNilKind = "LimitedSupportReasonNil"

// LimitedSupportReason represents the values of the 'limited_support_reason' type.
//
// A reason that a cluster is in limited support.
type LimitedSupportReason struct {
	bitmap_           uint32
	id                string
	href              string
	creationTimestamp time.Time
	details           string
	detectionType     DetectionType
	override          *LimitedSupportReasonOverride
	summary           string
	template          *LimitedSupportReasonTemplate
}

// Kind returns the name of the type of the object.
func (o *LimitedSupportReason) Kind() string {
	if o == nil {
		return LimitedSupportReasonNilKind
	}
	if o.bitmap_&1 != 0 {
		return LimitedSupportReasonLinkKind
	}
	return LimitedSupportReasonKind
}

// Link returns true iif this is a link.
func (o *LimitedSupportReason) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *LimitedSupportReason) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *LimitedSupportReason) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *LimitedSupportReason) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *LimitedSupportReason) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *LimitedSupportReason) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The time the reason was detected.
func (o *LimitedSupportReason) CreationTimestamp() time.Time {
	if o != nil && o.bitmap_&8 != 0 {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// The time the reason was detected.
func (o *LimitedSupportReason) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.creationTimestamp
	}
	return
}

// Details returns the value of the 'details' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// URL with a link to a detailed description of the reason.
func (o *LimitedSupportReason) Details() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.details
	}
	return ""
}

// GetDetails returns the value of the 'details' attribute and
// a flag indicating if the attribute has a value.
//
// URL with a link to a detailed description of the reason.
func (o *LimitedSupportReason) GetDetails() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.details
	}
	return
}

// DetectionType returns the value of the 'detection_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if the reason was detected automatically or manually.
// When creating a new reason this field should be empty or "manual".
func (o *LimitedSupportReason) DetectionType() DetectionType {
	if o != nil && o.bitmap_&32 != 0 {
		return o.detectionType
	}
	return DetectionType("")
}

// GetDetectionType returns the value of the 'detection_type' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if the reason was detected automatically or manually.
// When creating a new reason this field should be empty or "manual".
func (o *LimitedSupportReason) GetDetectionType() (value DetectionType, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.detectionType
	}
	return
}

// Override returns the value of the 'override' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates if the override is enabled
func (o *LimitedSupportReason) Override() *LimitedSupportReasonOverride {
	if o != nil && o.bitmap_&64 != 0 {
		return o.override
	}
	return nil
}

// GetOverride returns the value of the 'override' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates if the override is enabled
func (o *LimitedSupportReason) GetOverride() (value *LimitedSupportReasonOverride, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.override
	}
	return
}

// Summary returns the value of the 'summary' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Summary of the reason.
func (o *LimitedSupportReason) Summary() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.summary
	}
	return ""
}

// GetSummary returns the value of the 'summary' attribute and
// a flag indicating if the attribute has a value.
//
// Summary of the reason.
func (o *LimitedSupportReason) GetSummary() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.summary
	}
	return
}

// Template returns the value of the 'template' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional link to a template with summary and details.
func (o *LimitedSupportReason) Template() *LimitedSupportReasonTemplate {
	if o != nil && o.bitmap_&256 != 0 {
		return o.template
	}
	return nil
}

// GetTemplate returns the value of the 'template' attribute and
// a flag indicating if the attribute has a value.
//
// Optional link to a template with summary and details.
func (o *LimitedSupportReason) GetTemplate() (value *LimitedSupportReasonTemplate, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.template
	}
	return
}

// LimitedSupportReasonListKind is the name of the type used to represent list of objects of
// type 'limited_support_reason'.
const LimitedSupportReasonListKind = "LimitedSupportReasonList"

// LimitedSupportReasonListLinkKind is the name of the type used to represent links to list
// of objects of type 'limited_support_reason'.
const LimitedSupportReasonListLinkKind = "LimitedSupportReasonListLink"

// LimitedSupportReasonNilKind is the name of the type used to nil lists of objects of
// type 'limited_support_reason'.
const LimitedSupportReasonListNilKind = "LimitedSupportReasonListNil"

// LimitedSupportReasonList is a list of values of the 'limited_support_reason' type.
type LimitedSupportReasonList struct {
	href  string
	link  bool
	items []*LimitedSupportReason
}

// Kind returns the name of the type of the object.
func (l *LimitedSupportReasonList) Kind() string {
	if l == nil {
		return LimitedSupportReasonListNilKind
	}
	if l.link {
		return LimitedSupportReasonListLinkKind
	}
	return LimitedSupportReasonListKind
}

// Link returns true iif this is a link.
func (l *LimitedSupportReasonList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *LimitedSupportReasonList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *LimitedSupportReasonList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *LimitedSupportReasonList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *LimitedSupportReasonList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *LimitedSupportReasonList) Get(i int) *LimitedSupportReason {
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
func (l *LimitedSupportReasonList) Slice() []*LimitedSupportReason {
	var slice []*LimitedSupportReason
	if l == nil {
		slice = make([]*LimitedSupportReason, 0)
	} else {
		slice = make([]*LimitedSupportReason, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *LimitedSupportReasonList) Each(f func(item *LimitedSupportReason) bool) {
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
func (l *LimitedSupportReasonList) Range(f func(index int, item *LimitedSupportReason) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
