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

// Contains the necessary attributes to support etcd data encryption with customer managed keys
// for Azure based clusters.
type AzureEtcdDataEncryptionCustomerManagedBuilder struct {
	fieldSet_      []bool
	encryptionType string
	kms            *AzureKmsEncryptionBuilder
}

// NewAzureEtcdDataEncryptionCustomerManaged creates a new builder of 'azure_etcd_data_encryption_customer_managed' objects.
func NewAzureEtcdDataEncryptionCustomerManaged() *AzureEtcdDataEncryptionCustomerManagedBuilder {
	return &AzureEtcdDataEncryptionCustomerManagedBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureEtcdDataEncryptionCustomerManagedBuilder) Empty() bool {
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

// EncryptionType sets the value of the 'encryption_type' attribute to the given value.
func (b *AzureEtcdDataEncryptionCustomerManagedBuilder) EncryptionType(value string) *AzureEtcdDataEncryptionCustomerManagedBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.encryptionType = value
	b.fieldSet_[0] = true
	return b
}

// Kms sets the value of the 'kms' attribute to the given value.
//
// Contains the necessary attributes to support KMS encryption for Azure based clusters.
func (b *AzureEtcdDataEncryptionCustomerManagedBuilder) Kms(value *AzureKmsEncryptionBuilder) *AzureEtcdDataEncryptionCustomerManagedBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.kms = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureEtcdDataEncryptionCustomerManagedBuilder) Copy(object *AzureEtcdDataEncryptionCustomerManaged) *AzureEtcdDataEncryptionCustomerManagedBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.encryptionType = object.encryptionType
	if object.kms != nil {
		b.kms = NewAzureKmsEncryption().Copy(object.kms)
	} else {
		b.kms = nil
	}
	return b
}

// Build creates a 'azure_etcd_data_encryption_customer_managed' object using the configuration stored in the builder.
func (b *AzureEtcdDataEncryptionCustomerManagedBuilder) Build() (object *AzureEtcdDataEncryptionCustomerManaged, err error) {
	object = new(AzureEtcdDataEncryptionCustomerManaged)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.encryptionType = b.encryptionType
	if b.kms != nil {
		object.kms, err = b.kms.Build()
		if err != nil {
			return
		}
	}
	return
}
