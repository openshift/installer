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

// FeatureReviewResponse represents the values of the 'feature_review_response' type.
//
// Representation of a feature review response
type FeatureReviewResponse struct {
	bitmap_   uint32
	featureID string
	enabled   bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *FeatureReviewResponse) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines whether the feature can be toggled
func (o *FeatureReviewResponse) Enabled() bool {
	if o != nil && o.bitmap_&1 != 0 {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Defines whether the feature can be toggled
func (o *FeatureReviewResponse) GetEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.enabled
	}
	return
}

// FeatureID returns the value of the 'feature_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines the feature id which can be toggled
func (o *FeatureReviewResponse) FeatureID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.featureID
	}
	return ""
}

// GetFeatureID returns the value of the 'feature_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Defines the feature id which can be toggled
func (o *FeatureReviewResponse) GetFeatureID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.featureID
	}
	return
}

// FeatureReviewResponseListKind is the name of the type used to represent list of objects of
// type 'feature_review_response'.
const FeatureReviewResponseListKind = "FeatureReviewResponseList"

// FeatureReviewResponseListLinkKind is the name of the type used to represent links to list
// of objects of type 'feature_review_response'.
const FeatureReviewResponseListLinkKind = "FeatureReviewResponseListLink"

// FeatureReviewResponseNilKind is the name of the type used to nil lists of objects of
// type 'feature_review_response'.
const FeatureReviewResponseListNilKind = "FeatureReviewResponseListNil"

// FeatureReviewResponseList is a list of values of the 'feature_review_response' type.
type FeatureReviewResponseList struct {
	href  string
	link  bool
	items []*FeatureReviewResponse
}

// Len returns the length of the list.
func (l *FeatureReviewResponseList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *FeatureReviewResponseList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *FeatureReviewResponseList) Get(i int) *FeatureReviewResponse {
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
func (l *FeatureReviewResponseList) Slice() []*FeatureReviewResponse {
	var slice []*FeatureReviewResponse
	if l == nil {
		slice = make([]*FeatureReviewResponse, 0)
	} else {
		slice = make([]*FeatureReviewResponse, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *FeatureReviewResponseList) Each(f func(item *FeatureReviewResponse) bool) {
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
func (l *FeatureReviewResponseList) Range(f func(index int, item *FeatureReviewResponse) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
