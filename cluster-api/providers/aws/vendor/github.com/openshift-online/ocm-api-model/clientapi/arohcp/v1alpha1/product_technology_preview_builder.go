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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	time "time"
)

// Representation of a product technology preview.
type ProductTechnologyPreviewBuilder struct {
	fieldSet_      []bool
	id             string
	href           string
	additionalText string
	endDate        time.Time
	startDate      time.Time
}

// NewProductTechnologyPreview creates a new builder of 'product_technology_preview' objects.
func NewProductTechnologyPreview() *ProductTechnologyPreviewBuilder {
	return &ProductTechnologyPreviewBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ProductTechnologyPreviewBuilder) Link(value bool) *ProductTechnologyPreviewBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ProductTechnologyPreviewBuilder) ID(value string) *ProductTechnologyPreviewBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ProductTechnologyPreviewBuilder) HREF(value string) *ProductTechnologyPreviewBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ProductTechnologyPreviewBuilder) Empty() bool {
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

// AdditionalText sets the value of the 'additional_text' attribute to the given value.
func (b *ProductTechnologyPreviewBuilder) AdditionalText(value string) *ProductTechnologyPreviewBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.additionalText = value
	b.fieldSet_[3] = true
	return b
}

// EndDate sets the value of the 'end_date' attribute to the given value.
func (b *ProductTechnologyPreviewBuilder) EndDate(value time.Time) *ProductTechnologyPreviewBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.endDate = value
	b.fieldSet_[4] = true
	return b
}

// StartDate sets the value of the 'start_date' attribute to the given value.
func (b *ProductTechnologyPreviewBuilder) StartDate(value time.Time) *ProductTechnologyPreviewBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.startDate = value
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ProductTechnologyPreviewBuilder) Copy(object *ProductTechnologyPreview) *ProductTechnologyPreviewBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.additionalText = object.additionalText
	b.endDate = object.endDate
	b.startDate = object.startDate
	return b
}

// Build creates a 'product_technology_preview' object using the configuration stored in the builder.
func (b *ProductTechnologyPreviewBuilder) Build() (object *ProductTechnologyPreview, err error) {
	object = new(ProductTechnologyPreview)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.additionalText = b.additionalText
	object.endDate = b.endDate
	object.startDate = b.startDate
	return
}
