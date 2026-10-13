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

// Representation of an addon catalog source object used by addon versions.
type AdditionalCatalogSourceBuilder struct {
	fieldSet_ []bool
	id        string
	image     string
	name      string
	enabled   bool
}

// NewAdditionalCatalogSource creates a new builder of 'additional_catalog_source' objects.
func NewAdditionalCatalogSource() *AdditionalCatalogSourceBuilder {
	return &AdditionalCatalogSourceBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AdditionalCatalogSourceBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *AdditionalCatalogSourceBuilder) ID(value string) *AdditionalCatalogSourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.id = value
	b.fieldSet_[0] = true
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AdditionalCatalogSourceBuilder) Enabled(value bool) *AdditionalCatalogSourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.enabled = value
	b.fieldSet_[1] = true
	return b
}

// Image sets the value of the 'image' attribute to the given value.
func (b *AdditionalCatalogSourceBuilder) Image(value string) *AdditionalCatalogSourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.image = value
	b.fieldSet_[2] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AdditionalCatalogSourceBuilder) Name(value string) *AdditionalCatalogSourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.name = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AdditionalCatalogSourceBuilder) Copy(object *AdditionalCatalogSource) *AdditionalCatalogSourceBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.enabled = object.enabled
	b.image = object.image
	b.name = object.name
	return b
}

// Build creates a 'additional_catalog_source' object using the configuration stored in the builder.
func (b *AdditionalCatalogSourceBuilder) Build() (object *AdditionalCatalogSource, err error) {
	object = new(AdditionalCatalogSource)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.id = b.id
	object.enabled = b.enabled
	object.image = b.image
	object.name = b.name
	return
}
