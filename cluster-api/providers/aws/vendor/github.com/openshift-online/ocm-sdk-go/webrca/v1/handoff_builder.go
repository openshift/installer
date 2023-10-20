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

// HandoffBuilder contains the data and logic needed to build 'handoff' objects.
//
// Definition of a Web RCA handoff.
type HandoffBuilder struct {
	bitmap_     uint32
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
	return &HandoffBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *HandoffBuilder) Link(value bool) *HandoffBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *HandoffBuilder) ID(value string) *HandoffBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *HandoffBuilder) HREF(value string) *HandoffBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *HandoffBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *HandoffBuilder) CreatedAt(value time.Time) *HandoffBuilder {
	b.createdAt = value
	b.bitmap_ |= 8
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *HandoffBuilder) DeletedAt(value time.Time) *HandoffBuilder {
	b.deletedAt = value
	b.bitmap_ |= 16
	return b
}

// HandoffFrom sets the value of the 'handoff_from' attribute to the given value.
//
// Definition of a Web RCA user.
func (b *HandoffBuilder) HandoffFrom(value *UserBuilder) *HandoffBuilder {
	b.handoffFrom = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// HandoffTo sets the value of the 'handoff_to' attribute to the given value.
//
// Definition of a Web RCA user.
func (b *HandoffBuilder) HandoffTo(value *UserBuilder) *HandoffBuilder {
	b.handoffTo = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// HandoffType sets the value of the 'handoff_type' attribute to the given value.
func (b *HandoffBuilder) HandoffType(value string) *HandoffBuilder {
	b.handoffType = value
	b.bitmap_ |= 128
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *HandoffBuilder) UpdatedAt(value time.Time) *HandoffBuilder {
	b.updatedAt = value
	b.bitmap_ |= 256
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *HandoffBuilder) Copy(object *Handoff) *HandoffBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
