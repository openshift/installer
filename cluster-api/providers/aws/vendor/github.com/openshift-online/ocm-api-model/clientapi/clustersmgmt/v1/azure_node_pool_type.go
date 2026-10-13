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

// AzureNodePool represents the values of the 'azure_node_pool' type.
//
// Representation of azure node pool specific parameters.
type AzureNodePool struct {
	fieldSet_        []bool
	vmSize           string
	encryptionAtHost *AzureNodePoolEncryptionAtHost
	osDisk           *AzureNodePoolOsDisk
	resourceName     string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AzureNodePool) Empty() bool {
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

// VMSize returns the value of the 'VM_size' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Virtual Machine size identifier used for the
// Nodes of the Node Pool.
// Availability of VM sizes are dependent on the Azure Location
// of the parent Cluster.
// Required during creation.
func (o *AzureNodePool) VMSize() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
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
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.vmSize
	}
	return
}

// EncryptionAtHost returns the value of the 'encryption_at_host' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// EncryptionAtHost contains Encryption At Host disk encryption configuration.
// When enabled, it enhances Azure Disk Storage Server-Side Encryption to ensure that all temporary disks
// and disk caches are encrypted at rest and flow encrypted to the Storage clusters.
// If not specified, Encryption at Host is not enabled.
// Immutable.
func (o *AzureNodePool) EncryptionAtHost() *AzureNodePoolEncryptionAtHost {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.encryptionAtHost
	}
	return nil
}

// GetEncryptionAtHost returns the value of the 'encryption_at_host' attribute and
// a flag indicating if the attribute has a value.
//
// EncryptionAtHost contains Encryption At Host disk encryption configuration.
// When enabled, it enhances Azure Disk Storage Server-Side Encryption to ensure that all temporary disks
// and disk caches are encrypted at rest and flow encrypted to the Storage clusters.
// If not specified, Encryption at Host is not enabled.
// Immutable.
func (o *AzureNodePool) GetEncryptionAtHost() (value *AzureNodePoolEncryptionAtHost, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.encryptionAtHost
	}
	return
}

// OsDisk returns the value of the 'os_disk' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The configuration for the OS disk used by the nodes in the Node Pool.
func (o *AzureNodePool) OsDisk() *AzureNodePoolOsDisk {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.osDisk
	}
	return nil
}

// GetOsDisk returns the value of the 'os_disk' attribute and
// a flag indicating if the attribute has a value.
//
// The configuration for the OS disk used by the nodes in the Node Pool.
func (o *AzureNodePool) GetOsDisk() (value *AzureNodePoolOsDisk, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.osDisk
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
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
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
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
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

// Items sets the items of the list.
func (l *AzureNodePoolList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AzureNodePoolList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AzureNodePoolList) SetItems(items []*AzureNodePool) {
	l.items = items
}

// Items returns the items of the list.
func (l *AzureNodePoolList) Items() []*AzureNodePool {
	if l == nil {
		return nil
	}
	return l.items
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
