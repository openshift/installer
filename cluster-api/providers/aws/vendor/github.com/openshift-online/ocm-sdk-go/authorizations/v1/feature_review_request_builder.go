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

// FeatureReviewRequestBuilder contains the data and logic needed to build 'feature_review_request' objects.
//
// Representation of a feature review
type FeatureReviewRequestBuilder struct {
	bitmap_         uint32
	accountUsername string
	feature         string
}

// NewFeatureReviewRequest creates a new builder of 'feature_review_request' objects.
func NewFeatureReviewRequest() *FeatureReviewRequestBuilder {
	return &FeatureReviewRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *FeatureReviewRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AccountUsername sets the value of the 'account_username' attribute to the given value.
func (b *FeatureReviewRequestBuilder) AccountUsername(value string) *FeatureReviewRequestBuilder {
	b.accountUsername = value
	b.bitmap_ |= 1
	return b
}

// Feature sets the value of the 'feature' attribute to the given value.
func (b *FeatureReviewRequestBuilder) Feature(value string) *FeatureReviewRequestBuilder {
	b.feature = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *FeatureReviewRequestBuilder) Copy(object *FeatureReviewRequest) *FeatureReviewRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.accountUsername = object.accountUsername
	b.feature = object.feature
	return b
}

// Build creates a 'feature_review_request' object using the configuration stored in the builder.
func (b *FeatureReviewRequestBuilder) Build() (object *FeatureReviewRequest, err error) {
	object = new(FeatureReviewRequest)
	object.bitmap_ = b.bitmap_
	object.accountUsername = b.accountUsername
	object.feature = b.feature
	return
}
