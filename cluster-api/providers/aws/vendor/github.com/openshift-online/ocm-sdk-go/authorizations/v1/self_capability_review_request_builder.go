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

// SelfCapabilityReviewRequestBuilder contains the data and logic needed to build 'self_capability_review_request' objects.
//
// Representation of a capability review.
type SelfCapabilityReviewRequestBuilder struct {
	bitmap_         uint32
	accountUsername string
	capability      string
	clusterID       string
	organizationID  string
	resourceType    string
	subscriptionID  string
	type_           string
}

// NewSelfCapabilityReviewRequest creates a new builder of 'self_capability_review_request' objects.
func NewSelfCapabilityReviewRequest() *SelfCapabilityReviewRequestBuilder {
	return &SelfCapabilityReviewRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SelfCapabilityReviewRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AccountUsername sets the value of the 'account_username' attribute to the given value.
func (b *SelfCapabilityReviewRequestBuilder) AccountUsername(value string) *SelfCapabilityReviewRequestBuilder {
	b.accountUsername = value
	b.bitmap_ |= 1
	return b
}

// Capability sets the value of the 'capability' attribute to the given value.
func (b *SelfCapabilityReviewRequestBuilder) Capability(value string) *SelfCapabilityReviewRequestBuilder {
	b.capability = value
	b.bitmap_ |= 2
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *SelfCapabilityReviewRequestBuilder) ClusterID(value string) *SelfCapabilityReviewRequestBuilder {
	b.clusterID = value
	b.bitmap_ |= 4
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *SelfCapabilityReviewRequestBuilder) OrganizationID(value string) *SelfCapabilityReviewRequestBuilder {
	b.organizationID = value
	b.bitmap_ |= 8
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *SelfCapabilityReviewRequestBuilder) ResourceType(value string) *SelfCapabilityReviewRequestBuilder {
	b.resourceType = value
	b.bitmap_ |= 16
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *SelfCapabilityReviewRequestBuilder) SubscriptionID(value string) *SelfCapabilityReviewRequestBuilder {
	b.subscriptionID = value
	b.bitmap_ |= 32
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *SelfCapabilityReviewRequestBuilder) Type(value string) *SelfCapabilityReviewRequestBuilder {
	b.type_ = value
	b.bitmap_ |= 64
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SelfCapabilityReviewRequestBuilder) Copy(object *SelfCapabilityReviewRequest) *SelfCapabilityReviewRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.accountUsername = object.accountUsername
	b.capability = object.capability
	b.clusterID = object.clusterID
	b.organizationID = object.organizationID
	b.resourceType = object.resourceType
	b.subscriptionID = object.subscriptionID
	b.type_ = object.type_
	return b
}

// Build creates a 'self_capability_review_request' object using the configuration stored in the builder.
func (b *SelfCapabilityReviewRequestBuilder) Build() (object *SelfCapabilityReviewRequest, err error) {
	object = new(SelfCapabilityReviewRequest)
	object.bitmap_ = b.bitmap_
	object.accountUsername = b.accountUsername
	object.capability = b.capability
	object.clusterID = b.clusterID
	object.organizationID = b.organizationID
	object.resourceType = b.resourceType
	object.subscriptionID = b.subscriptionID
	object.type_ = b.type_
	return
}
