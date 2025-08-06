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

// ReleaseImageDetails represents the values of the 'release_image_details' type.
type ReleaseImageDetails struct {
	fieldSet_         []bool
	availableUpgrades []string
	releaseImage      string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ReleaseImageDetails) Empty() bool {
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

// AvailableUpgrades returns the value of the 'available_upgrades' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AvailableUpgrades is the list of versions this version can be upgraded to.
func (o *ReleaseImageDetails) AvailableUpgrades() []string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.availableUpgrades
	}
	return nil
}

// GetAvailableUpgrades returns the value of the 'available_upgrades' attribute and
// a flag indicating if the attribute has a value.
//
// AvailableUpgrades is the list of versions this version can be upgraded to.
func (o *ReleaseImageDetails) GetAvailableUpgrades() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.availableUpgrades
	}
	return
}

// ReleaseImage returns the value of the 'release_image' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ReleaseImage contains the URI of Openshift release image.
func (o *ReleaseImageDetails) ReleaseImage() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.releaseImage
	}
	return ""
}

// GetReleaseImage returns the value of the 'release_image' attribute and
// a flag indicating if the attribute has a value.
//
// ReleaseImage contains the URI of Openshift release image.
func (o *ReleaseImageDetails) GetReleaseImage() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.releaseImage
	}
	return
}

// ReleaseImageDetailsListKind is the name of the type used to represent list of objects of
// type 'release_image_details'.
const ReleaseImageDetailsListKind = "ReleaseImageDetailsList"

// ReleaseImageDetailsListLinkKind is the name of the type used to represent links to list
// of objects of type 'release_image_details'.
const ReleaseImageDetailsListLinkKind = "ReleaseImageDetailsListLink"

// ReleaseImageDetailsNilKind is the name of the type used to nil lists of objects of
// type 'release_image_details'.
const ReleaseImageDetailsListNilKind = "ReleaseImageDetailsListNil"

// ReleaseImageDetailsList is a list of values of the 'release_image_details' type.
type ReleaseImageDetailsList struct {
	href  string
	link  bool
	items []*ReleaseImageDetails
}

// Len returns the length of the list.
func (l *ReleaseImageDetailsList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ReleaseImageDetailsList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ReleaseImageDetailsList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ReleaseImageDetailsList) SetItems(items []*ReleaseImageDetails) {
	l.items = items
}

// Items returns the items of the list.
func (l *ReleaseImageDetailsList) Items() []*ReleaseImageDetails {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ReleaseImageDetailsList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ReleaseImageDetailsList) Get(i int) *ReleaseImageDetails {
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
func (l *ReleaseImageDetailsList) Slice() []*ReleaseImageDetails {
	var slice []*ReleaseImageDetails
	if l == nil {
		slice = make([]*ReleaseImageDetails, 0)
	} else {
		slice = make([]*ReleaseImageDetails, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ReleaseImageDetailsList) Each(f func(item *ReleaseImageDetails) bool) {
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
func (l *ReleaseImageDetailsList) Range(f func(index int, item *ReleaseImageDetails) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
