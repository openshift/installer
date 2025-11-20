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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

// Identifies computing resources
type ResourceBuilder struct {
	fieldSet_            []bool
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
	return &ResourceBuilder{
		fieldSet_: make([]bool, 9),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ResourceBuilder) Link(value bool) *ResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ResourceBuilder) ID(value string) *ResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ResourceBuilder) HREF(value string) *ResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ResourceBuilder) Empty() bool {
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

// BYOC sets the value of the 'BYOC' attribute to the given value.
func (b *ResourceBuilder) BYOC(value bool) *ResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.byoc = value
	b.fieldSet_[3] = true
	return b
}

// SKU sets the value of the 'SKU' attribute to the given value.
func (b *ResourceBuilder) SKU(value string) *ResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.sku = value
	b.fieldSet_[4] = true
	return b
}

// Allowed sets the value of the 'allowed' attribute to the given value.
func (b *ResourceBuilder) Allowed(value int) *ResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.allowed = value
	b.fieldSet_[5] = true
	return b
}

// AvailabilityZoneType sets the value of the 'availability_zone_type' attribute to the given value.
func (b *ResourceBuilder) AvailabilityZoneType(value string) *ResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.availabilityZoneType = value
	b.fieldSet_[6] = true
	return b
}

// ResourceName sets the value of the 'resource_name' attribute to the given value.
func (b *ResourceBuilder) ResourceName(value string) *ResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.resourceName = value
	b.fieldSet_[7] = true
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *ResourceBuilder) ResourceType(value string) *ResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.resourceType = value
	b.fieldSet_[8] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ResourceBuilder) Copy(object *Resource) *ResourceBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.byoc = b.byoc
	object.sku = b.sku
	object.allowed = b.allowed
	object.availabilityZoneType = b.availabilityZoneType
	object.resourceName = b.resourceName
	object.resourceType = b.resourceType
	return
}
