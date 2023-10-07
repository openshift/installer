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

// ProductBuilder contains the data and logic needed to build 'product' objects.
//
// Representation of an product that can be selected as a cluster type.
type ProductBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	name    string
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

// Name sets the value of the 'name' attribute to the given value.
func (b *ProductBuilder) Name(value string) *ProductBuilder {
	b.name = value
	b.bitmap_ |= 8
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
	b.name = object.name
	return b
}

// Build creates a 'product' object using the configuration stored in the builder.
func (b *ProductBuilder) Build() (object *Product, err error) {
	object = new(Product)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.name = b.name
	return
}
