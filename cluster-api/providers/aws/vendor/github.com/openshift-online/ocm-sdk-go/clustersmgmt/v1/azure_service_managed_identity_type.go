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

// AzureServiceManagedIdentity represents the values of the 'azure_service_managed_identity' type.
//
// Represents the information associated to an Azure User-Assigned
// Managed Identity whose purpose is to perform service level actions.
type AzureServiceManagedIdentity struct {
	bitmap_     uint32
	clientID    string
	principalID string
	resourceID  string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AzureServiceManagedIdentity) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ClientID returns the value of the 'client_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Client ID associated to the Azure User-Assigned Managed Identity.
// Readonly.
func (o *AzureServiceManagedIdentity) ClientID() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.clientID
	}
	return ""
}

// GetClientID returns the value of the 'client_ID' attribute and
// a flag indicating if the attribute has a value.
//
// The Client ID associated to the Azure User-Assigned Managed Identity.
// Readonly.
func (o *AzureServiceManagedIdentity) GetClientID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.clientID
	}
	return
}

// PrincipalID returns the value of the 'principal_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Principal ID associated to the Azure User-Assigned Managed Identity.
// Readonly.
func (o *AzureServiceManagedIdentity) PrincipalID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.principalID
	}
	return ""
}

// GetPrincipalID returns the value of the 'principal_ID' attribute and
// a flag indicating if the attribute has a value.
//
// The Principal ID associated to the Azure User-Assigned Managed Identity.
// Readonly.
func (o *AzureServiceManagedIdentity) GetPrincipalID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.principalID
	}
	return
}

// ResourceID returns the value of the 'resource_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Resource ID of the Azure User-Assigned Managed
// Identity. The managed identity represented must exist before
// creating the cluster.
// The Azure Resource Group Name specified as part of the Resource ID
// must belong to the Azure Subscription specified in `.azure.subscription_id`,
// and in the same Azure location as the cluster's region.
// The Azure Resource Group Name specified as part of the Resource ID
// must be a different Resource Group Name than the one specified in
// `.azure.managed_resource_group_name`.
// The Azure Resource Group Name specified as part of the Resource ID
// can be the same, or a different one than the one specified in
// `.azure.resource_group_name`.
// Required during creation.
// Immutable.
func (o *AzureServiceManagedIdentity) ResourceID() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.resourceID
	}
	return ""
}

// GetResourceID returns the value of the 'resource_ID' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure Resource ID of the Azure User-Assigned Managed
// Identity. The managed identity represented must exist before
// creating the cluster.
// The Azure Resource Group Name specified as part of the Resource ID
// must belong to the Azure Subscription specified in `.azure.subscription_id`,
// and in the same Azure location as the cluster's region.
// The Azure Resource Group Name specified as part of the Resource ID
// must be a different Resource Group Name than the one specified in
// `.azure.managed_resource_group_name`.
// The Azure Resource Group Name specified as part of the Resource ID
// can be the same, or a different one than the one specified in
// `.azure.resource_group_name`.
// Required during creation.
// Immutable.
func (o *AzureServiceManagedIdentity) GetResourceID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.resourceID
	}
	return
}

// AzureServiceManagedIdentityListKind is the name of the type used to represent list of objects of
// type 'azure_service_managed_identity'.
const AzureServiceManagedIdentityListKind = "AzureServiceManagedIdentityList"

// AzureServiceManagedIdentityListLinkKind is the name of the type used to represent links to list
// of objects of type 'azure_service_managed_identity'.
const AzureServiceManagedIdentityListLinkKind = "AzureServiceManagedIdentityListLink"

// AzureServiceManagedIdentityNilKind is the name of the type used to nil lists of objects of
// type 'azure_service_managed_identity'.
const AzureServiceManagedIdentityListNilKind = "AzureServiceManagedIdentityListNil"

// AzureServiceManagedIdentityList is a list of values of the 'azure_service_managed_identity' type.
type AzureServiceManagedIdentityList struct {
	href  string
	link  bool
	items []*AzureServiceManagedIdentity
}

// Len returns the length of the list.
func (l *AzureServiceManagedIdentityList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AzureServiceManagedIdentityList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AzureServiceManagedIdentityList) Get(i int) *AzureServiceManagedIdentity {
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
func (l *AzureServiceManagedIdentityList) Slice() []*AzureServiceManagedIdentity {
	var slice []*AzureServiceManagedIdentity
	if l == nil {
		slice = make([]*AzureServiceManagedIdentity, 0)
	} else {
		slice = make([]*AzureServiceManagedIdentity, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AzureServiceManagedIdentityList) Each(f func(item *AzureServiceManagedIdentity) bool) {
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
func (l *AzureServiceManagedIdentityList) Range(f func(index int, item *AzureServiceManagedIdentity) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
