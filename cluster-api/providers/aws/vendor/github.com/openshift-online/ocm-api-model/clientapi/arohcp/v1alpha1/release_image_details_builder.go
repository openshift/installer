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

type ReleaseImageDetailsBuilder struct {
	fieldSet_         []bool
	availableUpgrades []string
	releaseImage      string
}

// NewReleaseImageDetails creates a new builder of 'release_image_details' objects.
func NewReleaseImageDetails() *ReleaseImageDetailsBuilder {
	return &ReleaseImageDetailsBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ReleaseImageDetailsBuilder) Empty() bool {
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

// AvailableUpgrades sets the value of the 'available_upgrades' attribute to the given values.
func (b *ReleaseImageDetailsBuilder) AvailableUpgrades(values ...string) *ReleaseImageDetailsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.availableUpgrades = make([]string, len(values))
	copy(b.availableUpgrades, values)
	b.fieldSet_[0] = true
	return b
}

// ReleaseImage sets the value of the 'release_image' attribute to the given value.
func (b *ReleaseImageDetailsBuilder) ReleaseImage(value string) *ReleaseImageDetailsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.releaseImage = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ReleaseImageDetailsBuilder) Copy(object *ReleaseImageDetails) *ReleaseImageDetailsBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.availableUpgrades != nil {
		object.availableUpgrades = make([]string, len(b.availableUpgrades))
		copy(object.availableUpgrades, b.availableUpgrades)
	}
	object.releaseImage = b.releaseImage
	return
}
