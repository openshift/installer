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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/jobqueue/v1

import (
	time "time"
)

// This struct is a job in a Job Queue.
type JobBuilder struct {
	fieldSet_   []bool
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
	return &JobBuilder{
		fieldSet_: make([]bool, 9),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *JobBuilder) Link(value bool) *JobBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *JobBuilder) ID(value string) *JobBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *JobBuilder) HREF(value string) *JobBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *JobBuilder) Empty() bool {
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

// AbandonedAt sets the value of the 'abandoned_at' attribute to the given value.
func (b *JobBuilder) AbandonedAt(value time.Time) *JobBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.abandonedAt = value
	b.fieldSet_[3] = true
	return b
}

// Arguments sets the value of the 'arguments' attribute to the given value.
func (b *JobBuilder) Arguments(value string) *JobBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.arguments = value
	b.fieldSet_[4] = true
	return b
}

// Attempts sets the value of the 'attempts' attribute to the given value.
func (b *JobBuilder) Attempts(value int) *JobBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.attempts = value
	b.fieldSet_[5] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *JobBuilder) CreatedAt(value time.Time) *JobBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.createdAt = value
	b.fieldSet_[6] = true
	return b
}

// ReceiptId sets the value of the 'receipt_id' attribute to the given value.
func (b *JobBuilder) ReceiptId(value string) *JobBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.receiptId = value
	b.fieldSet_[7] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *JobBuilder) UpdatedAt(value time.Time) *JobBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.updatedAt = value
	b.fieldSet_[8] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *JobBuilder) Copy(object *Job) *JobBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.abandonedAt = b.abandonedAt
	object.arguments = b.arguments
	object.attempts = b.attempts
	object.createdAt = b.createdAt
	object.receiptId = b.receiptId
	object.updatedAt = b.updatedAt
	return
}
