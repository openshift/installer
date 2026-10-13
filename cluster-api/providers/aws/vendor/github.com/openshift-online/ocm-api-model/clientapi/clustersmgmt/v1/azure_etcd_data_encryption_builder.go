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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// Contains the necessary attributes to support data encryption for Azure based clusters.
type AzureEtcdDataEncryptionBuilder struct {
	fieldSet_         []bool
	customerManaged   *AzureEtcdDataEncryptionCustomerManagedBuilder
	keyManagementMode string
}

// NewAzureEtcdDataEncryption creates a new builder of 'azure_etcd_data_encryption' objects.
func NewAzureEtcdDataEncryption() *AzureEtcdDataEncryptionBuilder {
	return &AzureEtcdDataEncryptionBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureEtcdDataEncryptionBuilder) Empty() bool {
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

// CustomerManaged sets the value of the 'customer_managed' attribute to the given value.
//
// Contains the necessary attributes to support etcd data encryption with customer managed keys
// for Azure based clusters.
func (b *AzureEtcdDataEncryptionBuilder) CustomerManaged(value *AzureEtcdDataEncryptionCustomerManagedBuilder) *AzureEtcdDataEncryptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.customerManaged = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// KeyManagementMode sets the value of the 'key_management_mode' attribute to the given value.
func (b *AzureEtcdDataEncryptionBuilder) KeyManagementMode(value string) *AzureEtcdDataEncryptionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.keyManagementMode = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureEtcdDataEncryptionBuilder) Copy(object *AzureEtcdDataEncryption) *AzureEtcdDataEncryptionBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.customerManaged != nil {
		b.customerManaged = NewAzureEtcdDataEncryptionCustomerManaged().Copy(object.customerManaged)
	} else {
		b.customerManaged = nil
	}
	b.keyManagementMode = object.keyManagementMode
	return b
}

// Build creates a 'azure_etcd_data_encryption' object using the configuration stored in the builder.
func (b *AzureEtcdDataEncryptionBuilder) Build() (object *AzureEtcdDataEncryption, err error) {
	object = new(AzureEtcdDataEncryption)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.customerManaged != nil {
		object.customerManaged, err = b.customerManaged.Build()
		if err != nil {
			return
		}
	}
	object.keyManagementMode = b.keyManagementMode
	return
}
