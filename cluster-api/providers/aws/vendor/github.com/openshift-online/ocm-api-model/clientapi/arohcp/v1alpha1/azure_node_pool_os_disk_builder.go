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

// Defines the configuration of a Node Pool's OS disk.
type AzureNodePoolOsDiskBuilder struct {
	fieldSet_                  []bool
	persistence                string
	sizeGibibytes              int
	sseEncryptionSetResourceId string
	storageAccountType         string
}

// NewAzureNodePoolOsDisk creates a new builder of 'azure_node_pool_os_disk' objects.
func NewAzureNodePoolOsDisk() *AzureNodePoolOsDiskBuilder {
	return &AzureNodePoolOsDiskBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureNodePoolOsDiskBuilder) Empty() bool {
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

// Persistence sets the value of the 'persistence' attribute to the given value.
func (b *AzureNodePoolOsDiskBuilder) Persistence(value string) *AzureNodePoolOsDiskBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.persistence = value
	b.fieldSet_[0] = true
	return b
}

// SizeGibibytes sets the value of the 'size_gibibytes' attribute to the given value.
func (b *AzureNodePoolOsDiskBuilder) SizeGibibytes(value int) *AzureNodePoolOsDiskBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.sizeGibibytes = value
	b.fieldSet_[1] = true
	return b
}

// SseEncryptionSetResourceId sets the value of the 'sse_encryption_set_resource_id' attribute to the given value.
func (b *AzureNodePoolOsDiskBuilder) SseEncryptionSetResourceId(value string) *AzureNodePoolOsDiskBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.sseEncryptionSetResourceId = value
	b.fieldSet_[2] = true
	return b
}

// StorageAccountType sets the value of the 'storage_account_type' attribute to the given value.
func (b *AzureNodePoolOsDiskBuilder) StorageAccountType(value string) *AzureNodePoolOsDiskBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.storageAccountType = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureNodePoolOsDiskBuilder) Copy(object *AzureNodePoolOsDisk) *AzureNodePoolOsDiskBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.persistence = object.persistence
	b.sizeGibibytes = object.sizeGibibytes
	b.sseEncryptionSetResourceId = object.sseEncryptionSetResourceId
	b.storageAccountType = object.storageAccountType
	return b
}

// Build creates a 'azure_node_pool_os_disk' object using the configuration stored in the builder.
func (b *AzureNodePoolOsDiskBuilder) Build() (object *AzureNodePoolOsDisk, err error) {
	object = new(AzureNodePoolOsDisk)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.persistence = b.persistence
	object.sizeGibibytes = b.sizeGibibytes
	object.sseEncryptionSetResourceId = b.sseEncryptionSetResourceId
	object.storageAccountType = b.storageAccountType
	return
}
