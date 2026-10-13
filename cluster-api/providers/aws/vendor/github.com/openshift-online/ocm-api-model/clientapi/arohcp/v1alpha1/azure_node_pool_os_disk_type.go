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

// AzureNodePoolOsDisk represents the values of the 'azure_node_pool_os_disk' type.
//
// Defines the configuration of a Node Pool's OS disk.
type AzureNodePoolOsDisk struct {
	fieldSet_                  []bool
	persistence                string
	sizeGibibytes              int
	sseEncryptionSetResourceId string
	storageAccountType         string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AzureNodePoolOsDisk) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}
	for _, set := range o.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// Persistence returns the value of the 'persistence' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Specifies the OS Disk persistence for the OS Disks of the Nodes in the Node Pool.
// Valid values are:
// * persistent
// * ephemeral
// If not specified, Persistent OS Disks are used.
func (o *AzureNodePoolOsDisk) Persistence() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.persistence
	}
	return ""
}

// GetPersistence returns the value of the 'persistence' attribute and
// a flag indicating if the attribute has a value.
//
// Specifies the OS Disk persistence for the OS Disks of the Nodes in the Node Pool.
// Valid values are:
// * persistent
// * ephemeral
// If not specified, Persistent OS Disks are used.
func (o *AzureNodePoolOsDisk) GetPersistence() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.persistence
	}
	return
}

// SizeGibibytes returns the value of the 'size_gibibytes' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The size in GiB to assign to the OS disks of the
// Nodes in the Node Pool. The property
// is the number of bytes x 1024^3.
// If not specified, OS disk size is 64 GiB.
func (o *AzureNodePoolOsDisk) SizeGibibytes() int {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.sizeGibibytes
	}
	return 0
}

// GetSizeGibibytes returns the value of the 'size_gibibytes' attribute and
// a flag indicating if the attribute has a value.
//
// The size in GiB to assign to the OS disks of the
// Nodes in the Node Pool. The property
// is the number of bytes x 1024^3.
// If not specified, OS disk size is 64 GiB.
func (o *AzureNodePoolOsDisk) GetSizeGibibytes() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.sizeGibibytes
	}
	return
}

// SseEncryptionSetResourceId returns the value of the 'sse_encryption_set_resource_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Resource ID of a pre-existing Azure Disk Encryption Set (DES).
// When provided, Server-Side Encryption (SSE) on the OS Disks of the Nodes of the Node Pool
// is performed using the provided Disk Encryption Set.
// It must be located in the same Azure location as the parent Cluster.
// It must be located in the same Azure Subscription as the parent Cluster.
// The Azure Resource Group Name specified as part of it must be a different resource group name
// than the one specified in the parent Cluster's `managed_resource_group_name`.
// The Azure Resource Group Name specified as part of it can be the same, or a different one
// than the one specified in the parent Cluster's `resource_group_name`.
// If not specified, Server-Side Encryption (SSE) on the OS Disks of the Nodes of the Node Pool
// is performed with platform managed keys.
func (o *AzureNodePoolOsDisk) SseEncryptionSetResourceId() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.sseEncryptionSetResourceId
	}
	return ""
}

// GetSseEncryptionSetResourceId returns the value of the 'sse_encryption_set_resource_id' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure Resource ID of a pre-existing Azure Disk Encryption Set (DES).
// When provided, Server-Side Encryption (SSE) on the OS Disks of the Nodes of the Node Pool
// is performed using the provided Disk Encryption Set.
// It must be located in the same Azure location as the parent Cluster.
// It must be located in the same Azure Subscription as the parent Cluster.
// The Azure Resource Group Name specified as part of it must be a different resource group name
// than the one specified in the parent Cluster's `managed_resource_group_name`.
// The Azure Resource Group Name specified as part of it can be the same, or a different one
// than the one specified in the parent Cluster's `resource_group_name`.
// If not specified, Server-Side Encryption (SSE) on the OS Disks of the Nodes of the Node Pool
// is performed with platform managed keys.
func (o *AzureNodePoolOsDisk) GetSseEncryptionSetResourceId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.sseEncryptionSetResourceId
	}
	return
}

// StorageAccountType returns the value of the 'storage_account_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The disk storage account type to use for the OS disks of the Nodes in the
// Node Pool. Valid values are:
// * Standard_LRS: HDD
// * StandardSSD_LRS: Standard SSD
// * Premium_LRS: Premium SDD
// * UltraSSD_LRS: Ultra SDD
//
// If not specified, `Premium_LRS` is used.
func (o *AzureNodePoolOsDisk) StorageAccountType() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.storageAccountType
	}
	return ""
}

// GetStorageAccountType returns the value of the 'storage_account_type' attribute and
// a flag indicating if the attribute has a value.
//
// The disk storage account type to use for the OS disks of the Nodes in the
// Node Pool. Valid values are:
// * Standard_LRS: HDD
// * StandardSSD_LRS: Standard SSD
// * Premium_LRS: Premium SDD
// * UltraSSD_LRS: Ultra SDD
//
// If not specified, `Premium_LRS` is used.
func (o *AzureNodePoolOsDisk) GetStorageAccountType() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.storageAccountType
	}
	return
}

// AzureNodePoolOsDiskListKind is the name of the type used to represent list of objects of
// type 'azure_node_pool_os_disk'.
const AzureNodePoolOsDiskListKind = "AzureNodePoolOsDiskList"

// AzureNodePoolOsDiskListLinkKind is the name of the type used to represent links to list
// of objects of type 'azure_node_pool_os_disk'.
const AzureNodePoolOsDiskListLinkKind = "AzureNodePoolOsDiskListLink"

// AzureNodePoolOsDiskNilKind is the name of the type used to nil lists of objects of
// type 'azure_node_pool_os_disk'.
const AzureNodePoolOsDiskListNilKind = "AzureNodePoolOsDiskListNil"

// AzureNodePoolOsDiskList is a list of values of the 'azure_node_pool_os_disk' type.
type AzureNodePoolOsDiskList struct {
	href  string
	link  bool
	items []*AzureNodePoolOsDisk
}

// Len returns the length of the list.
func (l *AzureNodePoolOsDiskList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AzureNodePoolOsDiskList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AzureNodePoolOsDiskList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AzureNodePoolOsDiskList) SetItems(items []*AzureNodePoolOsDisk) {
	l.items = items
}

// Items returns the items of the list.
func (l *AzureNodePoolOsDiskList) Items() []*AzureNodePoolOsDisk {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AzureNodePoolOsDiskList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AzureNodePoolOsDiskList) Get(i int) *AzureNodePoolOsDisk {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *AzureNodePoolOsDiskList) Slice() []*AzureNodePoolOsDisk {
	var slice []*AzureNodePoolOsDisk
	if l == nil {
		slice = make([]*AzureNodePoolOsDisk, 0)
	} else {
		slice = make([]*AzureNodePoolOsDisk, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AzureNodePoolOsDiskList) Each(f func(item *AzureNodePoolOsDisk) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *AzureNodePoolOsDiskList) Range(f func(index int, item *AzureNodePoolOsDisk) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
