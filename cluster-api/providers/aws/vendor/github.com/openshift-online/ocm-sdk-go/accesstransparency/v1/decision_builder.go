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

package v1 // github.com/openshift-online/ocm-sdk-go/accesstransparency/v1

import (
	time "time"
)

// DecisionBuilder contains the data and logic needed to build 'decision' objects.
//
// Representation of an decision.
type DecisionBuilder struct {
	bitmap_       uint32
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
	return &DecisionBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *DecisionBuilder) Link(value bool) *DecisionBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *DecisionBuilder) ID(value string) *DecisionBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *DecisionBuilder) HREF(value string) *DecisionBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *DecisionBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *DecisionBuilder) CreatedAt(value time.Time) *DecisionBuilder {
	b.createdAt = value
	b.bitmap_ |= 8
	return b
}

// DecidedBy sets the value of the 'decided_by' attribute to the given value.
func (b *DecisionBuilder) DecidedBy(value string) *DecisionBuilder {
	b.decidedBy = value
	b.bitmap_ |= 16
	return b
}

// Decision sets the value of the 'decision' attribute to the given value.
//
// Possible decisions to a decision status.
func (b *DecisionBuilder) Decision(value DecisionDecision) *DecisionBuilder {
	b.decision = value
	b.bitmap_ |= 32
	return b
}

// Justification sets the value of the 'justification' attribute to the given value.
func (b *DecisionBuilder) Justification(value string) *DecisionBuilder {
	b.justification = value
	b.bitmap_ |= 64
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *DecisionBuilder) UpdatedAt(value time.Time) *DecisionBuilder {
	b.updatedAt = value
	b.bitmap_ |= 128
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *DecisionBuilder) Copy(object *Decision) *DecisionBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	object.createdAt = b.createdAt
	object.decidedBy = b.decidedBy
	object.decision = b.decision
	object.justification = b.justification
	object.updatedAt = b.updatedAt
	return
}
