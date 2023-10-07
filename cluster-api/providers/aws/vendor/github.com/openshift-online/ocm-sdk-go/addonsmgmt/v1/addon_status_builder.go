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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

// AddonStatusBuilder contains the data and logic needed to build 'addon_status' objects.
//
// Representation of an addon status.
type AddonStatusBuilder struct {
	bitmap_          uint32
	id               string
	href             string
	addonId          string
	correlationID    string
	statusConditions []*AddonStatusConditionBuilder
}

// NewAddonStatus creates a new builder of 'addon_status' objects.
func NewAddonStatus() *AddonStatusBuilder {
	return &AddonStatusBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AddonStatusBuilder) Link(value bool) *AddonStatusBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AddonStatusBuilder) ID(value string) *AddonStatusBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AddonStatusBuilder) HREF(value string) *AddonStatusBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonStatusBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// AddonId sets the value of the 'addon_id' attribute to the given value.
func (b *AddonStatusBuilder) AddonId(value string) *AddonStatusBuilder {
	b.addonId = value
	b.bitmap_ |= 8
	return b
}

// CorrelationID sets the value of the 'correlation_ID' attribute to the given value.
func (b *AddonStatusBuilder) CorrelationID(value string) *AddonStatusBuilder {
	b.correlationID = value
	b.bitmap_ |= 16
	return b
}

// StatusConditions sets the value of the 'status_conditions' attribute to the given values.
func (b *AddonStatusBuilder) StatusConditions(values ...*AddonStatusConditionBuilder) *AddonStatusBuilder {
	b.statusConditions = make([]*AddonStatusConditionBuilder, len(values))
	copy(b.statusConditions, values)
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonStatusBuilder) Copy(object *AddonStatus) *AddonStatusBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.addonId = object.addonId
	b.correlationID = object.correlationID
	if object.statusConditions != nil {
		b.statusConditions = make([]*AddonStatusConditionBuilder, len(object.statusConditions))
		for i, v := range object.statusConditions {
			b.statusConditions[i] = NewAddonStatusCondition().Copy(v)
		}
	} else {
		b.statusConditions = nil
	}
	return b
}

// Build creates a 'addon_status' object using the configuration stored in the builder.
func (b *AddonStatusBuilder) Build() (object *AddonStatus, err error) {
	object = new(AddonStatus)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.addonId = b.addonId
	object.correlationID = b.correlationID
	if b.statusConditions != nil {
		object.statusConditions = make([]*AddonStatusCondition, len(b.statusConditions))
		for i, v := range b.statusConditions {
			object.statusConditions[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
