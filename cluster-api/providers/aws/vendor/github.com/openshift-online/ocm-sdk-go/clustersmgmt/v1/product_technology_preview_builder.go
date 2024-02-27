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

// ProductTechnologyPreviewBuilder contains the data and logic needed to build 'product_technology_preview' objects.
//
// Representation of a product technology preview.
type ProductTechnologyPreviewBuilder struct {
	bitmap_        uint32
	id             string
	href           string
	additionalText string
	endDate        time.Time
	startDate      time.Time
}

// NewProductTechnologyPreview creates a new builder of 'product_technology_preview' objects.
func NewProductTechnologyPreview() *ProductTechnologyPreviewBuilder {
	return &ProductTechnologyPreviewBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ProductTechnologyPreviewBuilder) Link(value bool) *ProductTechnologyPreviewBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ProductTechnologyPreviewBuilder) ID(value string) *ProductTechnologyPreviewBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ProductTechnologyPreviewBuilder) HREF(value string) *ProductTechnologyPreviewBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ProductTechnologyPreviewBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// AdditionalText sets the value of the 'additional_text' attribute to the given value.
func (b *ProductTechnologyPreviewBuilder) AdditionalText(value string) *ProductTechnologyPreviewBuilder {
	b.additionalText = value
	b.bitmap_ |= 8
	return b
}

// EndDate sets the value of the 'end_date' attribute to the given value.
func (b *ProductTechnologyPreviewBuilder) EndDate(value time.Time) *ProductTechnologyPreviewBuilder {
	b.endDate = value
	b.bitmap_ |= 16
	return b
}

// StartDate sets the value of the 'start_date' attribute to the given value.
func (b *ProductTechnologyPreviewBuilder) StartDate(value time.Time) *ProductTechnologyPreviewBuilder {
	b.startDate = value
	b.bitmap_ |= 32
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ProductTechnologyPreviewBuilder) Copy(object *ProductTechnologyPreview) *ProductTechnologyPreviewBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	object.additionalText = b.additionalText
	object.endDate = b.endDate
	object.startDate = b.startDate
	return
}
