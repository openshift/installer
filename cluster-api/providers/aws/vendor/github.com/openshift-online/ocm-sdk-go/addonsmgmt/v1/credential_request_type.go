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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

// CredentialRequest represents the values of the 'credential_request' type.
//
// Contains the necessary attributes to allow each operator to access the necessary AWS resources
type CredentialRequest struct {
	bitmap_           uint32
	name              string
	namespace         string
	policyPermissions []string
	serviceAccount    string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *CredentialRequest) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the credentials secret used to access cloud resources
func (o *CredentialRequest) Name() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the credentials secret used to access cloud resources
func (o *CredentialRequest) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.name
	}
	return
}

// Namespace returns the value of the 'namespace' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Namespace where the credentials secret lives in the cluster
func (o *CredentialRequest) Namespace() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.namespace
	}
	return ""
}

// GetNamespace returns the value of the 'namespace' attribute and
// a flag indicating if the attribute has a value.
//
// Namespace where the credentials secret lives in the cluster
func (o *CredentialRequest) GetNamespace() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.namespace
	}
	return
}

// PolicyPermissions returns the value of the 'policy_permissions' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of policy permissions needed to access cloud resources
func (o *CredentialRequest) PolicyPermissions() []string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.policyPermissions
	}
	return nil
}

// GetPolicyPermissions returns the value of the 'policy_permissions' attribute and
// a flag indicating if the attribute has a value.
//
// List of policy permissions needed to access cloud resources
func (o *CredentialRequest) GetPolicyPermissions() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.policyPermissions
	}
	return
}

// ServiceAccount returns the value of the 'service_account' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Service account name to use when authenticating
func (o *CredentialRequest) ServiceAccount() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.serviceAccount
	}
	return ""
}

// GetServiceAccount returns the value of the 'service_account' attribute and
// a flag indicating if the attribute has a value.
//
// Service account name to use when authenticating
func (o *CredentialRequest) GetServiceAccount() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.serviceAccount
	}
	return
}

// CredentialRequestListKind is the name of the type used to represent list of objects of
// type 'credential_request'.
const CredentialRequestListKind = "CredentialRequestList"

// CredentialRequestListLinkKind is the name of the type used to represent links to list
// of objects of type 'credential_request'.
const CredentialRequestListLinkKind = "CredentialRequestListLink"

// CredentialRequestNilKind is the name of the type used to nil lists of objects of
// type 'credential_request'.
const CredentialRequestListNilKind = "CredentialRequestListNil"

// CredentialRequestList is a list of values of the 'credential_request' type.
type CredentialRequestList struct {
	href  string
	link  bool
	items []*CredentialRequest
}

// Len returns the length of the list.
func (l *CredentialRequestList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *CredentialRequestList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *CredentialRequestList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *CredentialRequestList) SetItems(items []*CredentialRequest) {
	l.items = items
}

// Items returns the items of the list.
func (l *CredentialRequestList) Items() []*CredentialRequest {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *CredentialRequestList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *CredentialRequestList) Get(i int) *CredentialRequest {
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
func (l *CredentialRequestList) Slice() []*CredentialRequest {
	var slice []*CredentialRequest
	if l == nil {
		slice = make([]*CredentialRequest, 0)
	} else {
		slice = make([]*CredentialRequest, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *CredentialRequestList) Each(f func(item *CredentialRequest) bool) {
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
func (l *CredentialRequestList) Range(f func(index int, item *CredentialRequest) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
