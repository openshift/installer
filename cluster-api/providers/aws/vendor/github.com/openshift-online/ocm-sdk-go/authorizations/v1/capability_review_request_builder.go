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

// CapabilityReviewRequestBuilder contains the data and logic needed to build 'capability_review_request' objects.
//
// Representation of a capability review.
type CapabilityReviewRequestBuilder struct {
	bitmap_         uint32
	accountUsername string
	capability      string
	clusterID       string
	organizationID  string
	resourceType    string
	subscriptionID  string
	type_           string
}

// NewCapabilityReviewRequest creates a new builder of 'capability_review_request' objects.
func NewCapabilityReviewRequest() *CapabilityReviewRequestBuilder {
	return &CapabilityReviewRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CapabilityReviewRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AccountUsername sets the value of the 'account_username' attribute to the given value.
func (b *CapabilityReviewRequestBuilder) AccountUsername(value string) *CapabilityReviewRequestBuilder {
	b.accountUsername = value
	b.bitmap_ |= 1
	return b
}

// Capability sets the value of the 'capability' attribute to the given value.
func (b *CapabilityReviewRequestBuilder) Capability(value string) *CapabilityReviewRequestBuilder {
	b.capability = value
	b.bitmap_ |= 2
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *CapabilityReviewRequestBuilder) ClusterID(value string) *CapabilityReviewRequestBuilder {
	b.clusterID = value
	b.bitmap_ |= 4
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *CapabilityReviewRequestBuilder) OrganizationID(value string) *CapabilityReviewRequestBuilder {
	b.organizationID = value
	b.bitmap_ |= 8
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *CapabilityReviewRequestBuilder) ResourceType(value string) *CapabilityReviewRequestBuilder {
	b.resourceType = value
	b.bitmap_ |= 16
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *CapabilityReviewRequestBuilder) SubscriptionID(value string) *CapabilityReviewRequestBuilder {
	b.subscriptionID = value
	b.bitmap_ |= 32
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *CapabilityReviewRequestBuilder) Type(value string) *CapabilityReviewRequestBuilder {
	b.type_ = value
	b.bitmap_ |= 64
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CapabilityReviewRequestBuilder) Copy(object *CapabilityReviewRequest) *CapabilityReviewRequestBuilder {
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

// Build creates a 'capability_review_request' object using the configuration stored in the builder.
func (b *CapabilityReviewRequestBuilder) Build() (object *CapabilityReviewRequest, err error) {
	object = new(CapabilityReviewRequest)
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
