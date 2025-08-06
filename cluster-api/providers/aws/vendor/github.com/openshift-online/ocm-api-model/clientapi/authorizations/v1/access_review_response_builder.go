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

// Representation of an access review response
type AccessReviewResponseBuilder struct {
	fieldSet_       []bool
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
	return &AccessReviewResponseBuilder{
		fieldSet_: make([]bool, 10),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccessReviewResponseBuilder) Empty() bool {
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
func (b *AccessReviewResponseBuilder) AccountUsername(value string) *AccessReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.accountUsername = value
	b.fieldSet_[0] = true
	return b
}

// Action sets the value of the 'action' attribute to the given value.
func (b *AccessReviewResponseBuilder) Action(value string) *AccessReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.action = value
	b.fieldSet_[1] = true
	return b
}

// Allowed sets the value of the 'allowed' attribute to the given value.
func (b *AccessReviewResponseBuilder) Allowed(value bool) *AccessReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.allowed = value
	b.fieldSet_[2] = true
	return b
}

// ClusterID sets the value of the 'cluster_ID' attribute to the given value.
func (b *AccessReviewResponseBuilder) ClusterID(value string) *AccessReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.clusterID = value
	b.fieldSet_[3] = true
	return b
}

// ClusterUUID sets the value of the 'cluster_UUID' attribute to the given value.
func (b *AccessReviewResponseBuilder) ClusterUUID(value string) *AccessReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.clusterUUID = value
	b.fieldSet_[4] = true
	return b
}

// IsOCMInternal sets the value of the 'is_OCM_internal' attribute to the given value.
func (b *AccessReviewResponseBuilder) IsOCMInternal(value bool) *AccessReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.isOCMInternal = value
	b.fieldSet_[5] = true
	return b
}

// OrganizationID sets the value of the 'organization_ID' attribute to the given value.
func (b *AccessReviewResponseBuilder) OrganizationID(value string) *AccessReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.organizationID = value
	b.fieldSet_[6] = true
	return b
}

// Reason sets the value of the 'reason' attribute to the given value.
func (b *AccessReviewResponseBuilder) Reason(value string) *AccessReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.reason = value
	b.fieldSet_[7] = true
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *AccessReviewResponseBuilder) ResourceType(value string) *AccessReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.resourceType = value
	b.fieldSet_[8] = true
	return b
}

// SubscriptionID sets the value of the 'subscription_ID' attribute to the given value.
func (b *AccessReviewResponseBuilder) SubscriptionID(value string) *AccessReviewResponseBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.subscriptionID = value
	b.fieldSet_[9] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccessReviewResponseBuilder) Copy(object *AccessReviewResponse) *AccessReviewResponseBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
