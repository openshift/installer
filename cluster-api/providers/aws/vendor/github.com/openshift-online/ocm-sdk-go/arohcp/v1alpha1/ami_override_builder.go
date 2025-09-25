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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
)

// AMIOverrideBuilder contains the data and logic needed to build 'AMI_override' objects.
//
// AMIOverride specifies what Amazon Machine Image should be used for a particular product and region.
type AMIOverrideBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	ami     string
	product *v1.ProductBuilder
	region  *v1.CloudRegionBuilder
}

// NewAMIOverride creates a new builder of 'AMI_override' objects.
func NewAMIOverride() *AMIOverrideBuilder {
	return &AMIOverrideBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *AMIOverrideBuilder) Link(value bool) *AMIOverrideBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *AMIOverrideBuilder) ID(value string) *AMIOverrideBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *AMIOverrideBuilder) HREF(value string) *AMIOverrideBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AMIOverrideBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// AMI sets the value of the 'AMI' attribute to the given value.
func (b *AMIOverrideBuilder) AMI(value string) *AMIOverrideBuilder {
	b.ami = value
	b.bitmap_ |= 8
	return b
}

// Product sets the value of the 'product' attribute to the given value.
//
// Representation of an product that can be selected as a cluster type.
func (b *AMIOverrideBuilder) Product(value *v1.ProductBuilder) *AMIOverrideBuilder {
	b.product = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// Region sets the value of the 'region' attribute to the given value.
//
// Description of a region of a cloud provider.
func (b *AMIOverrideBuilder) Region(value *v1.CloudRegionBuilder) *AMIOverrideBuilder {
	b.region = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AMIOverrideBuilder) Copy(object *AMIOverride) *AMIOverrideBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.ami = object.ami
	if object.product != nil {
		b.product = v1.NewProduct().Copy(object.product)
	} else {
		b.product = nil
	}
	if object.region != nil {
		b.region = v1.NewCloudRegion().Copy(object.region)
	} else {
		b.region = nil
	}
	return b
}

// Build creates a 'AMI_override' object using the configuration stored in the builder.
func (b *AMIOverrideBuilder) Build() (object *AMIOverride, err error) {
	object = new(AMIOverride)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.ami = b.ami
	if b.product != nil {
		object.product, err = b.product.Build()
		if err != nil {
			return
		}
	}
	if b.region != nil {
		object.region, err = b.region.Build()
		if err != nil {
			return
		}
	}
	return
}
