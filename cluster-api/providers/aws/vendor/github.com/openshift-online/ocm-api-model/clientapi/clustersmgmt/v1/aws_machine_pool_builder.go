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

// Representation of aws machine pool specific parameters.
type AWSMachinePoolBuilder struct {
	fieldSet_                  []bool
	id                         string
	href                       string
	additionalSecurityGroupIds []string
	availabilityZoneTypes      map[string]string
	spotMarketOptions          *AWSSpotMarketOptionsBuilder
	subnetOutposts             map[string]string
	tags                       map[string]string
}

// NewAWSMachinePool creates a new builder of 'AWS_machine_pool' objects.
func NewAWSMachinePool() *AWSMachinePoolBuilder {
	return &AWSMachinePoolBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *AWSMachinePoolBuilder) Link(value bool) *AWSMachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *AWSMachinePoolBuilder) ID(value string) *AWSMachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *AWSMachinePoolBuilder) HREF(value string) *AWSMachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSMachinePoolBuilder) Empty() bool {
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
func (b *AWSMachinePoolBuilder) AdditionalSecurityGroupIds(values ...string) *AWSMachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.additionalSecurityGroupIds = make([]string, len(values))
	copy(b.additionalSecurityGroupIds, values)
	b.fieldSet_[3] = true
	return b
}

// AvailabilityZoneTypes sets the value of the 'availability_zone_types' attribute to the given value.
func (b *AWSMachinePoolBuilder) AvailabilityZoneTypes(value map[string]string) *AWSMachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.availabilityZoneTypes = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// SpotMarketOptions sets the value of the 'spot_market_options' attribute to the given value.
//
// Spot market options for AWS machine pool.
func (b *AWSMachinePoolBuilder) SpotMarketOptions(value *AWSSpotMarketOptionsBuilder) *AWSMachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.spotMarketOptions = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// SubnetOutposts sets the value of the 'subnet_outposts' attribute to the given value.
func (b *AWSMachinePoolBuilder) SubnetOutposts(value map[string]string) *AWSMachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.subnetOutposts = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// Tags sets the value of the 'tags' attribute to the given value.
func (b *AWSMachinePoolBuilder) Tags(value map[string]string) *AWSMachinePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.tags = value
	if value != nil {
		b.fieldSet_[7] = true
	} else {
		b.fieldSet_[7] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSMachinePoolBuilder) Copy(object *AWSMachinePool) *AWSMachinePoolBuilder {
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
	if object.spotMarketOptions != nil {
		b.spotMarketOptions = NewAWSSpotMarketOptions().Copy(object.spotMarketOptions)
	} else {
		b.spotMarketOptions = nil
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

// Build creates a 'AWS_machine_pool' object using the configuration stored in the builder.
func (b *AWSMachinePoolBuilder) Build() (object *AWSMachinePool, err error) {
	object = new(AWSMachinePool)
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
	if b.spotMarketOptions != nil {
		object.spotMarketOptions, err = b.spotMarketOptions.Build()
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
