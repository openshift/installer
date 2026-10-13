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

// TokenClaimValidationRule represents the values of the 'token_claim_validation_rule' type.
//
// The rule that is applied to validate token claims to authenticate users.
type TokenClaimValidationRule struct {
	fieldSet_     []bool
	claim         string
	requiredValue string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *TokenClaimValidationRule) Empty() bool {
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

// Claim returns the value of the 'claim' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Claim is a name of a required claim.
func (o *TokenClaimValidationRule) Claim() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.claim
	}
	return ""
}

// GetClaim returns the value of the 'claim' attribute and
// a flag indicating if the attribute has a value.
//
// Claim is a name of a required claim.
func (o *TokenClaimValidationRule) GetClaim() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.claim
	}
	return
}

// RequiredValue returns the value of the 'required_value' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// RequiredValue is the required value for the claim.
func (o *TokenClaimValidationRule) RequiredValue() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.requiredValue
	}
	return ""
}

// GetRequiredValue returns the value of the 'required_value' attribute and
// a flag indicating if the attribute has a value.
//
// RequiredValue is the required value for the claim.
func (o *TokenClaimValidationRule) GetRequiredValue() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.requiredValue
	}
	return
}

// TokenClaimValidationRuleListKind is the name of the type used to represent list of objects of
// type 'token_claim_validation_rule'.
const TokenClaimValidationRuleListKind = "TokenClaimValidationRuleList"

// TokenClaimValidationRuleListLinkKind is the name of the type used to represent links to list
// of objects of type 'token_claim_validation_rule'.
const TokenClaimValidationRuleListLinkKind = "TokenClaimValidationRuleListLink"

// TokenClaimValidationRuleNilKind is the name of the type used to nil lists of objects of
// type 'token_claim_validation_rule'.
const TokenClaimValidationRuleListNilKind = "TokenClaimValidationRuleListNil"

// TokenClaimValidationRuleList is a list of values of the 'token_claim_validation_rule' type.
type TokenClaimValidationRuleList struct {
	href  string
	link  bool
	items []*TokenClaimValidationRule
}

// Len returns the length of the list.
func (l *TokenClaimValidationRuleList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *TokenClaimValidationRuleList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *TokenClaimValidationRuleList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *TokenClaimValidationRuleList) SetItems(items []*TokenClaimValidationRule) {
	l.items = items
}

// Items returns the items of the list.
func (l *TokenClaimValidationRuleList) Items() []*TokenClaimValidationRule {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *TokenClaimValidationRuleList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *TokenClaimValidationRuleList) Get(i int) *TokenClaimValidationRule {
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
func (l *TokenClaimValidationRuleList) Slice() []*TokenClaimValidationRule {
	var slice []*TokenClaimValidationRule
	if l == nil {
		slice = make([]*TokenClaimValidationRule, 0)
	} else {
		slice = make([]*TokenClaimValidationRule, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *TokenClaimValidationRuleList) Each(f func(item *TokenClaimValidationRule) bool) {
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
func (l *TokenClaimValidationRuleList) Range(f func(index int, item *TokenClaimValidationRule) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
