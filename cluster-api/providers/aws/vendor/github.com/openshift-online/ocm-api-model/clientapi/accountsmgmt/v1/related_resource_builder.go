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

// Resource which can be provisioned using the allowed quota.
type RelatedResourceBuilder struct {
	fieldSet_            []bool
	byoc                 string
	availabilityZoneType string
	billingModel         string
	cloudProvider        string
	cost                 int
	product              string
	resourceName         string
	resourceType         string
}

// NewRelatedResource creates a new builder of 'related_resource' objects.
func NewRelatedResource() *RelatedResourceBuilder {
	return &RelatedResourceBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RelatedResourceBuilder) Empty() bool {
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

// BYOC sets the value of the 'BYOC' attribute to the given value.
func (b *RelatedResourceBuilder) BYOC(value string) *RelatedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.byoc = value
	b.fieldSet_[0] = true
	return b
}

// AvailabilityZoneType sets the value of the 'availability_zone_type' attribute to the given value.
func (b *RelatedResourceBuilder) AvailabilityZoneType(value string) *RelatedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.availabilityZoneType = value
	b.fieldSet_[1] = true
	return b
}

// BillingModel sets the value of the 'billing_model' attribute to the given value.
func (b *RelatedResourceBuilder) BillingModel(value string) *RelatedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.billingModel = value
	b.fieldSet_[2] = true
	return b
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
func (b *RelatedResourceBuilder) CloudProvider(value string) *RelatedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.cloudProvider = value
	b.fieldSet_[3] = true
	return b
}

// Cost sets the value of the 'cost' attribute to the given value.
func (b *RelatedResourceBuilder) Cost(value int) *RelatedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.cost = value
	b.fieldSet_[4] = true
	return b
}

// Product sets the value of the 'product' attribute to the given value.
func (b *RelatedResourceBuilder) Product(value string) *RelatedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.product = value
	b.fieldSet_[5] = true
	return b
}

// ResourceName sets the value of the 'resource_name' attribute to the given value.
func (b *RelatedResourceBuilder) ResourceName(value string) *RelatedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.resourceName = value
	b.fieldSet_[6] = true
	return b
}

// ResourceType sets the value of the 'resource_type' attribute to the given value.
func (b *RelatedResourceBuilder) ResourceType(value string) *RelatedResourceBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.resourceType = value
	b.fieldSet_[7] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RelatedResourceBuilder) Copy(object *RelatedResource) *RelatedResourceBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.byoc = object.byoc
	b.availabilityZoneType = object.availabilityZoneType
	b.billingModel = object.billingModel
	b.cloudProvider = object.cloudProvider
	b.cost = object.cost
	b.product = object.product
	b.resourceName = object.resourceName
	b.resourceType = object.resourceType
	return b
}

// Build creates a 'related_resource' object using the configuration stored in the builder.
func (b *RelatedResourceBuilder) Build() (object *RelatedResource, err error) {
	object = new(RelatedResource)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.byoc = b.byoc
	object.availabilityZoneType = b.availabilityZoneType
	object.billingModel = b.billingModel
	object.cloudProvider = b.cloudProvider
	object.cost = b.cost
	object.product = b.product
	object.resourceName = b.resourceName
	object.resourceType = b.resourceType
	return
}
