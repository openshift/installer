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

// ProductMinimalVersionBuilder contains the data and logic needed to build 'product_minimal_version' objects.
//
// Representation of a product minimal version.
type ProductMinimalVersionBuilder struct {
	bitmap_   uint32
	id        string
	href      string
	rosaCli   string
	startDate time.Time
}

// NewProductMinimalVersion creates a new builder of 'product_minimal_version' objects.
func NewProductMinimalVersion() *ProductMinimalVersionBuilder {
	return &ProductMinimalVersionBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ProductMinimalVersionBuilder) Link(value bool) *ProductMinimalVersionBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ProductMinimalVersionBuilder) ID(value string) *ProductMinimalVersionBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ProductMinimalVersionBuilder) HREF(value string) *ProductMinimalVersionBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ProductMinimalVersionBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// RosaCli sets the value of the 'rosa_cli' attribute to the given value.
func (b *ProductMinimalVersionBuilder) RosaCli(value string) *ProductMinimalVersionBuilder {
	b.rosaCli = value
	b.bitmap_ |= 8
	return b
}

// StartDate sets the value of the 'start_date' attribute to the given value.
func (b *ProductMinimalVersionBuilder) StartDate(value time.Time) *ProductMinimalVersionBuilder {
	b.startDate = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ProductMinimalVersionBuilder) Copy(object *ProductMinimalVersion) *ProductMinimalVersionBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.rosaCli = object.rosaCli
	b.startDate = object.startDate
	return b
}

// Build creates a 'product_minimal_version' object using the configuration stored in the builder.
func (b *ProductMinimalVersionBuilder) Build() (object *ProductMinimalVersion, err error) {
	object = new(ProductMinimalVersion)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.rosaCli = b.rosaCli
	object.startDate = b.startDate
	return
}
