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

// Representation of a feature review
type FeatureReviewRequestBuilder struct {
	fieldSet_       []bool
	accountUsername string
	clusterId       string
	feature         string
	organizationId  string
}

// NewFeatureReviewRequest creates a new builder of 'feature_review_request' objects.
func NewFeatureReviewRequest() *FeatureReviewRequestBuilder {
	return &FeatureReviewRequestBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *FeatureReviewRequestBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// AccountUsername sets the value of the 'account_username' attribute to the given value.
func (b *FeatureReviewRequestBuilder) AccountUsername(value string) *FeatureReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.accountUsername = value
	b.fieldSet_[0] = true
	return b
}

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *FeatureReviewRequestBuilder) ClusterId(value string) *FeatureReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.clusterId = value
	b.fieldSet_[1] = true
	return b
}

// Feature sets the value of the 'feature' attribute to the given value.
func (b *FeatureReviewRequestBuilder) Feature(value string) *FeatureReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.feature = value
	b.fieldSet_[2] = true
	return b
}

// OrganizationId sets the value of the 'organization_id' attribute to the given value.
func (b *FeatureReviewRequestBuilder) OrganizationId(value string) *FeatureReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.organizationId = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *FeatureReviewRequestBuilder) Copy(object *FeatureReviewRequest) *FeatureReviewRequestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.accountUsername = object.accountUsername
	b.clusterId = object.clusterId
	b.feature = object.feature
	b.organizationId = object.organizationId
	return b
}

// Build creates a 'feature_review_request' object using the configuration stored in the builder.
func (b *FeatureReviewRequestBuilder) Build() (object *FeatureReviewRequest, err error) {
	object = new(FeatureReviewRequest)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.accountUsername = b.accountUsername
	object.clusterId = b.clusterId
	object.feature = b.feature
	object.organizationId = b.organizationId
	return
}
