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

// AMIOverrideKind is the name of the type used to represent objects
// of type 'AMI_override'.
const AMIOverrideKind = "AMIOverride"

// AMIOverrideLinkKind is the name of the type used to represent links
// to objects of type 'AMI_override'.
const AMIOverrideLinkKind = "AMIOverrideLink"

// AMIOverrideNilKind is the name of the type used to nil references
// to objects of type 'AMI_override'.
const AMIOverrideNilKind = "AMIOverrideNil"

// AMIOverride represents the values of the 'AMI_override' type.
//
// AMIOverride specifies what Amazon Machine Image should be used for a particular product and region.
type AMIOverride struct {
	bitmap_ uint32
	id      string
	href    string
	ami     string
	product *Product
	region  *CloudRegion
}

// Kind returns the name of the type of the object.
func (o *AMIOverride) Kind() string {
	if o == nil {
		return AMIOverrideNilKind
	}
	if o.bitmap_&1 != 0 {
		return AMIOverrideLinkKind
	}
	return AMIOverrideKind
}

// Link returns true iif this is a link.
func (o *AMIOverride) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *AMIOverride) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AMIOverride) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AMIOverride) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AMIOverride) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AMIOverride) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// AMI returns the value of the 'AMI' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AMI is the id of the Amazon Machine Image.
func (o *AMIOverride) AMI() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.ami
	}
	return ""
}

// GetAMI returns the value of the 'AMI' attribute and
// a flag indicating if the attribute has a value.
//
// AMI is the id of the Amazon Machine Image.
func (o *AMIOverride) GetAMI() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.ami
	}
	return
}

// Product returns the value of the 'product' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the product type.
func (o *AMIOverride) Product() *Product {
	if o != nil && o.bitmap_&16 != 0 {
		return o.product
	}
	return nil
}

// GetProduct returns the value of the 'product' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the product type.
func (o *AMIOverride) GetProduct() (value *Product, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.product
	}
	return
}

// Region returns the value of the 'region' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the cloud provider region.
func (o *AMIOverride) Region() *CloudRegion {
	if o != nil && o.bitmap_&32 != 0 {
		return o.region
	}
	return nil
}

// GetRegion returns the value of the 'region' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the cloud provider region.
func (o *AMIOverride) GetRegion() (value *CloudRegion, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.region
	}
	return
}

// AMIOverrideListKind is the name of the type used to represent list of objects of
// type 'AMI_override'.
const AMIOverrideListKind = "AMIOverrideList"

// AMIOverrideListLinkKind is the name of the type used to represent links to list
// of objects of type 'AMI_override'.
const AMIOverrideListLinkKind = "AMIOverrideListLink"

// AMIOverrideNilKind is the name of the type used to nil lists of objects of
// type 'AMI_override'.
const AMIOverrideListNilKind = "AMIOverrideListNil"

// AMIOverrideList is a list of values of the 'AMI_override' type.
type AMIOverrideList struct {
	href  string
	link  bool
	items []*AMIOverride
}

// Kind returns the name of the type of the object.
func (l *AMIOverrideList) Kind() string {
	if l == nil {
		return AMIOverrideListNilKind
	}
	if l.link {
		return AMIOverrideListLinkKind
	}
	return AMIOverrideListKind
}

// Link returns true iif this is a link.
func (l *AMIOverrideList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AMIOverrideList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AMIOverrideList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AMIOverrideList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AMIOverrideList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AMIOverrideList) Get(i int) *AMIOverride {
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
func (l *AMIOverrideList) Slice() []*AMIOverride {
	var slice []*AMIOverride
	if l == nil {
		slice = make([]*AMIOverride, 0)
	} else {
		slice = make([]*AMIOverride, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AMIOverrideList) Each(f func(item *AMIOverride) bool) {
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
func (l *AMIOverrideList) Range(f func(index int, item *AMIOverride) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
