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

import (
	time "time"
)

// Representation of a manifestwork.
type ManifestBuilder struct {
	fieldSet_         []bool
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
	return &ManifestBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ManifestBuilder) Link(value bool) *ManifestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ManifestBuilder) ID(value string) *ManifestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ManifestBuilder) HREF(value string) *ManifestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ManifestBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ManifestBuilder) CreationTimestamp(value time.Time) *ManifestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.creationTimestamp = value
	b.fieldSet_[3] = true
	return b
}

// LiveResource sets the value of the 'live_resource' attribute to the given value.
func (b *ManifestBuilder) LiveResource(value interface{}) *ManifestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.liveResource = value
	b.fieldSet_[4] = true
	return b
}

// Spec sets the value of the 'spec' attribute to the given value.
func (b *ManifestBuilder) Spec(value interface{}) *ManifestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.spec = value
	b.fieldSet_[5] = true
	return b
}

// UpdatedTimestamp sets the value of the 'updated_timestamp' attribute to the given value.
func (b *ManifestBuilder) UpdatedTimestamp(value time.Time) *ManifestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.updatedTimestamp = value
	b.fieldSet_[6] = true
	return b
}

// Workloads sets the value of the 'workloads' attribute to the given values.
func (b *ManifestBuilder) Workloads(values ...interface{}) *ManifestBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.workloads = make([]interface{}, len(values))
	copy(b.workloads, values)
	b.fieldSet_[7] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ManifestBuilder) Copy(object *Manifest) *ManifestBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
