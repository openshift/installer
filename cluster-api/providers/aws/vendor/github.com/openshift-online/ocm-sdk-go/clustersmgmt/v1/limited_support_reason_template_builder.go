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

// LimitedSupportReasonTemplateBuilder contains the data and logic needed to build 'limited_support_reason_template' objects.
//
// A template for cluster limited support reason.
type LimitedSupportReasonTemplateBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	details string
	summary string
}

// NewLimitedSupportReasonTemplate creates a new builder of 'limited_support_reason_template' objects.
func NewLimitedSupportReasonTemplate() *LimitedSupportReasonTemplateBuilder {
	return &LimitedSupportReasonTemplateBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *LimitedSupportReasonTemplateBuilder) Link(value bool) *LimitedSupportReasonTemplateBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *LimitedSupportReasonTemplateBuilder) ID(value string) *LimitedSupportReasonTemplateBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *LimitedSupportReasonTemplateBuilder) HREF(value string) *LimitedSupportReasonTemplateBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LimitedSupportReasonTemplateBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Details sets the value of the 'details' attribute to the given value.
func (b *LimitedSupportReasonTemplateBuilder) Details(value string) *LimitedSupportReasonTemplateBuilder {
	b.details = value
	b.bitmap_ |= 8
	return b
}

// Summary sets the value of the 'summary' attribute to the given value.
func (b *LimitedSupportReasonTemplateBuilder) Summary(value string) *LimitedSupportReasonTemplateBuilder {
	b.summary = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LimitedSupportReasonTemplateBuilder) Copy(object *LimitedSupportReasonTemplate) *LimitedSupportReasonTemplateBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.details = object.details
	b.summary = object.summary
	return b
}

// Build creates a 'limited_support_reason_template' object using the configuration stored in the builder.
func (b *LimitedSupportReasonTemplateBuilder) Build() (object *LimitedSupportReasonTemplate, err error) {
	object = new(LimitedSupportReasonTemplate)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.details = b.details
	object.summary = b.summary
	return
}
