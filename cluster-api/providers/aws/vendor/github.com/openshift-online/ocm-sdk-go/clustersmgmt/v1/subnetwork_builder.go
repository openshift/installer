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

// SubnetworkBuilder contains the data and logic needed to build 'subnetwork' objects.
//
// AWS subnetwork object to be used while installing a cluster
type SubnetworkBuilder struct {
	bitmap_          uint32
	availabilityZone string
	name             string
	subnetID         string
	public           bool
}

// NewSubnetwork creates a new builder of 'subnetwork' objects.
func NewSubnetwork() *SubnetworkBuilder {
	return &SubnetworkBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SubnetworkBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AvailabilityZone sets the value of the 'availability_zone' attribute to the given value.
func (b *SubnetworkBuilder) AvailabilityZone(value string) *SubnetworkBuilder {
	b.availabilityZone = value
	b.bitmap_ |= 1
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *SubnetworkBuilder) Name(value string) *SubnetworkBuilder {
	b.name = value
	b.bitmap_ |= 2
	return b
}

// Public sets the value of the 'public' attribute to the given value.
func (b *SubnetworkBuilder) Public(value bool) *SubnetworkBuilder {
	b.public = value
	b.bitmap_ |= 4
	return b
}

// SubnetID sets the value of the 'subnet_ID' attribute to the given value.
func (b *SubnetworkBuilder) SubnetID(value string) *SubnetworkBuilder {
	b.subnetID = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SubnetworkBuilder) Copy(object *Subnetwork) *SubnetworkBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.availabilityZone = object.availabilityZone
	b.name = object.name
	b.public = object.public
	b.subnetID = object.subnetID
	return b
}

// Build creates a 'subnetwork' object using the configuration stored in the builder.
func (b *SubnetworkBuilder) Build() (object *Subnetwork, err error) {
	object = new(Subnetwork)
	object.bitmap_ = b.bitmap_
	object.availabilityZone = b.availabilityZone
	object.name = b.name
	object.public = b.public
	object.subnetID = b.subnetID
	return
}
