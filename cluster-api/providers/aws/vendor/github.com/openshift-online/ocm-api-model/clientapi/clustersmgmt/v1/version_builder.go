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

import (
	time "time"
)

// Representation of an _OpenShift_ version.
type VersionBuilder struct {
	fieldSet_                 []bool
	id                        string
	href                      string
	availableUpgrades         []string
	channelGroup              string
	endOfLifeTimestamp        time.Time
	imageOverrides            *ImageOverridesBuilder
	rawID                     string
	releaseImage              string
	releaseImages             *ReleaseImagesBuilder
	gcpMarketplaceEnabled     bool
	rosaEnabled               bool
	default_                  bool
	enabled                   bool
	hostedControlPlaneDefault bool
	hostedControlPlaneEnabled bool
	wifEnabled                bool
}

// NewVersion creates a new builder of 'version' objects.
func NewVersion() *VersionBuilder {
	return &VersionBuilder{
		fieldSet_: make([]bool, 17),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *VersionBuilder) Link(value bool) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *VersionBuilder) ID(value string) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *VersionBuilder) HREF(value string) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *VersionBuilder) Empty() bool {
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

// GCPMarketplaceEnabled sets the value of the 'GCP_marketplace_enabled' attribute to the given value.
func (b *VersionBuilder) GCPMarketplaceEnabled(value bool) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.gcpMarketplaceEnabled = value
	b.fieldSet_[3] = true
	return b
}

// ROSAEnabled sets the value of the 'ROSA_enabled' attribute to the given value.
func (b *VersionBuilder) ROSAEnabled(value bool) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.rosaEnabled = value
	b.fieldSet_[4] = true
	return b
}

// AvailableUpgrades sets the value of the 'available_upgrades' attribute to the given values.
func (b *VersionBuilder) AvailableUpgrades(values ...string) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.availableUpgrades = make([]string, len(values))
	copy(b.availableUpgrades, values)
	b.fieldSet_[5] = true
	return b
}

// ChannelGroup sets the value of the 'channel_group' attribute to the given value.
func (b *VersionBuilder) ChannelGroup(value string) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.channelGroup = value
	b.fieldSet_[6] = true
	return b
}

// Default sets the value of the 'default' attribute to the given value.
func (b *VersionBuilder) Default(value bool) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.default_ = value
	b.fieldSet_[7] = true
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *VersionBuilder) Enabled(value bool) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.enabled = value
	b.fieldSet_[8] = true
	return b
}

// EndOfLifeTimestamp sets the value of the 'end_of_life_timestamp' attribute to the given value.
func (b *VersionBuilder) EndOfLifeTimestamp(value time.Time) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.endOfLifeTimestamp = value
	b.fieldSet_[9] = true
	return b
}

// HostedControlPlaneDefault sets the value of the 'hosted_control_plane_default' attribute to the given value.
func (b *VersionBuilder) HostedControlPlaneDefault(value bool) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.hostedControlPlaneDefault = value
	b.fieldSet_[10] = true
	return b
}

// HostedControlPlaneEnabled sets the value of the 'hosted_control_plane_enabled' attribute to the given value.
func (b *VersionBuilder) HostedControlPlaneEnabled(value bool) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.hostedControlPlaneEnabled = value
	b.fieldSet_[11] = true
	return b
}

// ImageOverrides sets the value of the 'image_overrides' attribute to the given value.
//
// ImageOverrides holds the lists of available images per cloud provider.
func (b *VersionBuilder) ImageOverrides(value *ImageOverridesBuilder) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.imageOverrides = value
	if value != nil {
		b.fieldSet_[12] = true
	} else {
		b.fieldSet_[12] = false
	}
	return b
}

// RawID sets the value of the 'raw_ID' attribute to the given value.
func (b *VersionBuilder) RawID(value string) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.rawID = value
	b.fieldSet_[13] = true
	return b
}

// ReleaseImage sets the value of the 'release_image' attribute to the given value.
func (b *VersionBuilder) ReleaseImage(value string) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.releaseImage = value
	b.fieldSet_[14] = true
	return b
}

// ReleaseImages sets the value of the 'release_images' attribute to the given value.
func (b *VersionBuilder) ReleaseImages(value *ReleaseImagesBuilder) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.releaseImages = value
	if value != nil {
		b.fieldSet_[15] = true
	} else {
		b.fieldSet_[15] = false
	}
	return b
}

// WifEnabled sets the value of the 'wif_enabled' attribute to the given value.
func (b *VersionBuilder) WifEnabled(value bool) *VersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 17)
	}
	b.wifEnabled = value
	b.fieldSet_[16] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *VersionBuilder) Copy(object *Version) *VersionBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.gcpMarketplaceEnabled = object.gcpMarketplaceEnabled
	b.rosaEnabled = object.rosaEnabled
	if object.availableUpgrades != nil {
		b.availableUpgrades = make([]string, len(object.availableUpgrades))
		copy(b.availableUpgrades, object.availableUpgrades)
	} else {
		b.availableUpgrades = nil
	}
	b.channelGroup = object.channelGroup
	b.default_ = object.default_
	b.enabled = object.enabled
	b.endOfLifeTimestamp = object.endOfLifeTimestamp
	b.hostedControlPlaneDefault = object.hostedControlPlaneDefault
	b.hostedControlPlaneEnabled = object.hostedControlPlaneEnabled
	if object.imageOverrides != nil {
		b.imageOverrides = NewImageOverrides().Copy(object.imageOverrides)
	} else {
		b.imageOverrides = nil
	}
	b.rawID = object.rawID
	b.releaseImage = object.releaseImage
	if object.releaseImages != nil {
		b.releaseImages = NewReleaseImages().Copy(object.releaseImages)
	} else {
		b.releaseImages = nil
	}
	b.wifEnabled = object.wifEnabled
	return b
}

// Build creates a 'version' object using the configuration stored in the builder.
func (b *VersionBuilder) Build() (object *Version, err error) {
	object = new(Version)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.gcpMarketplaceEnabled = b.gcpMarketplaceEnabled
	object.rosaEnabled = b.rosaEnabled
	if b.availableUpgrades != nil {
		object.availableUpgrades = make([]string, len(b.availableUpgrades))
		copy(object.availableUpgrades, b.availableUpgrades)
	}
	object.channelGroup = b.channelGroup
	object.default_ = b.default_
	object.enabled = b.enabled
	object.endOfLifeTimestamp = b.endOfLifeTimestamp
	object.hostedControlPlaneDefault = b.hostedControlPlaneDefault
	object.hostedControlPlaneEnabled = b.hostedControlPlaneEnabled
	if b.imageOverrides != nil {
		object.imageOverrides, err = b.imageOverrides.Build()
		if err != nil {
			return
		}
	}
	object.rawID = b.rawID
	object.releaseImage = b.releaseImage
	if b.releaseImages != nil {
		object.releaseImages, err = b.releaseImages.Build()
		if err != nil {
			return
		}
	}
	object.wifEnabled = b.wifEnabled
	return
}
