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

// Representation of an access request status.
type AccessRequestStatusBuilder struct {
	fieldSet_ []bool
	expiresAt time.Time
	state     AccessRequestState
}

// NewAccessRequestStatus creates a new builder of 'access_request_status' objects.
func NewAccessRequestStatus() *AccessRequestStatusBuilder {
	return &AccessRequestStatusBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AccessRequestStatusBuilder) Empty() bool {
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

// ExpiresAt sets the value of the 'expires_at' attribute to the given value.
func (b *AccessRequestStatusBuilder) ExpiresAt(value time.Time) *AccessRequestStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.expiresAt = value
	b.fieldSet_[0] = true
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Possible states to an access request status.
func (b *AccessRequestStatusBuilder) State(value AccessRequestState) *AccessRequestStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.state = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AccessRequestStatusBuilder) Copy(object *AccessRequestStatus) *AccessRequestStatusBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.expiresAt = object.expiresAt
	b.state = object.state
	return b
}

// Build creates a 'access_request_status' object using the configuration stored in the builder.
func (b *AccessRequestStatusBuilder) Build() (object *AccessRequestStatus, err error) {
	object = new(AccessRequestStatus)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.expiresAt = b.expiresAt
	object.state = b.state
	return
}
