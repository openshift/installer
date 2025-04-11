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

// WifServiceAccount represents the values of the 'wif_service_account' type.
type WifServiceAccount struct {
	bitmap_           uint32
	accessMethod      WifAccessMethod
	credentialRequest *WifCredentialRequest
	osdRole           string
	roles             []*WifRole
	serviceAccountId  string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *WifServiceAccount) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AccessMethod returns the value of the 'access_method' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifServiceAccount) AccessMethod() WifAccessMethod {
	if o != nil && o.bitmap_&1 != 0 {
		return o.accessMethod
	}
	return WifAccessMethod("")
}

// GetAccessMethod returns the value of the 'access_method' attribute and
// a flag indicating if the attribute has a value.
func (o *WifServiceAccount) GetAccessMethod() (value WifAccessMethod, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.accessMethod
	}
	return
}

// CredentialRequest returns the value of the 'credential_request' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifServiceAccount) CredentialRequest() *WifCredentialRequest {
	if o != nil && o.bitmap_&2 != 0 {
		return o.credentialRequest
	}
	return nil
}

// GetCredentialRequest returns the value of the 'credential_request' attribute and
// a flag indicating if the attribute has a value.
func (o *WifServiceAccount) GetCredentialRequest() (value *WifCredentialRequest, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.credentialRequest
	}
	return
}

// OsdRole returns the value of the 'osd_role' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifServiceAccount) OsdRole() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.osdRole
	}
	return ""
}

// GetOsdRole returns the value of the 'osd_role' attribute and
// a flag indicating if the attribute has a value.
func (o *WifServiceAccount) GetOsdRole() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.osdRole
	}
	return
}

// Roles returns the value of the 'roles' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifServiceAccount) Roles() []*WifRole {
	if o != nil && o.bitmap_&8 != 0 {
		return o.roles
	}
	return nil
}

// GetRoles returns the value of the 'roles' attribute and
// a flag indicating if the attribute has a value.
func (o *WifServiceAccount) GetRoles() (value []*WifRole, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.roles
	}
	return
}

// ServiceAccountId returns the value of the 'service_account_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifServiceAccount) ServiceAccountId() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.serviceAccountId
	}
	return ""
}

// GetServiceAccountId returns the value of the 'service_account_id' attribute and
// a flag indicating if the attribute has a value.
func (o *WifServiceAccount) GetServiceAccountId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.serviceAccountId
	}
	return
}

// WifServiceAccountListKind is the name of the type used to represent list of objects of
// type 'wif_service_account'.
const WifServiceAccountListKind = "WifServiceAccountList"

// WifServiceAccountListLinkKind is the name of the type used to represent links to list
// of objects of type 'wif_service_account'.
const WifServiceAccountListLinkKind = "WifServiceAccountListLink"

// WifServiceAccountNilKind is the name of the type used to nil lists of objects of
// type 'wif_service_account'.
const WifServiceAccountListNilKind = "WifServiceAccountListNil"

// WifServiceAccountList is a list of values of the 'wif_service_account' type.
type WifServiceAccountList struct {
	href  string
	link  bool
	items []*WifServiceAccount
}

// Len returns the length of the list.
func (l *WifServiceAccountList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *WifServiceAccountList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *WifServiceAccountList) Get(i int) *WifServiceAccount {
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
func (l *WifServiceAccountList) Slice() []*WifServiceAccount {
	var slice []*WifServiceAccount
	if l == nil {
		slice = make([]*WifServiceAccount, 0)
	} else {
		slice = make([]*WifServiceAccount, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *WifServiceAccountList) Each(f func(item *WifServiceAccount) bool) {
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
func (l *WifServiceAccountList) Range(f func(index int, item *WifServiceAccount) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
