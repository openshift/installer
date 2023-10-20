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

// CloudProviderDataBuilder contains the data and logic needed to build 'cloud_provider_data' objects.
//
// Description of a cloud provider data used for cloud provider inquiries.
type CloudProviderDataBuilder struct {
	bitmap_           uint32
	aws               *AWSBuilder
	gcp               *GCPBuilder
	availabilityZones []string
	keyLocation       string
	keyRingName       string
	region            *CloudRegionBuilder
	subnets           []string
	version           *VersionBuilder
}

// NewCloudProviderData creates a new builder of 'cloud_provider_data' objects.
func NewCloudProviderData() *CloudProviderDataBuilder {
	return &CloudProviderDataBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CloudProviderDataBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AWS sets the value of the 'AWS' attribute to the given value.
//
// _Amazon Web Services_ specific settings of a cluster.
func (b *CloudProviderDataBuilder) AWS(value *AWSBuilder) *CloudProviderDataBuilder {
	b.aws = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// GCP sets the value of the 'GCP' attribute to the given value.
//
// Google cloud platform settings of a cluster.
func (b *CloudProviderDataBuilder) GCP(value *GCPBuilder) *CloudProviderDataBuilder {
	b.gcp = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// AvailabilityZones sets the value of the 'availability_zones' attribute to the given values.
func (b *CloudProviderDataBuilder) AvailabilityZones(values ...string) *CloudProviderDataBuilder {
	b.availabilityZones = make([]string, len(values))
	copy(b.availabilityZones, values)
	b.bitmap_ |= 4
	return b
}

// KeyLocation sets the value of the 'key_location' attribute to the given value.
func (b *CloudProviderDataBuilder) KeyLocation(value string) *CloudProviderDataBuilder {
	b.keyLocation = value
	b.bitmap_ |= 8
	return b
}

// KeyRingName sets the value of the 'key_ring_name' attribute to the given value.
func (b *CloudProviderDataBuilder) KeyRingName(value string) *CloudProviderDataBuilder {
	b.keyRingName = value
	b.bitmap_ |= 16
	return b
}

// Region sets the value of the 'region' attribute to the given value.
//
// Description of a region of a cloud provider.
func (b *CloudProviderDataBuilder) Region(value *CloudRegionBuilder) *CloudProviderDataBuilder {
	b.region = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// Subnets sets the value of the 'subnets' attribute to the given values.
func (b *CloudProviderDataBuilder) Subnets(values ...string) *CloudProviderDataBuilder {
	b.subnets = make([]string, len(values))
	copy(b.subnets, values)
	b.bitmap_ |= 64
	return b
}

// Version sets the value of the 'version' attribute to the given value.
//
// Representation of an _OpenShift_ version.
func (b *CloudProviderDataBuilder) Version(value *VersionBuilder) *CloudProviderDataBuilder {
	b.version = value
	if value != nil {
		b.bitmap_ |= 128
	} else {
		b.bitmap_ &^= 128
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CloudProviderDataBuilder) Copy(object *CloudProviderData) *CloudProviderDataBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	return b
}

// Build creates a 'cloud_provider_data' object using the configuration stored in the builder.
func (b *CloudProviderDataBuilder) Build() (object *CloudProviderData, err error) {
	object = new(CloudProviderData)
	object.bitmap_ = b.bitmap_
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
	return
}
