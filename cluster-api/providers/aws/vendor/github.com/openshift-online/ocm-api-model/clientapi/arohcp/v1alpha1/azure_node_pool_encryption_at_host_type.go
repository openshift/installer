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

// AzureNodePoolEncryptionAtHost represents the values of the 'azure_node_pool_encryption_at_host' type.
//
// AzureNodePoolEncryptionAtHost defines the encryption setting for Encryption At Host.
// If not specified, Encryption at Host is not enabled.
type AzureNodePoolEncryptionAtHost struct {
	fieldSet_ []bool
	state     string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AzureNodePoolEncryptionAtHost) Empty() bool {
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

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// State indicates whether Encryption At Host is enabled.
// When enabled, it enhances Azure Disk Storage Server-Side Encryption to ensure that all temporary disks
// and disk caches are encrypted at rest and flow encrypted to the Storage clusters.
// Accepted values are: "disabled" or "enabled".
// If not specified, its value is "disabled", which indicates Encryption At Host is disabled.
// Immutable.
func (o *AzureNodePoolEncryptionAtHost) State() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.state
	}
	return ""
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
//
// State indicates whether Encryption At Host is enabled.
// When enabled, it enhances Azure Disk Storage Server-Side Encryption to ensure that all temporary disks
// and disk caches are encrypted at rest and flow encrypted to the Storage clusters.
// Accepted values are: "disabled" or "enabled".
// If not specified, its value is "disabled", which indicates Encryption At Host is disabled.
// Immutable.
func (o *AzureNodePoolEncryptionAtHost) GetState() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.state
	}
	return
}

// AzureNodePoolEncryptionAtHostListKind is the name of the type used to represent list of objects of
// type 'azure_node_pool_encryption_at_host'.
const AzureNodePoolEncryptionAtHostListKind = "AzureNodePoolEncryptionAtHostList"

// AzureNodePoolEncryptionAtHostListLinkKind is the name of the type used to represent links to list
// of objects of type 'azure_node_pool_encryption_at_host'.
const AzureNodePoolEncryptionAtHostListLinkKind = "AzureNodePoolEncryptionAtHostListLink"

// AzureNodePoolEncryptionAtHostNilKind is the name of the type used to nil lists of objects of
// type 'azure_node_pool_encryption_at_host'.
const AzureNodePoolEncryptionAtHostListNilKind = "AzureNodePoolEncryptionAtHostListNil"

// AzureNodePoolEncryptionAtHostList is a list of values of the 'azure_node_pool_encryption_at_host' type.
type AzureNodePoolEncryptionAtHostList struct {
	href  string
	link  bool
	items []*AzureNodePoolEncryptionAtHost
}

// Len returns the length of the list.
func (l *AzureNodePoolEncryptionAtHostList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AzureNodePoolEncryptionAtHostList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AzureNodePoolEncryptionAtHostList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AzureNodePoolEncryptionAtHostList) SetItems(items []*AzureNodePoolEncryptionAtHost) {
	l.items = items
}

// Items returns the items of the list.
func (l *AzureNodePoolEncryptionAtHostList) Items() []*AzureNodePoolEncryptionAtHost {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AzureNodePoolEncryptionAtHostList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AzureNodePoolEncryptionAtHostList) Get(i int) *AzureNodePoolEncryptionAtHost {
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
func (l *AzureNodePoolEncryptionAtHostList) Slice() []*AzureNodePoolEncryptionAtHost {
	var slice []*AzureNodePoolEncryptionAtHost
	if l == nil {
		slice = make([]*AzureNodePoolEncryptionAtHost, 0)
	} else {
		slice = make([]*AzureNodePoolEncryptionAtHost, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AzureNodePoolEncryptionAtHostList) Each(f func(item *AzureNodePoolEncryptionAtHost) bool) {
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
func (l *AzureNodePoolEncryptionAtHostList) Range(f func(index int, item *AzureNodePoolEncryptionAtHost) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
