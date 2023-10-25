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

// SyncsetBuilder contains the data and logic needed to build 'syncset' objects.
//
// Representation of a syncset.
type SyncsetBuilder struct {
	bitmap_   uint32
	id        string
	href      string
	resources []interface{}
}

// NewSyncset creates a new builder of 'syncset' objects.
func NewSyncset() *SyncsetBuilder {
	return &SyncsetBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *SyncsetBuilder) Link(value bool) *SyncsetBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *SyncsetBuilder) ID(value string) *SyncsetBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *SyncsetBuilder) HREF(value string) *SyncsetBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *SyncsetBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Resources sets the value of the 'resources' attribute to the given values.
func (b *SyncsetBuilder) Resources(values ...interface{}) *SyncsetBuilder {
	b.resources = make([]interface{}, len(values))
	copy(b.resources, values)
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *SyncsetBuilder) Copy(object *Syncset) *SyncsetBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.resources != nil {
		b.resources = make([]interface{}, len(object.resources))
		copy(b.resources, object.resources)
	} else {
		b.resources = nil
	}
	return b
}

// Build creates a 'syncset' object using the configuration stored in the builder.
func (b *SyncsetBuilder) Build() (object *Syncset, err error) {
	object = new(Syncset)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.resources != nil {
		object.resources = make([]interface{}, len(b.resources))
		copy(object.resources, b.resources)
	}
	return
}
