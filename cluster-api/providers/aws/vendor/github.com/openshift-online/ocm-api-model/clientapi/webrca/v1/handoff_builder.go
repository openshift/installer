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

// Definition of a Web RCA handoff.
type HandoffBuilder struct {
	fieldSet_   []bool
	id          string
	href        string
	createdAt   time.Time
	deletedAt   time.Time
	handoffFrom *UserBuilder
	handoffTo   *UserBuilder
	handoffType string
	updatedAt   time.Time
}

// NewHandoff creates a new builder of 'handoff' objects.
func NewHandoff() *HandoffBuilder {
	return &HandoffBuilder{
		fieldSet_: make([]bool, 9),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *HandoffBuilder) Link(value bool) *HandoffBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *HandoffBuilder) ID(value string) *HandoffBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *HandoffBuilder) HREF(value string) *HandoffBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *HandoffBuilder) Empty() bool {
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
func (b *HandoffBuilder) CreatedAt(value time.Time) *HandoffBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.createdAt = value
	b.fieldSet_[3] = true
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *HandoffBuilder) DeletedAt(value time.Time) *HandoffBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.deletedAt = value
	b.fieldSet_[4] = true
	return b
}

// HandoffFrom sets the value of the 'handoff_from' attribute to the given value.
//
// Definition of a Web RCA user.
func (b *HandoffBuilder) HandoffFrom(value *UserBuilder) *HandoffBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.handoffFrom = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// HandoffTo sets the value of the 'handoff_to' attribute to the given value.
//
// Definition of a Web RCA user.
func (b *HandoffBuilder) HandoffTo(value *UserBuilder) *HandoffBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.handoffTo = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// HandoffType sets the value of the 'handoff_type' attribute to the given value.
func (b *HandoffBuilder) HandoffType(value string) *HandoffBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.handoffType = value
	b.fieldSet_[7] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *HandoffBuilder) UpdatedAt(value time.Time) *HandoffBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.updatedAt = value
	b.fieldSet_[8] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *HandoffBuilder) Copy(object *Handoff) *HandoffBuilder {
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
	if object.handoffFrom != nil {
		b.handoffFrom = NewUser().Copy(object.handoffFrom)
	} else {
		b.handoffFrom = nil
	}
	if object.handoffTo != nil {
		b.handoffTo = NewUser().Copy(object.handoffTo)
	} else {
		b.handoffTo = nil
	}
	b.handoffType = object.handoffType
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'handoff' object using the configuration stored in the builder.
func (b *HandoffBuilder) Build() (object *Handoff, err error) {
	object = new(Handoff)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.createdAt = b.createdAt
	object.deletedAt = b.deletedAt
	if b.handoffFrom != nil {
		object.handoffFrom, err = b.handoffFrom.Build()
		if err != nil {
			return
		}
	}
	if b.handoffTo != nil {
		object.handoffTo, err = b.handoffTo.Build()
		if err != nil {
			return
		}
	}
	object.handoffType = b.handoffType
	object.updatedAt = b.updatedAt
	return
}
