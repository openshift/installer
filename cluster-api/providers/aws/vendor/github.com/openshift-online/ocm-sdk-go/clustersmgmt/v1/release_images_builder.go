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

// ReleaseImagesBuilder contains the data and logic needed to build 'release_images' objects.
type ReleaseImagesBuilder struct {
	bitmap_ uint32
	arm64   *ReleaseImageDetailsBuilder
	multi   *ReleaseImageDetailsBuilder
}

// NewReleaseImages creates a new builder of 'release_images' objects.
func NewReleaseImages() *ReleaseImagesBuilder {
	return &ReleaseImagesBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ReleaseImagesBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ARM64 sets the value of the 'ARM64' attribute to the given value.
func (b *ReleaseImagesBuilder) ARM64(value *ReleaseImageDetailsBuilder) *ReleaseImagesBuilder {
	b.arm64 = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// Multi sets the value of the 'multi' attribute to the given value.
func (b *ReleaseImagesBuilder) Multi(value *ReleaseImageDetailsBuilder) *ReleaseImagesBuilder {
	b.multi = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ReleaseImagesBuilder) Copy(object *ReleaseImages) *ReleaseImagesBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.arm64 != nil {
		b.arm64 = NewReleaseImageDetails().Copy(object.arm64)
	} else {
		b.arm64 = nil
	}
	if object.multi != nil {
		b.multi = NewReleaseImageDetails().Copy(object.multi)
	} else {
		b.multi = nil
	}
	return b
}

// Build creates a 'release_images' object using the configuration stored in the builder.
func (b *ReleaseImagesBuilder) Build() (object *ReleaseImages, err error) {
	object = new(ReleaseImages)
	object.bitmap_ = b.bitmap_
	if b.arm64 != nil {
		object.arm64, err = b.arm64.Build()
		if err != nil {
			return
		}
	}
	if b.multi != nil {
		object.multi, err = b.multi.Build()
		if err != nil {
			return
		}
	}
	return
}
