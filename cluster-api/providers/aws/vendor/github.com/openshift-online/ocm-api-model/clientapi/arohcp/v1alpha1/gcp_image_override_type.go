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

import (
	v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
)

// GCPImageOverrideKind is the name of the type used to represent objects
// of type 'GCP_image_override'.
const GCPImageOverrideKind = "GCPImageOverride"

// GCPImageOverrideLinkKind is the name of the type used to represent links
// to objects of type 'GCP_image_override'.
const GCPImageOverrideLinkKind = "GCPImageOverrideLink"

// GCPImageOverrideNilKind is the name of the type used to nil references
// to objects of type 'GCP_image_override'.
const GCPImageOverrideNilKind = "GCPImageOverrideNil"

// GCPImageOverride represents the values of the 'GCP_image_override' type.
//
// GcpImageOverride specifies what a GCP VM Image should be used for a particular product and billing model
type GCPImageOverride struct {
	fieldSet_    []bool
	id           string
	href         string
	billingModel *v1.BillingModelItem
	imageID      string
	product      *v1.Product
	projectID    string
}

// Kind returns the name of the type of the object.
func (o *GCPImageOverride) Kind() string {
	if o == nil {
		return GCPImageOverrideNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return GCPImageOverrideLinkKind
	}
	return GCPImageOverrideKind
}

// Link returns true if this is a link.
func (o *GCPImageOverride) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *GCPImageOverride) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *GCPImageOverride) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *GCPImageOverride) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *GCPImageOverride) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *GCPImageOverride) Empty() bool {
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

// BillingModel returns the value of the 'billing_model' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the billing model.
func (o *GCPImageOverride) BillingModel() *v1.BillingModelItem {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.billingModel
	}
	return nil
}

// GetBillingModel returns the value of the 'billing_model' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the billing model.
func (o *GCPImageOverride) GetBillingModel() (value *v1.BillingModelItem, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.billingModel
	}
	return
}

// ImageID returns the value of the 'image_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ImageID is the id of the Google Cloud Platform image.
func (o *GCPImageOverride) ImageID() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.imageID
	}
	return ""
}

// GetImageID returns the value of the 'image_ID' attribute and
// a flag indicating if the attribute has a value.
//
// ImageID is the id of the Google Cloud Platform image.
func (o *GCPImageOverride) GetImageID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.imageID
	}
	return
}

// Product returns the value of the 'product' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the product type.
func (o *GCPImageOverride) Product() *v1.Product {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.product
	}
	return nil
}

// GetProduct returns the value of the 'product' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the product type.
func (o *GCPImageOverride) GetProduct() (value *v1.Product, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.product
	}
	return
}

// ProjectID returns the value of the 'project_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ProjectID is the id of the Google Cloud Platform project that hosts the image.
func (o *GCPImageOverride) ProjectID() string {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.projectID
	}
	return ""
}

// GetProjectID returns the value of the 'project_ID' attribute and
// a flag indicating if the attribute has a value.
//
// ProjectID is the id of the Google Cloud Platform project that hosts the image.
func (o *GCPImageOverride) GetProjectID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.projectID
	}
	return
}

// GCPImageOverrideListKind is the name of the type used to represent list of objects of
// type 'GCP_image_override'.
const GCPImageOverrideListKind = "GCPImageOverrideList"

// GCPImageOverrideListLinkKind is the name of the type used to represent links to list
// of objects of type 'GCP_image_override'.
const GCPImageOverrideListLinkKind = "GCPImageOverrideListLink"

// GCPImageOverrideNilKind is the name of the type used to nil lists of objects of
// type 'GCP_image_override'.
const GCPImageOverrideListNilKind = "GCPImageOverrideListNil"

// GCPImageOverrideList is a list of values of the 'GCP_image_override' type.
type GCPImageOverrideList struct {
	href  string
	link  bool
	items []*GCPImageOverride
}

// Kind returns the name of the type of the object.
func (l *GCPImageOverrideList) Kind() string {
	if l == nil {
		return GCPImageOverrideListNilKind
	}
	if l.link {
		return GCPImageOverrideListLinkKind
	}
	return GCPImageOverrideListKind
}

// Link returns true iif this is a link.
func (l *GCPImageOverrideList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *GCPImageOverrideList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *GCPImageOverrideList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *GCPImageOverrideList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *GCPImageOverrideList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *GCPImageOverrideList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *GCPImageOverrideList) SetItems(items []*GCPImageOverride) {
	l.items = items
}

// Items returns the items of the list.
func (l *GCPImageOverrideList) Items() []*GCPImageOverride {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *GCPImageOverrideList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *GCPImageOverrideList) Get(i int) *GCPImageOverride {
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
func (l *GCPImageOverrideList) Slice() []*GCPImageOverride {
	var slice []*GCPImageOverride
	if l == nil {
		slice = make([]*GCPImageOverride, 0)
	} else {
		slice = make([]*GCPImageOverride, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *GCPImageOverrideList) Each(f func(item *GCPImageOverride) bool) {
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
func (l *GCPImageOverrideList) Range(f func(index int, item *GCPImageOverride) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
