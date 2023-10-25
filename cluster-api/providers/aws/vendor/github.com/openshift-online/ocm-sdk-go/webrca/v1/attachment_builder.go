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

// AttachmentBuilder contains the data and logic needed to build 'attachment' objects.
//
// Definition of a Web RCA attachment.
type AttachmentBuilder struct {
	bitmap_     uint32
	id          string
	href        string
	contentType string
	createdAt   time.Time
	creator     *UserBuilder
	deletedAt   time.Time
	event       *EventBuilder
	fileSize    int
	name        string
	updatedAt   time.Time
}

// NewAttachment creates a new builder of 'attachment' objects.
func NewAttachment() *AttachmentBuilder {
	return &AttachmentBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AttachmentBuilder) Link(value bool) *AttachmentBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AttachmentBuilder) ID(value string) *AttachmentBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AttachmentBuilder) HREF(value string) *AttachmentBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AttachmentBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// ContentType sets the value of the 'content_type' attribute to the given value.
func (b *AttachmentBuilder) ContentType(value string) *AttachmentBuilder {
	b.contentType = value
	b.bitmap_ |= 8
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *AttachmentBuilder) CreatedAt(value time.Time) *AttachmentBuilder {
	b.createdAt = value
	b.bitmap_ |= 16
	return b
}

// Creator sets the value of the 'creator' attribute to the given value.
//
// Definition of a Web RCA user.
func (b *AttachmentBuilder) Creator(value *UserBuilder) *AttachmentBuilder {
	b.creator = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *AttachmentBuilder) DeletedAt(value time.Time) *AttachmentBuilder {
	b.deletedAt = value
	b.bitmap_ |= 64
	return b
}

// Event sets the value of the 'event' attribute to the given value.
//
// Definition of a Web RCA event.
func (b *AttachmentBuilder) Event(value *EventBuilder) *AttachmentBuilder {
	b.event = value
	if value != nil {
		b.bitmap_ |= 128
	} else {
		b.bitmap_ &^= 128
	}
	return b
}

// FileSize sets the value of the 'file_size' attribute to the given value.
func (b *AttachmentBuilder) FileSize(value int) *AttachmentBuilder {
	b.fileSize = value
	b.bitmap_ |= 256
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AttachmentBuilder) Name(value string) *AttachmentBuilder {
	b.name = value
	b.bitmap_ |= 512
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *AttachmentBuilder) UpdatedAt(value time.Time) *AttachmentBuilder {
	b.updatedAt = value
	b.bitmap_ |= 1024
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AttachmentBuilder) Copy(object *Attachment) *AttachmentBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.contentType = object.contentType
	b.createdAt = object.createdAt
	if object.creator != nil {
		b.creator = NewUser().Copy(object.creator)
	} else {
		b.creator = nil
	}
	b.deletedAt = object.deletedAt
	if object.event != nil {
		b.event = NewEvent().Copy(object.event)
	} else {
		b.event = nil
	}
	b.fileSize = object.fileSize
	b.name = object.name
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'attachment' object using the configuration stored in the builder.
func (b *AttachmentBuilder) Build() (object *Attachment, err error) {
	object = new(Attachment)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.contentType = b.contentType
	object.createdAt = b.createdAt
	if b.creator != nil {
		object.creator, err = b.creator.Build()
		if err != nil {
			return
		}
	}
	object.deletedAt = b.deletedAt
	if b.event != nil {
		object.event, err = b.event.Build()
		if err != nil {
			return
		}
	}
	object.fileSize = b.fileSize
	object.name = b.name
	object.updatedAt = b.updatedAt
	return
}
