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

// TokenAuthorizationResponse represents the values of the 'token_authorization_response' type.
type TokenAuthorizationResponse struct {
	bitmap_ uint32
	account *Account
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *TokenAuthorizationResponse) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Account returns the value of the 'account' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *TokenAuthorizationResponse) Account() *Account {
	if o != nil && o.bitmap_&1 != 0 {
		return o.account
	}
	return nil
}

// GetAccount returns the value of the 'account' attribute and
// a flag indicating if the attribute has a value.
func (o *TokenAuthorizationResponse) GetAccount() (value *Account, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.account
	}
	return
}

// TokenAuthorizationResponseListKind is the name of the type used to represent list of objects of
// type 'token_authorization_response'.
const TokenAuthorizationResponseListKind = "TokenAuthorizationResponseList"

// TokenAuthorizationResponseListLinkKind is the name of the type used to represent links to list
// of objects of type 'token_authorization_response'.
const TokenAuthorizationResponseListLinkKind = "TokenAuthorizationResponseListLink"

// TokenAuthorizationResponseNilKind is the name of the type used to nil lists of objects of
// type 'token_authorization_response'.
const TokenAuthorizationResponseListNilKind = "TokenAuthorizationResponseListNil"

// TokenAuthorizationResponseList is a list of values of the 'token_authorization_response' type.
type TokenAuthorizationResponseList struct {
	href  string
	link  bool
	items []*TokenAuthorizationResponse
}

// Len returns the length of the list.
func (l *TokenAuthorizationResponseList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *TokenAuthorizationResponseList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *TokenAuthorizationResponseList) Get(i int) *TokenAuthorizationResponse {
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
func (l *TokenAuthorizationResponseList) Slice() []*TokenAuthorizationResponse {
	var slice []*TokenAuthorizationResponse
	if l == nil {
		slice = make([]*TokenAuthorizationResponse, 0)
	} else {
		slice = make([]*TokenAuthorizationResponse, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *TokenAuthorizationResponseList) Each(f func(item *TokenAuthorizationResponse) bool) {
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
func (l *TokenAuthorizationResponseList) Range(f func(index int, item *TokenAuthorizationResponse) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
