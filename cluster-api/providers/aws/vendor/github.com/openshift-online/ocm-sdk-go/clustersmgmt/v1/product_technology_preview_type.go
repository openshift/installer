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

// ProductTechnologyPreviewKind is the name of the type used to represent objects
// of type 'product_technology_preview'.
const ProductTechnologyPreviewKind = "ProductTechnologyPreview"

// ProductTechnologyPreviewLinkKind is the name of the type used to represent links
// to objects of type 'product_technology_preview'.
const ProductTechnologyPreviewLinkKind = "ProductTechnologyPreviewLink"

// ProductTechnologyPreviewNilKind is the name of the type used to nil references
// to objects of type 'product_technology_preview'.
const ProductTechnologyPreviewNilKind = "ProductTechnologyPreviewNil"

// ProductTechnologyPreview represents the values of the 'product_technology_preview' type.
//
// Representation of a product technology preview.
type ProductTechnologyPreview struct {
	bitmap_        uint32
	id             string
	href           string
	additionalText string
	endDate        time.Time
	startDate      time.Time
}

// Kind returns the name of the type of the object.
func (o *ProductTechnologyPreview) Kind() string {
	if o == nil {
		return ProductTechnologyPreviewNilKind
	}
	if o.bitmap_&1 != 0 {
		return ProductTechnologyPreviewLinkKind
	}
	return ProductTechnologyPreviewKind
}

// Link returns true if this is a link.
func (o *ProductTechnologyPreview) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *ProductTechnologyPreview) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ProductTechnologyPreview) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ProductTechnologyPreview) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ProductTechnologyPreview) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ProductTechnologyPreview) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// AdditionalText returns the value of the 'additional_text' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Message associated with this technology preview.
func (o *ProductTechnologyPreview) AdditionalText() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.additionalText
	}
	return ""
}

// GetAdditionalText returns the value of the 'additional_text' attribute and
// a flag indicating if the attribute has a value.
//
// Message associated with this technology preview.
func (o *ProductTechnologyPreview) GetAdditionalText() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.additionalText
	}
	return
}

// EndDate returns the value of the 'end_date' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The end date for this technology preview.
func (o *ProductTechnologyPreview) EndDate() time.Time {
	if o != nil && o.bitmap_&16 != 0 {
		return o.endDate
	}
	return time.Time{}
}

// GetEndDate returns the value of the 'end_date' attribute and
// a flag indicating if the attribute has a value.
//
// The end date for this technology preview.
func (o *ProductTechnologyPreview) GetEndDate() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.endDate
	}
	return
}

// StartDate returns the value of the 'start_date' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The start date for this technology preview.
func (o *ProductTechnologyPreview) StartDate() time.Time {
	if o != nil && o.bitmap_&32 != 0 {
		return o.startDate
	}
	return time.Time{}
}

// GetStartDate returns the value of the 'start_date' attribute and
// a flag indicating if the attribute has a value.
//
// The start date for this technology preview.
func (o *ProductTechnologyPreview) GetStartDate() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.startDate
	}
	return
}

// ProductTechnologyPreviewListKind is the name of the type used to represent list of objects of
// type 'product_technology_preview'.
const ProductTechnologyPreviewListKind = "ProductTechnologyPreviewList"

// ProductTechnologyPreviewListLinkKind is the name of the type used to represent links to list
// of objects of type 'product_technology_preview'.
const ProductTechnologyPreviewListLinkKind = "ProductTechnologyPreviewListLink"

// ProductTechnologyPreviewNilKind is the name of the type used to nil lists of objects of
// type 'product_technology_preview'.
const ProductTechnologyPreviewListNilKind = "ProductTechnologyPreviewListNil"

// ProductTechnologyPreviewList is a list of values of the 'product_technology_preview' type.
type ProductTechnologyPreviewList struct {
	href  string
	link  bool
	items []*ProductTechnologyPreview
}

// Kind returns the name of the type of the object.
func (l *ProductTechnologyPreviewList) Kind() string {
	if l == nil {
		return ProductTechnologyPreviewListNilKind
	}
	if l.link {
		return ProductTechnologyPreviewListLinkKind
	}
	return ProductTechnologyPreviewListKind
}

// Link returns true iif this is a link.
func (l *ProductTechnologyPreviewList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ProductTechnologyPreviewList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ProductTechnologyPreviewList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ProductTechnologyPreviewList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ProductTechnologyPreviewList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ProductTechnologyPreviewList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ProductTechnologyPreviewList) SetItems(items []*ProductTechnologyPreview) {
	l.items = items
}

// Items returns the items of the list.
func (l *ProductTechnologyPreviewList) Items() []*ProductTechnologyPreview {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ProductTechnologyPreviewList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ProductTechnologyPreviewList) Get(i int) *ProductTechnologyPreview {
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
func (l *ProductTechnologyPreviewList) Slice() []*ProductTechnologyPreview {
	var slice []*ProductTechnologyPreview
	if l == nil {
		slice = make([]*ProductTechnologyPreview, 0)
	} else {
		slice = make([]*ProductTechnologyPreview, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ProductTechnologyPreviewList) Each(f func(item *ProductTechnologyPreview) bool) {
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
func (l *ProductTechnologyPreviewList) Range(f func(index int, item *ProductTechnologyPreview) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
