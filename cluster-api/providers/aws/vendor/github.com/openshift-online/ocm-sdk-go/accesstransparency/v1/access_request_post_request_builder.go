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

package v1 // github.com/openshift-online/ocm-sdk-go/accesstransparency/v1

// AccessRequestPostRequestBuilder contains the data and logic needed to build 'access_request_post_request' objects.
//
// Representation of an access request post request.
type AccessRequestPostRequestBuilder struct {
	bitmap_               uint32
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
	return &AccessRequestPostRequestBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccessRequestPostRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) ClusterId(value string) *AccessRequestPostRequestBuilder {
	b.clusterId = value
	b.bitmap_ |= 1
	return b
}

// Deadline sets the value of the 'deadline' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) Deadline(value string) *AccessRequestPostRequestBuilder {
	b.deadline = value
	b.bitmap_ |= 2
	return b
}

// Duration sets the value of the 'duration' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) Duration(value string) *AccessRequestPostRequestBuilder {
	b.duration = value
	b.bitmap_ |= 4
	return b
}

// InternalSupportCaseId sets the value of the 'internal_support_case_id' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) InternalSupportCaseId(value string) *AccessRequestPostRequestBuilder {
	b.internalSupportCaseId = value
	b.bitmap_ |= 8
	return b
}

// Justification sets the value of the 'justification' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) Justification(value string) *AccessRequestPostRequestBuilder {
	b.justification = value
	b.bitmap_ |= 16
	return b
}

// SubscriptionId sets the value of the 'subscription_id' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) SubscriptionId(value string) *AccessRequestPostRequestBuilder {
	b.subscriptionId = value
	b.bitmap_ |= 32
	return b
}

// SupportCaseId sets the value of the 'support_case_id' attribute to the given value.
func (b *AccessRequestPostRequestBuilder) SupportCaseId(value string) *AccessRequestPostRequestBuilder {
	b.supportCaseId = value
	b.bitmap_ |= 64
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccessRequestPostRequestBuilder) Copy(object *AccessRequestPostRequest) *AccessRequestPostRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	object.clusterId = b.clusterId
	object.deadline = b.deadline
	object.duration = b.duration
	object.internalSupportCaseId = b.internalSupportCaseId
	object.justification = b.justification
	object.subscriptionId = b.subscriptionId
	object.supportCaseId = b.supportCaseId
	return
}
