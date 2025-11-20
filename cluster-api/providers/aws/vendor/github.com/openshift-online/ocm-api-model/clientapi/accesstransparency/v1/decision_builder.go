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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accesstransparency/v1

import (
	time "time"
)

// Representation of an decision.
type DecisionBuilder struct {
	fieldSet_     []bool
	id            string
	href          string
	createdAt     time.Time
	decidedBy     string
	decision      DecisionDecision
	justification string
	updatedAt     time.Time
}

// NewDecision creates a new builder of 'decision' objects.
func NewDecision() *DecisionBuilder {
	return &DecisionBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *DecisionBuilder) Link(value bool) *DecisionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *DecisionBuilder) ID(value string) *DecisionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *DecisionBuilder) HREF(value string) *DecisionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *DecisionBuilder) Empty() bool {
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
func (b *DecisionBuilder) CreatedAt(value time.Time) *DecisionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.createdAt = value
	b.fieldSet_[3] = true
	return b
}

// DecidedBy sets the value of the 'decided_by' attribute to the given value.
func (b *DecisionBuilder) DecidedBy(value string) *DecisionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.decidedBy = value
	b.fieldSet_[4] = true
	return b
}

// Decision sets the value of the 'decision' attribute to the given value.
//
// Possible decisions to a decision status.
func (b *DecisionBuilder) Decision(value DecisionDecision) *DecisionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.decision = value
	b.fieldSet_[5] = true
	return b
}

// Justification sets the value of the 'justification' attribute to the given value.
func (b *DecisionBuilder) Justification(value string) *DecisionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.justification = value
	b.fieldSet_[6] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *DecisionBuilder) UpdatedAt(value time.Time) *DecisionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.updatedAt = value
	b.fieldSet_[7] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *DecisionBuilder) Copy(object *Decision) *DecisionBuilder {
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
	b.decidedBy = object.decidedBy
	b.decision = object.decision
	b.justification = object.justification
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'decision' object using the configuration stored in the builder.
func (b *DecisionBuilder) Build() (object *Decision, err error) {
	object = new(Decision)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.createdAt = b.createdAt
	object.decidedBy = b.decidedBy
	object.decision = b.decision
	object.justification = b.justification
	object.updatedAt = b.updatedAt
	return
}
