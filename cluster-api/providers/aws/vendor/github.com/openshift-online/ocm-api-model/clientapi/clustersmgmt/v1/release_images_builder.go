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

type ReleaseImagesBuilder struct {
	fieldSet_ []bool
	arm64     *ReleaseImageDetailsBuilder
	multi     *ReleaseImageDetailsBuilder
}

// NewReleaseImages creates a new builder of 'release_images' objects.
func NewReleaseImages() *ReleaseImagesBuilder {
	return &ReleaseImagesBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ReleaseImagesBuilder) Empty() bool {
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

// ARM64 sets the value of the 'ARM64' attribute to the given value.
func (b *ReleaseImagesBuilder) ARM64(value *ReleaseImageDetailsBuilder) *ReleaseImagesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.arm64 = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// Multi sets the value of the 'multi' attribute to the given value.
func (b *ReleaseImagesBuilder) Multi(value *ReleaseImageDetailsBuilder) *ReleaseImagesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.multi = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ReleaseImagesBuilder) Copy(object *ReleaseImages) *ReleaseImagesBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
