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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

// InstanceIAMRoles represents the values of the 'instance_IAM_roles' type.
//
// Contains the necessary attributes to support role-based authentication on AWS.
type InstanceIAMRoles struct {
	bitmap_       uint32
	masterRoleARN string
	workerRoleARN string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *InstanceIAMRoles) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// MasterRoleARN returns the value of the 'master_role_ARN' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The IAM role ARN that will be attached to master instances
func (o *InstanceIAMRoles) MasterRoleARN() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.masterRoleARN
	}
	return ""
}

// GetMasterRoleARN returns the value of the 'master_role_ARN' attribute and
// a flag indicating if the attribute has a value.
//
// The IAM role ARN that will be attached to master instances
func (o *InstanceIAMRoles) GetMasterRoleARN() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.masterRoleARN
	}
	return
}

// WorkerRoleARN returns the value of the 'worker_role_ARN' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The IAM role ARN that will be attached to worker instances
func (o *InstanceIAMRoles) WorkerRoleARN() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.workerRoleARN
	}
	return ""
}

// GetWorkerRoleARN returns the value of the 'worker_role_ARN' attribute and
// a flag indicating if the attribute has a value.
//
// The IAM role ARN that will be attached to worker instances
func (o *InstanceIAMRoles) GetWorkerRoleARN() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.workerRoleARN
	}
	return
}

// InstanceIAMRolesListKind is the name of the type used to represent list of objects of
// type 'instance_IAM_roles'.
const InstanceIAMRolesListKind = "InstanceIAMRolesList"

// InstanceIAMRolesListLinkKind is the name of the type used to represent links to list
// of objects of type 'instance_IAM_roles'.
const InstanceIAMRolesListLinkKind = "InstanceIAMRolesListLink"

// InstanceIAMRolesNilKind is the name of the type used to nil lists of objects of
// type 'instance_IAM_roles'.
const InstanceIAMRolesListNilKind = "InstanceIAMRolesListNil"

// InstanceIAMRolesList is a list of values of the 'instance_IAM_roles' type.
type InstanceIAMRolesList struct {
	href  string
	link  bool
	items []*InstanceIAMRoles
}

// Len returns the length of the list.
func (l *InstanceIAMRolesList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *InstanceIAMRolesList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *InstanceIAMRolesList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *InstanceIAMRolesList) SetItems(items []*InstanceIAMRoles) {
	l.items = items
}

// Items returns the items of the list.
func (l *InstanceIAMRolesList) Items() []*InstanceIAMRoles {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *InstanceIAMRolesList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *InstanceIAMRolesList) Get(i int) *InstanceIAMRoles {
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
func (l *InstanceIAMRolesList) Slice() []*InstanceIAMRoles {
	var slice []*InstanceIAMRoles
	if l == nil {
		slice = make([]*InstanceIAMRoles, 0)
	} else {
		slice = make([]*InstanceIAMRoles, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *InstanceIAMRolesList) Each(f func(item *InstanceIAMRoles) bool) {
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
func (l *InstanceIAMRolesList) Range(f func(index int, item *InstanceIAMRoles) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
