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

// AzureOperatorsAuthenticationManagedIdentities represents the values of the 'azure_operators_authentication_managed_identities' type.
//
// Represents the information related to Azure User-Assigned managed identities
// needed to perform Operators authentication based on Azure User-Assigned
// Managed Identities
type AzureOperatorsAuthenticationManagedIdentities struct {
	bitmap_                                uint32
	controlPlaneOperatorsManagedIdentities map[string]*AzureControlPlaneManagedIdentity
	dataPlaneOperatorsManagedIdentities    map[string]*AzureDataPlaneManagedIdentity
	managedIdentitiesDataPlaneIdentityUrl  string
	serviceManagedIdentity                 *AzureServiceManagedIdentity
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AzureOperatorsAuthenticationManagedIdentities) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ControlPlaneOperatorsManagedIdentities returns the value of the 'control_plane_operators_managed_identities' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The set of Azure User-Assigned Managed Identities leveraged for the
// Control Plane operators of the cluster. The set of required managed
// identities is dependent on the Cluster's OpenShift version.
// Immutable
func (o *AzureOperatorsAuthenticationManagedIdentities) ControlPlaneOperatorsManagedIdentities() map[string]*AzureControlPlaneManagedIdentity {
	if o != nil && o.bitmap_&1 != 0 {
		return o.controlPlaneOperatorsManagedIdentities
	}
	return nil
}

// GetControlPlaneOperatorsManagedIdentities returns the value of the 'control_plane_operators_managed_identities' attribute and
// a flag indicating if the attribute has a value.
//
// The set of Azure User-Assigned Managed Identities leveraged for the
// Control Plane operators of the cluster. The set of required managed
// identities is dependent on the Cluster's OpenShift version.
// Immutable
func (o *AzureOperatorsAuthenticationManagedIdentities) GetControlPlaneOperatorsManagedIdentities() (value map[string]*AzureControlPlaneManagedIdentity, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.controlPlaneOperatorsManagedIdentities
	}
	return
}

// DataPlaneOperatorsManagedIdentities returns the value of the 'data_plane_operators_managed_identities' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The set of Azure User-Assigned Managed Identities leveraged for the
// Data Plane operators of the cluster. The set of required managed
// identities is dependent on the Cluster's OpenShift version.
// Immutable.
func (o *AzureOperatorsAuthenticationManagedIdentities) DataPlaneOperatorsManagedIdentities() map[string]*AzureDataPlaneManagedIdentity {
	if o != nil && o.bitmap_&2 != 0 {
		return o.dataPlaneOperatorsManagedIdentities
	}
	return nil
}

// GetDataPlaneOperatorsManagedIdentities returns the value of the 'data_plane_operators_managed_identities' attribute and
// a flag indicating if the attribute has a value.
//
// The set of Azure User-Assigned Managed Identities leveraged for the
// Data Plane operators of the cluster. The set of required managed
// identities is dependent on the Cluster's OpenShift version.
// Immutable.
func (o *AzureOperatorsAuthenticationManagedIdentities) GetDataPlaneOperatorsManagedIdentities() (value map[string]*AzureDataPlaneManagedIdentity, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.dataPlaneOperatorsManagedIdentities
	}
	return
}

// ManagedIdentitiesDataPlaneIdentityUrl returns the value of the 'managed_identities_data_plane_identity_url' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Managed Identities Data Plane Identity URL associated with the
// cluster. It is the URL that will be used to communicate with the
// Managed Identities Resource Provider (MI RP).
// Required during creation.
// Immutable.
func (o *AzureOperatorsAuthenticationManagedIdentities) ManagedIdentitiesDataPlaneIdentityUrl() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.managedIdentitiesDataPlaneIdentityUrl
	}
	return ""
}

