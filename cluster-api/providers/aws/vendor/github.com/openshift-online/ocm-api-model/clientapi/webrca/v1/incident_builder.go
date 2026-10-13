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

// Definition of a Web RCA incident.
type IncidentBuilder struct {
	fieldSet_            []bool
	id                   string
	href                 string
	createdAt            time.Time
	creatorId            string
	deletedAt            time.Time
	description          string
	externalCoordination []string
	incidentId           string
	incidentType         string
	lastUpdated          time.Time
	primaryTeam          string
	severity             string
	status               string
	summary              string
	updatedAt            time.Time
	workedAt             time.Time
}

// NewIncident creates a new builder of 'incident' objects.
func NewIncident() *IncidentBuilder {
	return &IncidentBuilder{
		fieldSet_: make([]bool, 17),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *IncidentBuilder) Link(value bool) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *IncidentBuilder) ID(value string) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *IncidentBuilder) HREF(value string) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *IncidentBuilder) Empty() bool {
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
func (b *IncidentBuilder) CreatedAt(value time.Time) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.createdAt = value
	b.fieldSet_[3] = true
	return b
}

// CreatorId sets the value of the 'creator_id' attribute to the given value.
func (b *IncidentBuilder) CreatorId(value string) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.creatorId = value
	b.fieldSet_[4] = true
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *IncidentBuilder) DeletedAt(value time.Time) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.deletedAt = value
	b.fieldSet_[5] = true
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *IncidentBuilder) Description(value string) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.description = value
	b.fieldSet_[6] = true
	return b
}

// ExternalCoordination sets the value of the 'external_coordination' attribute to the given values.
func (b *IncidentBuilder) ExternalCoordination(values ...string) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.externalCoordination = make([]string, len(values))
	copy(b.externalCoordination, values)
	b.fieldSet_[7] = true
	return b
}

// IncidentId sets the value of the 'incident_id' attribute to the given value.
func (b *IncidentBuilder) IncidentId(value string) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.incidentId = value
	b.fieldSet_[8] = true
	return b
}

// IncidentType sets the value of the 'incident_type' attribute to the given value.
func (b *IncidentBuilder) IncidentType(value string) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.incidentType = value
	b.fieldSet_[9] = true
	return b
}

// LastUpdated sets the value of the 'last_updated' attribute to the given value.
func (b *IncidentBuilder) LastUpdated(value time.Time) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.lastUpdated = value
	b.fieldSet_[10] = true
	return b
}

// PrimaryTeam sets the value of the 'primary_team' attribute to the given value.
func (b *IncidentBuilder) PrimaryTeam(value string) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.primaryTeam = value
	b.fieldSet_[11] = true
	return b
}

// Severity sets the value of the 'severity' attribute to the given value.
func (b *IncidentBuilder) Severity(value string) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.severity = value
	b.fieldSet_[12] = true
	return b
}

// Status sets the value of the 'status' attribute to the given value.
func (b *IncidentBuilder) Status(value string) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.status = value
	b.fieldSet_[13] = true
	return b
}

// Summary sets the value of the 'summary' attribute to the given value.
func (b *IncidentBuilder) Summary(value string) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.summary = value
	b.fieldSet_[14] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *IncidentBuilder) UpdatedAt(value time.Time) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.updatedAt = value
	b.fieldSet_[15] = true
	return b
}

// WorkedAt sets the value of the 'worked_at' attribute to the given value.
func (b *IncidentBuilder) WorkedAt(value time.Time) *IncidentBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.workedAt = value
	b.fieldSet_[16] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *IncidentBuilder) Copy(object *Incident) *IncidentBuilder {
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
	b.creatorId = object.creatorId
	b.deletedAt = object.deletedAt
	b.description = object.description
	if object.externalCoordination != nil {
		b.externalCoordination = make([]string, len(object.externalCoordination))
		copy(b.externalCoordination, object.externalCoordination)
	} else {
		b.externalCoordination = nil
	}
	b.incidentId = object.incidentId
	b.incidentType = object.incidentType
	b.lastUpdated = object.lastUpdated
	b.primaryTeam = object.primaryTeam
	b.severity = object.severity
	b.status = object.status
	b.summary = object.summary
	b.updatedAt = object.updatedAt
	b.workedAt = object.workedAt
	return b
}

// Build creates a 'incident' object using the configuration stored in the builder.
func (b *IncidentBuilder) Build() (object *Incident, err error) {
	object = new(Incident)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.createdAt = b.createdAt
	object.creatorId = b.creatorId
	object.deletedAt = b.deletedAt
	object.description = b.description
	if b.externalCoordination != nil {
		object.externalCoordination = make([]string, len(b.externalCoordination))
		copy(object.externalCoordination, b.externalCoordination)
	}
	object.incidentId = b.incidentId
	object.incidentType = b.incidentType
	object.lastUpdated = b.lastUpdated
	object.primaryTeam = b.primaryTeam
	object.severity = b.severity
	object.status = b.status
	object.summary = b.summary
	object.updatedAt = b.updatedAt
	object.workedAt = b.workedAt
	return
}
