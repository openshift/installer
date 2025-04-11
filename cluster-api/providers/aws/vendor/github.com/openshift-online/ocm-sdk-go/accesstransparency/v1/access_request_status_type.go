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

package v1 // github.com/openshift-online/ocm-sdk-go/accesstransparency/v1

import (
	time "time"
)

// AccessRequestStatus represents the values of the 'access_request_status' type.
//
// Representation of an access request status.
type AccessRequestStatus struct {
	bitmap_   uint32
	expiresAt time.Time
	state     AccessRequestState
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AccessRequestStatus) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ExpiresAt returns the value of the 'expires_at' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the access request will expire, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *AccessRequestStatus) ExpiresAt() time.Time {
	if o != nil && o.bitmap_&1 != 0 {
		return o.expiresAt
	}
	return time.Time{}
}

// GetExpiresAt returns the value of the 'expires_at' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the access request will expire, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *AccessRequestStatus) GetExpiresAt() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.expiresAt
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Current state of the Access Request.
func (o *AccessRequestStatus) State() AccessRequestState {
	if o != nil && o.bitmap_&2 != 0 {
		return o.state
	}
	return AccessRequestState("")
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
//
// Current state of the Access Request.
func (o *AccessRequestStatus) GetState() (value AccessRequestState, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.state
	}
	return
}

// AccessRequestStatusListKind is the name of the type used to represent list of objects of
// type 'access_request_status'.
const AccessRequestStatusListKind = "AccessRequestStatusList"

// AccessRequestStatusListLinkKind is the name of the type used to represent links to list
// of objects of type 'access_request_status'.
const AccessRequestStatusListLinkKind = "AccessRequestStatusListLink"

// AccessRequestStatusNilKind is the name of the type used to nil lists of objects of
// type 'access_request_status'.
const AccessRequestStatusListNilKind = "AccessRequestStatusListNil"

// AccessRequestStatusList is a list of values of the 'access_request_status' type.
type AccessRequestStatusList struct {
	href  string
	link  bool
	items []*AccessRequestStatus
}

// Len returns the length of the list.
func (l *AccessRequestStatusList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AccessRequestStatusList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AccessRequestStatusList) Get(i int) *AccessRequestStatus {
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
func (l *AccessRequestStatusList) Slice() []*AccessRequestStatus {
	var slice []*AccessRequestStatus
	if l == nil {
		slice = make([]*AccessRequestStatus, 0)
	} else {
		slice = make([]*AccessRequestStatus, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AccessRequestStatusList) Each(f func(item *AccessRequestStatus) bool) {
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
func (l *AccessRequestStatusList) Range(f func(index int, item *AccessRequestStatus) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
