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

// FollowUpBuilder contains the data and logic needed to build 'follow_up' objects.
//
// Definition of a Web RCA event.
type FollowUpBuilder struct {
	bitmap_      uint32
	id           string
	href         string
	createdAt    time.Time
	deletedAt    time.Time
	followUpType string
	incident     *IncidentBuilder
	owner        string
	priority     string
	status       string
	title        string
	updatedAt    time.Time
	url          string
	workedAt     time.Time
	archived     bool
	done         bool
}

// NewFollowUp creates a new builder of 'follow_up' objects.
func NewFollowUp() *FollowUpBuilder {
	return &FollowUpBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *FollowUpBuilder) Link(value bool) *FollowUpBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *FollowUpBuilder) ID(value string) *FollowUpBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *FollowUpBuilder) HREF(value string) *FollowUpBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *FollowUpBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Archived sets the value of the 'archived' attribute to the given value.
func (b *FollowUpBuilder) Archived(value bool) *FollowUpBuilder {
	b.archived = value
	b.bitmap_ |= 8
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *FollowUpBuilder) CreatedAt(value time.Time) *FollowUpBuilder {
	b.createdAt = value
	b.bitmap_ |= 16
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *FollowUpBuilder) DeletedAt(value time.Time) *FollowUpBuilder {
	b.deletedAt = value
	b.bitmap_ |= 32
	return b
}

// Done sets the value of the 'done' attribute to the given value.
func (b *FollowUpBuilder) Done(value bool) *FollowUpBuilder {
	b.done = value
	b.bitmap_ |= 64
	return b
}

// FollowUpType sets the value of the 'follow_up_type' attribute to the given value.
func (b *FollowUpBuilder) FollowUpType(value string) *FollowUpBuilder {
	b.followUpType = value
	b.bitmap_ |= 128
	return b
}

// Incident sets the value of the 'incident' attribute to the given value.
//
// Definition of a Web RCA incident.
func (b *FollowUpBuilder) Incident(value *IncidentBuilder) *FollowUpBuilder {
	b.incident = value
	if value != nil {
		b.bitmap_ |= 256
	} else {
		b.bitmap_ &^= 256
	}
	return b
}

// Owner sets the value of the 'owner' attribute to the given value.
func (b *FollowUpBuilder) Owner(value string) *FollowUpBuilder {
	b.owner = value
	b.bitmap_ |= 512
	return b
}

// Priority sets the value of the 'priority' attribute to the given value.
func (b *FollowUpBuilder) Priority(value string) *FollowUpBuilder {
	b.priority = value
	b.bitmap_ |= 1024
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *FollowUpBuilder) Status(value string) *FollowUpBuilder {
	b.status = value
	b.bitmap_ |= 2048
	return b
}

// Title sets the value of the 'title' attribute to the given value.
func (b *FollowUpBuilder) Title(value string) *FollowUpBuilder {
	b.title = value
	b.bitmap_ |= 4096
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *FollowUpBuilder) UpdatedAt(value time.Time) *FollowUpBuilder {
	b.updatedAt = value
	b.bitmap_ |= 8192
	return b
}

// Url sets the value of the 'url' attribute to the given value.
func (b *FollowUpBuilder) Url(value string) *FollowUpBuilder {
	b.url = value
	b.bitmap_ |= 16384
	return b
}

// WorkedAt sets the value of the 'worked_at' attribute to the given value.
func (b *FollowUpBuilder) WorkedAt(value time.Time) *FollowUpBuilder {
	b.workedAt = value
	b.bitmap_ |= 32768
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *FollowUpBuilder) Copy(object *FollowUp) *FollowUpBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.archived = object.archived
	b.createdAt = object.createdAt
	b.deletedAt = object.deletedAt
	b.done = object.done
	b.followUpType = object.followUpType
	if object.incident != nil {
		b.incident = NewIncident().Copy(object.incident)
	} else {
		b.incident = nil
	}
	b.owner = object.owner
	b.priority = object.priority
	b.status = object.status
	b.title = object.title
	b.updatedAt = object.updatedAt
	b.url = object.url
	b.workedAt = object.workedAt
	return b
}

// Build creates a 'follow_up' object using the configuration stored in the builder.
func (b *FollowUpBuilder) Build() (object *FollowUp, err error) {
	object = new(FollowUp)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.archived = b.archived
	object.createdAt = b.createdAt
	object.deletedAt = b.deletedAt
	object.done = b.done
	object.followUpType = b.followUpType
	if b.incident != nil {
		object.incident, err = b.incident.Build()
		if err != nil {
			return
		}
	}
	object.owner = b.owner
	object.priority = b.priority
	object.status = b.status
	object.title = b.title
	object.updatedAt = b.updatedAt
	object.url = b.url
	object.workedAt = b.workedAt
	return
}
