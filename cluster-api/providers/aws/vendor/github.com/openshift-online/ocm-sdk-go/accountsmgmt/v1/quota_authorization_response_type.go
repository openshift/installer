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

// QuotaAuthorizationResponse represents the values of the 'quota_authorization_response' type.
type QuotaAuthorizationResponse struct {
	bitmap_         uint32
	excessResources []*ReservedResource
	subscription    *Subscription
	allowed         bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *QuotaAuthorizationResponse) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Allowed returns the value of the 'allowed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaAuthorizationResponse) Allowed() bool {
	if o != nil && o.bitmap_&1 != 0 {
		return o.allowed
	}
	return false
}

// GetAllowed returns the value of the 'allowed' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaAuthorizationResponse) GetAllowed() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.allowed
	}
	return
}

// ExcessResources returns the value of the 'excess_resources' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaAuthorizationResponse) ExcessResources() []*ReservedResource {
	if o != nil && o.bitmap_&2 != 0 {
		return o.excessResources
	}
	return nil
}

// GetExcessResources returns the value of the 'excess_resources' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaAuthorizationResponse) GetExcessResources() (value []*ReservedResource, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.excessResources
	}
	return
}

// Subscription returns the value of the 'subscription' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *QuotaAuthorizationResponse) Subscription() *Subscription {
	if o != nil && o.bitmap_&4 != 0 {
		return o.subscription
	}
	return nil
}

// GetSubscription returns the value of the 'subscription' attribute and
// a flag indicating if the attribute has a value.
func (o *QuotaAuthorizationResponse) GetSubscription() (value *Subscription, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.subscription
	}
	return
}

// QuotaAuthorizationResponseListKind is the name of the type used to represent list of objects of
// type 'quota_authorization_response'.
const QuotaAuthorizationResponseListKind = "QuotaAuthorizationResponseList"

// QuotaAuthorizationResponseListLinkKind is the name of the type used to represent links to list
// of objects of type 'quota_authorization_response'.
const QuotaAuthorizationResponseListLinkKind = "QuotaAuthorizationResponseListLink"

// QuotaAuthorizationResponseNilKind is the name of the type used to nil lists of objects of
// type 'quota_authorization_response'.
const QuotaAuthorizationResponseListNilKind = "QuotaAuthorizationResponseListNil"

// QuotaAuthorizationResponseList is a list of values of the 'quota_authorization_response' type.
type QuotaAuthorizationResponseList struct {
	href  string
	link  bool
	items []*QuotaAuthorizationResponse
}

// Len returns the length of the list.
func (l *QuotaAuthorizationResponseList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *QuotaAuthorizationResponseList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *QuotaAuthorizationResponseList) Get(i int) *QuotaAuthorizationResponse {
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
func (l *QuotaAuthorizationResponseList) Slice() []*QuotaAuthorizationResponse {
	var slice []*QuotaAuthorizationResponse
	if l == nil {
		slice = make([]*QuotaAuthorizationResponse, 0)
	} else {
		slice = make([]*QuotaAuthorizationResponse, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *QuotaAuthorizationResponseList) Each(f func(item *QuotaAuthorizationResponse) bool) {
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
func (l *QuotaAuthorizationResponseList) Range(f func(index int, item *QuotaAuthorizationResponse) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
