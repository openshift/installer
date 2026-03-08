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

// ImageMirror represents a container image mirror configuration for a cluster.
// This enables Day 2 image mirroring configuration for ROSA HCP clusters using
// HyperShift's native imageContentSources mechanism.
type ImageMirrorBuilder struct {
	fieldSet_           []bool
	id                  string
	href                string
	creationTimestamp   time.Time
	lastUpdateTimestamp time.Time
	mirrors             []string
	source              string
	type_               string
}

// NewImageMirror creates a new builder of 'image_mirror' objects.
func NewImageMirror() *ImageMirrorBuilder {
	return &ImageMirrorBuilder{
		fieldSet_: make([]bool, 8),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *ImageMirrorBuilder) Link(value bool) *ImageMirrorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *ImageMirrorBuilder) ID(value string) *ImageMirrorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *ImageMirrorBuilder) HREF(value string) *ImageMirrorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ImageMirrorBuilder) Empty() bool {
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
func (b *ImageMirrorBuilder) CreationTimestamp(value time.Time) *ImageMirrorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.creationTimestamp = value
	b.fieldSet_[3] = true
	return b
}

// LastUpdateTimestamp sets the value of the 'last_update_timestamp' attribute to the given value.
func (b *ImageMirrorBuilder) LastUpdateTimestamp(value time.Time) *ImageMirrorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.lastUpdateTimestamp = value
	b.fieldSet_[4] = true
	return b
}

// Mirrors sets the value of the 'mirrors' attribute to the given values.
func (b *ImageMirrorBuilder) Mirrors(values ...string) *ImageMirrorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.mirrors = make([]string, len(values))
	copy(b.mirrors, values)
	b.fieldSet_[5] = true
	return b
}

// Source sets the value of the 'source' attribute to the given value.
func (b *ImageMirrorBuilder) Source(value string) *ImageMirrorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.source = value
	b.fieldSet_[6] = true
	return b
}

// Type sets the value of the 'type' attribute to the given value.
func (b *ImageMirrorBuilder) Type(value string) *ImageMirrorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 8)
	}
	b.type_ = value
	b.fieldSet_[7] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ImageMirrorBuilder) Copy(object *ImageMirror) *ImageMirrorBuilder {
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
	b.lastUpdateTimestamp = object.lastUpdateTimestamp
	if object.mirrors != nil {
		b.mirrors = make([]string, len(object.mirrors))
		copy(b.mirrors, object.mirrors)
	} else {
		b.mirrors = nil
	}
	b.source = object.source
	b.type_ = object.type_
	return b
}

// Build creates a 'image_mirror' object using the configuration stored in the builder.
func (b *ImageMirrorBuilder) Build() (object *ImageMirror, err error) {
	object = new(ImageMirror)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.creationTimestamp = b.creationTimestamp
	object.lastUpdateTimestamp = b.lastUpdateTimestamp
	if b.mirrors != nil {
		object.mirrors = make([]string, len(b.mirrors))
		copy(object.mirrors, b.mirrors)
	}
	object.source = b.source
	object.type_ = b.type_
	return
}
