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

// Description of a region of a cloud provider.
type CloudRegionBuilder struct {
	fieldSet_          []bool
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
	return &CloudRegionBuilder{
		fieldSet_: make([]bool, 13),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *CloudRegionBuilder) Link(value bool) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *CloudRegionBuilder) ID(value string) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *CloudRegionBuilder) HREF(value string) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CloudRegionBuilder) Empty() bool {
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

// CCSOnly sets the value of the 'CCS_only' attribute to the given value.
func (b *CloudRegionBuilder) CCSOnly(value bool) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.ccsOnly = value
	b.fieldSet_[3] = true
	return b
}

// KMSLocationID sets the value of the 'KMS_location_ID' attribute to the given value.
func (b *CloudRegionBuilder) KMSLocationID(value string) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.kmsLocationID = value
	b.fieldSet_[4] = true
	return b
}

// KMSLocationName sets the value of the 'KMS_location_name' attribute to the given value.
func (b *CloudRegionBuilder) KMSLocationName(value string) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.kmsLocationName = value
	b.fieldSet_[5] = true
	return b
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
//
// Cloud provider.
func (b *CloudRegionBuilder) CloudProvider(value *CloudProviderBuilder) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.cloudProvider = value
	if value != nil {
		b.fieldSet_[6] = true
	} else {
		b.fieldSet_[6] = false
	}
	return b
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *CloudRegionBuilder) DisplayName(value string) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.displayName = value
	b.fieldSet_[7] = true
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *CloudRegionBuilder) Enabled(value bool) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.enabled = value
	b.fieldSet_[8] = true
	return b
}

// GovCloud sets the value of the 'gov_cloud' attribute to the given value.
func (b *CloudRegionBuilder) GovCloud(value bool) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.govCloud = value
	b.fieldSet_[9] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *CloudRegionBuilder) Name(value string) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.name = value
	b.fieldSet_[10] = true
	return b
}

// SupportsHypershift sets the value of the 'supports_hypershift' attribute to the given value.
func (b *CloudRegionBuilder) SupportsHypershift(value bool) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.supportsHypershift = value
	b.fieldSet_[11] = true
	return b
}

// SupportsMultiAZ sets the value of the 'supports_multi_AZ' attribute to the given value.
func (b *CloudRegionBuilder) SupportsMultiAZ(value bool) *CloudRegionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 13)
	}
	b.supportsMultiAZ = value
	b.fieldSet_[12] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CloudRegionBuilder) Copy(object *CloudRegion) *CloudRegionBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
