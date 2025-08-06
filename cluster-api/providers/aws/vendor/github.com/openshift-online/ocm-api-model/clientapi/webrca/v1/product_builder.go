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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/webrca/v1

import (
	time "time"
)

// Definition of a Web RCA product.
type ProductBuilder struct {
	fieldSet_   []bool
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
	return &ProductBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ProductBuilder) Link(value bool) *ProductBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ProductBuilder) ID(value string) *ProductBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ProductBuilder) HREF(value string) *ProductBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ProductBuilder) Empty() bool {
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

// CreatedAt sets the value of the 'created_at' attribute to the given value.
func (b *ProductBuilder) CreatedAt(value time.Time) *ProductBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.createdAt = value
	b.fieldSet_[3] = true
	return b
}

// DeletedAt sets the value of the 'deleted_at' attribute to the given value.
func (b *ProductBuilder) DeletedAt(value time.Time) *ProductBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.deletedAt = value
	b.fieldSet_[4] = true
	return b
}

// ProductId sets the value of the 'product_id' attribute to the given value.
func (b *ProductBuilder) ProductId(value string) *ProductBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.productId = value
	b.fieldSet_[5] = true
	return b
}

// ProductName sets the value of the 'product_name' attribute to the given value.
func (b *ProductBuilder) ProductName(value string) *ProductBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.productName = value
	b.fieldSet_[6] = true
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ProductBuilder) UpdatedAt(value time.Time) *ProductBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.updatedAt = value
	b.fieldSet_[7] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ProductBuilder) Copy(object *Product) *ProductBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.createdAt = b.createdAt
	object.deletedAt = b.deletedAt
	object.productId = b.productId
	object.productName = b.productName
	object.updatedAt = b.updatedAt
	return
}
