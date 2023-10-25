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

// RootVolumeBuilder contains the data and logic needed to build 'root_volume' objects.
//
// Root volume capabilities.
type RootVolumeBuilder struct {
	bitmap_ uint32
	aws     *AWSVolumeBuilder
	gcp     *GCPVolumeBuilder
}

// NewRootVolume creates a new builder of 'root_volume' objects.
func NewRootVolume() *RootVolumeBuilder {
	return &RootVolumeBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RootVolumeBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AWS sets the value of the 'AWS' attribute to the given value.
//
// Holds settings for an AWS storage volume.
func (b *RootVolumeBuilder) AWS(value *AWSVolumeBuilder) *RootVolumeBuilder {
	b.aws = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// GCP sets the value of the 'GCP' attribute to the given value.
//
// Holds settings for an GCP storage volume.
func (b *RootVolumeBuilder) GCP(value *GCPVolumeBuilder) *RootVolumeBuilder {
	b.gcp = value
	if value != nil {
		b.bitmap_ |= 2
	} else {
		b.bitmap_ &^= 2
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RootVolumeBuilder) Copy(object *RootVolume) *RootVolumeBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
