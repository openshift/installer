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

import (
	time "time"
)

// AccessRequestBuilder contains the data and logic needed to build 'access_request' objects.
//
// Representation of an access request.
type AccessRequestBuilder struct {
	bitmap_               uint32
	id                    string
	href                  string
	clusterId             string
	createdAt             time.Time
	deadline              string
	deadlineAt            time.Time
	decisions             []*DecisionBuilder
	duration              string
	internalSupportCaseId string
	justification         string
	organizationId        string
	requestedBy           string
	status                *AccessRequestStatusBuilder
	subscriptionId        string
	supportCaseId         string
	updatedAt             time.Time
}

// NewAccessRequest creates a new builder of 'access_request' objects.
func NewAccessRequest() *AccessRequestBuilder {
	return &AccessRequestBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AccessRequestBuilder) Link(value bool) *AccessRequestBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AccessRequestBuilder) ID(value string) *AccessRequestBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AccessRequestBuilder) HREF(value string) *AccessRequestBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccessRequestBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *AccessRequestBuilder) ClusterId(value string) *AccessRequestBuilder {
	b.clusterId = value
	b.bitmap_ |= 8
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *AccessRequestBuilder) CreatedAt(value time.Time) *AccessRequestBuilder {
	b.createdAt = value
	b.bitmap_ |= 16
	return b
}

// Deadline sets the value of the 'deadline' attribute to the given value.
func (b *AccessRequestBuilder) Deadline(value string) *AccessRequestBuilder {
	b.deadline = value
	b.bitmap_ |= 32
	return b
}

// DeadlineAt sets the value of the 'deadline_at' attribute to the given value.
func (b *AccessRequestBuilder) DeadlineAt(value time.Time) *AccessRequestBuilder {
	b.deadlineAt = value
	b.bitmap_ |= 64
	return b
}

// Decisions sets the value of the 'decisions' attribute to the given values.
func (b *AccessRequestBuilder) Decisions(values ...*DecisionBuilder) *AccessRequestBuilder {
	b.decisions = make([]*DecisionBuilder, len(values))
	copy(b.decisions, values)
	b.bitmap_ |= 128
	return b
}

// Duration sets the value of the 'duration' attribute to the given value.
func (b *AccessRequestBuilder) Duration(value string) *AccessRequestBuilder {
	b.duration = value
	b.bitmap_ |= 256
	return b
}

// InternalSupportCaseId sets the value of the 'internal_support_case_id' attribute to the given value.
func (b *AccessRequestBuilder) InternalSupportCaseId(value string) *AccessRequestBuilder {
	b.internalSupportCaseId = value
	b.bitmap_ |= 512
	return b
}

// Justification sets the value of the 'justification' attribute to the given value.
func (b *AccessRequestBuilder) Justification(value string) *AccessRequestBuilder {
	b.justification = value
	b.bitmap_ |= 1024
	return b
}

// OrganizationId sets the value of the 'organization_id' attribute to the given value.
func (b *AccessRequestBuilder) OrganizationId(value string) *AccessRequestBuilder {
	b.organizationId = value
	b.bitmap_ |= 2048
	return b
}

// RequestedBy sets the value of the 'requested_by' attribute to the given value.
func (b *AccessRequestBuilder) RequestedBy(value string) *AccessRequestBuilder {
	b.requestedBy = value
	b.bitmap_ |= 4096
	return b
}

// Status sets the value of the 'status' attribute to the given value.
//
// Representation of an access request status.
func (b *AccessRequestBuilder) Status(value *AccessRequestStatusBuilder) *AccessRequestBuilder {
	b.status = value
	if value != nil {
		b.bitmap_ |= 8192
	} else {
		b.bitmap_ &^= 8192
	}
	return b
}

// SubscriptionId sets the value of the 'subscription_id' attribute to the given value.
func (b *AccessRequestBuilder) SubscriptionId(value string) *AccessRequestBuilder {
	b.subscriptionId = value
	b.bitmap_ |= 16384
	return b
}

// SupportCaseId sets the value of the 'support_case_id' attribute to the given value.
func (b *AccessRequestBuilder) SupportCaseId(value string) *AccessRequestBuilder {
	b.supportCaseId = value
	b.bitmap_ |= 32768
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *AccessRequestBuilder) UpdatedAt(value time.Time) *AccessRequestBuilder {
	b.updatedAt = value
	b.bitmap_ |= 65536
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccessRequestBuilder) Copy(object *AccessRequest) *AccessRequestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.clusterId = object.clusterId
	b.createdAt = object.createdAt
	b.deadline = object.deadline
	b.deadlineAt = object.deadlineAt
	if object.decisions != nil {
		b.decisions = make([]*DecisionBuilder, len(object.decisions))
		for i, v := range object.decisions {
			b.decisions[i] = NewDecision().Copy(v)
		}
	} else {
		b.decisions = nil
	}
	b.duration = object.duration
	b.internalSupportCaseId = object.internalSupportCaseId
	b.justification = object.justification
	b.organizationId = object.organizationId
	b.requestedBy = object.requestedBy
	if object.status != nil {
		b.status = NewAccessRequestStatus().Copy(object.status)
	} else {
		b.status = nil
	}
	b.subscriptionId = object.subscriptionId
	b.supportCaseId = object.supportCaseId
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'access_request' object using the configuration stored in the builder.
func (b *AccessRequestBuilder) Build() (object *AccessRequest, err error) {
	object = new(AccessRequest)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.clusterId = b.clusterId
	object.createdAt = b.createdAt
	object.deadline = b.deadline
	object.deadlineAt = b.deadlineAt
	if b.decisions != nil {
		object.decisions = make([]*Decision, len(b.decisions))
		for i, v := range b.decisions {
			object.decisions[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.duration = b.duration
	object.internalSupportCaseId = b.internalSupportCaseId
	object.justification = b.justification
	object.organizationId = b.organizationId
	object.requestedBy = b.requestedBy
	if b.status != nil {
		object.status, err = b.status.Build()
		if err != nil {
			return
		}
	}
	object.subscriptionId = b.subscriptionId
	object.supportCaseId = b.supportCaseId
	object.updatedAt = b.updatedAt
	return
}
