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

// Definition of a Web RCA event.
type FollowUpBuilder struct {
	fieldSet_    []bool
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
	return &FollowUpBuilder{
		fieldSet_: make([]bool, 16),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *FollowUpBuilder) Link(value bool) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *FollowUpBuilder) ID(value string) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *FollowUpBuilder) HREF(value string) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *FollowUpBuilder) Empty() bool {
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

// Archived sets the value of the 'archived' attribute to the given value.
func (b *FollowUpBuilder) Archived(value bool) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.archived = value
	b.fieldSet_[3] = true
	return b
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *FollowUpBuilder) CreatedAt(value time.Time) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.createdAt = value
	b.fieldSet_[4] = true
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *FollowUpBuilder) DeletedAt(value time.Time) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.deletedAt = value
	b.fieldSet_[5] = true
	return b
}

// Done sets the value of the 'done' attribute to the given value.
func (b *FollowUpBuilder) Done(value bool) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.done = value
	b.fieldSet_[6] = true
	return b
}

// FollowUpType sets the value of the 'follow_up_type' attribute to the given value.
func (b *FollowUpBuilder) FollowUpType(value string) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.followUpType = value
	b.fieldSet_[7] = true
	return b
}

// Incident sets the value of the 'incident' attribute to the given value.
//
// Definition of a Web RCA incident.
func (b *FollowUpBuilder) Incident(value *IncidentBuilder) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.incident = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// Owner sets the value of the 'owner' attribute to the given value.
func (b *FollowUpBuilder) Owner(value string) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.owner = value
	b.fieldSet_[9] = true
	return b
}

// Priority sets the value of the 'priority' attribute to the given value.
func (b *FollowUpBuilder) Priority(value string) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.priority = value
	b.fieldSet_[10] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *FollowUpBuilder) Status(value string) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.status = value
	b.fieldSet_[11] = true
	return b
}

// Title sets the value of the 'title' attribute to the given value.
func (b *FollowUpBuilder) Title(value string) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.title = value
	b.fieldSet_[12] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *FollowUpBuilder) UpdatedAt(value time.Time) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.updatedAt = value
	b.fieldSet_[13] = true
	return b
}

// Url sets the value of the 'url' attribute to the given value.
func (b *FollowUpBuilder) Url(value string) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.url = value
	b.fieldSet_[14] = true
	return b
}

// WorkedAt sets the value of the 'worked_at' attribute to the given value.
func (b *FollowUpBuilder) WorkedAt(value time.Time) *FollowUpBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.workedAt = value
	b.fieldSet_[15] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *FollowUpBuilder) Copy(object *FollowUp) *FollowUpBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
