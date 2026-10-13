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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

import (
	time "time"
)

// Representation of check running before the cluster is provisioned.
type InflightCheckBuilder struct {
	fieldSet_ []bool
	id        string
	href      string
	details   interface{}
	endedAt   time.Time
	name      string
	restarts  int
	startedAt time.Time
	state     InflightCheckState
}

// NewInflightCheck creates a new builder of 'inflight_check' objects.
func NewInflightCheck() *InflightCheckBuilder {
	return &InflightCheckBuilder{
		fieldSet_: make([]bool, 9),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *InflightCheckBuilder) Link(value bool) *InflightCheckBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *InflightCheckBuilder) ID(value string) *InflightCheckBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *InflightCheckBuilder) HREF(value string) *InflightCheckBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *InflightCheckBuilder) Empty() bool {
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

// Details sets the value of the 'details' attribute to the given value.
func (b *InflightCheckBuilder) Details(value interface{}) *InflightCheckBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.details = value
	b.fieldSet_[3] = true
	return b
}

// EndedAt sets the value of the 'ended_at' attribute to the given value.
func (b *InflightCheckBuilder) EndedAt(value time.Time) *InflightCheckBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.endedAt = value
	b.fieldSet_[4] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *InflightCheckBuilder) Name(value string) *InflightCheckBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.name = value
	b.fieldSet_[5] = true
	return b
}

// Restarts sets the value of the 'restarts' attribute to the given value.
func (b *InflightCheckBuilder) Restarts(value int) *InflightCheckBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.restarts = value
	b.fieldSet_[6] = true
	return b
}

// StartedAt sets the value of the 'started_at' attribute to the given value.
func (b *InflightCheckBuilder) StartedAt(value time.Time) *InflightCheckBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.startedAt = value
	b.fieldSet_[7] = true
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// State of an inflight check.
func (b *InflightCheckBuilder) State(value InflightCheckState) *InflightCheckBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.state = value
	b.fieldSet_[8] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *InflightCheckBuilder) Copy(object *InflightCheck) *InflightCheckBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.details = object.details
	b.endedAt = object.endedAt
	b.name = object.name
	b.restarts = object.restarts
	b.startedAt = object.startedAt
	b.state = object.state
	return b
}

// Build creates a 'inflight_check' object using the configuration stored in the builder.
func (b *InflightCheckBuilder) Build() (object *InflightCheck, err error) {
	object = new(InflightCheck)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.details = b.details
	object.endedAt = b.endedAt
	object.name = b.name
	object.restarts = b.restarts
	object.startedAt = b.startedAt
	object.state = b.state
	return
}
