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

// TermsReviewRequest represents the values of the 'terms_review_request' type.
//
// Representation of Red Hat's Terms and Conditions for using OpenShift Dedicated and Amazon Red Hat OpenShift [Terms]
// review requests.
type TermsReviewRequest struct {
	bitmap_            uint32
	accountUsername    string
	eventCode          string
	siteCode           string
	checkOptionalTerms bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *TermsReviewRequest) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AccountUsername returns the value of the 'account_username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines the username of the account of which Terms is being reviewed.
func (o *TermsReviewRequest) AccountUsername() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.accountUsername
	}
	return ""
}

// GetAccountUsername returns the value of the 'account_username' attribute and
// a flag indicating if the attribute has a value.
//
// Defines the username of the account of which Terms is being reviewed.
func (o *TermsReviewRequest) GetAccountUsername() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.accountUsername
	}
	return
}

// CheckOptionalTerms returns the value of the 'check_optional_terms' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// If false, only `terms_required` will be checked
func (o *TermsReviewRequest) CheckOptionalTerms() bool {
	if o != nil && o.bitmap_&2 != 0 {
		return o.checkOptionalTerms
	}
	return false
}

// GetCheckOptionalTerms returns the value of the 'check_optional_terms' attribute and
// a flag indicating if the attribute has a value.
//
// If false, only `terms_required` will be checked
func (o *TermsReviewRequest) GetCheckOptionalTerms() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.checkOptionalTerms
	}
	return
}

// EventCode returns the value of the 'event_code' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines the event code of the terms being checked
func (o *TermsReviewRequest) EventCode() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.eventCode
	}
	return ""
}

// GetEventCode returns the value of the 'event_code' attribute and
// a flag indicating if the attribute has a value.
//
// Defines the event code of the terms being checked
func (o *TermsReviewRequest) GetEventCode() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.eventCode
	}
	return
}

// SiteCode returns the value of the 'site_code' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines the site code of the terms being checked
func (o *TermsReviewRequest) SiteCode() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.siteCode
	}
	return ""
}

// GetSiteCode returns the value of the 'site_code' attribute and
// a flag indicating if the attribute has a value.
//
// Defines the site code of the terms being checked
func (o *TermsReviewRequest) GetSiteCode() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.siteCode
	}
	return
}

// TermsReviewRequestListKind is the name of the type used to represent list of objects of
// type 'terms_review_request'.
const TermsReviewRequestListKind = "TermsReviewRequestList"

// TermsReviewRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'terms_review_request'.
const TermsReviewRequestListLinkKind = "TermsReviewRequestListLink"

// TermsReviewRequestNilKind is the name of the type used to nil lists of objects of
// type 'terms_review_request'.
const TermsReviewRequestListNilKind = "TermsReviewRequestListNil"

// TermsReviewRequestList is a list of values of the 'terms_review_request' type.
type TermsReviewRequestList struct {
	href  string
	link  bool
	items []*TermsReviewRequest
}

// Len returns the length of the list.
func (l *TermsReviewRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *TermsReviewRequestList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *TermsReviewRequestList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *TermsReviewRequestList) SetItems(items []*TermsReviewRequest) {
	l.items = items
}

// Items returns the items of the list.
func (l *TermsReviewRequestList) Items() []*TermsReviewRequest {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *TermsReviewRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *TermsReviewRequestList) Get(i int) *TermsReviewRequest {
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
func (l *TermsReviewRequestList) Slice() []*TermsReviewRequest {
	var slice []*TermsReviewRequest
	if l == nil {
		slice = make([]*TermsReviewRequest, 0)
	} else {
		slice = make([]*TermsReviewRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *TermsReviewRequestList) Each(f func(item *TermsReviewRequest) bool) {
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
func (l *TermsReviewRequestList) Range(f func(index int, item *TermsReviewRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
