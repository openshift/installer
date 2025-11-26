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

import (
	time "time"
)

// Representation of an access request.
type AccessRequestBuilder struct {
	fieldSet_             []bool
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
	return &AccessRequestBuilder{
		fieldSet_: make([]bool, 17),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AccessRequestBuilder) Link(value bool) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AccessRequestBuilder) ID(value string) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AccessRequestBuilder) HREF(value string) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccessRequestBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *AccessRequestBuilder) ClusterId(value string) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.clusterId = value
	b.fieldSet_[3] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *AccessRequestBuilder) CreatedAt(value time.Time) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.createdAt = value
	b.fieldSet_[4] = true
	return b
}

// Deadline sets the value of the 'deadline' attribute to the given value.
func (b *AccessRequestBuilder) Deadline(value string) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.deadline = value
	b.fieldSet_[5] = true
	return b
}

// DeadlineAt sets the value of the 'deadline_at' attribute to the given value.
func (b *AccessRequestBuilder) DeadlineAt(value time.Time) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.deadlineAt = value
	b.fieldSet_[6] = true
	return b
}

// Decisions sets the value of the 'decisions' attribute to the given values.
func (b *AccessRequestBuilder) Decisions(values ...*DecisionBuilder) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.decisions = make([]*DecisionBuilder, len(values))
	copy(b.decisions, values)
	b.fieldSet_[7] = true
	return b
}

// Duration sets the value of the 'duration' attribute to the given value.
func (b *AccessRequestBuilder) Duration(value string) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.duration = value
	b.fieldSet_[8] = true
	return b
}

// InternalSupportCaseId sets the value of the 'internal_support_case_id' attribute to the given value.
func (b *AccessRequestBuilder) InternalSupportCaseId(value string) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.internalSupportCaseId = value
	b.fieldSet_[9] = true
	return b
}

// Justification sets the value of the 'justification' attribute to the given value.
func (b *AccessRequestBuilder) Justification(value string) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.justification = value
	b.fieldSet_[10] = true
	return b
}

// OrganizationId sets the value of the 'organization_id' attribute to the given value.
func (b *AccessRequestBuilder) OrganizationId(value string) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.organizationId = value
	b.fieldSet_[11] = true
	return b
}

// RequestedBy sets the value of the 'requested_by' attribute to the given value.
func (b *AccessRequestBuilder) RequestedBy(value string) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.requestedBy = value
	b.fieldSet_[12] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
//
// Representation of an access request status.
func (b *AccessRequestBuilder) Status(value *AccessRequestStatusBuilder) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.status = value
	if value != nil {
		b.fieldSet_[13] = true
	} else {
		b.fieldSet_[13] = false
	}
	return b
}

// SubscriptionId sets the value of the 'subscription_id' attribute to the given value.
func (b *AccessRequestBuilder) SubscriptionId(value string) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.subscriptionId = value
	b.fieldSet_[14] = true
	return b
}

// SupportCaseId sets the value of the 'support_case_id' attribute to the given value.
func (b *AccessRequestBuilder) SupportCaseId(value string) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.supportCaseId = value
	b.fieldSet_[15] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *AccessRequestBuilder) UpdatedAt(value time.Time) *AccessRequestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.updatedAt = value
	b.fieldSet_[16] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccessRequestBuilder) Copy(object *AccessRequest) *AccessRequestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
