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

// Representation of azure node pool specific parameters.
type AzureNodePoolBuilder struct {
	fieldSet_        []bool
	vmSize           string
	encryptionAtHost *AzureNodePoolEncryptionAtHostBuilder
	osDisk           *AzureNodePoolOsDiskBuilder
	resourceName     string
}

// NewAzureNodePool creates a new builder of 'azure_node_pool' objects.
func NewAzureNodePool() *AzureNodePoolBuilder {
	return &AzureNodePoolBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureNodePoolBuilder) Empty() bool {
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

// VMSize sets the value of the 'VM_size' attribute to the given value.
func (b *AzureNodePoolBuilder) VMSize(value string) *AzureNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.vmSize = value
	b.fieldSet_[0] = true
	return b
}

// EncryptionAtHost sets the value of the 'encryption_at_host' attribute to the given value.
//
// AzureNodePoolEncryptionAtHost defines the encryption setting for Encryption At Host.
// If not specified, Encryption at Host is not enabled.
func (b *AzureNodePoolBuilder) EncryptionAtHost(value *AzureNodePoolEncryptionAtHostBuilder) *AzureNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.encryptionAtHost = value
	if value != nil {
		b.fieldSet_[1] = true
	} else {
		b.fieldSet_[1] = false
	}
	return b
}

// OsDisk sets the value of the 'os_disk' attribute to the given value.
//
// Defines the configuration of a Node Pool's OS disk.
func (b *AzureNodePoolBuilder) OsDisk(value *AzureNodePoolOsDiskBuilder) *AzureNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.osDisk = value
	if value != nil {
		b.fieldSet_[2] = true
	} else {
		b.fieldSet_[2] = false
	}
	return b
}

// ResourceName sets the value of the 'resource_name' attribute to the given value.
func (b *AzureNodePoolBuilder) ResourceName(value string) *AzureNodePoolBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.resourceName = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureNodePoolBuilder) Copy(object *AzureNodePool) *AzureNodePoolBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.vmSize = object.vmSize
	if object.encryptionAtHost != nil {
		b.encryptionAtHost = NewAzureNodePoolEncryptionAtHost().Copy(object.encryptionAtHost)
	} else {
		b.encryptionAtHost = nil
	}
	if object.osDisk != nil {
		b.osDisk = NewAzureNodePoolOsDisk().Copy(object.osDisk)
	} else {
		b.osDisk = nil
	}
	b.resourceName = object.resourceName
	return b
}

// Build creates a 'azure_node_pool' object using the configuration stored in the builder.
func (b *AzureNodePoolBuilder) Build() (object *AzureNodePool, err error) {
	object = new(AzureNodePool)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.vmSize = b.vmSize
	if b.encryptionAtHost != nil {
		object.encryptionAtHost, err = b.encryptionAtHost.Build()
		if err != nil {
			return
		}
	}
	if b.osDisk != nil {
		object.osDisk, err = b.osDisk.Build()
		if err != nil {
			return
		}
	}
	object.resourceName = b.resourceName
	return
}
