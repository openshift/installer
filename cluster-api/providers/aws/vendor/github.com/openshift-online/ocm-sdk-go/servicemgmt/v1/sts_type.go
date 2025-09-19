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

// STS represents the values of the 'STS' type.
//
// Contains the necessary attributes to support role-based authentication on AWS.
type STS struct {
	bitmap_            uint32
	oidcEndpointURL    string
	instanceIAMRoles   *InstanceIAMRoles
	operatorIAMRoles   []*OperatorIAMRole
	operatorRolePrefix string
	roleARN            string
	supportRoleARN     string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *STS) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// OIDCEndpointURL returns the value of the 'OIDC_endpoint_URL' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// URL of the location where OIDC configuration and keys are available
func (o *STS) OIDCEndpointURL() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.oidcEndpointURL
	}
	return ""
}

// GetOIDCEndpointURL returns the value of the 'OIDC_endpoint_URL' attribute and
// a flag indicating if the attribute has a value.
//
// URL of the location where OIDC configuration and keys are available
func (o *STS) GetOIDCEndpointURL() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.oidcEndpointURL
	}
	return
}

// InstanceIAMRoles returns the value of the 'instance_IAM_roles' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Instance IAM roles to use for the instance profiles of the master and worker instances
func (o *STS) InstanceIAMRoles() *InstanceIAMRoles {
	if o != nil && o.bitmap_&2 != 0 {
		return o.instanceIAMRoles
	}
	return nil
}

// GetInstanceIAMRoles returns the value of the 'instance_IAM_roles' attribute and
// a flag indicating if the attribute has a value.
//
// Instance IAM roles to use for the instance profiles of the master and worker instances
func (o *STS) GetInstanceIAMRoles() (value *InstanceIAMRoles, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.instanceIAMRoles
	}
	return
}

// OperatorIAMRoles returns the value of the 'operator_IAM_roles' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of roles necessary to access the AWS resources of the various operators used during installation
func (o *STS) OperatorIAMRoles() []*OperatorIAMRole {
	if o != nil && o.bitmap_&4 != 0 {
		return o.operatorIAMRoles
	}
	return nil
}

// GetOperatorIAMRoles returns the value of the 'operator_IAM_roles' attribute and
// a flag indicating if the attribute has a value.
//
// List of roles necessary to access the AWS resources of the various operators used during installation
func (o *STS) GetOperatorIAMRoles() (value []*OperatorIAMRole, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.operatorIAMRoles
	}
	return
}

// OperatorRolePrefix returns the value of the 'operator_role_prefix' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *STS) OperatorRolePrefix() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.operatorRolePrefix
	}
	return ""
}

// GetOperatorRolePrefix returns the value of the 'operator_role_prefix' attribute and
// a flag indicating if the attribute has a value.
func (o *STS) GetOperatorRolePrefix() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.operatorRolePrefix
	}
	return
}

// RoleARN returns the value of the 'role_ARN' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ARN of the AWS role to assume when installing the cluster
func (o *STS) RoleARN() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.roleARN
	}
	return ""
}

// GetRoleARN returns the value of the 'role_ARN' attribute and
// a flag indicating if the attribute has a value.
//
// ARN of the AWS role to assume when installing the cluster
func (o *STS) GetRoleARN() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.roleARN
	}
	return
}

// SupportRoleARN returns the value of the 'support_role_ARN' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ARN of the AWS role used by SREs to access the cluster AWS account in order to provide support
func (o *STS) SupportRoleARN() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.supportRoleARN
	}
	return ""
}

// GetSupportRoleARN returns the value of the 'support_role_ARN' attribute and
// a flag indicating if the attribute has a value.
//
// ARN of the AWS role used by SREs to access the cluster AWS account in order to provide support
func (o *STS) GetSupportRoleARN() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.supportRoleARN
	}
	return
}

// STSListKind is the name of the type used to represent list of objects of
// type 'STS'.
const STSListKind = "STSList"

// STSListLinkKind is the name of the type used to represent links to list
// of objects of type 'STS'.
const STSListLinkKind = "STSListLink"

// STSNilKind is the name of the type used to nil lists of objects of
// type 'STS'.
const STSListNilKind = "STSListNil"

// STSList is a list of values of the 'STS' type.
type STSList struct {
	href  string
	link  bool
	items []*STS
}

// Len returns the length of the list.
func (l *STSList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *STSList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *STSList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *STSList) SetItems(items []*STS) {
	l.items = items
}

// Items returns the items of the list.
func (l *STSList) Items() []*STS {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *STSList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *STSList) Get(i int) *STS {
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
func (l *STSList) Slice() []*STS {
	var slice []*STS
	if l == nil {
		slice = make([]*STS, 0)
	} else {
		slice = make([]*STS, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *STSList) Each(f func(item *STS) bool) {
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
func (l *STSList) Range(f func(index int, item *STS) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
