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

// Representation of the limited support reason override.
type LimitedSupportReasonOverrideBuilder struct {
	fieldSet_ []bool
	id        string
	href      string
	enabled   bool
}

// NewLimitedSupportReasonOverride creates a new builder of 'limited_support_reason_override' objects.
func NewLimitedSupportReasonOverride() *LimitedSupportReasonOverrideBuilder {
	return &LimitedSupportReasonOverrideBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *LimitedSupportReasonOverrideBuilder) Link(value bool) *LimitedSupportReasonOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *LimitedSupportReasonOverrideBuilder) ID(value string) *LimitedSupportReasonOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *LimitedSupportReasonOverrideBuilder) HREF(value string) *LimitedSupportReasonOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LimitedSupportReasonOverrideBuilder) Empty() bool {
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

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *LimitedSupportReasonOverrideBuilder) Enabled(value bool) *LimitedSupportReasonOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.enabled = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LimitedSupportReasonOverrideBuilder) Copy(object *LimitedSupportReasonOverride) *LimitedSupportReasonOverrideBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.enabled = object.enabled
	return b
}

// Build creates a 'limited_support_reason_override' object using the configuration stored in the builder.
func (b *LimitedSupportReasonOverrideBuilder) Build() (object *LimitedSupportReasonOverride, err error) {
	object = new(LimitedSupportReasonOverride)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.enabled = b.enabled
	return
}
