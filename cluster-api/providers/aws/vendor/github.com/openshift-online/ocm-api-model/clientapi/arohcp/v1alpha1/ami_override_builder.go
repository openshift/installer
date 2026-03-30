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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

// AMIOverride specifies what Amazon Machine Image should be used for a particular product and region.
type AMIOverrideBuilder struct {
	fieldSet_ []bool
	id        string
	href      string
	ami       string
	product   *ProductBuilder
	region    *CloudRegionBuilder
}

// NewAMIOverride creates a new builder of 'AMI_override' objects.
func NewAMIOverride() *AMIOverrideBuilder {
	return &AMIOverrideBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AMIOverrideBuilder) Link(value bool) *AMIOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AMIOverrideBuilder) ID(value string) *AMIOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AMIOverrideBuilder) HREF(value string) *AMIOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AMIOverrideBuilder) Empty() bool {
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

// AMI sets the value of the 'AMI' attribute to the given value.
func (b *AMIOverrideBuilder) AMI(value string) *AMIOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.ami = value
	b.fieldSet_[3] = true
	return b
}

// Product sets the value of the 'product' attribute to the given value.
//
// Representation of an product that can be selected as a cluster type.
func (b *AMIOverrideBuilder) Product(value *ProductBuilder) *AMIOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.product = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// Region sets the value of the 'region' attribute to the given value.
//
// Description of a region of a cloud provider.
func (b *AMIOverrideBuilder) Region(value *CloudRegionBuilder) *AMIOverrideBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.region = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AMIOverrideBuilder) Copy(object *AMIOverride) *AMIOverrideBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.ami = object.ami
	if object.product != nil {
		b.product = NewProduct().Copy(object.product)
	} else {
		b.product = nil
	}
	if object.region != nil {
		b.region = NewCloudRegion().Copy(object.region)
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
