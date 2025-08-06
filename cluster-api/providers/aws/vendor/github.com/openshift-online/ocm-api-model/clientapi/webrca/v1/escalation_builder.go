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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/webrca/v1

import (
	time "time"
)

// Definition of a Web RCA escalation.
type EscalationBuilder struct {
	fieldSet_ []bool
	id        string
	href      string
	createdAt time.Time
	deletedAt time.Time
	updatedAt time.Time
	user      *UserBuilder
}

// NewEscalation creates a new builder of 'escalation' objects.
func NewEscalation() *EscalationBuilder {
	return &EscalationBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *EscalationBuilder) Link(value bool) *EscalationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *EscalationBuilder) ID(value string) *EscalationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *EscalationBuilder) HREF(value string) *EscalationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *EscalationBuilder) Empty() bool {
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

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *EscalationBuilder) CreatedAt(value time.Time) *EscalationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.createdAt = value
	b.fieldSet_[3] = true
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *EscalationBuilder) DeletedAt(value time.Time) *EscalationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.deletedAt = value
	b.fieldSet_[4] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *EscalationBuilder) UpdatedAt(value time.Time) *EscalationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.updatedAt = value
	b.fieldSet_[5] = true
	return b
}

// User sets the value of the 'user' attribute to the given value.
//
// Definition of a Web RCA user.
func (b *EscalationBuilder) User(value *UserBuilder) *EscalationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.user = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *EscalationBuilder) Copy(object *Escalation) *EscalationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.createdAt = object.createdAt
	b.deletedAt = object.deletedAt
	b.updatedAt = object.updatedAt
	if object.user != nil {
		b.user = NewUser().Copy(object.user)
	} else {
		b.user = nil
	}
	return b
}

// Build creates a 'escalation' object using the configuration stored in the builder.
func (b *EscalationBuilder) Build() (object *Escalation, err error) {
	object = new(Escalation)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.createdAt = b.createdAt
	object.deletedAt = b.deletedAt
	object.updatedAt = b.updatedAt
	if b.user != nil {
		object.user, err = b.user.Build()
		if err != nil {
			return
		}
	}
	return
}