// GetManagedIdentitiesDataPlaneIdentityUrl returns the value of the 'managed_identities_data_plane_identity_url' attribute and
// a flag indicating if the attribute has a value.
//
// The Managed Identities Data Plane Identity URL associated with the
// cluster. It is the URL that will be used to communicate with the
// Managed Identities Resource Provider (MI RP).
// Required during creation.
// Immutable.
func (o *AzureOperatorsAuthenticationManagedIdentities) GetManagedIdentitiesDataPlaneIdentityUrl() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.managedIdentitiesDataPlaneIdentityUrl
	}
	return
}

// ServiceManagedIdentity returns the value of the 'service_managed_identity' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure User-Assigned Managed Identity used to perform service
// level actions. Specifically:
//   - Add Federated Identity Credentials to the identities in
//     `data_plane_operators_managed_identities` that belong to Data
//     Plane Cluster Operators
//   - Perform permissions validation for the BYOVNet related resources
//     associated to the Cluster
//
// Required during creation.
// Immutable.
func (o *AzureOperatorsAuthenticationManagedIdentities) ServiceManagedIdentity() *AzureServiceManagedIdentity {
	if o != nil && o.bitmap_&8 != 0 {
		return o.serviceManagedIdentity
	}
	return nil
}

// GetServiceManagedIdentity returns the value of the 'service_managed_identity' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure User-Assigned Managed Identity used to perform service
// level actions. Specifically:
//   - Add Federated Identity Credentials to the identities in
//     `data_plane_operators_managed_identities` that belong to Data
//     Plane Cluster Operators
//   - Perform permissions validation for the BYOVNet related resources
//     associated to the Cluster
//
// Required during creation.
// Immutable.
func (o *AzureOperatorsAuthenticationManagedIdentities) GetServiceManagedIdentity() (value *AzureServiceManagedIdentity, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.serviceManagedIdentity
	}
	return
}

// AzureOperatorsAuthenticationManagedIdentitiesListKind is the name of the type used to represent list of objects of
// type 'azure_operators_authentication_managed_identities'.
const AzureOperatorsAuthenticationManagedIdentitiesListKind = "AzureOperatorsAuthenticationManagedIdentitiesList"

// AzureOperatorsAuthenticationManagedIdentitiesListLinkKind is the name of the type used to represent links to list
// of objects of type 'azure_operators_authentication_managed_identities'.
const AzureOperatorsAuthenticationManagedIdentitiesListLinkKind = "AzureOperatorsAuthenticationManagedIdentitiesListLink"

// AzureOperatorsAuthenticationManagedIdentitiesNilKind is the name of the type used to nil lists of objects of
// type 'azure_operators_authentication_managed_identities'.
const AzureOperatorsAuthenticationManagedIdentitiesListNilKind = "AzureOperatorsAuthenticationManagedIdentitiesListNil"

// AzureOperatorsAuthenticationManagedIdentitiesList is a list of values of the 'azure_operators_authentication_managed_identities' type.
type AzureOperatorsAuthenticationManagedIdentitiesList struct {
	href  string
	link  bool
	items []*AzureOperatorsAuthenticationManagedIdentities
}

// Len returns the length of the list.
func (l *AzureOperatorsAuthenticationManagedIdentitiesList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *AzureOperatorsAuthenticationManagedIdentitiesList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AzureOperatorsAuthenticationManagedIdentitiesList) Get(i int) *AzureOperatorsAuthenticationManagedIdentities {
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
func (l *AzureOperatorsAuthenticationManagedIdentitiesList) Slice() []*AzureOperatorsAuthenticationManagedIdentities {
	var slice []*AzureOperatorsAuthenticationManagedIdentities
	if l == nil {
		slice = make([]*AzureOperatorsAuthenticationManagedIdentities, 0)
	} else {
		slice = make([]*AzureOperatorsAuthenticationManagedIdentities, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AzureOperatorsAuthenticationManagedIdentitiesList) Each(f func(item *AzureOperatorsAuthenticationManagedIdentities) bool) {
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
func (l *AzureOperatorsAuthenticationManagedIdentitiesList) Range(f func(index int, item *AzureOperatorsAuthenticationManagedIdentities) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
