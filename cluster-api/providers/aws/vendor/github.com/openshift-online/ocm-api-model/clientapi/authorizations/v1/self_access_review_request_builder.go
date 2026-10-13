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

// Representation of an access review performed against oneself
type SelfAccessReviewRequestBuilder struct {
	fieldSet_      []bool
	action         string
	clusterID      string
	clusterUUID    string
	organizationID string
	resourceType   string
	subscriptionID string
}

// NewSelfAccessReviewRequest creates a new builder of 'self_access_review_request' objects.
func NewSelfAccessReviewRequest() *SelfAccessReviewRequestBuilder {
	return &SelfAccessReviewRequestBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SelfAccessReviewRequestBuilder) Empty() bool {
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

// Action sets the value of the 'action' attribute to the given value.
func (b *SelfAccessReviewRequestBuilder) Action(value string) *SelfAccessReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.action = value
	b.fieldSet_[0] = true
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *SelfAccessReviewRequestBuilder) ClusterID(value string) *SelfAccessReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.clusterID = value
	b.fieldSet_[1] = true
	return b
}

// ClusterUUID sets the value of the 'cluster_UUID' attribute to the given value.
func (b *SelfAccessReviewRequestBuilder) ClusterUUID(value string) *SelfAccessReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.clusterUUID = value
	b.fieldSet_[2] = true
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *SelfAccessReviewRequestBuilder) OrganizationID(value string) *SelfAccessReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.organizationID = value
	b.fieldSet_[3] = true
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *SelfAccessReviewRequestBuilder) ResourceType(value string) *SelfAccessReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.resourceType = value
	b.fieldSet_[4] = true
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *SelfAccessReviewRequestBuilder) SubscriptionID(value string) *SelfAccessReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.subscriptionID = value
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SelfAccessReviewRequestBuilder) Copy(object *SelfAccessReviewRequest) *SelfAccessReviewRequestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.action = b.action
	object.clusterID = b.clusterID
	object.clusterUUID = b.clusterUUID
	object.organizationID = b.organizationID
	object.resourceType = b.resourceType
	object.subscriptionID = b.subscriptionID
	return
}
