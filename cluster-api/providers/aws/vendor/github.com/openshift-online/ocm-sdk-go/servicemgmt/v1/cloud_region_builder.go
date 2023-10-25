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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

// CloudRegionBuilder contains the data and logic needed to build 'cloud_region' objects.
//
// Description of a region of a cloud provider.
type CloudRegionBuilder struct {
	bitmap_ uint32
	id      string
}

// NewCloudRegion creates a new builder of 'cloud_region' objects.
func NewCloudRegion() *CloudRegionBuilder {
	return &CloudRegionBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *CloudRegionBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *CloudRegionBuilder) ID(value string) *CloudRegionBuilder {
	b.id = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *CloudRegionBuilder) Copy(object *CloudRegion) *CloudRegionBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	return b
}

// Build creates a 'cloud_region' object using the configuration stored in the builder.
func (b *CloudRegionBuilder) Build() (object *CloudRegion, err error) {
	object = new(CloudRegion)
	object.bitmap_ = b.bitmap_
	object.id = b.id
	return
}
