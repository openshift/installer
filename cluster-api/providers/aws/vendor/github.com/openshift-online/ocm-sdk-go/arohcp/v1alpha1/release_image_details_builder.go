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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// ReleaseImageDetailsBuilder contains the data and logic needed to build 'release_image_details' objects.
type ReleaseImageDetailsBuilder struct {
	bitmap_           uint32
	availableUpgrades []string
	releaseImage      string
}

// NewReleaseImageDetails creates a new builder of 'release_image_details' objects.
func NewReleaseImageDetails() *ReleaseImageDetailsBuilder {
	return &ReleaseImageDetailsBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ReleaseImageDetailsBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AvailableUpgrades sets the value of the 'available_upgrades' attribute to the given values.
func (b *ReleaseImageDetailsBuilder) AvailableUpgrades(values ...string) *ReleaseImageDetailsBuilder {
	b.availableUpgrades = make([]string, len(values))
	copy(b.availableUpgrades, values)
	b.bitmap_ |= 1
	return b
}

// ReleaseImage sets the value of the 'release_image' attribute to the given value.
func (b *ReleaseImageDetailsBuilder) ReleaseImage(value string) *ReleaseImageDetailsBuilder {
	b.releaseImage = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ReleaseImageDetailsBuilder) Copy(object *ReleaseImageDetails) *ReleaseImageDetailsBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.availableUpgrades != nil {
		b.availableUpgrades = make([]string, len(object.availableUpgrades))
		copy(b.availableUpgrades, object.availableUpgrades)
	} else {
		b.availableUpgrades = nil
	}
	b.releaseImage = object.releaseImage
	return b
}

// Build creates a 'release_image_details' object using the configuration stored in the builder.
func (b *ReleaseImageDetailsBuilder) Build() (object *ReleaseImageDetails, err error) {
	object = new(ReleaseImageDetails)
	object.bitmap_ = b.bitmap_
	if b.availableUpgrades != nil {
		object.availableUpgrades = make([]string, len(b.availableUpgrades))
		copy(object.availableUpgrades, b.availableUpgrades)
	}
	object.releaseImage = b.releaseImage
	return
}
