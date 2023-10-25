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

// EventBuilder contains the data and logic needed to build 'event' objects.
//
// Definition of a Web RCA event.
type EventBuilder struct {
	bitmap_              uint32
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
	return &EventBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *EventBuilder) Link(value bool) *EventBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *EventBuilder) ID(value string) *EventBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *EventBuilder) HREF(value string) *EventBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *EventBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *EventBuilder) CreatedAt(value time.Time) *EventBuilder {
	b.createdAt = value
	b.bitmap_ |= 8
	return b
}

// Creator sets the value of the 'creator' attribute to the given value.
//
// Definition of a Web RCA user.
func (b *EventBuilder) Creator(value *UserBuilder) *EventBuilder {
	b.creator = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *EventBuilder) DeletedAt(value time.Time) *EventBuilder {
	b.deletedAt = value
	b.bitmap_ |= 32
	return b
}

// Escalation sets the value of the 'escalation' attribute to the given value.
//
// Definition of a Web RCA escalation.
func (b *EventBuilder) Escalation(value *EscalationBuilder) *EventBuilder {
	b.escalation = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// EventType sets the value of the 'event_type' attribute to the given value.
func (b *EventBuilder) EventType(value string) *EventBuilder {
	b.eventType = value
	b.bitmap_ |= 128
	return b
}

// ExternalReferenceUrl sets the value of the 'external_reference_url' attribute to the given value.
func (b *EventBuilder) ExternalReferenceUrl(value string) *EventBuilder {
	b.externalReferenceUrl = value
	b.bitmap_ |= 256
	return b
}

// FollowUp sets the value of the 'follow_up' attribute to the given value.
//
// Definition of a Web RCA event.
func (b *EventBuilder) FollowUp(value *FollowUpBuilder) *EventBuilder {
	b.followUp = value
	if value != nil {
		b.bitmap_ |= 512
	} else {
		b.bitmap_ &^= 512
	}
	return b
}

// FollowUpChange sets the value of the 'follow_up_change' attribute to the given value.
//
// Definition of a Web RCA event.
func (b *EventBuilder) FollowUpChange(value *FollowUpChangeBuilder) *EventBuilder {
	b.followUpChange = value
	if value != nil {
		b.bitmap_ |= 1024
	} else {
		b.bitmap_ &^= 1024
	}
	return b
}

// Handoff sets the value of the 'handoff' attribute to the given value.
//
// Definition of a Web RCA handoff.
func (b *EventBuilder) Handoff(value *HandoffBuilder) *EventBuilder {
	b.handoff = value
	if value != nil {
		b.bitmap_ |= 2048
	} else {
		b.bitmap_ &^= 2048
	}
	return b
}

// Incident sets the value of the 'incident' attribute to the given value.
//
// Definition of a Web RCA incident.
func (b *EventBuilder) Incident(value *IncidentBuilder) *EventBuilder {
	b.incident = value
	if value != nil {
		b.bitmap_ |= 4096
	} else {
		b.bitmap_ &^= 4096
	}
	return b
}

// Note sets the value of the 'note' attribute to the given value.
func (b *EventBuilder) Note(value string) *EventBuilder {
	b.note = value
	b.bitmap_ |= 8192
	return b
}

// StatusChange sets the value of the 'status_change' attribute to the given value.
//
// Definition of a Web RCA event.
func (b *EventBuilder) StatusChange(value *StatusChangeBuilder) *EventBuilder {
	b.statusChange = value
	if value != nil {
		b.bitmap_ |= 16384
	} else {
		b.bitmap_ &^= 16384
	}
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *EventBuilder) UpdatedAt(value time.Time) *EventBuilder {
	b.updatedAt = value
	b.bitmap_ |= 32768
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *EventBuilder) Copy(object *Event) *EventBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
