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

// Holds settings for an AWS storage volume.
type AWSVolumeBuilder struct {
	fieldSet_ []bool
	iops      int
	size      int
}

// NewAWSVolume creates a new builder of 'AWS_volume' objects.
func NewAWSVolume() *AWSVolumeBuilder {
	return &AWSVolumeBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AWSVolumeBuilder) Empty() bool {
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

// IOPS sets the value of the 'IOPS' attribute to the given value.
func (b *AWSVolumeBuilder) IOPS(value int) *AWSVolumeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.iops = value
	b.fieldSet_[0] = true
	return b
}

// Size sets the value of the 'size' attribute to the given value.
func (b *AWSVolumeBuilder) Size(value int) *AWSVolumeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.size = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AWSVolumeBuilder) Copy(object *AWSVolume) *AWSVolumeBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.iops = object.iops
	b.size = object.size
	return b
}

// Build creates a 'AWS_volume' object using the configuration stored in the builder.
func (b *AWSVolumeBuilder) Build() (object *AWSVolume, err error) {
	object = new(AWSVolume)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.iops = b.iops
	object.size = b.size
	return
}
