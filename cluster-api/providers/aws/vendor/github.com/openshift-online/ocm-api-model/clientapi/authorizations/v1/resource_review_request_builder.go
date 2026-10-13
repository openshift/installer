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

// Request to perform a resource access review.
type ResourceReviewRequestBuilder struct {
	fieldSet_                   []bool
	accountUsername             string
	action                      string
	excludeSubscriptionStatuses []SubscriptionStatus
	resourceType                string
	reduceClusterList           bool
}

// NewResourceReviewRequest creates a new builder of 'resource_review_request' objects.
func NewResourceReviewRequest() *ResourceReviewRequestBuilder {
	return &ResourceReviewRequestBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ResourceReviewRequestBuilder) Empty() bool {
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
func (b *ResourceReviewRequestBuilder) AccountUsername(value string) *ResourceReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.accountUsername = value
	b.fieldSet_[0] = true
	return b
}

// Action sets the value of the 'action' attribute to the given value.
func (b *ResourceReviewRequestBuilder) Action(value string) *ResourceReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.action = value
	b.fieldSet_[1] = true
	return b
}

// ExcludeSubscriptionStatuses sets the value of the 'exclude_subscription_statuses' attribute to the given values.
func (b *ResourceReviewRequestBuilder) ExcludeSubscriptionStatuses(values ...SubscriptionStatus) *ResourceReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.excludeSubscriptionStatuses = make([]SubscriptionStatus, len(values))
	copy(b.excludeSubscriptionStatuses, values)
	b.fieldSet_[2] = true
	return b
}

// ReduceClusterList sets the value of the 'reduce_cluster_list' attribute to the given value.
func (b *ResourceReviewRequestBuilder) ReduceClusterList(value bool) *ResourceReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.reduceClusterList = value
	b.fieldSet_[3] = true
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *ResourceReviewRequestBuilder) ResourceType(value string) *ResourceReviewRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.resourceType = value
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ResourceReviewRequestBuilder) Copy(object *ResourceReviewRequest) *ResourceReviewRequestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.accountUsername = object.accountUsername
	b.action = object.action
	if object.excludeSubscriptionStatuses != nil {
		b.excludeSubscriptionStatuses = make([]SubscriptionStatus, len(object.excludeSubscriptionStatuses))
		copy(b.excludeSubscriptionStatuses, object.excludeSubscriptionStatuses)
	} else {
		b.excludeSubscriptionStatuses = nil
	}
	b.reduceClusterList = object.reduceClusterList
	b.resourceType = object.resourceType
	return b
}

// Build creates a 'resource_review_request' object using the configuration stored in the builder.
func (b *ResourceReviewRequestBuilder) Build() (object *ResourceReviewRequest, err error) {
	object = new(ResourceReviewRequest)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.accountUsername = b.accountUsername
	object.action = b.action
	if b.excludeSubscriptionStatuses != nil {
		object.excludeSubscriptionStatuses = make([]SubscriptionStatus, len(b.excludeSubscriptionStatuses))
		copy(object.excludeSubscriptionStatuses, b.excludeSubscriptionStatuses)
	}
	object.reduceClusterList = b.reduceClusterList
	object.resourceType = b.resourceType
	return
}
