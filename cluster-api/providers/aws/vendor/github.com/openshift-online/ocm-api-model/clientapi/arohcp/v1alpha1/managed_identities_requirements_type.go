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

// ManagedIdentitiesRequirementsKind is the name of the type used to represent objects
// of type 'managed_identities_requirements'.
const ManagedIdentitiesRequirementsKind = "ManagedIdentitiesRequirements"

// ManagedIdentitiesRequirementsLinkKind is the name of the type used to represent links
// to objects of type 'managed_identities_requirements'.
const ManagedIdentitiesRequirementsLinkKind = "ManagedIdentitiesRequirementsLink"

// ManagedIdentitiesRequirementsNilKind is the name of the type used to nil references
// to objects of type 'managed_identities_requirements'.
const ManagedIdentitiesRequirementsNilKind = "ManagedIdentitiesRequirementsNil"

// ManagedIdentitiesRequirements represents the values of the 'managed_identities_requirements' type.
//
// Representation of managed identities requirements.
// When creating ARO-HCP Clusters, the end-users will need to pre-create the set of Managed Identities
// required by the clusters.
// The set of Managed Identities that the end-users need to precreate is not static and depends on
// several factors:
// (1) The OpenShift version of the cluster being created.
// (2) The functionalities that are being enabled for the cluster. Some Managed Identities are not
// always required but become required if a given functionality is enabled.
// Additionally, the Managed Identities that the end-users will need to precreate will have to have a
// set of required permissions assigned to them which also have to be returned to the end users.
type ManagedIdentitiesRequirements struct {
	fieldSet_                       []bool
	id                              string
	href                            string
	controlPlaneOperatorsIdentities []*ControlPlaneOperatorIdentityRequirement
	dataPlaneOperatorsIdentities    []*DataPlaneOperatorIdentityRequirement
}

// Kind returns the name of the type of the object.
func (o *ManagedIdentitiesRequirements) Kind() string {
	if o == nil {
		return ManagedIdentitiesRequirementsNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ManagedIdentitiesRequirementsLinkKind
	}
	return ManagedIdentitiesRequirementsKind
}

// Link returns true if this is a link.
func (o *ManagedIdentitiesRequirements) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *ManagedIdentitiesRequirements) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ManagedIdentitiesRequirements) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ManagedIdentitiesRequirements) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ManagedIdentitiesRequirements) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ManagedIdentitiesRequirements) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}

	// Check all fields except the link flag (index 0)
	for i := 1; i < len(o.fieldSet_); i++ {
		if o.fieldSet_[i] {
			return false
		}
	}
	return true
}

// ControlPlaneOperatorsIdentities returns the value of the 'control_plane_operators_identities' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The control plane operators managed identities requirements
func (o *ManagedIdentitiesRequirements) ControlPlaneOperatorsIdentities() []*ControlPlaneOperatorIdentityRequirement {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.controlPlaneOperatorsIdentities
	}
	return nil
}

// GetControlPlaneOperatorsIdentities returns the value of the 'control_plane_operators_identities' attribute and
// a flag indicating if the attribute has a value.
//
// The control plane operators managed identities requirements
func (o *ManagedIdentitiesRequirements) GetControlPlaneOperatorsIdentities() (value []*ControlPlaneOperatorIdentityRequirement, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.controlPlaneOperatorsIdentities
	}
	return
}

// DataPlaneOperatorsIdentities returns the value of the 'data_plane_operators_identities' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The data plane operators managed identities requirements
func (o *ManagedIdentitiesRequirements) DataPlaneOperatorsIdentities() []*DataPlaneOperatorIdentityRequirement {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.dataPlaneOperatorsIdentities
	}
	return nil
}

// GetDataPlaneOperatorsIdentities returns the value of the 'data_plane_operators_identities' attribute and
// a flag indicating if the attribute has a value.
//
// The data plane operators managed identities requirements
func (o *ManagedIdentitiesRequirements) GetDataPlaneOperatorsIdentities() (value []*DataPlaneOperatorIdentityRequirement, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.dataPlaneOperatorsIdentities
	}
	return
}

// ManagedIdentitiesRequirementsListKind is the name of the type used to represent list of objects of
// type 'managed_identities_requirements'.
const ManagedIdentitiesRequirementsListKind = "ManagedIdentitiesRequirementsList"

// ManagedIdentitiesRequirementsListLinkKind is the name of the type used to represent links to list
// of objects of type 'managed_identities_requirements'.
const ManagedIdentitiesRequirementsListLinkKind = "ManagedIdentitiesRequirementsListLink"

// ManagedIdentitiesRequirementsNilKind is the name of the type used to nil lists of objects of
// type 'managed_identities_requirements'.
const ManagedIdentitiesRequirementsListNilKind = "ManagedIdentitiesRequirementsListNil"

// ManagedIdentitiesRequirementsList is a list of values of the 'managed_identities_requirements' type.
type ManagedIdentitiesRequirementsList struct {
	href  string
	link  bool
	items []*ManagedIdentitiesRequirements
}

// Kind returns the name of the type of the object.
func (l *ManagedIdentitiesRequirementsList) Kind() string {
	if l == nil {
		return ManagedIdentitiesRequirementsListNilKind
	}
	if l.link {
		return ManagedIdentitiesRequirementsListLinkKind
	}
	return ManagedIdentitiesRequirementsListKind
}

// Link returns true iif this is a link.
func (l *ManagedIdentitiesRequirementsList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ManagedIdentitiesRequirementsList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ManagedIdentitiesRequirementsList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ManagedIdentitiesRequirementsList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ManagedIdentitiesRequirementsList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ManagedIdentitiesRequirementsList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ManagedIdentitiesRequirementsList) SetItems(items []*ManagedIdentitiesRequirements) {
	l.items = items
}

// Items returns the items of the list.
func (l *ManagedIdentitiesRequirementsList) Items() []*ManagedIdentitiesRequirements {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ManagedIdentitiesRequirementsList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ManagedIdentitiesRequirementsList) Get(i int) *ManagedIdentitiesRequirements {
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
func (l *ManagedIdentitiesRequirementsList) Slice() []*ManagedIdentitiesRequirements {
	var slice []*ManagedIdentitiesRequirements
	if l == nil {
		slice = make([]*ManagedIdentitiesRequirements, 0)
	} else {
		slice = make([]*ManagedIdentitiesRequirements, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ManagedIdentitiesRequirementsList) Each(f func(item *ManagedIdentitiesRequirements) bool) {
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
func (l *ManagedIdentitiesRequirementsList) Range(f func(index int, item *ManagedIdentitiesRequirements) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
