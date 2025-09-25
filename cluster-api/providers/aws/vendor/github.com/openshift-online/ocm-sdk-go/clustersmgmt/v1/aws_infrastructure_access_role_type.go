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

// AWSInfrastructureAccessRoleKind is the name of the type used to represent objects
// of type 'AWS_infrastructure_access_role'.
const AWSInfrastructureAccessRoleKind = "AWSInfrastructureAccessRole"

// AWSInfrastructureAccessRoleLinkKind is the name of the type used to represent links
// to objects of type 'AWS_infrastructure_access_role'.
const AWSInfrastructureAccessRoleLinkKind = "AWSInfrastructureAccessRoleLink"

// AWSInfrastructureAccessRoleNilKind is the name of the type used to nil references
// to objects of type 'AWS_infrastructure_access_role'.
const AWSInfrastructureAccessRoleNilKind = "AWSInfrastructureAccessRoleNil"

// AWSInfrastructureAccessRole represents the values of the 'AWS_infrastructure_access_role' type.
//
// A set of acces permissions for AWS resources
type AWSInfrastructureAccessRole struct {
	bitmap_     uint32
	id          string
	href        string
	description string
	displayName string
	state       AWSInfrastructureAccessRoleState
}

// Kind returns the name of the type of the object.
func (o *AWSInfrastructureAccessRole) Kind() string {
	if o == nil {
		return AWSInfrastructureAccessRoleNilKind
	}
	if o.bitmap_&1 != 0 {
		return AWSInfrastructureAccessRoleLinkKind
	}
	return AWSInfrastructureAccessRoleKind
}

// Link returns true if this is a link.
func (o *AWSInfrastructureAccessRole) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *AWSInfrastructureAccessRole) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AWSInfrastructureAccessRole) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AWSInfrastructureAccessRole) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AWSInfrastructureAccessRole) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWSInfrastructureAccessRole) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Description returns the value of the 'description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Description of the role.
func (o *AWSInfrastructureAccessRole) Description() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.description
	}
	return ""
}

// GetDescription returns the value of the 'description' attribute and
// a flag indicating if the attribute has a value.
//
// Description of the role.
func (o *AWSInfrastructureAccessRole) GetDescription() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.description
	}
	return
}

// DisplayName returns the value of the 'display_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Human friendly identifier of the role, for example `Read only`.
func (o *AWSInfrastructureAccessRole) DisplayName() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.displayName
	}
	return ""
}

// GetDisplayName returns the value of the 'display_name' attribute and
// a flag indicating if the attribute has a value.
//
// Human friendly identifier of the role, for example `Read only`.
func (o *AWSInfrastructureAccessRole) GetDisplayName() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.displayName
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// State of the role.
func (o *AWSInfrastructureAccessRole) State() AWSInfrastructureAccessRoleState {
	if o != nil && o.bitmap_&32 != 0 {
		return o.state
	}
	return AWSInfrastructureAccessRoleState("")
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
//
// State of the role.
func (o *AWSInfrastructureAccessRole) GetState() (value AWSInfrastructureAccessRoleState, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.state
	}
	return
}

// AWSInfrastructureAccessRoleListKind is the name of the type used to represent list of objects of
// type 'AWS_infrastructure_access_role'.
const AWSInfrastructureAccessRoleListKind = "AWSInfrastructureAccessRoleList"

// AWSInfrastructureAccessRoleListLinkKind is the name of the type used to represent links to list
// of objects of type 'AWS_infrastructure_access_role'.
const AWSInfrastructureAccessRoleListLinkKind = "AWSInfrastructureAccessRoleListLink"

// AWSInfrastructureAccessRoleNilKind is the name of the type used to nil lists of objects of
// type 'AWS_infrastructure_access_role'.
const AWSInfrastructureAccessRoleListNilKind = "AWSInfrastructureAccessRoleListNil"

// AWSInfrastructureAccessRoleList is a list of values of the 'AWS_infrastructure_access_role' type.
type AWSInfrastructureAccessRoleList struct {
	href  string
	link  bool
	items []*AWSInfrastructureAccessRole
}

// Kind returns the name of the type of the object.
func (l *AWSInfrastructureAccessRoleList) Kind() string {
	if l == nil {
		return AWSInfrastructureAccessRoleListNilKind
	}
	if l.link {
		return AWSInfrastructureAccessRoleListLinkKind
	}
	return AWSInfrastructureAccessRoleListKind
}

// Link returns true iif this is a link.
func (l *AWSInfrastructureAccessRoleList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AWSInfrastructureAccessRoleList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AWSInfrastructureAccessRoleList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AWSInfrastructureAccessRoleList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AWSInfrastructureAccessRoleList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AWSInfrastructureAccessRoleList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AWSInfrastructureAccessRoleList) SetItems(items []*AWSInfrastructureAccessRole) {
	l.items = items
}

// Items returns the items of the list.
func (l *AWSInfrastructureAccessRoleList) Items() []*AWSInfrastructureAccessRole {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AWSInfrastructureAccessRoleList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AWSInfrastructureAccessRoleList) Get(i int) *AWSInfrastructureAccessRole {
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
func (l *AWSInfrastructureAccessRoleList) Slice() []*AWSInfrastructureAccessRole {
	var slice []*AWSInfrastructureAccessRole
	if l == nil {
		slice = make([]*AWSInfrastructureAccessRole, 0)
	} else {
		slice = make([]*AWSInfrastructureAccessRole, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AWSInfrastructureAccessRoleList) Each(f func(item *AWSInfrastructureAccessRole) bool) {
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
func (l *AWSInfrastructureAccessRoleList) Range(f func(index int, item *AWSInfrastructureAccessRole) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
