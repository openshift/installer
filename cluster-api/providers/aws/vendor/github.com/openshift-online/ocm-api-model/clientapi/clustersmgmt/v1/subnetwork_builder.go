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

// AWS subnetwork object to be used while installing a cluster
type SubnetworkBuilder struct {
	fieldSet_        []bool
	cidrBlock        string
	availabilityZone string
	name             string
	subnetID         string
	public           bool
	redHatManaged    bool
}

// NewSubnetwork creates a new builder of 'subnetwork' objects.
func NewSubnetwork() *SubnetworkBuilder {
	return &SubnetworkBuilder{
		fieldSet_: make([]bool, 6),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SubnetworkBuilder) Empty() bool {
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

// CIDRBlock sets the value of the 'CIDR_block' attribute to the given value.
func (b *SubnetworkBuilder) CIDRBlock(value string) *SubnetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.cidrBlock = value
	b.fieldSet_[0] = true
	return b
}

// AvailabilityZone sets the value of the 'availability_zone' attribute to the given value.
func (b *SubnetworkBuilder) AvailabilityZone(value string) *SubnetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.availabilityZone = value
	b.fieldSet_[1] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *SubnetworkBuilder) Name(value string) *SubnetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.name = value
	b.fieldSet_[2] = true
	return b
}

// Public sets the value of the 'public' attribute to the given value.
func (b *SubnetworkBuilder) Public(value bool) *SubnetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.public = value
	b.fieldSet_[3] = true
	return b
}

// RedHatManaged sets the value of the 'red_hat_managed' attribute to the given value.
func (b *SubnetworkBuilder) RedHatManaged(value bool) *SubnetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.redHatManaged = value
	b.fieldSet_[4] = true
	return b
}

// SubnetID sets the value of the 'subnet_ID' attribute to the given value.
func (b *SubnetworkBuilder) SubnetID(value string) *SubnetworkBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 6)
	}
	b.subnetID = value
	b.fieldSet_[5] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SubnetworkBuilder) Copy(object *Subnetwork) *SubnetworkBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.cidrBlock = object.cidrBlock
	b.availabilityZone = object.availabilityZone
	b.name = object.name
	b.public = object.public
	b.redHatManaged = object.redHatManaged
	b.subnetID = object.subnetID
	return b
}

// Build creates a 'subnetwork' object using the configuration stored in the builder.
func (b *SubnetworkBuilder) Build() (object *Subnetwork, err error) {
	object = new(Subnetwork)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.cidrBlock = b.cidrBlock
	object.availabilityZone = b.availabilityZone
	object.name = b.name
	object.public = b.public
	object.redHatManaged = b.redHatManaged
	object.subnetID = b.subnetID
	return
}
