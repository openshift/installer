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

// AWSVolumeBuilder contains the data and logic needed to build 'AWS_volume' objects.
//
// Holds settings for an AWS storage volume.
type AWSVolumeBuilder struct {
	bitmap_ uint32
	iops    int
	size    int
}

// NewAWSVolume creates a new builder of 'AWS_volume' objects.
func NewAWSVolume() *AWSVolumeBuilder {
	return &AWSVolumeBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSVolumeBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// IOPS sets the value of the 'IOPS' attribute to the given value.
func (b *AWSVolumeBuilder) IOPS(value int) *AWSVolumeBuilder {
	b.iops = value
	b.bitmap_ |= 1
	return b
}

// Size sets the value of the 'size' attribute to the given value.
func (b *AWSVolumeBuilder) Size(value int) *AWSVolumeBuilder {
	b.size = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSVolumeBuilder) Copy(object *AWSVolume) *AWSVolumeBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.iops = object.iops
	b.size = object.size
	return b
}

// Build creates a 'AWS_volume' object using the configuration stored in the builder.
func (b *AWSVolumeBuilder) Build() (object *AWSVolume, err error) {
	object = new(AWSVolume)
	object.bitmap_ = b.bitmap_
	object.iops = b.iops
	object.size = b.size
	return
}
