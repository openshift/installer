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

// CloudVPCBuilder contains the data and logic needed to build 'cloud_VPC' objects.
//
// Description of a cloud provider virtual private cloud.
type CloudVPCBuilder struct {
	bitmap_    uint32
	awsSubnets []*SubnetworkBuilder
	id         string
	name       string
	subnets    []string
}

// NewCloudVPC creates a new builder of 'cloud_VPC' objects.
func NewCloudVPC() *CloudVPCBuilder {
	return &CloudVPCBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CloudVPCBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AWSSubnets sets the value of the 'AWS_subnets' attribute to the given values.
func (b *CloudVPCBuilder) AWSSubnets(values ...*SubnetworkBuilder) *CloudVPCBuilder {
	b.awsSubnets = make([]*SubnetworkBuilder, len(values))
	copy(b.awsSubnets, values)
	b.bitmap_ |= 1
	return b
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *CloudVPCBuilder) ID(value string) *CloudVPCBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *CloudVPCBuilder) Name(value string) *CloudVPCBuilder {
	b.name = value
	b.bitmap_ |= 4
	return b
}

// Subnets sets the value of the 'subnets' attribute to the given values.
func (b *CloudVPCBuilder) Subnets(values ...string) *CloudVPCBuilder {
	b.subnets = make([]string, len(values))
	copy(b.subnets, values)
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CloudVPCBuilder) Copy(object *CloudVPC) *CloudVPCBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.awsSubnets != nil {
		b.awsSubnets = make([]*SubnetworkBuilder, len(object.awsSubnets))
		for i, v := range object.awsSubnets {
			b.awsSubnets[i] = NewSubnetwork().Copy(v)
		}
	} else {
		b.awsSubnets = nil
	}
	b.id = object.id
	b.name = object.name
	if object.subnets != nil {
		b.subnets = make([]string, len(object.subnets))
		copy(b.subnets, object.subnets)
	} else {
		b.subnets = nil
	}
	return b
}

// Build creates a 'cloud_VPC' object using the configuration stored in the builder.
func (b *CloudVPCBuilder) Build() (object *CloudVPC, err error) {
	object = new(CloudVPC)
	object.bitmap_ = b.bitmap_
	if b.awsSubnets != nil {
		object.awsSubnets = make([]*Subnetwork, len(b.awsSubnets))
		for i, v := range b.awsSubnets {
			object.awsSubnets[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.id = b.id
	object.name = b.name
	if b.subnets != nil {
		object.subnets = make([]string, len(b.subnets))
		copy(object.subnets, b.subnets)
	}
	return
}
