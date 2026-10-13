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

// Description of a cloud provider virtual private cloud.
type CloudVPCBuilder struct {
	fieldSet_         []bool
	awsSecurityGroups []*SecurityGroupBuilder
	awsSubnets        []*SubnetworkBuilder
	cidrBlock         string
	id                string
	name              string
	subnets           []string
	redHatManaged     bool
}

// NewCloudVPC creates a new builder of 'cloud_VPC' objects.
func NewCloudVPC() *CloudVPCBuilder {
	return &CloudVPCBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CloudVPCBuilder) Empty() bool {
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

// AWSSecurityGroups sets the value of the 'AWS_security_groups' attribute to the given values.
func (b *CloudVPCBuilder) AWSSecurityGroups(values ...*SecurityGroupBuilder) *CloudVPCBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.awsSecurityGroups = make([]*SecurityGroupBuilder, len(values))
	copy(b.awsSecurityGroups, values)
	b.fieldSet_[0] = true
	return b
}

// AWSSubnets sets the value of the 'AWS_subnets' attribute to the given values.
func (b *CloudVPCBuilder) AWSSubnets(values ...*SubnetworkBuilder) *CloudVPCBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.awsSubnets = make([]*SubnetworkBuilder, len(values))
	copy(b.awsSubnets, values)
	b.fieldSet_[1] = true
	return b
}

// CIDRBlock sets the value of the 'CIDR_block' attribute to the given value.
func (b *CloudVPCBuilder) CIDRBlock(value string) *CloudVPCBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.cidrBlock = value
	b.fieldSet_[2] = true
	return b
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *CloudVPCBuilder) ID(value string) *CloudVPCBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.id = value
	b.fieldSet_[3] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *CloudVPCBuilder) Name(value string) *CloudVPCBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.name = value
	b.fieldSet_[4] = true
	return b
}

// RedHatManaged sets the value of the 'red_hat_managed' attribute to the given value.
func (b *CloudVPCBuilder) RedHatManaged(value bool) *CloudVPCBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.redHatManaged = value
	b.fieldSet_[5] = true
	return b
}

// Subnets sets the value of the 'subnets' attribute to the given values.
func (b *CloudVPCBuilder) Subnets(values ...string) *CloudVPCBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.subnets = make([]string, len(values))
	copy(b.subnets, values)
	b.fieldSet_[6] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CloudVPCBuilder) Copy(object *CloudVPC) *CloudVPCBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.awsSecurityGroups != nil {
		b.awsSecurityGroups = make([]*SecurityGroupBuilder, len(object.awsSecurityGroups))
		for i, v := range object.awsSecurityGroups {
			b.awsSecurityGroups[i] = NewSecurityGroup().Copy(v)
		}
	} else {
		b.awsSecurityGroups = nil
	}
	if object.awsSubnets != nil {
		b.awsSubnets = make([]*SubnetworkBuilder, len(object.awsSubnets))
		for i, v := range object.awsSubnets {
			b.awsSubnets[i] = NewSubnetwork().Copy(v)
		}
	} else {
		b.awsSubnets = nil
	}
	b.cidrBlock = object.cidrBlock
	b.id = object.id
	b.name = object.name
	b.redHatManaged = object.redHatManaged
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.awsSecurityGroups != nil {
		object.awsSecurityGroups = make([]*SecurityGroup, len(b.awsSecurityGroups))
		for i, v := range b.awsSecurityGroups {
			object.awsSecurityGroups[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.awsSubnets != nil {
		object.awsSubnets = make([]*Subnetwork, len(b.awsSubnets))
		for i, v := range b.awsSubnets {
			object.awsSubnets[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	object.cidrBlock = b.cidrBlock
	object.id = b.id
	object.name = b.name
	object.redHatManaged = b.redHatManaged
	if b.subnets != nil {
		object.subnets = make([]string, len(b.subnets))
		copy(object.subnets, b.subnets)
	}
	return
}
