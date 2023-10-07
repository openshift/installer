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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

import (
	time "time"
)

// EscalationBuilder contains the data and logic needed to build 'escalation' objects.
//
// Definition of a Web RCA escalation.
type EscalationBuilder struct {
	bitmap_   uint32
	id        string
	href      string
	createdAt time.Time
	deletedAt time.Time
	updatedAt time.Time
	user      *UserBuilder
}

// NewEscalation creates a new builder of 'escalation' objects.
func NewEscalation() *EscalationBuilder {
	return &EscalationBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *EscalationBuilder) Link(value bool) *EscalationBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *EscalationBuilder) ID(value string) *EscalationBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *EscalationBuilder) HREF(value string) *EscalationBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *EscalationBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *EscalationBuilder) CreatedAt(value time.Time) *EscalationBuilder {
	b.createdAt = value
	b.bitmap_ |= 8
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *EscalationBuilder) DeletedAt(value time.Time) *EscalationBuilder {
	b.deletedAt = value
	b.bitmap_ |= 16
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *EscalationBuilder) UpdatedAt(value time.Time) *EscalationBuilder {
	b.updatedAt = value
	b.bitmap_ |= 32
	return b
}

// User sets the value of the 'user' attribute to the given value.
//
// Definition of a Web RCA user.
func (b *EscalationBuilder) User(value *UserBuilder) *EscalationBuilder {
	b.user = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *EscalationBuilder) Copy(object *Escalation) *EscalationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
