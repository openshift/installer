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

package v1 // github.com/openshift-online/ocm-sdk-go/webrca/v1

import (
	time "time"
)

// ProductBuilder contains the data and logic needed to build 'product' objects.
//
// Definition of a Web RCA product.
type ProductBuilder struct {
	bitmap_     uint32
	id          string
	href        string
	createdAt   time.Time
	deletedAt   time.Time
	productId   string
	productName string
	updatedAt   time.Time
}

// NewProduct creates a new builder of 'product' objects.
func NewProduct() *ProductBuilder {
	return &ProductBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ProductBuilder) Link(value bool) *ProductBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ProductBuilder) ID(value string) *ProductBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ProductBuilder) HREF(value string) *ProductBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ProductBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *ProductBuilder) CreatedAt(value time.Time) *ProductBuilder {
	b.createdAt = value
	b.bitmap_ |= 8
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *ProductBuilder) DeletedAt(value time.Time) *ProductBuilder {
	b.deletedAt = value
	b.bitmap_ |= 16
	return b
}

// ProductId sets the value of the 'product_id' attribute to the given value.
func (b *ProductBuilder) ProductId(value string) *ProductBuilder {
	b.productId = value
	b.bitmap_ |= 32
	return b
}

// ProductName sets the value of the 'product_name' attribute to the given value.
func (b *ProductBuilder) ProductName(value string) *ProductBuilder {
	b.productName = value
	b.bitmap_ |= 64
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ProductBuilder) UpdatedAt(value time.Time) *ProductBuilder {
	b.updatedAt = value
	b.bitmap_ |= 128
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ProductBuilder) Copy(object *Product) *ProductBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.createdAt = object.createdAt
	b.deletedAt = object.deletedAt
	b.productId = object.productId
	b.productName = object.productName
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'product' object using the configuration stored in the builder.
func (b *ProductBuilder) Build() (object *Product, err error) {
	object = new(Product)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.createdAt = b.createdAt
	object.deletedAt = b.deletedAt
	object.productId = b.productId
	object.productName = b.productName
	object.updatedAt = b.updatedAt
	return
}
