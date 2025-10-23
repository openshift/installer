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

// Definition of a Web RCA notification.
type NotificationBuilder struct {
	fieldSet_ []bool
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
	return &NotificationBuilder{
		fieldSet_: make([]bool, 10),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *NotificationBuilder) Link(value bool) *NotificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *NotificationBuilder) ID(value string) *NotificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *NotificationBuilder) HREF(value string) *NotificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NotificationBuilder) Empty() bool {
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

// Checked sets the value of the 'checked' attribute to the given value.
func (b *NotificationBuilder) Checked(value bool) *NotificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.checked = value
	b.fieldSet_[3] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *NotificationBuilder) CreatedAt(value time.Time) *NotificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.createdAt = value
	b.fieldSet_[4] = true
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *NotificationBuilder) DeletedAt(value time.Time) *NotificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.deletedAt = value
	b.fieldSet_[5] = true
	return b
}

// Incident sets the value of the 'incident' attribute to the given value.
//
// Definition of a Web RCA incident.
func (b *NotificationBuilder) Incident(value *IncidentBuilder) *NotificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.incident = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *NotificationBuilder) Name(value string) *NotificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.name = value
	b.fieldSet_[7] = true
	return b
}

// Rank sets the value of the 'rank' attribute to the given value.
func (b *NotificationBuilder) Rank(value int) *NotificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.rank = value
	b.fieldSet_[8] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *NotificationBuilder) UpdatedAt(value time.Time) *NotificationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.updatedAt = value
	b.fieldSet_[9] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NotificationBuilder) Copy(object *Notification) *NotificationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
