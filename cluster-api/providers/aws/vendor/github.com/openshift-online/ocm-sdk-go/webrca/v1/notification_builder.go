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

// NotificationBuilder contains the data and logic needed to build 'notification' objects.
//
// Definition of a Web RCA notification.
type NotificationBuilder struct {
	bitmap_   uint32
	id        string
	href      string
	createdAt time.Time
	deletedAt time.Time
	incident  *IncidentBuilder
	name      string
	rank      int
	updatedAt time.Time
	checked   bool
}

// NewNotification creates a new builder of 'notification' objects.
func NewNotification() *NotificationBuilder {
	return &NotificationBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *NotificationBuilder) Link(value bool) *NotificationBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *NotificationBuilder) ID(value string) *NotificationBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *NotificationBuilder) HREF(value string) *NotificationBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NotificationBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Checked sets the value of the 'checked' attribute to the given value.
func (b *NotificationBuilder) Checked(value bool) *NotificationBuilder {
	b.checked = value
	b.bitmap_ |= 8
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *NotificationBuilder) CreatedAt(value time.Time) *NotificationBuilder {
	b.createdAt = value
	b.bitmap_ |= 16
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *NotificationBuilder) DeletedAt(value time.Time) *NotificationBuilder {
	b.deletedAt = value
	b.bitmap_ |= 32
	return b
}

// Incident sets the value of the 'incident' attribute to the given value.
//
// Definition of a Web RCA incident.
func (b *NotificationBuilder) Incident(value *IncidentBuilder) *NotificationBuilder {
	b.incident = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *NotificationBuilder) Name(value string) *NotificationBuilder {
	b.name = value
	b.bitmap_ |= 128
	return b
}

// Rank sets the value of the 'rank' attribute to the given value.
func (b *NotificationBuilder) Rank(value int) *NotificationBuilder {
	b.rank = value
	b.bitmap_ |= 256
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *NotificationBuilder) UpdatedAt(value time.Time) *NotificationBuilder {
	b.updatedAt = value
	b.bitmap_ |= 512
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NotificationBuilder) Copy(object *Notification) *NotificationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.checked = object.checked
	b.createdAt = object.createdAt
	b.deletedAt = object.deletedAt
	if object.incident != nil {
		b.incident = NewIncident().Copy(object.incident)
	} else {
		b.incident = nil
	}
	b.name = object.name
	b.rank = object.rank
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'notification' object using the configuration stored in the builder.
func (b *NotificationBuilder) Build() (object *Notification, err error) {
	object = new(Notification)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.checked = b.checked
	object.createdAt = b.createdAt
	object.deletedAt = b.deletedAt
	if b.incident != nil {
		object.incident, err = b.incident.Build()
		if err != nil {
			return
		}
	}
	object.name = b.name
	object.rank = b.rank
	object.updatedAt = b.updatedAt
	return
}
