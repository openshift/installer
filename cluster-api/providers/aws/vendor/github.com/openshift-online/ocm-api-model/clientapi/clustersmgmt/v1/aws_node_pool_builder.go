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

// Representation of aws node pool specific parameters.
type AWSNodePoolBuilder struct {
	fieldSet_                  []bool
	id                         string
	href                       string
	additionalSecurityGroupIds []string
	availabilityZoneTypes      map[string]string
	capacityReservation        *AWSCapacityReservationBuilder
	ec2MetadataHttpTokens      Ec2MetadataHttpTokens
	instanceProfile            string
	instanceType               string
	rootVolume                 *AWSVolumeBuilder
	subnetOutposts             map[string]string
	tags                       map[string]string
}

// NewAWSNodePool creates a new builder of 'AWS_node_pool' objects.
func NewAWSNodePool() *AWSNodePoolBuilder {
	return &AWSNodePoolBuilder{
		fieldSet_: make([]bool, 12),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AWSNodePoolBuilder) Link(value bool) *AWSNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AWSNodePoolBuilder) ID(value string) *AWSNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AWSNodePoolBuilder) HREF(value string) *AWSNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSNodePoolBuilder) Empty() bool {
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

// AdditionalSecurityGroupIds sets the value of the 'additional_security_group_ids' attribute to the given values.
func (b *AWSNodePoolBuilder) AdditionalSecurityGroupIds(values ...string) *AWSNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.additionalSecurityGroupIds = make([]string, len(values))
	copy(b.additionalSecurityGroupIds, values)
	b.fieldSet_[3] = true
	return b
}

// AvailabilityZoneTypes sets the value of the 'availability_zone_types' attribute to the given value.
func (b *AWSNodePoolBuilder) AvailabilityZoneTypes(value map[string]string) *AWSNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.availabilityZoneTypes = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// CapacityReservation sets the value of the 'capacity_reservation' attribute to the given value.
//
// AWS Capacity Reservation specification.
func (b *AWSNodePoolBuilder) CapacityReservation(value *AWSCapacityReservationBuilder) *AWSNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.capacityReservation = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// Ec2MetadataHttpTokens sets the value of the 'ec_2_metadata_http_tokens' attribute to the given value.
//
// Which Ec2MetadataHttpTokens to use for metadata service interaction options for EC2 instances
func (b *AWSNodePoolBuilder) Ec2MetadataHttpTokens(value Ec2MetadataHttpTokens) *AWSNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.ec2MetadataHttpTokens = value
	b.fieldSet_[6] = true
	return b
}

// InstanceProfile sets the value of the 'instance_profile' attribute to the given value.
func (b *AWSNodePoolBuilder) InstanceProfile(value string) *AWSNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.instanceProfile = value
	b.fieldSet_[7] = true
	return b
}

// InstanceType sets the value of the 'instance_type' attribute to the given value.
func (b *AWSNodePoolBuilder) InstanceType(value string) *AWSNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.instanceType = value
	b.fieldSet_[8] = true
	return b
}

// RootVolume sets the value of the 'root_volume' attribute to the given value.
//
// Holds settings for an AWS storage volume.
func (b *AWSNodePoolBuilder) RootVolume(value *AWSVolumeBuilder) *AWSNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.rootVolume = value
	if value != nil {
		b.fieldSet_[9] = true
	} else {
		b.fieldSet_[9] = false
	}
	return b
}

// SubnetOutposts sets the value of the 'subnet_outposts' attribute to the given value.
func (b *AWSNodePoolBuilder) SubnetOutposts(value map[string]string) *AWSNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.subnetOutposts = value
	if value != nil {
		b.fieldSet_[10] = true
	} else {
		b.fieldSet_[10] = false
	}
	return b
}

// Tags sets the value of the 'tags' attribute to the given value.
func (b *AWSNodePoolBuilder) Tags(value map[string]string) *AWSNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.tags = value
	if value != nil {
		b.fieldSet_[11] = true
	} else {
		b.fieldSet_[11] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSNodePoolBuilder) Copy(object *AWSNodePool) *AWSNodePoolBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	if object.additionalSecurityGroupIds != nil {
		b.additionalSecurityGroupIds = make([]string, len(object.additionalSecurityGroupIds))
		copy(b.additionalSecurityGroupIds, object.additionalSecurityGroupIds)
	} else {
		b.additionalSecurityGroupIds = nil
	}
	if len(object.availabilityZoneTypes) > 0 {
		b.availabilityZoneTypes = map[string]string{}
		for k, v := range object.availabilityZoneTypes {
			b.availabilityZoneTypes[k] = v
		}
	} else {
		b.availabilityZoneTypes = nil
	}
	if object.capacityReservation != nil {
		b.capacityReservation = NewAWSCapacityReservation().Copy(object.capacityReservation)
	} else {
		b.capacityReservation = nil
	}
	b.ec2MetadataHttpTokens = object.ec2MetadataHttpTokens
	b.instanceProfile = object.instanceProfile
	b.instanceType = object.instanceType
	if object.rootVolume != nil {
		b.rootVolume = NewAWSVolume().Copy(object.rootVolume)
	} else {
		b.rootVolume = nil
	}
	if len(object.subnetOutposts) > 0 {
		b.subnetOutposts = map[string]string{}
		for k, v := range object.subnetOutposts {
			b.subnetOutposts[k] = v
		}
	} else {
		b.subnetOutposts = nil
	}
	if len(object.tags) > 0 {
		b.tags = map[string]string{}
		for k, v := range object.tags {
			b.tags[k] = v
		}
	} else {
		b.tags = nil
	}
	return b
}

// Build creates a 'AWS_node_pool' object using the configuration stored in the builder.
func (b *AWSNodePoolBuilder) Build() (object *AWSNodePool, err error) {
	object = new(AWSNodePool)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.additionalSecurityGroupIds != nil {
		object.additionalSecurityGroupIds = make([]string, len(b.additionalSecurityGroupIds))
		copy(object.additionalSecurityGroupIds, b.additionalSecurityGroupIds)
	}
	if b.availabilityZoneTypes != nil {
		object.availabilityZoneTypes = make(map[string]string)
		for k, v := range b.availabilityZoneTypes {
			object.availabilityZoneTypes[k] = v
		}
	}
	if b.capacityReservation != nil {
		object.capacityReservation, err = b.capacityReservation.Build()
		if err != nil {
			return
		}
	}
	object.ec2MetadataHttpTokens = b.ec2MetadataHttpTokens
	object.instanceProfile = b.instanceProfile
	object.instanceType = b.instanceType
	if b.rootVolume != nil {
		object.rootVolume, err = b.rootVolume.Build()
		if err != nil {
			return
		}
	}
	if b.subnetOutposts != nil {
		object.subnetOutposts = make(map[string]string)
		for k, v := range b.subnetOutposts {
			object.subnetOutposts[k] = v
		}
	}
	if b.tags != nil {
		object.tags = make(map[string]string)
		for k, v := range b.tags {
			object.tags[k] = v
		}
	}
	return
}
