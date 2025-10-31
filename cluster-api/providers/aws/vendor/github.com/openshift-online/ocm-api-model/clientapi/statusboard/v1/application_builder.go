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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/statusboard/v1

import (
	time "time"
)

// Definition of a Status Board application.
type ApplicationBuilder struct {
	fieldSet_ []bool
	id        string
	href      string
	createdAt time.Time
	fullname  string
	metadata  interface{}
	name      string
	owners    []*OwnerBuilder
	product   *ProductBuilder
	updatedAt time.Time
}

// NewApplication creates a new builder of 'application' objects.
func NewApplication() *ApplicationBuilder {
	return &ApplicationBuilder{
		fieldSet_: make([]bool, 10),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ApplicationBuilder) Link(value bool) *ApplicationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ApplicationBuilder) ID(value string) *ApplicationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ApplicationBuilder) HREF(value string) *ApplicationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ApplicationBuilder) Empty() bool {
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
func (b *ApplicationBuilder) CreatedAt(value time.Time) *ApplicationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.createdAt = value
	b.fieldSet_[3] = true
	return b
}

// Fullname sets the value of the 'fullname' attribute to the given value.
func (b *ApplicationBuilder) Fullname(value string) *ApplicationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.fullname = value
	b.fieldSet_[4] = true
	return b
}

// Metadata sets the value of the 'metadata' attribute to the given value.
func (b *ApplicationBuilder) Metadata(value interface{}) *ApplicationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.metadata = value
	b.fieldSet_[5] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ApplicationBuilder) Name(value string) *ApplicationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.name = value
	b.fieldSet_[6] = true
	return b
}

// Owners sets the value of the 'owners' attribute to the given values.
func (b *ApplicationBuilder) Owners(values ...*OwnerBuilder) *ApplicationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.owners = make([]*OwnerBuilder, len(values))
	copy(b.owners, values)
	b.fieldSet_[7] = true
	return b
}

// Product sets the value of the 'product' attribute to the given value.
//
// Definition of a Status Board product.
func (b *ApplicationBuilder) Product(value *ProductBuilder) *ApplicationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.product = value
	if value != nil {
		b.fieldSet_[8] = true
	} else {
		b.fieldSet_[8] = false
	}
	return b
}

// UpdatedAt sets the value of the 'updated_at' attribute to the given value.
func (b *ApplicationBuilder) UpdatedAt(value time.Time) *ApplicationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 10)
	}
	b.updatedAt = value
	b.fieldSet_[9] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ApplicationBuilder) Copy(object *Application) *ApplicationBuilder {
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
	b.fullname = object.fullname
	b.metadata = object.metadata
	b.name = object.name
	if object.owners != nil {
		b.owners = make([]*OwnerBuilder, len(object.owners))
		for i, v := range object.owners {
			b.owners[i] = NewOwner().Copy(v)
		}
	} else {
		b.owners = nil
	}
	if object.product != nil {
		b.product = NewProduct().Copy(object.product)
	} else {
		b.product = nil
	}
	b.updatedAt = object.updatedAt
	return b
}

// Build creates a 'application' object using the configuration stored in the builder.
func (b *ApplicationBuilder) Build() (object *Application, err error) {
	object = new(Application)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.createdAt = b.createdAt
	object.fullname = b.fullname
	object.metadata = b.metadata
	object.name = b.name
	if b.owners != nil {
		object.owners = make([]*Owner, len(b.owners))
		for i, v := range b.owners {
			object.owners[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.product != nil {
		object.product, err = b.product.Build()
		if err != nil {
			return
		}
	}
	object.updatedAt = b.updatedAt
	return
}
