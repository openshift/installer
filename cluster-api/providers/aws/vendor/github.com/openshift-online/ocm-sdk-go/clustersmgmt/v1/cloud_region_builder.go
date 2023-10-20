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

// CloudRegionBuilder contains the data and logic needed to build 'cloud_region' objects.
//
// Description of a region of a cloud provider.
type CloudRegionBuilder struct {
	bitmap_            uint32
	id                 string
	href               string
	kmsLocationID      string
	kmsLocationName    string
	cloudProvider      *CloudProviderBuilder
	displayName        string
	name               string
	ccsOnly            bool
	enabled            bool
	govCloud           bool
	supportsHypershift bool
	supportsMultiAZ    bool
}

// NewCloudRegion creates a new builder of 'cloud_region' objects.
func NewCloudRegion() *CloudRegionBuilder {
	return &CloudRegionBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *CloudRegionBuilder) Link(value bool) *CloudRegionBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *CloudRegionBuilder) ID(value string) *CloudRegionBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *CloudRegionBuilder) HREF(value string) *CloudRegionBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CloudRegionBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CCSOnly sets the value of the 'CCS_only' attribute to the given value.
func (b *CloudRegionBuilder) CCSOnly(value bool) *CloudRegionBuilder {
	b.ccsOnly = value
	b.bitmap_ |= 8
	return b
}

// KMSLocationID sets the value of the 'KMS_location_ID' attribute to the given value.
func (b *CloudRegionBuilder) KMSLocationID(value string) *CloudRegionBuilder {
	b.kmsLocationID = value
	b.bitmap_ |= 16
	return b
}

// KMSLocationName sets the value of the 'KMS_location_name' attribute to the given value.
func (b *CloudRegionBuilder) KMSLocationName(value string) *CloudRegionBuilder {
	b.kmsLocationName = value
	b.bitmap_ |= 32
	return b
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
//
// Cloud provider.
func (b *CloudRegionBuilder) CloudProvider(value *CloudProviderBuilder) *CloudRegionBuilder {
	b.cloudProvider = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *CloudRegionBuilder) DisplayName(value string) *CloudRegionBuilder {
	b.displayName = value
	b.bitmap_ |= 128
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *CloudRegionBuilder) Enabled(value bool) *CloudRegionBuilder {
	b.enabled = value
	b.bitmap_ |= 256
	return b
}

// GovCloud sets the value of the 'gov_cloud' attribute to the given value.
func (b *CloudRegionBuilder) GovCloud(value bool) *CloudRegionBuilder {
	b.govCloud = value
	b.bitmap_ |= 512
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *CloudRegionBuilder) Name(value string) *CloudRegionBuilder {
	b.name = value
	b.bitmap_ |= 1024
	return b
}

// SupportsHypershift sets the value of the 'supports_hypershift' attribute to the given value.
func (b *CloudRegionBuilder) SupportsHypershift(value bool) *CloudRegionBuilder {
	b.supportsHypershift = value
	b.bitmap_ |= 2048
	return b
}

// SupportsMultiAZ sets the value of the 'supports_multi_AZ' attribute to the given value.
func (b *CloudRegionBuilder) SupportsMultiAZ(value bool) *CloudRegionBuilder {
	b.supportsMultiAZ = value
	b.bitmap_ |= 4096
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CloudRegionBuilder) Copy(object *CloudRegion) *CloudRegionBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.ccsOnly = object.ccsOnly
	b.kmsLocationID = object.kmsLocationID
	b.kmsLocationName = object.kmsLocationName
	if object.cloudProvider != nil {
		b.cloudProvider = NewCloudProvider().Copy(object.cloudProvider)
	} else {
		b.cloudProvider = nil
	}
	b.displayName = object.displayName
	b.enabled = object.enabled
	b.govCloud = object.govCloud
	b.name = object.name
	b.supportsHypershift = object.supportsHypershift
	b.supportsMultiAZ = object.supportsMultiAZ
	return b
}

// Build creates a 'cloud_region' object using the configuration stored in the builder.
func (b *CloudRegionBuilder) Build() (object *CloudRegion, err error) {
	object = new(CloudRegion)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.ccsOnly = b.ccsOnly
	object.kmsLocationID = b.kmsLocationID
	object.kmsLocationName = b.kmsLocationName
	if b.cloudProvider != nil {
		object.cloudProvider, err = b.cloudProvider.Build()
		if err != nil {
			return
		}
	}
	object.displayName = b.displayName
	object.enabled = b.enabled
	object.govCloud = b.govCloud
	object.name = b.name
	object.supportsHypershift = b.supportsHypershift
	object.supportsMultiAZ = b.supportsMultiAZ
	return
}
