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

package v1 // github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1

// ProvisionShardReferenceBuilder contains the data and logic needed to build 'provision_shard_reference' objects.
//
// Provision Shard Reference of the cluster.
type ProvisionShardReferenceBuilder struct {
	bitmap_ uint32
	href    string
	id      string
}

// NewProvisionShardReference creates a new builder of 'provision_shard_reference' objects.
func NewProvisionShardReference() *ProvisionShardReferenceBuilder {
	return &ProvisionShardReferenceBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ProvisionShardReferenceBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Href sets the value of the 'href' attribute to the given value.
func (b *ProvisionShardReferenceBuilder) Href(value string) *ProvisionShardReferenceBuilder {
	b.href = value
	b.bitmap_ |= 1
	return b
}

// Id sets the value of the 'id' attribute to the given value.
func (b *ProvisionShardReferenceBuilder) Id(value string) *ProvisionShardReferenceBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ProvisionShardReferenceBuilder) Copy(object *ProvisionShardReference) *ProvisionShardReferenceBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.href = object.href
	b.id = object.id
	return b
}

// Build creates a 'provision_shard_reference' object using the configuration stored in the builder.
func (b *ProvisionShardReferenceBuilder) Build() (object *ProvisionShardReference, err error) {
	object = new(ProvisionShardReference)
	object.bitmap_ = b.bitmap_
	object.href = b.href
	object.id = b.id
	return
}
