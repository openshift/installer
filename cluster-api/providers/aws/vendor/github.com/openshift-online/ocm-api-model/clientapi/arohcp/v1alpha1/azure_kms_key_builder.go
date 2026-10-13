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

// Contains the necessary attributes to support KMS encryption key for Azure based clusters
type AzureKmsKeyBuilder struct {
	fieldSet_    []bool
	keyName      string
	keyVaultName string
	keyVersion   string
}

// NewAzureKmsKey creates a new builder of 'azure_kms_key' objects.
func NewAzureKmsKey() *AzureKmsKeyBuilder {
	return &AzureKmsKeyBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureKmsKeyBuilder) Empty() bool {
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

// KeyName sets the value of the 'key_name' attribute to the given value.
func (b *AzureKmsKeyBuilder) KeyName(value string) *AzureKmsKeyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.keyName = value
	b.fieldSet_[0] = true
	return b
}

// KeyVaultName sets the value of the 'key_vault_name' attribute to the given value.
func (b *AzureKmsKeyBuilder) KeyVaultName(value string) *AzureKmsKeyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.keyVaultName = value
	b.fieldSet_[1] = true
	return b
}

// KeyVersion sets the value of the 'key_version' attribute to the given value.
func (b *AzureKmsKeyBuilder) KeyVersion(value string) *AzureKmsKeyBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.keyVersion = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureKmsKeyBuilder) Copy(object *AzureKmsKey) *AzureKmsKeyBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.keyName = object.keyName
	b.keyVaultName = object.keyVaultName
	b.keyVersion = object.keyVersion
	return b
}

// Build creates a 'azure_kms_key' object using the configuration stored in the builder.
func (b *AzureKmsKeyBuilder) Build() (object *AzureKmsKey, err error) {
	object = new(AzureKmsKey)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.keyName = b.keyName
	object.keyVaultName = b.keyVaultName
	object.keyVersion = b.keyVersion
	return
}
