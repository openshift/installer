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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accesstransparency/v1

// Representation of an access request post request.
type AccessRequestPostRequestBuilder struct {
	fieldSet_             []bool
	clusterId             string
	deadline              string
	duration              string
	internalSupportCaseId string
	justification         string
	subscriptionId        string
	supportCaseId         string
}

// NewAccessRequestPostRequest creates a new builder of 'access_request_post_request' objects.
func NewAccessRequestPostRequest() *AccessRequestPostRequestBuilder {
	return &AccessRequestPostRequestBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccessRequestPostRequestBuilder) Empty() bool {
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

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) ClusterId(value string) *AccessRequestPostRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.clusterId = value
	b.fieldSet_[0] = true
	return b
}

// Deadline sets the value of the 'deadline' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) Deadline(value string) *AccessRequestPostRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.deadline = value
	b.fieldSet_[1] = true
	return b
}

// Duration sets the value of the 'duration' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) Duration(value string) *AccessRequestPostRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.duration = value
	b.fieldSet_[2] = true
	return b
}

// InternalSupportCaseId sets the value of the 'internal_support_case_id' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) InternalSupportCaseId(value string) *AccessRequestPostRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.internalSupportCaseId = value
	b.fieldSet_[3] = true
	return b
}

// Justification sets the value of the 'justification' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) Justification(value string) *AccessRequestPostRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.justification = value
	b.fieldSet_[4] = true
	return b
}

// SubscriptionId sets the value of the 'subscription_id' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) SubscriptionId(value string) *AccessRequestPostRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.subscriptionId = value
	b.fieldSet_[5] = true
	return b
}

// SupportCaseId sets the value of the 'support_case_id' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) SupportCaseId(value string) *AccessRequestPostRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.supportCaseId = value
	b.fieldSet_[6] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccessRequestPostRequestBuilder) Copy(object *AccessRequestPostRequest) *AccessRequestPostRequestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.clusterId = object.clusterId
	b.deadline = object.deadline
	b.duration = object.duration
	b.internalSupportCaseId = object.internalSupportCaseId
	b.justification = object.justification
	b.subscriptionId = object.subscriptionId
	b.supportCaseId = object.supportCaseId
	return b
}

// Build creates a 'access_request_post_request' object using the configuration stored in the builder.
func (b *AccessRequestPostRequestBuilder) Build() (object *AccessRequestPostRequest, err error) {
	object = new(AccessRequestPostRequest)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.clusterId = b.clusterId
	object.deadline = b.deadline
	object.duration = b.duration
	object.internalSupportCaseId = b.internalSupportCaseId
	object.justification = b.justification
	object.subscriptionId = b.subscriptionId
	object.supportCaseId = b.supportCaseId
	return
}
