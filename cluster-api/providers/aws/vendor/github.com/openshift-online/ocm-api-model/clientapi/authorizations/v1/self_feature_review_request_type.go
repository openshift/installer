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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/authorizations/v1

// SelfFeatureReviewRequest represents the values of the 'self_feature_review_request' type.
//
// Representation of a feature review performed against oneself
type SelfFeatureReviewRequest struct {
	fieldSet_ []bool
	feature   string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *SelfFeatureReviewRequest) Empty() bool {
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

// Feature returns the value of the 'feature' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates the feature which can be toggled
func (o *SelfFeatureReviewRequest) Feature() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.feature
	}
	return ""
}

// GetFeature returns the value of the 'feature' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates the feature which can be toggled
func (o *SelfFeatureReviewRequest) GetFeature() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.feature
	}
	return
}

// SelfFeatureReviewRequestListKind is the name of the type used to represent list of objects of
// type 'self_feature_review_request'.
const SelfFeatureReviewRequestListKind = "SelfFeatureReviewRequestList"

// SelfFeatureReviewRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'self_feature_review_request'.
const SelfFeatureReviewRequestListLinkKind = "SelfFeatureReviewRequestListLink"

// SelfFeatureReviewRequestNilKind is the name of the type used to nil lists of objects of
// type 'self_feature_review_request'.
const SelfFeatureReviewRequestListNilKind = "SelfFeatureReviewRequestListNil"

// SelfFeatureReviewRequestList is a list of values of the 'self_feature_review_request' type.
type SelfFeatureReviewRequestList struct {
	href  string
	link  bool
	items []*SelfFeatureReviewRequest
}

// Len returns the length of the list.
func (l *SelfFeatureReviewRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *SelfFeatureReviewRequestList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *SelfFeatureReviewRequestList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *SelfFeatureReviewRequestList) SetItems(items []*SelfFeatureReviewRequest) {
	l.items = items
}

// Items returns the items of the list.
func (l *SelfFeatureReviewRequestList) Items() []*SelfFeatureReviewRequest {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *SelfFeatureReviewRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SelfFeatureReviewRequestList) Get(i int) *SelfFeatureReviewRequest {
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
func (l *SelfFeatureReviewRequestList) Slice() []*SelfFeatureReviewRequest {
	var slice []*SelfFeatureReviewRequest
	if l == nil {
		slice = make([]*SelfFeatureReviewRequest, 0)
	} else {
		slice = make([]*SelfFeatureReviewRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SelfFeatureReviewRequestList) Each(f func(item *SelfFeatureReviewRequest) bool) {
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
func (l *SelfFeatureReviewRequestList) Range(f func(index int, item *SelfFeatureReviewRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
