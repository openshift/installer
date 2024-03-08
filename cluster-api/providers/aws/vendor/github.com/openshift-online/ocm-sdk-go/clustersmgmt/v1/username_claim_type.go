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

// UsernameClaim represents the values of the 'username_claim' type.
//
// The username claim mapping.
type UsernameClaim struct {
	bitmap_      uint32
	claim        string
	prefix       string
	prefixPolicy string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *UsernameClaim) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Claim returns the value of the 'claim' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The claim used in the token.
func (o *UsernameClaim) Claim() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.claim
	}
	return ""
}

// GetClaim returns the value of the 'claim' attribute and
// a flag indicating if the attribute has a value.
//
// The claim used in the token.
func (o *UsernameClaim) GetClaim() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.claim
	}
	return
}

// Prefix returns the value of the 'prefix' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// A prefix contatenated in the claim (Optional).
func (o *UsernameClaim) Prefix() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.prefix
	}
	return ""
}

// GetPrefix returns the value of the 'prefix' attribute and
// a flag indicating if the attribute has a value.
//
// A prefix contatenated in the claim (Optional).
func (o *UsernameClaim) GetPrefix() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.prefix
	}
	return
}

// PrefixPolicy returns the value of the 'prefix_policy' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// PrefixPolicy specifies how a prefix should apply.
//
// By default, claims other than `email` will be prefixed with the issuer URL to
// prevent naming clashes with other plugins.
//
// Set to "NoPrefix" to disable prefixing.
func (o *UsernameClaim) PrefixPolicy() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.prefixPolicy
	}
	return ""
}

// GetPrefixPolicy returns the value of the 'prefix_policy' attribute and
// a flag indicating if the attribute has a value.
//
// PrefixPolicy specifies how a prefix should apply.
//
// By default, claims other than `email` will be prefixed with the issuer URL to
// prevent naming clashes with other plugins.
//
// Set to "NoPrefix" to disable prefixing.
func (o *UsernameClaim) GetPrefixPolicy() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.prefixPolicy
	}
	return
}

// UsernameClaimListKind is the name of the type used to represent list of objects of
// type 'username_claim'.
const UsernameClaimListKind = "UsernameClaimList"

// UsernameClaimListLinkKind is the name of the type used to represent links to list
// of objects of type 'username_claim'.
const UsernameClaimListLinkKind = "UsernameClaimListLink"

// UsernameClaimNilKind is the name of the type used to nil lists of objects of
// type 'username_claim'.
const UsernameClaimListNilKind = "UsernameClaimListNil"

// UsernameClaimList is a list of values of the 'username_claim' type.
type UsernameClaimList struct {
	href  string
	link  bool
	items []*UsernameClaim
}

// Len returns the length of the list.
func (l *UsernameClaimList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *UsernameClaimList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *UsernameClaimList) Get(i int) *UsernameClaim {
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
func (l *UsernameClaimList) Slice() []*UsernameClaim {
	var slice []*UsernameClaim
	if l == nil {
		slice = make([]*UsernameClaim, 0)
	} else {
		slice = make([]*UsernameClaim, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *UsernameClaimList) Each(f func(item *UsernameClaim) bool) {
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
func (l *UsernameClaimList) Range(f func(index int, item *UsernameClaim) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
