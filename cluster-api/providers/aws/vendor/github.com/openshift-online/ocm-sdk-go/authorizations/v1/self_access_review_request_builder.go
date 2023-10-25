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

// SelfAccessReviewRequestBuilder contains the data and logic needed to build 'self_access_review_request' objects.
//
// Representation of an access review performed against oneself
type SelfAccessReviewRequestBuilder struct {
	bitmap_        uint32
	action         string
	clusterID      string
	clusterUUID    string
	organizationID string
	resourceType   string
	subscriptionID string
}

// NewSelfAccessReviewRequest creates a new builder of 'self_access_review_request' objects.
func NewSelfAccessReviewRequest() *SelfAccessReviewRequestBuilder {
	return &SelfAccessReviewRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SelfAccessReviewRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Action sets the value of the 'action' attribute to the given value.
func (b *SelfAccessReviewRequestBuilder) Action(value string) *SelfAccessReviewRequestBuilder {
	b.action = value
	b.bitmap_ |= 1
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *SelfAccessReviewRequestBuilder) ClusterID(value string) *SelfAccessReviewRequestBuilder {
	b.clusterID = value
	b.bitmap_ |= 2
	return b
}

// ClusterUUID sets the value of the 'cluster_UUID' attribute to the given value.
func (b *SelfAccessReviewRequestBuilder) ClusterUUID(value string) *SelfAccessReviewRequestBuilder {
	b.clusterUUID = value
	b.bitmap_ |= 4
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *SelfAccessReviewRequestBuilder) OrganizationID(value string) *SelfAccessReviewRequestBuilder {
	b.organizationID = value
	b.bitmap_ |= 8
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *SelfAccessReviewRequestBuilder) ResourceType(value string) *SelfAccessReviewRequestBuilder {
	b.resourceType = value
	b.bitmap_ |= 16
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *SelfAccessReviewRequestBuilder) SubscriptionID(value string) *SelfAccessReviewRequestBuilder {
	b.subscriptionID = value
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SelfAccessReviewRequestBuilder) Copy(object *SelfAccessReviewRequest) *SelfAccessReviewRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.action = object.action
	b.clusterID = object.clusterID
	b.clusterUUID = object.clusterUUID
	b.organizationID = object.organizationID
	b.resourceType = object.resourceType
	b.subscriptionID = object.subscriptionID
	return b
}

// Build creates a 'self_access_review_request' object using the configuration stored in the builder.
func (b *SelfAccessReviewRequestBuilder) Build() (object *SelfAccessReviewRequest, err error) {
	object = new(SelfAccessReviewRequest)
	object.bitmap_ = b.bitmap_
	object.action = b.action
	object.clusterID = b.clusterID
	object.clusterUUID = b.clusterUUID
	object.organizationID = b.organizationID
	object.resourceType = b.resourceType
	object.subscriptionID = b.subscriptionID
	return
}
