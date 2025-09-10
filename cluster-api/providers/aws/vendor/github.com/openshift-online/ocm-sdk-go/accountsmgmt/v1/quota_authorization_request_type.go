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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// QuotaAuthorizationRequest represents the values of the 'quota_authorization_request' type.
type QuotaAuthorizationRequest struct {
	bitmap_          uint32
	accountUsername  string
	availabilityZone string
	displayName      string
	productID        string
	productCategory  string
	quotaVersion     string
	resources        []*ReservedResource
	reserve          bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *QuotaAuthorizationRequest) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AccountUsername returns the value of the 'account_username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaAuthorizationRequest) AccountUsername() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.accountUsername
	}
	return ""
}

// GetAccountUsername returns the value of the 'account_username' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaAuthorizationRequest) GetAccountUsername() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.accountUsername
	}
	return
}

// AvailabilityZone returns the value of the 'availability_zone' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaAuthorizationRequest) AvailabilityZone() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.availabilityZone
	}
	return ""
}

// GetAvailabilityZone returns the value of the 'availability_zone' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaAuthorizationRequest) GetAvailabilityZone() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.availabilityZone
	}
	return
}

// DisplayName returns the value of the 'display_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaAuthorizationRequest) DisplayName() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.displayName
	}
	return ""
}

// GetDisplayName returns the value of the 'display_name' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaAuthorizationRequest) GetDisplayName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.displayName
	}
	return
}

// ProductID returns the value of the 'product_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaAuthorizationRequest) ProductID() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.productID
	}
	return ""
}

// GetProductID returns the value of the 'product_ID' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaAuthorizationRequest) GetProductID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.productID
	}
	return
}

// ProductCategory returns the value of the 'product_category' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaAuthorizationRequest) ProductCategory() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.productCategory
	}
	return ""
}

// GetProductCategory returns the value of the 'product_category' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaAuthorizationRequest) GetProductCategory() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.productCategory
	}
	return
}

// QuotaVersion returns the value of the 'quota_version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaAuthorizationRequest) QuotaVersion() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.quotaVersion
	}
	return ""
}

// GetQuotaVersion returns the value of the 'quota_version' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaAuthorizationRequest) GetQuotaVersion() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.quotaVersion
	}
	return
}

// Reserve returns the value of the 'reserve' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaAuthorizationRequest) Reserve() bool {
	if o != nil && o.bitmap_&64 != 0 {
		return o.reserve
	}
	return false
}

// GetReserve returns the value of the 'reserve' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaAuthorizationRequest) GetReserve() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.reserve
	}
	return
}

// Resources returns the value of the 'resources' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaAuthorizationRequest) Resources() []*ReservedResource {
	if o != nil && o.bitmap_&128 != 0 {
		return o.resources
	}
	return nil
}

// GetResources returns the value of the 'resources' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaAuthorizationRequest) GetResources() (value []*ReservedResource, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.resources
	}
	return
}

// QuotaAuthorizationRequestListKind is the name of the type used to represent list of objects of
// type 'quota_authorization_request'.
const QuotaAuthorizationRequestListKind = "QuotaAuthorizationRequestList"

// QuotaAuthorizationRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'quota_authorization_request'.
const QuotaAuthorizationRequestListLinkKind = "QuotaAuthorizationRequestListLink"

// QuotaAuthorizationRequestNilKind is the name of the type used to nil lists of objects of
// type 'quota_authorization_request'.
const QuotaAuthorizationRequestListNilKind = "QuotaAuthorizationRequestListNil"

// QuotaAuthorizationRequestList is a list of values of the 'quota_authorization_request' type.
type QuotaAuthorizationRequestList struct {
	href  string
	link  bool
	items []*QuotaAuthorizationRequest
}

// Len returns the length of the list.
func (l *QuotaAuthorizationRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *QuotaAuthorizationRequestList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *QuotaAuthorizationRequestList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *QuotaAuthorizationRequestList) SetItems(items []*QuotaAuthorizationRequest) {
	l.items = items
}

// Items returns the items of the list.
func (l *QuotaAuthorizationRequestList) Items() []*QuotaAuthorizationRequest {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *QuotaAuthorizationRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *QuotaAuthorizationRequestList) Get(i int) *QuotaAuthorizationRequest {
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
func (l *QuotaAuthorizationRequestList) Slice() []*QuotaAuthorizationRequest {
	var slice []*QuotaAuthorizationRequest
	if l == nil {
		slice = make([]*QuotaAuthorizationRequest, 0)
	} else {
		slice = make([]*QuotaAuthorizationRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *QuotaAuthorizationRequestList) Each(f func(item *QuotaAuthorizationRequest) bool) {
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
func (l *QuotaAuthorizationRequestList) Range(f func(index int, item *QuotaAuthorizationRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
