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

// AWSSTSRole represents the values of the 'AWSSTS_role' type.
//
// Representation of an sts role for a rosa cluster
type AWSSTSRole struct {
	bitmap_            uint32
	roleARN            string
	roleType           string
	roleVersion        string
	hcpManagedPolicies bool
	isAdmin            bool
	managedPolicies    bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWSSTSRole) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// HcpManagedPolicies returns the value of the 'hcp_managed_policies' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Does this Role have HCP Managed Policies?
func (o *AWSSTSRole) HcpManagedPolicies() bool {
	if o != nil && o.bitmap_&1 != 0 {
		return o.hcpManagedPolicies
	}
	return false
}

// GetHcpManagedPolicies returns the value of the 'hcp_managed_policies' attribute and
// a flag indicating if the attribute has a value.
//
// Does this Role have HCP Managed Policies?
func (o *AWSSTSRole) GetHcpManagedPolicies() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.hcpManagedPolicies
	}
	return
}

// IsAdmin returns the value of the 'is_admin' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Does this role have Admin permission?
func (o *AWSSTSRole) IsAdmin() bool {
	if o != nil && o.bitmap_&2 != 0 {
		return o.isAdmin
	}
	return false
}

// GetIsAdmin returns the value of the 'is_admin' attribute and
// a flag indicating if the attribute has a value.
//
// Does this role have Admin permission?
func (o *AWSSTSRole) GetIsAdmin() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.isAdmin
	}
	return
}

// ManagedPolicies returns the value of the 'managed_policies' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Does this Role have Managed Policies?
func (o *AWSSTSRole) ManagedPolicies() bool {
	if o != nil && o.bitmap_&4 != 0 {
		return o.managedPolicies
	}
	return false
}

// GetManagedPolicies returns the value of the 'managed_policies' attribute and
// a flag indicating if the attribute has a value.
//
// Does this Role have Managed Policies?
func (o *AWSSTSRole) GetManagedPolicies() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.managedPolicies
	}
	return
}

// RoleARN returns the value of the 'role_ARN' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The AWS ARN for this Role
func (o *AWSSTSRole) RoleARN() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.roleARN
	}
	return ""
}

// GetRoleARN returns the value of the 'role_ARN' attribute and
// a flag indicating if the attribute has a value.
//
// The AWS ARN for this Role
func (o *AWSSTSRole) GetRoleARN() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.roleARN
	}
	return
}

// RoleType returns the value of the 'role_type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The type of this Role
func (o *AWSSTSRole) RoleType() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.roleType
	}
	return ""
}

// GetRoleType returns the value of the 'role_type' attribute and
// a flag indicating if the attribute has a value.
//
// The type of this Role
func (o *AWSSTSRole) GetRoleType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.roleType
	}
	return
}

// RoleVersion returns the value of the 'role_version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Openshift Version for this Role
func (o *AWSSTSRole) RoleVersion() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.roleVersion
	}
	return ""
}

// GetRoleVersion returns the value of the 'role_version' attribute and
// a flag indicating if the attribute has a value.
//
// The Openshift Version for this Role
func (o *AWSSTSRole) GetRoleVersion() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.roleVersion
	}
	return
}

// AWSSTSRoleListKind is the name of the type used to represent list of objects of
// type 'AWSSTS_role'.
const AWSSTSRoleListKind = "AWSSTSRoleList"

// AWSSTSRoleListLinkKind is the name of the type used to represent links to list
// of objects of type 'AWSSTS_role'.
const AWSSTSRoleListLinkKind = "AWSSTSRoleListLink"

// AWSSTSRoleNilKind is the name of the type used to nil lists of objects of
// type 'AWSSTS_role'.
const AWSSTSRoleListNilKind = "AWSSTSRoleListNil"

// AWSSTSRoleList is a list of values of the 'AWSSTS_role' type.
type AWSSTSRoleList struct {
	href  string
	link  bool
	items []*AWSSTSRole
}

// Len returns the length of the list.
func (l *AWSSTSRoleList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AWSSTSRoleList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AWSSTSRoleList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AWSSTSRoleList) SetItems(items []*AWSSTSRole) {
	l.items = items
}

// Items returns the items of the list.
func (l *AWSSTSRoleList) Items() []*AWSSTSRole {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AWSSTSRoleList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AWSSTSRoleList) Get(i int) *AWSSTSRole {
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
func (l *AWSSTSRoleList) Slice() []*AWSSTSRole {
	var slice []*AWSSTSRole
	if l == nil {
		slice = make([]*AWSSTSRole, 0)
	} else {
		slice = make([]*AWSSTSRole, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AWSSTSRoleList) Each(f func(item *AWSSTSRole) bool) {
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
func (l *AWSSTSRoleList) Range(f func(index int, item *AWSSTSRole) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
