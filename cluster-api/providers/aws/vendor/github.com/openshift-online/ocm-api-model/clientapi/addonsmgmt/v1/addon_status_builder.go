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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1

// Representation of an addon status.
type AddonStatusBuilder struct {
	fieldSet_        []bool
	id               string
	href             string
	addonId          string
	correlationID    string
	statusConditions []*AddonStatusConditionBuilder
	version          string
}

// NewAddonStatus creates a new builder of 'addon_status' objects.
func NewAddonStatus() *AddonStatusBuilder {
	return &AddonStatusBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AddonStatusBuilder) Link(value bool) *AddonStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AddonStatusBuilder) ID(value string) *AddonStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AddonStatusBuilder) HREF(value string) *AddonStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonStatusBuilder) Empty() bool {
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

// AddonId sets the value of the 'addon_id' attribute to the given value.
func (b *AddonStatusBuilder) AddonId(value string) *AddonStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.addonId = value
	b.fieldSet_[3] = true
	return b
}

// CorrelationID sets the value of the 'correlation_ID' attribute to the given value.
func (b *AddonStatusBuilder) CorrelationID(value string) *AddonStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.correlationID = value
	b.fieldSet_[4] = true
	return b
}

// StatusConditions sets the value of the 'status_conditions' attribute to the given values.
func (b *AddonStatusBuilder) StatusConditions(values ...*AddonStatusConditionBuilder) *AddonStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.statusConditions = make([]*AddonStatusConditionBuilder, len(values))
	copy(b.statusConditions, values)
	b.fieldSet_[5] = true
	return b
}

// Version sets the value of the 'version' attribute to the given value.
func (b *AddonStatusBuilder) Version(value string) *AddonStatusBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.version = value
	b.fieldSet_[6] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonStatusBuilder) Copy(object *AddonStatus) *AddonStatusBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	b.version = object.version
	return b
}

// Build creates a 'addon_status' object using the configuration stored in the builder.
func (b *AddonStatusBuilder) Build() (object *AddonStatus, err error) {
	object = new(AddonStatus)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
	object.version = b.version
	return
}
