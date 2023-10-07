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

package v1 // github.com/openshift-online/ocm-sdk-go/jobqueue/v1

import (
	time "time"
)

// QueueBuilder contains the data and logic needed to build 'queue' objects.
type QueueBuilder struct {
	bitmap_     uint32
	id          string
	href        string
	createdAt   time.Time
	maxAttempts int
	maxRunTime  int
	name        string
	updatedAt   time.Time
}

// NewQueue creates a new builder of 'queue' objects.
func NewQueue() *QueueBuilder {
	return &QueueBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *QueueBuilder) Link(value bool) *QueueBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *QueueBuilder) ID(value string) *QueueBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *QueueBuilder) HREF(value string) *QueueBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *QueueBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *QueueBuilder) CreatedAt(value time.Time) *QueueBuilder {
	b.createdAt = value
	b.bitmap_ |= 8
	return b
}

// MaxAttempts sets the value of the 'max_attempts' attribute to the given value.
func (b *QueueBuilder) MaxAttempts(value int) *QueueBuilder {
	b.maxAttempts = value
	b.bitmap_ |= 16
	return b
}

// MaxRunTime sets the value of the 'max_run_time' attribute to the given value.
func (b *QueueBuilder) MaxRunTime(value int) *QueueBuilder {
	b.maxRunTime = value
	b.bitmap_ |= 32
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *QueueBuilder) Name(value string) *QueueBuilder {
	b.name = value
	b.bitmap_ |= 64
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *QueueBuilder) UpdatedAt(value time.Time) *QueueBuilder {
	b.updatedAt = value
	b.bitmap_ |= 128
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *QueueBuilder) Copy(object *Queue) *QueueBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.createdAt = object.createdAt
	b.maxAttempts = object.maxAttempts
	b.maxRunTime = object.maxRunTime
	b.name = object.name
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'queue' object using the configuration stored in the builder.
func (b *QueueBuilder) Build() (object *Queue, err error) {
	object = new(Queue)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.createdAt = b.createdAt
	object.maxAttempts = b.maxAttempts
	object.maxRunTime = b.maxRunTime
	object.name = b.name
	object.updatedAt = b.updatedAt
	return
}
