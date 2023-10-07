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

package v1 // github.com/openshift-online/ocm-sdk-go/authorizations/v1

// TermsReviewResponse represents the values of the 'terms_review_response' type.
//
// Representation of Red Hat's Terms and Conditions for using OpenShift Dedicated and Amazon Red Hat OpenShift [Terms]
// review response.
type TermsReviewResponse struct {
	bitmap_        uint32
	accountId      string
	organizationID string
	redirectUrl    string
	termsAvailable bool
	termsRequired  bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *TermsReviewResponse) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AccountId returns the value of the 'account_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Account ID of requesting user.
func (o *TermsReviewResponse) AccountId() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.accountId
	}
	return ""
}

// GetAccountId returns the value of the 'account_id' attribute and
// a flag indicating if the attribute has a value.
//
// Account ID of requesting user.
func (o *TermsReviewResponse) GetAccountId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.accountId
	}
	return
}

// OrganizationID returns the value of the 'organization_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates which Organization the user belongs to.
func (o *TermsReviewResponse) OrganizationID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.organizationID
	}
	return ""
}

// GetOrganizationID returns the value of the 'organization_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates which Organization the user belongs to.
func (o *TermsReviewResponse) GetOrganizationID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.organizationID
	}
	return
}

// RedirectUrl returns the value of the 'redirect_url' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional URL to Red Hat's Terms and Conditions Application if the user has either required or available Terms
// needs to acknowledge.
func (o *TermsReviewResponse) RedirectUrl() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.redirectUrl
	}
	return ""
}

// GetRedirectUrl returns the value of the 'redirect_url' attribute and
// a flag indicating if the attribute has a value.
//
// Optional URL to Red Hat's Terms and Conditions Application if the user has either required or available Terms
// needs to acknowledge.
func (o *TermsReviewResponse) GetRedirectUrl() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.redirectUrl
	}
	return
}

// TermsAvailable returns the value of the 'terms_available' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines whether Terms are available.
func (o *TermsReviewResponse) TermsAvailable() bool {
	if o != nil && o.bitmap_&8 != 0 {
		return o.termsAvailable
	}
	return false
}

// GetTermsAvailable returns the value of the 'terms_available' attribute and
// a flag indicating if the attribute has a value.
//
// Defines whether Terms are available.
func (o *TermsReviewResponse) GetTermsAvailable() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.termsAvailable
	}
	return
}

// TermsRequired returns the value of the 'terms_required' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines whether user is required to accept Terms before using OpenShift Dedicated and Amazon Red Hat OpenShift.
func (o *TermsReviewResponse) TermsRequired() bool {
	if o != nil && o.bitmap_&16 != 0 {
		return o.termsRequired
	}
	return false
}

// GetTermsRequired returns the value of the 'terms_required' attribute and
// a flag indicating if the attribute has a value.
//
// Defines whether user is required to accept Terms before using OpenShift Dedicated and Amazon Red Hat OpenShift.
func (o *TermsReviewResponse) GetTermsRequired() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.termsRequired
	}
	return
}

// TermsReviewResponseListKind is the name of the type used to represent list of objects of
// type 'terms_review_response'.
const TermsReviewResponseListKind = "TermsReviewResponseList"

// TermsReviewResponseListLinkKind is the name of the type used to represent links to list
// of objects of type 'terms_review_response'.
const TermsReviewResponseListLinkKind = "TermsReviewResponseListLink"

// TermsReviewResponseNilKind is the name of the type used to nil lists of objects of
// type 'terms_review_response'.
const TermsReviewResponseListNilKind = "TermsReviewResponseListNil"

// TermsReviewResponseList is a list of values of the 'terms_review_response' type.
type TermsReviewResponseList struct {
	href  string
	link  bool
	items []*TermsReviewResponse
}

// Len returns the length of the list.
func (l *TermsReviewResponseList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *TermsReviewResponseList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *TermsReviewResponseList) Get(i int) *TermsReviewResponse {
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
func (l *TermsReviewResponseList) Slice() []*TermsReviewResponse {
	var slice []*TermsReviewResponse
	if l == nil {
		slice = make([]*TermsReviewResponse, 0)
	} else {
		slice = make([]*TermsReviewResponse, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *TermsReviewResponseList) Each(f func(item *TermsReviewResponse) bool) {
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
func (l *TermsReviewResponseList) Range(f func(index int, item *TermsReviewResponse) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
