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

// OperatorIAMRole represents the values of the 'operator_IAM_role' type.
//
// Contains the necessary attributes to allow each operator to access the necessary AWS resources
type OperatorIAMRole struct {
	bitmap_        uint32
	id             string
	name           string
	namespace      string
	roleARN        string
	serviceAccount string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *OperatorIAMRole) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Randomly-generated ID to identify the operator role
func (o *OperatorIAMRole) ID() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// Randomly-generated ID to identify the operator role
func (o *OperatorIAMRole) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.id
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the credentials secret used to access cloud resources
func (o *OperatorIAMRole) Name() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the credentials secret used to access cloud resources
func (o *OperatorIAMRole) GetName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.name
	}
	return
}

// Namespace returns the value of the 'namespace' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Namespace where the credentials secret lives in the cluster
func (o *OperatorIAMRole) Namespace() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.namespace
	}
	return ""
}

// GetNamespace returns the value of the 'namespace' attribute and
// a flag indicating if the attribute has a value.
//
// Namespace where the credentials secret lives in the cluster
func (o *OperatorIAMRole) GetNamespace() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.namespace
	}
	return
}

// RoleARN returns the value of the 'role_ARN' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Role to assume when accessing AWS resources
func (o *OperatorIAMRole) RoleARN() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.roleARN
	}
	return ""
}

// GetRoleARN returns the value of the 'role_ARN' attribute and
// a flag indicating if the attribute has a value.
//
// Role to assume when accessing AWS resources
func (o *OperatorIAMRole) GetRoleARN() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.roleARN
	}
	return
}

// ServiceAccount returns the value of the 'service_account' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Service account name to use when authenticating
func (o *OperatorIAMRole) ServiceAccount() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.serviceAccount
	}
	return ""
}

// GetServiceAccount returns the value of the 'service_account' attribute and
// a flag indicating if the attribute has a value.
//
// Service account name to use when authenticating
func (o *OperatorIAMRole) GetServiceAccount() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.serviceAccount
	}
	return
}

// OperatorIAMRoleListKind is the name of the type used to represent list of objects of
// type 'operator_IAM_role'.
const OperatorIAMRoleListKind = "OperatorIAMRoleList"

// OperatorIAMRoleListLinkKind is the name of the type used to represent links to list
// of objects of type 'operator_IAM_role'.
const OperatorIAMRoleListLinkKind = "OperatorIAMRoleListLink"

// OperatorIAMRoleNilKind is the name of the type used to nil lists of objects of
// type 'operator_IAM_role'.
const OperatorIAMRoleListNilKind = "OperatorIAMRoleListNil"

// OperatorIAMRoleList is a list of values of the 'operator_IAM_role' type.
type OperatorIAMRoleList struct {
	href  string
	link  bool
	items []*OperatorIAMRole
}

// Len returns the length of the list.
func (l *OperatorIAMRoleList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *OperatorIAMRoleList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *OperatorIAMRoleList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *OperatorIAMRoleList) SetItems(items []*OperatorIAMRole) {
	l.items = items
}

// Items returns the items of the list.
func (l *OperatorIAMRoleList) Items() []*OperatorIAMRole {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *OperatorIAMRoleList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *OperatorIAMRoleList) Get(i int) *OperatorIAMRole {
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
func (l *OperatorIAMRoleList) Slice() []*OperatorIAMRole {
	var slice []*OperatorIAMRole
	if l == nil {
		slice = make([]*OperatorIAMRole, 0)
	} else {
		slice = make([]*OperatorIAMRole, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *OperatorIAMRoleList) Each(f func(item *OperatorIAMRole) bool) {
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
func (l *OperatorIAMRoleList) Range(f func(index int, item *OperatorIAMRole) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
