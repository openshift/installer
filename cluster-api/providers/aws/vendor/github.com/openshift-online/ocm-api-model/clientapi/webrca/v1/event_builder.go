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
type EventBuilder struct {
	fieldSet_            []bool
	id                   string
	href                 string
	createdAt            time.Time
	creator              *UserBuilder
	deletedAt            time.Time
	escalation           *EscalationBuilder
	eventType            string
	externalReferenceUrl string
	followUp             *FollowUpBuilder
	followUpChange       *FollowUpChangeBuilder
	handoff              *HandoffBuilder
	incident             *IncidentBuilder
	note                 string
	statusChange         *StatusChangeBuilder
	updatedAt            time.Time
}

// NewEvent creates a new builder of 'event' objects.
func NewEvent() *EventBuilder {
	return &EventBuilder{
		fieldSet_: make([]bool, 16),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *EventBuilder) Link(value bool) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *EventBuilder) ID(value string) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *EventBuilder) HREF(value string) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *EventBuilder) Empty() bool {
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

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *EventBuilder) CreatedAt(value time.Time) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.createdAt = value
	b.fieldSet_[3] = true
	return b
}

// Creator sets the value of the 'creator' attribute to the given value.
//
// Definition of a Web RCA user.
func (b *EventBuilder) Creator(value *UserBuilder) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.creator = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *EventBuilder) DeletedAt(value time.Time) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.deletedAt = value
	b.fieldSet_[5] = true
	return b
}

// Escalation sets the value of the 'escalation' attribute to the given value.
//
// Definition of a Web RCA escalation.
func (b *EventBuilder) Escalation(value *EscalationBuilder) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.escalation = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// EventType sets the value of the 'event_type' attribute to the given value.
func (b *EventBuilder) EventType(value string) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.eventType = value
	b.fieldSet_[7] = true
	return b
}

// ExternalReferenceUrl sets the value of the 'external_reference_url' attribute to the given value.
func (b *EventBuilder) ExternalReferenceUrl(value string) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.externalReferenceUrl = value
	b.fieldSet_[8] = true
	return b
}

// FollowUp sets the value of the 'follow_up' attribute to the given value.
//
// Definition of a Web RCA event.
func (b *EventBuilder) FollowUp(value *FollowUpBuilder) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.followUp = value
	if value != nil {
		b.fieldSet_[9] = true
	} else {
		b.fieldSet_[9] = false
	}
	return b
}

// FollowUpChange sets the value of the 'follow_up_change' attribute to the given value.
//
// Definition of a Web RCA event.
func (b *EventBuilder) FollowUpChange(value *FollowUpChangeBuilder) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.followUpChange = value
	if value != nil {
		b.fieldSet_[10] = true
	} else {
		b.fieldSet_[10] = false
	}
	return b
}

// Handoff sets the value of the 'handoff' attribute to the given value.
//
// Definition of a Web RCA handoff.
func (b *EventBuilder) Handoff(value *HandoffBuilder) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.handoff = value
	if value != nil {
		b.fieldSet_[11] = true
	} else {
		b.fieldSet_[11] = false
	}
	return b
}

// Incident sets the value of the 'incident' attribute to the given value.
//
// Definition of a Web RCA incident.
func (b *EventBuilder) Incident(value *IncidentBuilder) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.incident = value
	if value != nil {
		b.fieldSet_[12] = true
	} else {
		b.fieldSet_[12] = false
	}
	return b
}

// Note sets the value of the 'note' attribute to the given value.
func (b *EventBuilder) Note(value string) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.note = value
	b.fieldSet_[13] = true
	return b
}

// StatusChange sets the value of the 'status_change' attribute to the given value.
//
// Definition of a Web RCA event.
func (b *EventBuilder) StatusChange(value *StatusChangeBuilder) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.statusChange = value
	if value != nil {
		b.fieldSet_[14] = true
	} else {
		b.fieldSet_[14] = false
	}
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *EventBuilder) UpdatedAt(value time.Time) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 16)
	}
	b.updatedAt = value
	b.fieldSet_[15] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *EventBuilder) Copy(object *Event) *EventBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.createdAt = object.createdAt
	if object.creator != nil {
		b.creator = NewUser().Copy(object.creator)
	} else {
		b.creator = nil
	}
	b.deletedAt = object.deletedAt
	if object.escalation != nil {
		b.escalation = NewEscalation().Copy(object.escalation)
	} else {
		b.escalation = nil
	}
	b.eventType = object.eventType
	b.externalReferenceUrl = object.externalReferenceUrl
	if object.followUp != nil {
		b.followUp = NewFollowUp().Copy(object.followUp)
	} else {
		b.followUp = nil
	}
	if object.followUpChange != nil {
		b.followUpChange = NewFollowUpChange().Copy(object.followUpChange)
	} else {
		b.followUpChange = nil
	}
	if object.handoff != nil {
		b.handoff = NewHandoff().Copy(object.handoff)
	} else {
		b.handoff = nil
	}
	if object.incident != nil {
		b.incident = NewIncident().Copy(object.incident)
	} else {
		b.incident = nil
	}
	b.note = object.note
	if object.statusChange != nil {
		b.statusChange = NewStatusChange().Copy(object.statusChange)
	} else {
		b.statusChange = nil
	}
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'event' object using the configuration stored in the builder.
func (b *EventBuilder) Build() (object *Event, err error) {
	object = new(Event)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.createdAt = b.createdAt
	if b.creator != nil {
		object.creator, err = b.creator.Build()
		if err != nil {
			return
		}
	}
	object.deletedAt = b.deletedAt
	if b.escalation != nil {
		object.escalation, err = b.escalation.Build()
		if err != nil {
			return
		}
	}
	object.eventType = b.eventType
	object.externalReferenceUrl = b.externalReferenceUrl
	if b.followUp != nil {
		object.followUp, err = b.followUp.Build()
		if err != nil {
			return
		}
	}
	if b.followUpChange != nil {
		object.followUpChange, err = b.followUpChange.Build()
		if err != nil {
			return
		}
	}
	if b.handoff != nil {
		object.handoff, err = b.handoff.Build()
		if err != nil {
			return
		}
	}
	if b.incident != nil {
		object.incident, err = b.incident.Build()
		if err != nil {
			return
		}
	}
	object.note = b.note
	if b.statusChange != nil {
		object.statusChange, err = b.statusChange.Build()
		if err != nil {
			return
		}
	}
	object.updatedAt = b.updatedAt
	return
}
