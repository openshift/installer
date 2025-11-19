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

// ControlPlaneOperatorIdentityRequirement represents the values of the 'control_plane_operator_identity_requirement' type.
type ControlPlaneOperatorIdentityRequirement struct {
	fieldSet_           []bool
	maxOpenShiftVersion string
	minOpenShiftVersion string
	operatorName        string
	required            string
	roleDefinitions     []*RoleDefinitionOperatorIdentityRequirement
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ControlPlaneOperatorIdentityRequirement) Empty() bool {
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

// MaxOpenShiftVersion returns the value of the 'max_open_shift_version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The field is a string and it is of format X.Y.
// Not specifying it indicates support for this operator in all Openshift versions,
// starting from min_openshift_version if min_openshift_version is defined.
func (o *ControlPlaneOperatorIdentityRequirement) MaxOpenShiftVersion() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.maxOpenShiftVersion
	}
	return ""
}

// GetMaxOpenShiftVersion returns the value of the 'max_open_shift_version' attribute and
// a flag indicating if the attribute has a value.
//
// The field is a string and it is of format X.Y.
// Not specifying it indicates support for this operator in all Openshift versions,
// starting from min_openshift_version if min_openshift_version is defined.
func (o *ControlPlaneOperatorIdentityRequirement) GetMaxOpenShiftVersion() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.maxOpenShiftVersion
	}
	return
}

// MinOpenShiftVersion returns the value of the 'min_open_shift_version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The field is a string and it is of format X.Y.
// Not specifying it indicates support for this operator in all Openshift versions,
// or up to max_openshift_version, if defined.
func (o *ControlPlaneOperatorIdentityRequirement) MinOpenShiftVersion() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.minOpenShiftVersion
	}
	return ""
}

// GetMinOpenShiftVersion returns the value of the 'min_open_shift_version' attribute and
// a flag indicating if the attribute has a value.
//
// The field is a string and it is of format X.Y.
// Not specifying it indicates support for this operator in all Openshift versions,
// or up to max_openshift_version, if defined.
func (o *ControlPlaneOperatorIdentityRequirement) GetMinOpenShiftVersion() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.minOpenShiftVersion
	}
	return
}

// OperatorName returns the value of the 'operator_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the control plane operator that needs the identity
func (o *ControlPlaneOperatorIdentityRequirement) OperatorName() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.operatorName
	}
	return ""
}

// GetOperatorName returns the value of the 'operator_name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the control plane operator that needs the identity
func (o *ControlPlaneOperatorIdentityRequirement) GetOperatorName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.operatorName
	}
	return
}

// Required returns the value of the 'required' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates whether the identity is always required or not.
// "always" means that the identity is always required
// "on_enablement" means that the identity is only required when a functionality
// that leverages the operator is enabled.
// Possible values are ("always", "on_enablement")
func (o *ControlPlaneOperatorIdentityRequirement) Required() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.required
	}
	return ""
}

// GetRequired returns the value of the 'required' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates whether the identity is always required or not.
// "always" means that the identity is always required
// "on_enablement" means that the identity is only required when a functionality
// that leverages the operator is enabled.
// Possible values are ("always", "on_enablement")
func (o *ControlPlaneOperatorIdentityRequirement) GetRequired() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.required
	}
	return
}

// RoleDefinitions returns the value of the 'role_definitions' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// A list of roles that are required by the operator
func (o *ControlPlaneOperatorIdentityRequirement) RoleDefinitions() []*RoleDefinitionOperatorIdentityRequirement {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.roleDefinitions
	}
	return nil
}

// GetRoleDefinitions returns the value of the 'role_definitions' attribute and
// a flag indicating if the attribute has a value.
//
// A list of roles that are required by the operator
func (o *ControlPlaneOperatorIdentityRequirement) GetRoleDefinitions() (value []*RoleDefinitionOperatorIdentityRequirement, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.roleDefinitions
	}
	return
}

// ControlPlaneOperatorIdentityRequirementListKind is the name of the type used to represent list of objects of
// type 'control_plane_operator_identity_requirement'.
const ControlPlaneOperatorIdentityRequirementListKind = "ControlPlaneOperatorIdentityRequirementList"

// ControlPlaneOperatorIdentityRequirementListLinkKind is the name of the type used to represent links to list
// of objects of type 'control_plane_operator_identity_requirement'.
const ControlPlaneOperatorIdentityRequirementListLinkKind = "ControlPlaneOperatorIdentityRequirementListLink"

// ControlPlaneOperatorIdentityRequirementNilKind is the name of the type used to nil lists of objects of
// type 'control_plane_operator_identity_requirement'.
const ControlPlaneOperatorIdentityRequirementListNilKind = "ControlPlaneOperatorIdentityRequirementListNil"

// ControlPlaneOperatorIdentityRequirementList is a list of values of the 'control_plane_operator_identity_requirement' type.
type ControlPlaneOperatorIdentityRequirementList struct {
	href  string
	link  bool
	items []*ControlPlaneOperatorIdentityRequirement
}

// Len returns the length of the list.
func (l *ControlPlaneOperatorIdentityRequirementList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ControlPlaneOperatorIdentityRequirementList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ControlPlaneOperatorIdentityRequirementList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ControlPlaneOperatorIdentityRequirementList) SetItems(items []*ControlPlaneOperatorIdentityRequirement) {
	l.items = items
}

// Items returns the items of the list.
func (l *ControlPlaneOperatorIdentityRequirementList) Items() []*ControlPlaneOperatorIdentityRequirement {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ControlPlaneOperatorIdentityRequirementList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ControlPlaneOperatorIdentityRequirementList) Get(i int) *ControlPlaneOperatorIdentityRequirement {
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
func (l *ControlPlaneOperatorIdentityRequirementList) Slice() []*ControlPlaneOperatorIdentityRequirement {
	var slice []*ControlPlaneOperatorIdentityRequirement
	if l == nil {
		slice = make([]*ControlPlaneOperatorIdentityRequirement, 0)
	} else {
		slice = make([]*ControlPlaneOperatorIdentityRequirement, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ControlPlaneOperatorIdentityRequirementList) Each(f func(item *ControlPlaneOperatorIdentityRequirement) bool) {
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
func (l *ControlPlaneOperatorIdentityRequirementList) Range(f func(index int, item *ControlPlaneOperatorIdentityRequirement) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
