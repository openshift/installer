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

// LimitedSupportReasonTemplateKind is the name of the type used to represent objects
// of type 'limited_support_reason_template'.
const LimitedSupportReasonTemplateKind = "LimitedSupportReasonTemplate"

// LimitedSupportReasonTemplateLinkKind is the name of the type used to represent links
// to objects of type 'limited_support_reason_template'.
const LimitedSupportReasonTemplateLinkKind = "LimitedSupportReasonTemplateLink"

// LimitedSupportReasonTemplateNilKind is the name of the type used to nil references
// to objects of type 'limited_support_reason_template'.
const LimitedSupportReasonTemplateNilKind = "LimitedSupportReasonTemplateNil"

// LimitedSupportReasonTemplate represents the values of the 'limited_support_reason_template' type.
//
// A template for cluster limited support reason.
type LimitedSupportReasonTemplate struct {
	bitmap_ uint32
	id      string
	href    string
	details string
	summary string
}

// Kind returns the name of the type of the object.
func (o *LimitedSupportReasonTemplate) Kind() string {
	if o == nil {
		return LimitedSupportReasonTemplateNilKind
	}
	if o.bitmap_&1 != 0 {
		return LimitedSupportReasonTemplateLinkKind
	}
	return LimitedSupportReasonTemplateKind
}

// Link returns true if this is a link.
func (o *LimitedSupportReasonTemplate) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *LimitedSupportReasonTemplate) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *LimitedSupportReasonTemplate) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *LimitedSupportReasonTemplate) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *LimitedSupportReasonTemplate) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *LimitedSupportReasonTemplate) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Details returns the value of the 'details' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// A detailed description of the reason.
func (o *LimitedSupportReasonTemplate) Details() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.details
	}
	return ""
}

// GetDetails returns the value of the 'details' attribute and
// a flag indicating if the attribute has a value.
//
// A detailed description of the reason.
func (o *LimitedSupportReasonTemplate) GetDetails() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.details
	}
	return
}

// Summary returns the value of the 'summary' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Summary of the reason.
func (o *LimitedSupportReasonTemplate) Summary() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.summary
	}
	return ""
}

// GetSummary returns the value of the 'summary' attribute and
// a flag indicating if the attribute has a value.
//
// Summary of the reason.
func (o *LimitedSupportReasonTemplate) GetSummary() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.summary
	}
	return
}

// LimitedSupportReasonTemplateListKind is the name of the type used to represent list of objects of
// type 'limited_support_reason_template'.
const LimitedSupportReasonTemplateListKind = "LimitedSupportReasonTemplateList"

// LimitedSupportReasonTemplateListLinkKind is the name of the type used to represent links to list
// of objects of type 'limited_support_reason_template'.
const LimitedSupportReasonTemplateListLinkKind = "LimitedSupportReasonTemplateListLink"

// LimitedSupportReasonTemplateNilKind is the name of the type used to nil lists of objects of
// type 'limited_support_reason_template'.
const LimitedSupportReasonTemplateListNilKind = "LimitedSupportReasonTemplateListNil"

// LimitedSupportReasonTemplateList is a list of values of the 'limited_support_reason_template' type.
type LimitedSupportReasonTemplateList struct {
	href  string
	link  bool
	items []*LimitedSupportReasonTemplate
}

// Kind returns the name of the type of the object.
func (l *LimitedSupportReasonTemplateList) Kind() string {
	if l == nil {
		return LimitedSupportReasonTemplateListNilKind
	}
	if l.link {
		return LimitedSupportReasonTemplateListLinkKind
	}
	return LimitedSupportReasonTemplateListKind
}

// Link returns true iif this is a link.
func (l *LimitedSupportReasonTemplateList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *LimitedSupportReasonTemplateList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *LimitedSupportReasonTemplateList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *LimitedSupportReasonTemplateList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *LimitedSupportReasonTemplateList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *LimitedSupportReasonTemplateList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *LimitedSupportReasonTemplateList) SetItems(items []*LimitedSupportReasonTemplate) {
	l.items = items
}

// Items returns the items of the list.
func (l *LimitedSupportReasonTemplateList) Items() []*LimitedSupportReasonTemplate {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *LimitedSupportReasonTemplateList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *LimitedSupportReasonTemplateList) Get(i int) *LimitedSupportReasonTemplate {
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
func (l *LimitedSupportReasonTemplateList) Slice() []*LimitedSupportReasonTemplate {
	var slice []*LimitedSupportReasonTemplate
	if l == nil {
		slice = make([]*LimitedSupportReasonTemplate, 0)
	} else {
		slice = make([]*LimitedSupportReasonTemplate, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *LimitedSupportReasonTemplateList) Each(f func(item *LimitedSupportReasonTemplate) bool) {
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
func (l *LimitedSupportReasonTemplateList) Range(f func(index int, item *LimitedSupportReasonTemplate) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
