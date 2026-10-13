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

// Description of a cloud provider data used for cloud provider inquiries.
type CloudProviderDataBuilder struct {
	fieldSet_         []bool
	aws               *AWSBuilder
	gcp               *GCPBuilder
	availabilityZones []string
	keyLocation       string
	keyRingName       string
	region            *CloudRegionBuilder
	subnets           []string
	version           *VersionBuilder
	vpcIds            []string
}

// NewCloudProviderData creates a new builder of 'cloud_provider_data' objects.
func NewCloudProviderData() *CloudProviderDataBuilder {
	return &CloudProviderDataBuilder{
		fieldSet_: make([]bool, 9),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CloudProviderDataBuilder) Empty() bool {
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

// AWS sets the value of the 'AWS' attribute to the given value.
//
// _Amazon Web Services_ specific settings of a cluster.
func (b *CloudProviderDataBuilder) AWS(value *AWSBuilder) *CloudProviderDataBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.aws = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// GCP sets the value of the 'GCP' attribute to the given value.
//
// Google cloud platform settings of a cluster.
func (b *CloudProviderDataBuilder) GCP(value *GCPBuilder) *CloudProviderDataBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.gcp = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// AvailabilityZones sets the value of the 'availability_zones' attribute to the given values.
func (b *CloudProviderDataBuilder) AvailabilityZones(values ...string) *CloudProviderDataBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.availabilityZones = make([]string, len(values))
	copy(b.availabilityZones, values)
	b.fieldSet_[2] = true
	return b
}

// KeyLocation sets the value of the 'key_location' attribute to the given value.
func (b *CloudProviderDataBuilder) KeyLocation(value string) *CloudProviderDataBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.keyLocation = value
	b.fieldSet_[3] = true
	return b
}

// KeyRingName sets the value of the 'key_ring_name' attribute to the given value.
func (b *CloudProviderDataBuilder) KeyRingName(value string) *CloudProviderDataBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.keyRingName = value
	b.fieldSet_[4] = true
	return b
}

// Region sets the value of the 'region' attribute to the given value.
//
// Description of a region of a cloud provider.
func (b *CloudProviderDataBuilder) Region(value *CloudRegionBuilder) *CloudProviderDataBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.region = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// Subnets sets the value of the 'subnets' attribute to the given values.
func (b *CloudProviderDataBuilder) Subnets(values ...string) *CloudProviderDataBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.subnets = make([]string, len(values))
	copy(b.subnets, values)
	b.fieldSet_[6] = true
	return b
}

// Version sets the value of the 'version' attribute to the given value.
//
// Representation of an _OpenShift_ version.
func (b *CloudProviderDataBuilder) Version(value *VersionBuilder) *CloudProviderDataBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.version = value
	if value != nil {
		b.fieldSet_[7] = true
	} else {
		b.fieldSet_[7] = false
	}
	return b
}

// VpcIds sets the value of the 'vpc_ids' attribute to the given values.
func (b *CloudProviderDataBuilder) VpcIds(values ...string) *CloudProviderDataBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 9)
	}
	b.vpcIds = make([]string, len(values))
	copy(b.vpcIds, values)
	b.fieldSet_[8] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CloudProviderDataBuilder) Copy(object *CloudProviderData) *CloudProviderDataBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.aws != nil {
		b.aws = NewAWS().Copy(object.aws)
	} else {
		b.aws = nil
	}
	if object.gcp != nil {
		b.gcp = NewGCP().Copy(object.gcp)
	} else {
		b.gcp = nil
	}
	if object.availabilityZones != nil {
		b.availabilityZones = make([]string, len(object.availabilityZones))
		copy(b.availabilityZones, object.availabilityZones)
	} else {
		b.availabilityZones = nil
	}
	b.keyLocation = object.keyLocation
	b.keyRingName = object.keyRingName
	if object.region != nil {
		b.region = NewCloudRegion().Copy(object.region)
	} else {
		b.region = nil
	}
	if object.subnets != nil {
		b.subnets = make([]string, len(object.subnets))
		copy(b.subnets, object.subnets)
	} else {
		b.subnets = nil
	}
	if object.version != nil {
		b.version = NewVersion().Copy(object.version)
	} else {
		b.version = nil
	}
	if object.vpcIds != nil {
		b.vpcIds = make([]string, len(object.vpcIds))
		copy(b.vpcIds, object.vpcIds)
	} else {
		b.vpcIds = nil
	}
	return b
}

// Build creates a 'cloud_provider_data' object using the configuration stored in the builder.
func (b *CloudProviderDataBuilder) Build() (object *CloudProviderData, err error) {
	object = new(CloudProviderData)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.aws != nil {
		object.aws, err = b.aws.Build()
		if err != nil {
			return
		}
	}
	if b.gcp != nil {
		object.gcp, err = b.gcp.Build()
		if err != nil {
			return
		}
	}
	if b.availabilityZones != nil {
		object.availabilityZones = make([]string, len(b.availabilityZones))
		copy(object.availabilityZones, b.availabilityZones)
	}
	object.keyLocation = b.keyLocation
	object.keyRingName = b.keyRingName
	if b.region != nil {
		object.region, err = b.region.Build()
		if err != nil {
			return
		}
	}
	if b.subnets != nil {
		object.subnets = make([]string, len(b.subnets))
		copy(object.subnets, b.subnets)
	}
	if b.version != nil {
		object.version, err = b.version.Build()
		if err != nil {
			return
		}
	}
	if b.vpcIds != nil {
		object.vpcIds = make([]string, len(b.vpcIds))
		copy(object.vpcIds, b.vpcIds)
	}
	return
}
