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

// StatusChangeBuilder contains the data and logic needed to build 'status_change' objects.
//
// Definition of a Web RCA event.
type StatusChangeBuilder struct {
	bitmap_   uint32
	id        string
	href      string
	createdAt time.Time
	deletedAt time.Time
	status    interface{}
	statusId  string
	updatedAt time.Time
}

// NewStatusChange creates a new builder of 'status_change' objects.
func NewStatusChange() *StatusChangeBuilder {
	return &StatusChangeBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *StatusChangeBuilder) Link(value bool) *StatusChangeBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *StatusChangeBuilder) ID(value string) *StatusChangeBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *StatusChangeBuilder) HREF(value string) *StatusChangeBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *StatusChangeBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *StatusChangeBuilder) CreatedAt(value time.Time) *StatusChangeBuilder {
	b.createdAt = value
	b.bitmap_ |= 8
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *StatusChangeBuilder) DeletedAt(value time.Time) *StatusChangeBuilder {
	b.deletedAt = value
	b.bitmap_ |= 16
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *StatusChangeBuilder) Status(value interface{}) *StatusChangeBuilder {
	b.status = value
	b.bitmap_ |= 32
	return b
}

// StatusId sets the value of the 'status_id' attribute to the given value.
func (b *StatusChangeBuilder) StatusId(value string) *StatusChangeBuilder {
	b.statusId = value
	b.bitmap_ |= 64
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *StatusChangeBuilder) UpdatedAt(value time.Time) *StatusChangeBuilder {
	b.updatedAt = value
	b.bitmap_ |= 128
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *StatusChangeBuilder) Copy(object *StatusChange) *StatusChangeBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.createdAt = object.createdAt
	b.deletedAt = object.deletedAt
	b.status = object.status
	b.statusId = object.statusId
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'status_change' object using the configuration stored in the builder.
func (b *StatusChangeBuilder) Build() (object *StatusChange, err error) {
	object = new(StatusChange)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.createdAt = b.createdAt
	object.deletedAt = b.deletedAt
	object.status = b.status
	object.statusId = b.statusId
	object.updatedAt = b.updatedAt
	return
}
