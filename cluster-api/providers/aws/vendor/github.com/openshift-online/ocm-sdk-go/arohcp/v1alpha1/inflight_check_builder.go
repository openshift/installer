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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	time "time"
)

// InflightCheckBuilder contains the data and logic needed to build 'inflight_check' objects.
//
// Representation of check running before the cluster is provisioned.
type InflightCheckBuilder struct {
	bitmap_   uint32
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
	return &InflightCheckBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *InflightCheckBuilder) Link(value bool) *InflightCheckBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *InflightCheckBuilder) ID(value string) *InflightCheckBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *InflightCheckBuilder) HREF(value string) *InflightCheckBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *InflightCheckBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Details sets the value of the 'details' attribute to the given value.
func (b *InflightCheckBuilder) Details(value interface{}) *InflightCheckBuilder {
	b.details = value
	b.bitmap_ |= 8
	return b
}

// EndedAt sets the value of the 'ended_at' attribute to the given value.
func (b *InflightCheckBuilder) EndedAt(value time.Time) *InflightCheckBuilder {
	b.endedAt = value
	b.bitmap_ |= 16
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *InflightCheckBuilder) Name(value string) *InflightCheckBuilder {
	b.name = value
	b.bitmap_ |= 32
	return b
}

// Restarts sets the value of the 'restarts' attribute to the given value.
func (b *InflightCheckBuilder) Restarts(value int) *InflightCheckBuilder {
	b.restarts = value
	b.bitmap_ |= 64
	return b
}

// StartedAt sets the value of the 'started_at' attribute to the given value.
func (b *InflightCheckBuilder) StartedAt(value time.Time) *InflightCheckBuilder {
	b.startedAt = value
	b.bitmap_ |= 128
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// State of an inflight check.
func (b *InflightCheckBuilder) State(value InflightCheckState) *InflightCheckBuilder {
	b.state = value
	b.bitmap_ |= 256
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *InflightCheckBuilder) Copy(object *InflightCheck) *InflightCheckBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	object.details = b.details
	object.endedAt = b.endedAt
	object.name = b.name
	object.restarts = b.restarts
	object.startedAt = b.startedAt
	object.state = b.state
	return
}
