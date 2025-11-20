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

// A reason that a cluster is in limited support.
type LimitedSupportReasonBuilder struct {
	fieldSet_         []bool
	id                string
	href              string
	creationTimestamp time.Time
	details           string
	detectionType     DetectionType
	override          *LimitedSupportReasonOverrideBuilder
	summary           string
	template          *LimitedSupportReasonTemplateBuilder
}

// NewLimitedSupportReason creates a new builder of 'limited_support_reason' objects.
func NewLimitedSupportReason() *LimitedSupportReasonBuilder {
	return &LimitedSupportReasonBuilder{
		fieldSet_: make([]bool, 9),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *LimitedSupportReasonBuilder) Link(value bool) *LimitedSupportReasonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *LimitedSupportReasonBuilder) ID(value string) *LimitedSupportReasonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *LimitedSupportReasonBuilder) HREF(value string) *LimitedSupportReasonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LimitedSupportReasonBuilder) Empty() bool {
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

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *LimitedSupportReasonBuilder) CreationTimestamp(value time.Time) *LimitedSupportReasonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.creationTimestamp = value
	b.fieldSet_[3] = true
	return b
}

// Details sets the value of the 'details' attribute to the given value.
func (b *LimitedSupportReasonBuilder) Details(value string) *LimitedSupportReasonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.details = value
	b.fieldSet_[4] = true
	return b
}

// DetectionType sets the value of the 'detection_type' attribute to the given value.
func (b *LimitedSupportReasonBuilder) DetectionType(value DetectionType) *LimitedSupportReasonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.detectionType = value
	b.fieldSet_[5] = true
	return b
}

// Override sets the value of the 'override' attribute to the given value.
//
// Representation of the limited support reason override.
func (b *LimitedSupportReasonBuilder) Override(value *LimitedSupportReasonOverrideBuilder) *LimitedSupportReasonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.override = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Summary sets the value of the 'summary' attribute to the given value.
func (b *LimitedSupportReasonBuilder) Summary(value string) *LimitedSupportReasonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.summary = value
	b.fieldSet_[7] = true
	return b
}

// Template sets the value of the 'template' attribute to the given value.
//
// A template for cluster limited support reason.
func (b *LimitedSupportReasonBuilder) Template(value *LimitedSupportReasonTemplateBuilder) *LimitedSupportReasonBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.template = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LimitedSupportReasonBuilder) Copy(object *LimitedSupportReason) *LimitedSupportReasonBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.creationTimestamp = object.creationTimestamp
	b.details = object.details
	b.detectionType = object.detectionType
	if object.override != nil {
		b.override = NewLimitedSupportReasonOverride().Copy(object.override)
	} else {
		b.override = nil
	}
	b.summary = object.summary
	if object.template != nil {
		b.template = NewLimitedSupportReasonTemplate().Copy(object.template)
	} else {
		b.template = nil
	}
	return b
}

// Build creates a 'limited_support_reason' object using the configuration stored in the builder.
func (b *LimitedSupportReasonBuilder) Build() (object *LimitedSupportReason, err error) {
	object = new(LimitedSupportReason)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.creationTimestamp = b.creationTimestamp
	object.details = b.details
	object.detectionType = b.detectionType
	if b.override != nil {
		object.override, err = b.override.Build()
		if err != nil {
			return
		}
	}
	object.summary = b.summary
	if b.template != nil {
		object.template, err = b.template.Build()
		if err != nil {
			return
		}
	}
	return
}
