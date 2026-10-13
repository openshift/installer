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

// Root volume capabilities.
type RootVolumeBuilder struct {
	fieldSet_ []bool
	aws       *AWSVolumeBuilder
	gcp       *GCPVolumeBuilder
}

// NewRootVolume creates a new builder of 'root_volume' objects.
func NewRootVolume() *RootVolumeBuilder {
	return &RootVolumeBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RootVolumeBuilder) Empty() bool {
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

// AWS sets the value of the 'AWS' attribute to the given value.
//
// Holds settings for an AWS storage volume.
func (b *RootVolumeBuilder) AWS(value *AWSVolumeBuilder) *RootVolumeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.aws = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// GCP sets the value of the 'GCP' attribute to the given value.
//
// Holds settings for an GCP storage volume.
func (b *RootVolumeBuilder) GCP(value *GCPVolumeBuilder) *RootVolumeBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.gcp = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RootVolumeBuilder) Copy(object *RootVolume) *RootVolumeBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.aws != nil {
		b.aws = NewAWSVolume().Copy(object.aws)
	} else {
		b.aws = nil
	}
	if object.gcp != nil {
		b.gcp = NewGCPVolume().Copy(object.gcp)
	} else {
		b.gcp = nil
	}
	return b
}

// Build creates a 'root_volume' object using the configuration stored in the builder.
func (b *RootVolumeBuilder) Build() (object *RootVolume, err error) {
	object = new(RootVolume)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
	return
}
