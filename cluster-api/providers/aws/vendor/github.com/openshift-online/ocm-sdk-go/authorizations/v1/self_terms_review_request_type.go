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

// SelfTermsReviewRequest represents the values of the 'self_terms_review_request' type.
//
// Representation of Red Hat's Terms and Conditions for using OpenShift Dedicated and Amazon Red Hat OpenShift [Terms]
// review requests.
type SelfTermsReviewRequest struct {
	bitmap_   uint32
	eventCode string
	siteCode  string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *SelfTermsReviewRequest) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// EventCode returns the value of the 'event_code' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines the event code of the terms being checked
func (o *SelfTermsReviewRequest) EventCode() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.eventCode
	}
	return ""
}

// GetEventCode returns the value of the 'event_code' attribute and
// a flag indicating if the attribute has a value.
//
// Defines the event code of the terms being checked
func (o *SelfTermsReviewRequest) GetEventCode() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.eventCode
	}
	return
}

// SiteCode returns the value of the 'site_code' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines the site code of the terms being checked
func (o *SelfTermsReviewRequest) SiteCode() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.siteCode
	}
	return ""
}

// GetSiteCode returns the value of the 'site_code' attribute and
// a flag indicating if the attribute has a value.
//
// Defines the site code of the terms being checked
func (o *SelfTermsReviewRequest) GetSiteCode() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.siteCode
	}
	return
}

// SelfTermsReviewRequestListKind is the name of the type used to represent list of objects of
// type 'self_terms_review_request'.
const SelfTermsReviewRequestListKind = "SelfTermsReviewRequestList"

// SelfTermsReviewRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'self_terms_review_request'.
const SelfTermsReviewRequestListLinkKind = "SelfTermsReviewRequestListLink"

// SelfTermsReviewRequestNilKind is the name of the type used to nil lists of objects of
// type 'self_terms_review_request'.
const SelfTermsReviewRequestListNilKind = "SelfTermsReviewRequestListNil"

// SelfTermsReviewRequestList is a list of values of the 'self_terms_review_request' type.
type SelfTermsReviewRequestList struct {
	href  string
	link  bool
	items []*SelfTermsReviewRequest
}

// Len returns the length of the list.
func (l *SelfTermsReviewRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *SelfTermsReviewRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SelfTermsReviewRequestList) Get(i int) *SelfTermsReviewRequest {
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
func (l *SelfTermsReviewRequestList) Slice() []*SelfTermsReviewRequest {
	var slice []*SelfTermsReviewRequest
	if l == nil {
		slice = make([]*SelfTermsReviewRequest, 0)
	} else {
		slice = make([]*SelfTermsReviewRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SelfTermsReviewRequestList) Each(f func(item *SelfTermsReviewRequest) bool) {
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
func (l *SelfTermsReviewRequestList) Range(f func(index int, item *SelfTermsReviewRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
