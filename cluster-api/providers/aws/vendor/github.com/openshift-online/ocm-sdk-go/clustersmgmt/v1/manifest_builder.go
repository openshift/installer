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

import (
	time "time"
)

// ManifestBuilder contains the data and logic needed to build 'manifest' objects.
//
// Representation of a manifestwork.
type ManifestBuilder struct {
	bitmap_           uint32
	id                string
	href              string
	creationTimestamp time.Time
	liveResource      interface{}
	spec              interface{}
	updatedTimestamp  time.Time
	workloads         []interface{}
}

// NewManifest creates a new builder of 'manifest' objects.
func NewManifest() *ManifestBuilder {
	return &ManifestBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ManifestBuilder) Link(value bool) *ManifestBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ManifestBuilder) ID(value string) *ManifestBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ManifestBuilder) HREF(value string) *ManifestBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ManifestBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ManifestBuilder) CreationTimestamp(value time.Time) *ManifestBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 8
	return b
}

// LiveResource sets the value of the 'live_resource' attribute to the given value.
func (b *ManifestBuilder) LiveResource(value interface{}) *ManifestBuilder {
	b.liveResource = value
	b.bitmap_ |= 16
	return b
}

// Spec sets the value of the 'spec' attribute to the given value.
func (b *ManifestBuilder) Spec(value interface{}) *ManifestBuilder {
	b.spec = value
	b.bitmap_ |= 32
	return b
}

// UpdatedTimestamp sets the value of the 'updated_timestamp' attribute to the given value.
func (b *ManifestBuilder) UpdatedTimestamp(value time.Time) *ManifestBuilder {
	b.updatedTimestamp = value
	b.bitmap_ |= 64
	return b
}

// Workloads sets the value of the 'workloads' attribute to the given values.
func (b *ManifestBuilder) Workloads(values ...interface{}) *ManifestBuilder {
	b.workloads = make([]interface{}, len(values))
	copy(b.workloads, values)
	b.bitmap_ |= 128
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ManifestBuilder) Copy(object *Manifest) *ManifestBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.creationTimestamp = object.creationTimestamp
	b.liveResource = object.liveResource
	b.spec = object.spec
	b.updatedTimestamp = object.updatedTimestamp
	if object.workloads != nil {
		b.workloads = make([]interface{}, len(object.workloads))
		copy(b.workloads, object.workloads)
	} else {
		b.workloads = nil
	}
	return b
}

// Build creates a 'manifest' object using the configuration stored in the builder.
func (b *ManifestBuilder) Build() (object *Manifest, err error) {
	object = new(Manifest)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.creationTimestamp = b.creationTimestamp
	object.liveResource = b.liveResource
	object.spec = b.spec
	object.updatedTimestamp = b.updatedTimestamp
	if b.workloads != nil {
		object.workloads = make([]interface{}, len(b.workloads))
		copy(object.workloads, b.workloads)
	}
	return
}
