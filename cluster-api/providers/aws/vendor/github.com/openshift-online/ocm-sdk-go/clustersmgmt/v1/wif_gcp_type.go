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

// WifGcp represents the values of the 'wif_gcp' type.
type WifGcp struct {
	bitmap_              uint32
	impersonatorEmail    string
	projectId            string
	projectNumber        string
	rolePrefix           string
	serviceAccounts      []*WifServiceAccount
	support              *WifSupport
	workloadIdentityPool *WifPool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *WifGcp) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ImpersonatorEmail returns the value of the 'impersonator_email' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// This is the service account email that OCM will use to access other SAs.
func (o *WifGcp) ImpersonatorEmail() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.impersonatorEmail
	}
	return ""
}

// GetImpersonatorEmail returns the value of the 'impersonator_email' attribute and
// a flag indicating if the attribute has a value.
//
// This is the service account email that OCM will use to access other SAs.
func (o *WifGcp) GetImpersonatorEmail() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.impersonatorEmail
	}
	return
}

// ProjectId returns the value of the 'project_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// This represents the GCP project ID in which the wif resources will be configured.
func (o *WifGcp) ProjectId() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.projectId
	}
	return ""
}

// GetProjectId returns the value of the 'project_id' attribute and
// a flag indicating if the attribute has a value.
//
// This represents the GCP project ID in which the wif resources will be configured.
func (o *WifGcp) GetProjectId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.projectId
	}
	return
}

// ProjectNumber returns the value of the 'project_number' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// This represents the GCP project number in which the wif resources will be configured.
func (o *WifGcp) ProjectNumber() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.projectNumber
	}
	return ""
}

// GetProjectNumber returns the value of the 'project_number' attribute and
// a flag indicating if the attribute has a value.
//
// This represents the GCP project number in which the wif resources will be configured.
func (o *WifGcp) GetProjectNumber() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.projectNumber
	}
	return
}

// RolePrefix returns the value of the 'role_prefix' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Prefix for naming GCP custom roles configured.
func (o *WifGcp) RolePrefix() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.rolePrefix
	}
	return ""
}

// GetRolePrefix returns the value of the 'role_prefix' attribute and
// a flag indicating if the attribute has a value.
//
// Prefix for naming GCP custom roles configured.
func (o *WifGcp) GetRolePrefix() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.rolePrefix
	}
	return
}

// ServiceAccounts returns the value of the 'service_accounts' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The list of service accounts and their associated roles that will need to be
// configured on the user's GCP project.
func (o *WifGcp) ServiceAccounts() []*WifServiceAccount {
	if o != nil && o.bitmap_&16 != 0 {
		return o.serviceAccounts
	}
	return nil
}

// GetServiceAccounts returns the value of the 'service_accounts' attribute and
// a flag indicating if the attribute has a value.
//
// The list of service accounts and their associated roles that will need to be
// configured on the user's GCP project.
func (o *WifGcp) GetServiceAccounts() (value []*WifServiceAccount, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.serviceAccounts
	}
	return
}

// Support returns the value of the 'support' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Defines the access configuration for support.
func (o *WifGcp) Support() *WifSupport {
	if o != nil && o.bitmap_&32 != 0 {
		return o.support
	}
	return nil
}

// GetSupport returns the value of the 'support' attribute and
// a flag indicating if the attribute has a value.
//
// Defines the access configuration for support.
func (o *WifGcp) GetSupport() (value *WifSupport, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.support
	}
	return
}

// WorkloadIdentityPool returns the value of the 'workload_identity_pool' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The workload identity configuration data that will be used to create the
// workload identity pool on the user's account.
func (o *WifGcp) WorkloadIdentityPool() *WifPool {
	if o != nil && o.bitmap_&64 != 0 {
		return o.workloadIdentityPool
	}
	return nil
}

// GetWorkloadIdentityPool returns the value of the 'workload_identity_pool' attribute and
// a flag indicating if the attribute has a value.
//
// The workload identity configuration data that will be used to create the
// workload identity pool on the user's account.
func (o *WifGcp) GetWorkloadIdentityPool() (value *WifPool, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.workloadIdentityPool
	}
	return
}

// WifGcpListKind is the name of the type used to represent list of objects of
// type 'wif_gcp'.
const WifGcpListKind = "WifGcpList"

// WifGcpListLinkKind is the name of the type used to represent links to list
// of objects of type 'wif_gcp'.
const WifGcpListLinkKind = "WifGcpListLink"

// WifGcpNilKind is the name of the type used to nil lists of objects of
// type 'wif_gcp'.
const WifGcpListNilKind = "WifGcpListNil"

// WifGcpList is a list of values of the 'wif_gcp' type.
type WifGcpList struct {
	href  string
	link  bool
	items []*WifGcp
}

// Len returns the length of the list.
func (l *WifGcpList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *WifGcpList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *WifGcpList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *WifGcpList) SetItems(items []*WifGcp) {
	l.items = items
}

// Items returns the items of the list.
func (l *WifGcpList) Items() []*WifGcp {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *WifGcpList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *WifGcpList) Get(i int) *WifGcp {
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
func (l *WifGcpList) Slice() []*WifGcp {
	var slice []*WifGcp
	if l == nil {
		slice = make([]*WifGcp, 0)
	} else {
		slice = make([]*WifGcp, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *WifGcpList) Each(f func(item *WifGcp) bool) {
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
func (l *WifGcpList) Range(f func(index int, item *WifGcp) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
