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

// JobBuilder contains the data and logic needed to build 'job' objects.
//
// This struct is a job in a Job Queue.
type JobBuilder struct {
	bitmap_     uint32
	id          string
	href        string
	abandonedAt time.Time
	arguments   string
	attempts    int
	createdAt   time.Time
	receiptId   string
	updatedAt   time.Time
}

// NewJob creates a new builder of 'job' objects.
func NewJob() *JobBuilder {
	return &JobBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *JobBuilder) Link(value bool) *JobBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *JobBuilder) ID(value string) *JobBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *JobBuilder) HREF(value string) *JobBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *JobBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// AbandonedAt sets the value of the 'abandoned_at' attribute to the given value.
func (b *JobBuilder) AbandonedAt(value time.Time) *JobBuilder {
	b.abandonedAt = value
	b.bitmap_ |= 8
	return b
}

// Arguments sets the value of the 'arguments' attribute to the given value.
func (b *JobBuilder) Arguments(value string) *JobBuilder {
	b.arguments = value
	b.bitmap_ |= 16
	return b
}

// Attempts sets the value of the 'attempts' attribute to the given value.
func (b *JobBuilder) Attempts(value int) *JobBuilder {
	b.attempts = value
	b.bitmap_ |= 32
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *JobBuilder) CreatedAt(value time.Time) *JobBuilder {
	b.createdAt = value
	b.bitmap_ |= 64
	return b
}

// ReceiptId sets the value of the 'receipt_id' attribute to the given value.
func (b *JobBuilder) ReceiptId(value string) *JobBuilder {
	b.receiptId = value
	b.bitmap_ |= 128
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *JobBuilder) UpdatedAt(value time.Time) *JobBuilder {
	b.updatedAt = value
	b.bitmap_ |= 256
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *JobBuilder) Copy(object *Job) *JobBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.abandonedAt = object.abandonedAt
	b.arguments = object.arguments
	b.attempts = object.attempts
	b.createdAt = object.createdAt
	b.receiptId = object.receiptId
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'job' object using the configuration stored in the builder.
func (b *JobBuilder) Build() (object *Job, err error) {
	object = new(Job)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.abandonedAt = b.abandonedAt
	object.arguments = b.arguments
	object.attempts = b.attempts
	object.createdAt = b.createdAt
	object.receiptId = b.receiptId
	object.updatedAt = b.updatedAt
	return
}
