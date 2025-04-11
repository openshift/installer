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

// AzureNodePool represents the values of the 'azure_node_pool' type.
//
// Representation of azure node pool specific parameters.
type AzureNodePool struct {
	bitmap_                  uint32
	osDiskSizeGibibytes      int
	osDiskStorageAccountType string
	vmSize                   string
	resourceName             string
	ephemeralOSDiskEnabled   bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AzureNodePool) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// OSDiskSizeGibibytes returns the value of the 'OS_disk_size_gibibytes' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The size in GiB to assign to the OS disks of the
// Nodes in the Node Pool. The property
// is the number of bytes x 1024^3.
// If not specified, OS disk size is 30 GiB.
func (o *AzureNodePool) OSDiskSizeGibibytes() int {
	if o != nil && o.bitmap_&1 != 0 {
		return o.osDiskSizeGibibytes
	}
	return 0
}

// GetOSDiskSizeGibibytes returns the value of the 'OS_disk_size_gibibytes' attribute and
// a flag indicating if the attribute has a value.
//
// The size in GiB to assign to the OS disks of the
// Nodes in the Node Pool. The property
// is the number of bytes x 1024^3.
// If not specified, OS disk size is 30 GiB.
func (o *AzureNodePool) GetOSDiskSizeGibibytes() (value int, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.osDiskSizeGibibytes
	}
	return
}

// OSDiskStorageAccountType returns the value of the 'OS_disk_storage_account_type' attribute, or
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
func (o *AzureNodePool) OSDiskStorageAccountType() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.osDiskStorageAccountType
	}
	return ""
}

// GetOSDiskStorageAccountType returns the value of the 'OS_disk_storage_account_type' attribute and
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
func (o *AzureNodePool) GetOSDiskStorageAccountType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.osDiskStorageAccountType
	}
	return
}

// VMSize returns the value of the 'VM_size' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Virtual Machine size identifier used for the
// Nodes of the Node Pool.
// Availability of VM sizes are dependent on the Azure Location
// of the parent Cluster.
// Required during creation.
func (o *AzureNodePool) VMSize() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.vmSize
	}
	return ""
}

// GetVMSize returns the value of the 'VM_size' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure Virtual Machine size identifier used for the
// Nodes of the Node Pool.
// Availability of VM sizes are dependent on the Azure Location
// of the parent Cluster.
// Required during creation.
func (o *AzureNodePool) GetVMSize() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.vmSize
	}
	return
}

// EphemeralOSDiskEnabled returns the value of the 'ephemeral_OS_disk_enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Enables Ephemeral OS Disks for the Nodes in the Node Pool.
// If not specified, no Ephemeral OS Disks are used.
func (o *AzureNodePool) EphemeralOSDiskEnabled() bool {
	if o != nil && o.bitmap_&8 != 0 {
		return o.ephemeralOSDiskEnabled
	}
	return false
}

// GetEphemeralOSDiskEnabled returns the value of the 'ephemeral_OS_disk_enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Enables Ephemeral OS Disks for the Nodes in the Node Pool.
// If not specified, no Ephemeral OS Disks are used.
func (o *AzureNodePool) GetEphemeralOSDiskEnabled() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.ephemeralOSDiskEnabled
	}
	return
}

// ResourceName returns the value of the 'resource_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ResourceName is the Azure Resource Name of the NodePool.
// ResourceName must be within the Azure Resource Group Name of the parent
// Cluster it belongs to.
// ResourceName must be located in the same Azure Location as the parent
// Cluster it belongs to.
// ResourceName must be located in the same Azure Subscription as the parent
// Cluster it belongs to.
// ResourceName must belong to the same Microsoft Entra Tenant ID as the parent
// Cluster it belongs to.
// Required during creation.
// Immutable.
func (o *AzureNodePool) ResourceName() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.resourceName
	}
	return ""
}

// GetResourceName returns the value of the 'resource_name' attribute and
// a flag indicating if the attribute has a value.
//
// ResourceName is the Azure Resource Name of the NodePool.
// ResourceName must be within the Azure Resource Group Name of the parent
// Cluster it belongs to.
// ResourceName must be located in the same Azure Location as the parent
// Cluster it belongs to.
// ResourceName must be located in the same Azure Subscription as the parent
// Cluster it belongs to.
// ResourceName must belong to the same Microsoft Entra Tenant ID as the parent
// Cluster it belongs to.
// Required during creation.
// Immutable.
func (o *AzureNodePool) GetResourceName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.resourceName
	}
	return
}

// AzureNodePoolListKind is the name of the type used to represent list of objects of
// type 'azure_node_pool'.
const AzureNodePoolListKind = "AzureNodePoolList"

// AzureNodePoolListLinkKind is the name of the type used to represent links to list
// of objects of type 'azure_node_pool'.
const AzureNodePoolListLinkKind = "AzureNodePoolListLink"

// AzureNodePoolNilKind is the name of the type used to nil lists of objects of
// type 'azure_node_pool'.
const AzureNodePoolListNilKind = "AzureNodePoolListNil"

// AzureNodePoolList is a list of values of the 'azure_node_pool' type.
type AzureNodePoolList struct {
	href  string
	link  bool
	items []*AzureNodePool
}

// Len returns the length of the list.
func (l *AzureNodePoolList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AzureNodePoolList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AzureNodePoolList) Get(i int) *AzureNodePool {
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
func (l *AzureNodePoolList) Slice() []*AzureNodePool {
	var slice []*AzureNodePool
	if l == nil {
		slice = make([]*AzureNodePool, 0)
	} else {
		slice = make([]*AzureNodePool, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AzureNodePoolList) Each(f func(item *AzureNodePool) bool) {
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
func (l *AzureNodePoolList) Range(f func(index int, item *AzureNodePool) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
