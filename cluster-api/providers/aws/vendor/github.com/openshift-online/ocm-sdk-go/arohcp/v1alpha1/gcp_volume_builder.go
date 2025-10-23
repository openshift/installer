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

// GCPVolumeBuilder contains the data and logic needed to build 'GCP_volume' objects.
//
// Holds settings for an GCP storage volume.
type GCPVolumeBuilder struct {
	bitmap_ uint32
	size    int
}

// NewGCPVolume creates a new builder of 'GCP_volume' objects.
func NewGCPVolume() *GCPVolumeBuilder {
	return &GCPVolumeBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GCPVolumeBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Size sets the value of the 'size' attribute to the given value.
func (b *GCPVolumeBuilder) Size(value int) *GCPVolumeBuilder {
	b.size = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GCPVolumeBuilder) Copy(object *GCPVolume) *GCPVolumeBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.size = object.size
	return b
}

// Build creates a 'GCP_volume' object using the configuration stored in the builder.
func (b *GCPVolumeBuilder) Build() (object *GCPVolume, err error) {
	object = new(GCPVolume)
	object.bitmap_ = b.bitmap_
	object.size = b.size
	return
}
