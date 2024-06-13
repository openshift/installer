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

// StorageQuotaBuilder contains the data and logic needed to build 'storage_quota' objects.
//
// Representation of a storage quota
type StorageQuotaBuilder struct {
	bitmap_ uint32
	unit    string
	value   float64
}

// NewStorageQuota creates a new builder of 'storage_quota' objects.
func NewStorageQuota() *StorageQuotaBuilder {
	return &StorageQuotaBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *StorageQuotaBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Unit sets the value of the 'unit' attribute to the given value.
func (b *StorageQuotaBuilder) Unit(value string) *StorageQuotaBuilder {
	b.unit = value
	b.bitmap_ |= 1
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *StorageQuotaBuilder) Value(value float64) *StorageQuotaBuilder {
	b.value = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *StorageQuotaBuilder) Copy(object *StorageQuota) *StorageQuotaBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.unit = object.unit
	b.value = object.value
	return b
}

// Build creates a 'storage_quota' object using the configuration stored in the builder.
func (b *StorageQuotaBuilder) Build() (object *StorageQuota, err error) {
	object = new(StorageQuota)
	object.bitmap_ = b.bitmap_
	object.unit = b.unit
	object.value = b.value
	return
}
