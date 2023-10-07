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

// AdditionalCatalogSourceBuilder contains the data and logic needed to build 'additional_catalog_source' objects.
//
// Representation of an addon catalog source object used by addon versions.
type AdditionalCatalogSourceBuilder struct {
	bitmap_ uint32
	id      string
	image   string
	name    string
	enabled bool
}

// NewAdditionalCatalogSource creates a new builder of 'additional_catalog_source' objects.
func NewAdditionalCatalogSource() *AdditionalCatalogSourceBuilder {
	return &AdditionalCatalogSourceBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AdditionalCatalogSourceBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *AdditionalCatalogSourceBuilder) ID(value string) *AdditionalCatalogSourceBuilder {
	b.id = value
	b.bitmap_ |= 1
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AdditionalCatalogSourceBuilder) Enabled(value bool) *AdditionalCatalogSourceBuilder {
	b.enabled = value
	b.bitmap_ |= 2
	return b
}

// Image sets the value of the 'image' attribute to the given value.
func (b *AdditionalCatalogSourceBuilder) Image(value string) *AdditionalCatalogSourceBuilder {
	b.image = value
	b.bitmap_ |= 4
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *AdditionalCatalogSourceBuilder) Name(value string) *AdditionalCatalogSourceBuilder {
	b.name = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AdditionalCatalogSourceBuilder) Copy(object *AdditionalCatalogSource) *AdditionalCatalogSourceBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.enabled = object.enabled
	b.image = object.image
	b.name = object.name
	return b
}

// Build creates a 'additional_catalog_source' object using the configuration stored in the builder.
func (b *AdditionalCatalogSourceBuilder) Build() (object *AdditionalCatalogSource, err error) {
	object = new(AdditionalCatalogSource)
	object.bitmap_ = b.bitmap_
	object.id = b.id
	object.enabled = b.enabled
	object.image = b.image
	object.name = b.name
	return
}
