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

// AccessReviewResponseBuilder contains the data and logic needed to build 'access_review_response' objects.
//
// Representation of an access review response
type AccessReviewResponseBuilder struct {
	bitmap_         uint32
	accountUsername string
	action          string
	clusterID       string
	clusterUUID     string
	organizationID  string
	reason          string
	resourceType    string
	subscriptionID  string
	allowed         bool
	isOCMInternal   bool
}

// NewAccessReviewResponse creates a new builder of 'access_review_response' objects.
func NewAccessReviewResponse() *AccessReviewResponseBuilder {
	return &AccessReviewResponseBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccessReviewResponseBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AccountUsername sets the value of the 'account_username' attribute to the given value.
func (b *AccessReviewResponseBuilder) AccountUsername(value string) *AccessReviewResponseBuilder {
	b.accountUsername = value
	b.bitmap_ |= 1
	return b
}

// Action sets the value of the 'action' attribute to the given value.
func (b *AccessReviewResponseBuilder) Action(value string) *AccessReviewResponseBuilder {
	b.action = value
	b.bitmap_ |= 2
	return b
}

// Allowed sets the value of the 'allowed' attribute to the given value.
func (b *AccessReviewResponseBuilder) Allowed(value bool) *AccessReviewResponseBuilder {
	b.allowed = value
	b.bitmap_ |= 4
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *AccessReviewResponseBuilder) ClusterID(value string) *AccessReviewResponseBuilder {
	b.clusterID = value
	b.bitmap_ |= 8
	return b
}

// ClusterUUID sets the value of the 'cluster_UUID' attribute to the given value.
func (b *AccessReviewResponseBuilder) ClusterUUID(value string) *AccessReviewResponseBuilder {
	b.clusterUUID = value
	b.bitmap_ |= 16
	return b
}

// IsOCMInternal sets the value of the 'is_OCM_internal' attribute to the given value.
func (b *AccessReviewResponseBuilder) IsOCMInternal(value bool) *AccessReviewResponseBuilder {
	b.isOCMInternal = value
	b.bitmap_ |= 32
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *AccessReviewResponseBuilder) OrganizationID(value string) *AccessReviewResponseBuilder {
	b.organizationID = value
	b.bitmap_ |= 64
	return b
}

// Reason sets the value of the 'reason' attribute to the given value.
func (b *AccessReviewResponseBuilder) Reason(value string) *AccessReviewResponseBuilder {
	b.reason = value
	b.bitmap_ |= 128
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *AccessReviewResponseBuilder) ResourceType(value string) *AccessReviewResponseBuilder {
	b.resourceType = value
	b.bitmap_ |= 256
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *AccessReviewResponseBuilder) SubscriptionID(value string) *AccessReviewResponseBuilder {
	b.subscriptionID = value
	b.bitmap_ |= 512
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccessReviewResponseBuilder) Copy(object *AccessReviewResponse) *AccessReviewResponseBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.accountUsername = object.accountUsername
	b.action = object.action
	b.allowed = object.allowed
	b.clusterID = object.clusterID
	b.clusterUUID = object.clusterUUID
	b.isOCMInternal = object.isOCMInternal
	b.organizationID = object.organizationID
	b.reason = object.reason
	b.resourceType = object.resourceType
	b.subscriptionID = object.subscriptionID
	return b
}

// Build creates a 'access_review_response' object using the configuration stored in the builder.
func (b *AccessReviewResponseBuilder) Build() (object *AccessReviewResponse, err error) {
	object = new(AccessReviewResponse)
	object.bitmap_ = b.bitmap_
	object.accountUsername = b.accountUsername
	object.action = b.action
	object.allowed = b.allowed
	object.clusterID = b.clusterID
	object.clusterUUID = b.clusterUUID
	object.isOCMInternal = b.isOCMInternal
	object.organizationID = b.organizationID
	object.reason = b.reason
	object.resourceType = b.resourceType
	object.subscriptionID = b.subscriptionID
	return
}
