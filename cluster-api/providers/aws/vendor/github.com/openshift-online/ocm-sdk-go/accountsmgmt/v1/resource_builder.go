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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

// ResourceBuilder contains the data and logic needed to build 'resource' objects.
//
// Identifies computing resources
type ResourceBuilder struct {
	bitmap_              uint32
	id                   string
	href                 string
	sku                  string
	allowed              int
	availabilityZoneType string
	resourceName         string
	resourceType         string
	byoc                 bool
}

// NewResource creates a new builder of 'resource' objects.
func NewResource() *ResourceBuilder {
	return &ResourceBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ResourceBuilder) Link(value bool) *ResourceBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ResourceBuilder) ID(value string) *ResourceBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ResourceBuilder) HREF(value string) *ResourceBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ResourceBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// BYOC sets the value of the 'BYOC' attribute to the given value.
func (b *ResourceBuilder) BYOC(value bool) *ResourceBuilder {
	b.byoc = value
	b.bitmap_ |= 8
	return b
}

// SKU sets the value of the 'SKU' attribute to the given value.
func (b *ResourceBuilder) SKU(value string) *ResourceBuilder {
	b.sku = value
	b.bitmap_ |= 16
	return b
}

// Allowed sets the value of the 'allowed' attribute to the given value.
func (b *ResourceBuilder) Allowed(value int) *ResourceBuilder {
	b.allowed = value
	b.bitmap_ |= 32
	return b
}

// AvailabilityZoneType sets the value of the 'availability_zone_type' attribute to the given value.
func (b *ResourceBuilder) AvailabilityZoneType(value string) *ResourceBuilder {
	b.availabilityZoneType = value
	b.bitmap_ |= 64
	return b
}

// ResourceName sets the value of the 'resource_name' attribute to the given value.
func (b *ResourceBuilder) ResourceName(value string) *ResourceBuilder {
	b.resourceName = value
	b.bitmap_ |= 128
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *ResourceBuilder) ResourceType(value string) *ResourceBuilder {
	b.resourceType = value
	b.bitmap_ |= 256
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ResourceBuilder) Copy(object *Resource) *ResourceBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.byoc = object.byoc
	b.sku = object.sku
	b.allowed = object.allowed
	b.availabilityZoneType = object.availabilityZoneType
	b.resourceName = object.resourceName
	b.resourceType = object.resourceType
	return b
}

// Build creates a 'resource' object using the configuration stored in the builder.
func (b *ResourceBuilder) Build() (object *Resource, err error) {
	object = new(Resource)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.byoc = b.byoc
	object.sku = b.sku
	object.allowed = b.allowed
	object.availabilityZoneType = b.availabilityZoneType
	object.resourceName = b.resourceName
	object.resourceType = b.resourceType
	return
}
