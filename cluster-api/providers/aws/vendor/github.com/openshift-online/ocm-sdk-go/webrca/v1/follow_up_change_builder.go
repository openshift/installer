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

// FollowUpChangeBuilder contains the data and logic needed to build 'follow_up_change' objects.
//
// Definition of a Web RCA event.
type FollowUpChangeBuilder struct {
	bitmap_   uint32
	id        string
	href      string
	createdAt time.Time
	deletedAt time.Time
	followUp  *FollowUpBuilder
	status    interface{}
	updatedAt time.Time
}

// NewFollowUpChange creates a new builder of 'follow_up_change' objects.
func NewFollowUpChange() *FollowUpChangeBuilder {
	return &FollowUpChangeBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *FollowUpChangeBuilder) Link(value bool) *FollowUpChangeBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *FollowUpChangeBuilder) ID(value string) *FollowUpChangeBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *FollowUpChangeBuilder) HREF(value string) *FollowUpChangeBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *FollowUpChangeBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *FollowUpChangeBuilder) CreatedAt(value time.Time) *FollowUpChangeBuilder {
	b.createdAt = value
	b.bitmap_ |= 8
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *FollowUpChangeBuilder) DeletedAt(value time.Time) *FollowUpChangeBuilder {
	b.deletedAt = value
	b.bitmap_ |= 16
	return b
}

// FollowUp sets the value of the 'follow_up' attribute to the given value.
//
// Definition of a Web RCA event.
func (b *FollowUpChangeBuilder) FollowUp(value *FollowUpBuilder) *FollowUpChangeBuilder {
	b.followUp = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *FollowUpChangeBuilder) Status(value interface{}) *FollowUpChangeBuilder {
	b.status = value
	b.bitmap_ |= 64
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *FollowUpChangeBuilder) UpdatedAt(value time.Time) *FollowUpChangeBuilder {
	b.updatedAt = value
	b.bitmap_ |= 128
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *FollowUpChangeBuilder) Copy(object *FollowUpChange) *FollowUpChangeBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.createdAt = object.createdAt
	b.deletedAt = object.deletedAt
	if object.followUp != nil {
		b.followUp = NewFollowUp().Copy(object.followUp)
	} else {
		b.followUp = nil
	}
	b.status = object.status
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'follow_up_change' object using the configuration stored in the builder.
func (b *FollowUpChangeBuilder) Build() (object *FollowUpChange, err error) {
	object = new(FollowUpChange)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.createdAt = b.createdAt
	object.deletedAt = b.deletedAt
	if b.followUp != nil {
		object.followUp, err = b.followUp.Build()
		if err != nil {
			return
		}
	}
	object.status = b.status
	object.updatedAt = b.updatedAt
	return
}
