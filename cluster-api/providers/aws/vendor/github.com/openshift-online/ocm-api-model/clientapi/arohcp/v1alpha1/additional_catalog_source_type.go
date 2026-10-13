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

// AdditionalCatalogSource represents the values of the 'additional_catalog_source' type.
//
// Representation of an addon catalog source object used by addon versions.
type AdditionalCatalogSource struct {
	fieldSet_ []bool
	id        string
	image     string
	name      string
	enabled   bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AdditionalCatalogSource) Empty() bool {
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

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ID of the additional catalog source
func (o *AdditionalCatalogSource) ID() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// ID of the additional catalog source
func (o *AdditionalCatalogSource) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.id
	}
	return
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates is this additional catalog source is enabled for the addon
func (o *AdditionalCatalogSource) Enabled() bool {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates is this additional catalog source is enabled for the addon
func (o *AdditionalCatalogSource) GetEnabled() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.enabled
	}
	return
}

// Image returns the value of the 'image' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Image of the additional catalog source.
func (o *AdditionalCatalogSource) Image() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.image
	}
	return ""
}

// GetImage returns the value of the 'image' attribute and
// a flag indicating if the attribute has a value.
//
// Image of the additional catalog source.
func (o *AdditionalCatalogSource) GetImage() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.image
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the additional catalog source.
func (o *AdditionalCatalogSource) Name() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the additional catalog source.
func (o *AdditionalCatalogSource) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.name
	}
	return
}

// AdditionalCatalogSourceListKind is the name of the type used to represent list of objects of
// type 'additional_catalog_source'.
const AdditionalCatalogSourceListKind = "AdditionalCatalogSourceList"

// AdditionalCatalogSourceListLinkKind is the name of the type used to represent links to list
// of objects of type 'additional_catalog_source'.
const AdditionalCatalogSourceListLinkKind = "AdditionalCatalogSourceListLink"

// AdditionalCatalogSourceNilKind is the name of the type used to nil lists of objects of
// type 'additional_catalog_source'.
const AdditionalCatalogSourceListNilKind = "AdditionalCatalogSourceListNil"

// AdditionalCatalogSourceList is a list of values of the 'additional_catalog_source' type.
type AdditionalCatalogSourceList struct {
	href  string
	link  bool
	items []*AdditionalCatalogSource
}

// Len returns the length of the list.
func (l *AdditionalCatalogSourceList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AdditionalCatalogSourceList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AdditionalCatalogSourceList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AdditionalCatalogSourceList) SetItems(items []*AdditionalCatalogSource) {
	l.items = items
}

// Items returns the items of the list.
func (l *AdditionalCatalogSourceList) Items() []*AdditionalCatalogSource {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AdditionalCatalogSourceList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AdditionalCatalogSourceList) Get(i int) *AdditionalCatalogSource {
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
func (l *AdditionalCatalogSourceList) Slice() []*AdditionalCatalogSource {
	var slice []*AdditionalCatalogSource
	if l == nil {
		slice = make([]*AdditionalCatalogSource, 0)
	} else {
		slice = make([]*AdditionalCatalogSource, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AdditionalCatalogSourceList) Each(f func(item *AdditionalCatalogSource) bool) {
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
func (l *AdditionalCatalogSourceList) Range(f func(index int, item *AdditionalCatalogSource) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
