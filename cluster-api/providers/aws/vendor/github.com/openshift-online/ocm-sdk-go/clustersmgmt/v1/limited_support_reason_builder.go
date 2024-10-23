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

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	time "time"
)

// LimitedSupportReasonBuilder contains the data and logic needed to build 'limited_support_reason' objects.
//
// A reason that a cluster is in limited support.
type LimitedSupportReasonBuilder struct {
	bitmap_           uint32
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
	return &LimitedSupportReasonBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *LimitedSupportReasonBuilder) Link(value bool) *LimitedSupportReasonBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *LimitedSupportReasonBuilder) ID(value string) *LimitedSupportReasonBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *LimitedSupportReasonBuilder) HREF(value string) *LimitedSupportReasonBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LimitedSupportReasonBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *LimitedSupportReasonBuilder) CreationTimestamp(value time.Time) *LimitedSupportReasonBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 8
	return b
}

// Details sets the value of the 'details' attribute to the given value.
func (b *LimitedSupportReasonBuilder) Details(value string) *LimitedSupportReasonBuilder {
	b.details = value
	b.bitmap_ |= 16
	return b
}

// DetectionType sets the value of the 'detection_type' attribute to the given value.
func (b *LimitedSupportReasonBuilder) DetectionType(value DetectionType) *LimitedSupportReasonBuilder {
	b.detectionType = value
	b.bitmap_ |= 32
	return b
}

// Override sets the value of the 'override' attribute to the given value.
//
// Representation of the limited support reason override.
func (b *LimitedSupportReasonBuilder) Override(value *LimitedSupportReasonOverrideBuilder) *LimitedSupportReasonBuilder {
	b.override = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// Summary sets the value of the 'summary' attribute to the given value.
func (b *LimitedSupportReasonBuilder) Summary(value string) *LimitedSupportReasonBuilder {
	b.summary = value
	b.bitmap_ |= 128
	return b
}

// Template sets the value of the 'template' attribute to the given value.
//
// A template for cluster limited support reason.
func (b *LimitedSupportReasonBuilder) Template(value *LimitedSupportReasonTemplateBuilder) *LimitedSupportReasonBuilder {
	b.template = value
	if value != nil {
		b.bitmap_ |= 256
	} else {
		b.bitmap_ &^= 256
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LimitedSupportReasonBuilder) Copy(object *LimitedSupportReason) *LimitedSupportReasonBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
