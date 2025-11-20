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

// TokenClaimMappings represents the values of the 'token_claim_mappings' type.
//
// The claim mappings defined for users and groups.
type TokenClaimMappings struct {
	fieldSet_ []bool
	groups    *GroupsClaim
	userName  *UsernameClaim
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *TokenClaimMappings) Empty() bool {
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

// Groups returns the value of the 'groups' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Groups is a name of the claim that should be used to construct groups for the cluster identity.
func (o *TokenClaimMappings) Groups() *GroupsClaim {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.groups
	}
	return nil
}

// GetGroups returns the value of the 'groups' attribute and
// a flag indicating if the attribute has a value.
//
// Groups is a name of the claim that should be used to construct groups for the cluster identity.
func (o *TokenClaimMappings) GetGroups() (value *GroupsClaim, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.groups
	}
	return
}

// UserName returns the value of the 'user_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Username is a name of the claim that should be used to construct usernames for the cluster identity.
func (o *TokenClaimMappings) UserName() *UsernameClaim {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.userName
	}
	return nil
}

// GetUserName returns the value of the 'user_name' attribute and
// a flag indicating if the attribute has a value.
//
// Username is a name of the claim that should be used to construct usernames for the cluster identity.
func (o *TokenClaimMappings) GetUserName() (value *UsernameClaim, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.userName
	}
	return
}

// TokenClaimMappingsListKind is the name of the type used to represent list of objects of
// type 'token_claim_mappings'.
const TokenClaimMappingsListKind = "TokenClaimMappingsList"

// TokenClaimMappingsListLinkKind is the name of the type used to represent links to list
// of objects of type 'token_claim_mappings'.
const TokenClaimMappingsListLinkKind = "TokenClaimMappingsListLink"

// TokenClaimMappingsNilKind is the name of the type used to nil lists of objects of
// type 'token_claim_mappings'.
const TokenClaimMappingsListNilKind = "TokenClaimMappingsListNil"

// TokenClaimMappingsList is a list of values of the 'token_claim_mappings' type.
type TokenClaimMappingsList struct {
	href  string
	link  bool
	items []*TokenClaimMappings
}

// Len returns the length of the list.
func (l *TokenClaimMappingsList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *TokenClaimMappingsList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *TokenClaimMappingsList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *TokenClaimMappingsList) SetItems(items []*TokenClaimMappings) {
	l.items = items
}

// Items returns the items of the list.
func (l *TokenClaimMappingsList) Items() []*TokenClaimMappings {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *TokenClaimMappingsList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *TokenClaimMappingsList) Get(i int) *TokenClaimMappings {
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
func (l *TokenClaimMappingsList) Slice() []*TokenClaimMappings {
	var slice []*TokenClaimMappings
	if l == nil {
		slice = make([]*TokenClaimMappings, 0)
	} else {
		slice = make([]*TokenClaimMappings, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *TokenClaimMappingsList) Each(f func(item *TokenClaimMappings) bool) {
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
func (l *TokenClaimMappingsList) Range(f func(index int, item *TokenClaimMappings) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
