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

// DataPlaneOperatorIdentityRequirement represents the values of the 'data_plane_operator_identity_requirement' type.
type DataPlaneOperatorIdentityRequirement struct {
	fieldSet_           []bool
	maxOpenShiftVersion string
	minOpenShiftVersion string
	operatorName        string
	required            string
	roleDefinitions     []*RoleDefinitionOperatorIdentityRequirement
	serviceAccounts     []*K8sServiceAccountOperatorIdentityRequirement
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *DataPlaneOperatorIdentityRequirement) Empty() bool {
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
// The field is a string and it is of format X.Y (e.g 4.18) where X and Y are major and
// minor segments of the OpenShift version respectively.
// Not specifying it indicates support for this operator in all Openshift versions,
// starting from min_openshift_version if min_openshift_version is defined.
func (o *DataPlaneOperatorIdentityRequirement) MaxOpenShiftVersion() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.maxOpenShiftVersion
	}
	return ""
}

// GetMaxOpenShiftVersion returns the value of the 'max_open_shift_version' attribute and
// a flag indicating if the attribute has a value.
//
// The field is a string and it is of format X.Y (e.g 4.18) where X and Y are major and
// minor segments of the OpenShift version respectively.
// Not specifying it indicates support for this operator in all Openshift versions,
// starting from min_openshift_version if min_openshift_version is defined.
func (o *DataPlaneOperatorIdentityRequirement) GetMaxOpenShiftVersion() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.maxOpenShiftVersion
	}
	return
}

// MinOpenShiftVersion returns the value of the 'min_open_shift_version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The field is a string and it is of format X.Y (e.g 4.18) where X and Y are major and
// minor segments of the OpenShift version respectively.
// Not specifying it indicates support for this operator in all Openshift versions,
// or up to max_openshift_version, if defined.
func (o *DataPlaneOperatorIdentityRequirement) MinOpenShiftVersion() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.minOpenShiftVersion
	}
	return ""
}

// GetMinOpenShiftVersion returns the value of the 'min_open_shift_version' attribute and
// a flag indicating if the attribute has a value.
//
// The field is a string and it is of format X.Y (e.g 4.18) where X and Y are major and
// minor segments of the OpenShift version respectively.
// Not specifying it indicates support for this operator in all Openshift versions,
// or up to max_openshift_version, if defined.
func (o *DataPlaneOperatorIdentityRequirement) GetMinOpenShiftVersion() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.minOpenShiftVersion
	}
	return
}

// OperatorName returns the value of the 'operator_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the data plane operator that needs the identity
func (o *DataPlaneOperatorIdentityRequirement) OperatorName() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.operatorName
	}
	return ""
}

// GetOperatorName returns the value of the 'operator_name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the data plane operator that needs the identity
func (o *DataPlaneOperatorIdentityRequirement) GetOperatorName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.operatorName
	}
	return
}

// Required returns the value of the 'required' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates whether the identity is always required or not
// "always" means that the identity is always required
// "on_enablement" means that the identity is only required when a functionality
// that leverages the operator is enabled.
// Possible values are ("always", "on_enablement")
func (o *DataPlaneOperatorIdentityRequirement) Required() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.required
	}
	return ""
}

// GetRequired returns the value of the 'required' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates whether the identity is always required or not
// "always" means that the identity is always required
// "on_enablement" means that the identity is only required when a functionality
// that leverages the operator is enabled.
// Possible values are ("always", "on_enablement")
func (o *DataPlaneOperatorIdentityRequirement) GetRequired() (value string, ok bool) {
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
func (o *DataPlaneOperatorIdentityRequirement) RoleDefinitions() []*RoleDefinitionOperatorIdentityRequirement {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.roleDefinitions
	}
	return nil
}

// GetRoleDefinitions returns the value of the 'role_definitions' attribute and
// a flag indicating if the attribute has a value.
//
// A list of roles that are required by the operator
func (o *DataPlaneOperatorIdentityRequirement) GetRoleDefinitions() (value []*RoleDefinitionOperatorIdentityRequirement, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.roleDefinitions
	}
	return
}

// ServiceAccounts returns the value of the 'service_accounts' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// It is a list of K8s ServiceAccounts leveraged by the operator.
// There must be at least a single service account specified.
// This information is needed to federate a managed identity to a k8s subject.
// There should be no duplicated "name:namespace" entries within this field.
func (o *DataPlaneOperatorIdentityRequirement) ServiceAccounts() []*K8sServiceAccountOperatorIdentityRequirement {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.serviceAccounts
	}
	return nil
}

// GetServiceAccounts returns the value of the 'service_accounts' attribute and
// a flag indicating if the attribute has a value.
//
// It is a list of K8s ServiceAccounts leveraged by the operator.
// There must be at least a single service account specified.
// This information is needed to federate a managed identity to a k8s subject.
// There should be no duplicated "name:namespace" entries within this field.
func (o *DataPlaneOperatorIdentityRequirement) GetServiceAccounts() (value []*K8sServiceAccountOperatorIdentityRequirement, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.serviceAccounts
	}
	return
}

// DataPlaneOperatorIdentityRequirementListKind is the name of the type used to represent list of objects of
// type 'data_plane_operator_identity_requirement'.
const DataPlaneOperatorIdentityRequirementListKind = "DataPlaneOperatorIdentityRequirementList"

// DataPlaneOperatorIdentityRequirementListLinkKind is the name of the type used to represent links to list
// of objects of type 'data_plane_operator_identity_requirement'.
const DataPlaneOperatorIdentityRequirementListLinkKind = "DataPlaneOperatorIdentityRequirementListLink"

// DataPlaneOperatorIdentityRequirementNilKind is the name of the type used to nil lists of objects of
// type 'data_plane_operator_identity_requirement'.
const DataPlaneOperatorIdentityRequirementListNilKind = "DataPlaneOperatorIdentityRequirementListNil"

// DataPlaneOperatorIdentityRequirementList is a list of values of the 'data_plane_operator_identity_requirement' type.
type DataPlaneOperatorIdentityRequirementList struct {
	href  string
	link  bool
	items []*DataPlaneOperatorIdentityRequirement
}

// Len returns the length of the list.
func (l *DataPlaneOperatorIdentityRequirementList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *DataPlaneOperatorIdentityRequirementList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *DataPlaneOperatorIdentityRequirementList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *DataPlaneOperatorIdentityRequirementList) SetItems(items []*DataPlaneOperatorIdentityRequirement) {
	l.items = items
}

// Items returns the items of the list.
func (l *DataPlaneOperatorIdentityRequirementList) Items() []*DataPlaneOperatorIdentityRequirement {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *DataPlaneOperatorIdentityRequirementList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *DataPlaneOperatorIdentityRequirementList) Get(i int) *DataPlaneOperatorIdentityRequirement {
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
func (l *DataPlaneOperatorIdentityRequirementList) Slice() []*DataPlaneOperatorIdentityRequirement {
	var slice []*DataPlaneOperatorIdentityRequirement
	if l == nil {
		slice = make([]*DataPlaneOperatorIdentityRequirement, 0)
	} else {
		slice = make([]*DataPlaneOperatorIdentityRequirement, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *DataPlaneOperatorIdentityRequirementList) Each(f func(item *DataPlaneOperatorIdentityRequirement) bool) {
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
func (l *DataPlaneOperatorIdentityRequirementList) Range(f func(index int, item *DataPlaneOperatorIdentityRequirement) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
