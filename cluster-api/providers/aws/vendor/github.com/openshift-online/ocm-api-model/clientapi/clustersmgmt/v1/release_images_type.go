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

// ReleaseImages represents the values of the 'release_images' type.
type ReleaseImages struct {
	fieldSet_ []bool
	arm64     *ReleaseImageDetails
	multi     *ReleaseImageDetails
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ReleaseImages) Empty() bool {
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

// ARM64 returns the value of the 'ARM64' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Arm64 will contain the reference for the arm64 image which will be used for cluster deployments
func (o *ReleaseImages) ARM64() *ReleaseImageDetails {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.arm64
	}
	return nil
}

// GetARM64 returns the value of the 'ARM64' attribute and
// a flag indicating if the attribute has a value.
//
// Arm64 will contain the reference for the arm64 image which will be used for cluster deployments
func (o *ReleaseImages) GetARM64() (value *ReleaseImageDetails, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.arm64
	}
	return
}

// Multi returns the value of the 'multi' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Multi will contain the reference for the multi image which will be used for cluster deployments
func (o *ReleaseImages) Multi() *ReleaseImageDetails {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.multi
	}
	return nil
}

// GetMulti returns the value of the 'multi' attribute and
// a flag indicating if the attribute has a value.
//
// Multi will contain the reference for the multi image which will be used for cluster deployments
func (o *ReleaseImages) GetMulti() (value *ReleaseImageDetails, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.multi
	}
	return
}

// ReleaseImagesListKind is the name of the type used to represent list of objects of
// type 'release_images'.
const ReleaseImagesListKind = "ReleaseImagesList"

// ReleaseImagesListLinkKind is the name of the type used to represent links to list
// of objects of type 'release_images'.
const ReleaseImagesListLinkKind = "ReleaseImagesListLink"

// ReleaseImagesNilKind is the name of the type used to nil lists of objects of
// type 'release_images'.
const ReleaseImagesListNilKind = "ReleaseImagesListNil"

// ReleaseImagesList is a list of values of the 'release_images' type.
type ReleaseImagesList struct {
	href  string
	link  bool
	items []*ReleaseImages
}

// Len returns the length of the list.
func (l *ReleaseImagesList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ReleaseImagesList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ReleaseImagesList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ReleaseImagesList) SetItems(items []*ReleaseImages) {
	l.items = items
}

// Items returns the items of the list.
func (l *ReleaseImagesList) Items() []*ReleaseImages {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ReleaseImagesList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ReleaseImagesList) Get(i int) *ReleaseImages {
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
func (l *ReleaseImagesList) Slice() []*ReleaseImages {
	var slice []*ReleaseImages
	if l == nil {
		slice = make([]*ReleaseImages, 0)
	} else {
		slice = make([]*ReleaseImages, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ReleaseImagesList) Each(f func(item *ReleaseImages) bool) {
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
func (l *ReleaseImagesList) Range(f func(index int, item *ReleaseImages) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
