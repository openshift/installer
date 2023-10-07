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

// AccessReviewRequestBuilder contains the data and logic needed to build 'access_review_request' objects.
//
// Representation of an access review
type AccessReviewRequestBuilder struct {
	bitmap_         uint32
	accountUsername string
	action          string
	clusterID       string
	clusterUUID     string
	organizationID  string
	resourceType    string
	subscriptionID  string
}

// NewAccessReviewRequest creates a new builder of 'access_review_request' objects.
func NewAccessReviewRequest() *AccessReviewRequestBuilder {
	return &AccessReviewRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccessReviewRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AccountUsername sets the value of the 'account_username' attribute to the given value.
func (b *AccessReviewRequestBuilder) AccountUsername(value string) *AccessReviewRequestBuilder {
	b.accountUsername = value
	b.bitmap_ |= 1
	return b
}

// Action sets the value of the 'action' attribute to the given value.
func (b *AccessReviewRequestBuilder) Action(value string) *AccessReviewRequestBuilder {
	b.action = value
	b.bitmap_ |= 2
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *AccessReviewRequestBuilder) ClusterID(value string) *AccessReviewRequestBuilder {
	b.clusterID = value
	b.bitmap_ |= 4
	return b
}

// ClusterUUID sets the value of the 'cluster_UUID' attribute to the given value.
func (b *AccessReviewRequestBuilder) ClusterUUID(value string) *AccessReviewRequestBuilder {
	b.clusterUUID = value
	b.bitmap_ |= 8
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *AccessReviewRequestBuilder) OrganizationID(value string) *AccessReviewRequestBuilder {
	b.organizationID = value
	b.bitmap_ |= 16
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *AccessReviewRequestBuilder) ResourceType(value string) *AccessReviewRequestBuilder {
	b.resourceType = value
	b.bitmap_ |= 32
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *AccessReviewRequestBuilder) SubscriptionID(value string) *AccessReviewRequestBuilder {
	b.subscriptionID = value
	b.bitmap_ |= 64
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccessReviewRequestBuilder) Copy(object *AccessReviewRequest) *AccessReviewRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.accountUsername = object.accountUsername
	b.action = object.action
	b.clusterID = object.clusterID
	b.clusterUUID = object.clusterUUID
	b.organizationID = object.organizationID
	b.resourceType = object.resourceType
	b.subscriptionID = object.subscriptionID
	return b
}

// Build creates a 'access_review_request' object using the configuration stored in the builder.
func (b *AccessReviewRequestBuilder) Build() (object *AccessReviewRequest, err error) {
	object = new(AccessReviewRequest)
	object.bitmap_ = b.bitmap_
	object.accountUsername = b.accountUsername
	object.action = b.action
	object.clusterID = b.clusterID
	object.clusterUUID = b.clusterUUID
	object.organizationID = b.organizationID
	object.resourceType = b.resourceType
	object.subscriptionID = b.subscriptionID
	return
}
