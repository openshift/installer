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

// FeatureReviewRequest represents the values of the 'feature_review_request' type.
//
// Representation of a feature review
type FeatureReviewRequest struct {
	bitmap_         uint32
	accountUsername string
	feature         string
	organizationId  string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *FeatureReviewRequest) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AccountUsername returns the value of the 'account_username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines the username of the account of which access is being reviewed
func (o *FeatureReviewRequest) AccountUsername() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.accountUsername
	}
	return ""
}

// GetAccountUsername returns the value of the 'account_username' attribute and
// a flag indicating if the attribute has a value.
//
// Defines the username of the account of which access is being reviewed
func (o *FeatureReviewRequest) GetAccountUsername() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.accountUsername
	}
	return
}

// Feature returns the value of the 'feature' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the feature which can be toggled
func (o *FeatureReviewRequest) Feature() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.feature
	}
	return ""
}

// GetFeature returns the value of the 'feature' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the feature which can be toggled
func (o *FeatureReviewRequest) GetFeature() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.feature
	}
	return
}

// OrganizationId returns the value of the 'organization_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines the organisation id of the account of which access is being reviewed
func (o *FeatureReviewRequest) OrganizationId() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.organizationId
	}
	return ""
}

// GetOrganizationId returns the value of the 'organization_id' attribute and
// a flag indicating if the attribute has a value.
//
// Defines the organisation id of the account of which access is being reviewed
func (o *FeatureReviewRequest) GetOrganizationId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.organizationId
	}
	return
}

// FeatureReviewRequestListKind is the name of the type used to represent list of objects of
// type 'feature_review_request'.
const FeatureReviewRequestListKind = "FeatureReviewRequestList"

// FeatureReviewRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'feature_review_request'.
const FeatureReviewRequestListLinkKind = "FeatureReviewRequestListLink"

// FeatureReviewRequestNilKind is the name of the type used to nil lists of objects of
// type 'feature_review_request'.
const FeatureReviewRequestListNilKind = "FeatureReviewRequestListNil"

// FeatureReviewRequestList is a list of values of the 'feature_review_request' type.
type FeatureReviewRequestList struct {
	href  string
	link  bool
	items []*FeatureReviewRequest
}

// Len returns the length of the list.
func (l *FeatureReviewRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *FeatureReviewRequestList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *FeatureReviewRequestList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *FeatureReviewRequestList) SetItems(items []*FeatureReviewRequest) {
	l.items = items
}

// Items returns the items of the list.
func (l *FeatureReviewRequestList) Items() []*FeatureReviewRequest {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *FeatureReviewRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *FeatureReviewRequestList) Get(i int) *FeatureReviewRequest {
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
func (l *FeatureReviewRequestList) Slice() []*FeatureReviewRequest {
	var slice []*FeatureReviewRequest
	if l == nil {
		slice = make([]*FeatureReviewRequest, 0)
	} else {
		slice = make([]*FeatureReviewRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *FeatureReviewRequestList) Each(f func(item *FeatureReviewRequest) bool) {
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
func (l *FeatureReviewRequestList) Range(f func(index int, item *FeatureReviewRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
