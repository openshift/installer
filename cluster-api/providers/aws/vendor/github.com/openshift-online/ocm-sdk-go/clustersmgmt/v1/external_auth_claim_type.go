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

// ExternalAuthClaim represents the values of the 'external_auth_claim' type.
//
// The claims and validation rules used in the configuration of the external authentication.
type ExternalAuthClaim struct {
	bitmap_         uint32
	mappings        *TokenClaimMappings
	validationRules []*TokenClaimValidationRule
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ExternalAuthClaim) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Mappings returns the value of the 'mappings' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Mapping describes rules on how to transform information from an ID token into a cluster identity.
func (o *ExternalAuthClaim) Mappings() *TokenClaimMappings {
	if o != nil && o.bitmap_&1 != 0 {
		return o.mappings
	}
	return nil
}

// GetMappings returns the value of the 'mappings' attribute and
// a flag indicating if the attribute has a value.
//
// Mapping describes rules on how to transform information from an ID token into a cluster identity.
func (o *ExternalAuthClaim) GetMappings() (value *TokenClaimMappings, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.mappings
	}
	return
}

// ValidationRules returns the value of the 'validation_rules' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ValidationRules are rules that are applied to validate token claims to authenticate users.
func (o *ExternalAuthClaim) ValidationRules() []*TokenClaimValidationRule {
	if o != nil && o.bitmap_&2 != 0 {
		return o.validationRules
	}
	return nil
}

// GetValidationRules returns the value of the 'validation_rules' attribute and
// a flag indicating if the attribute has a value.
//
// ValidationRules are rules that are applied to validate token claims to authenticate users.
func (o *ExternalAuthClaim) GetValidationRules() (value []*TokenClaimValidationRule, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.validationRules
	}
	return
}

// ExternalAuthClaimListKind is the name of the type used to represent list of objects of
// type 'external_auth_claim'.
const ExternalAuthClaimListKind = "ExternalAuthClaimList"

// ExternalAuthClaimListLinkKind is the name of the type used to represent links to list
// of objects of type 'external_auth_claim'.
const ExternalAuthClaimListLinkKind = "ExternalAuthClaimListLink"

// ExternalAuthClaimNilKind is the name of the type used to nil lists of objects of
// type 'external_auth_claim'.
const ExternalAuthClaimListNilKind = "ExternalAuthClaimListNil"

// ExternalAuthClaimList is a list of values of the 'external_auth_claim' type.
type ExternalAuthClaimList struct {
	href  string
	link  bool
	items []*ExternalAuthClaim
}

// Len returns the length of the list.
func (l *ExternalAuthClaimList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ExternalAuthClaimList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ExternalAuthClaimList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ExternalAuthClaimList) SetItems(items []*ExternalAuthClaim) {
	l.items = items
}

// Items returns the items of the list.
func (l *ExternalAuthClaimList) Items() []*ExternalAuthClaim {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ExternalAuthClaimList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ExternalAuthClaimList) Get(i int) *ExternalAuthClaim {
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
func (l *ExternalAuthClaimList) Slice() []*ExternalAuthClaim {
	var slice []*ExternalAuthClaim
	if l == nil {
		slice = make([]*ExternalAuthClaim, 0)
	} else {
		slice = make([]*ExternalAuthClaim, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ExternalAuthClaimList) Each(f func(item *ExternalAuthClaim) bool) {
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
func (l *ExternalAuthClaimList) Range(f func(index int, item *ExternalAuthClaim) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
