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

// AWSInfrastructureAccessRoleGrantKind is the name of the type used to represent objects
// of type 'AWS_infrastructure_access_role_grant'.
const AWSInfrastructureAccessRoleGrantKind = "AWSInfrastructureAccessRoleGrant"

// AWSInfrastructureAccessRoleGrantLinkKind is the name of the type used to represent links
// to objects of type 'AWS_infrastructure_access_role_grant'.
const AWSInfrastructureAccessRoleGrantLinkKind = "AWSInfrastructureAccessRoleGrantLink"

// AWSInfrastructureAccessRoleGrantNilKind is the name of the type used to nil references
// to objects of type 'AWS_infrastructure_access_role_grant'.
const AWSInfrastructureAccessRoleGrantNilKind = "AWSInfrastructureAccessRoleGrantNil"

// AWSInfrastructureAccessRoleGrant represents the values of the 'AWS_infrastructure_access_role_grant' type.
//
// Representation of an AWS infrastructure access role grant.
type AWSInfrastructureAccessRoleGrant struct {
	bitmap_          uint32
	id               string
	href             string
	consoleURL       string
	role             *AWSInfrastructureAccessRole
	state            AWSInfrastructureAccessRoleGrantState
	stateDescription string
	userARN          string
}

// Kind returns the name of the type of the object.
func (o *AWSInfrastructureAccessRoleGrant) Kind() string {
	if o == nil {
		return AWSInfrastructureAccessRoleGrantNilKind
	}
	if o.bitmap_&1 != 0 {
		return AWSInfrastructureAccessRoleGrantLinkKind
	}
	return AWSInfrastructureAccessRoleGrantKind
}

// Link returns true if this is a link.
func (o *AWSInfrastructureAccessRoleGrant) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *AWSInfrastructureAccessRoleGrant) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *AWSInfrastructureAccessRoleGrant) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *AWSInfrastructureAccessRoleGrant) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *AWSInfrastructureAccessRoleGrant) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AWSInfrastructureAccessRoleGrant) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// ConsoleURL returns the value of the 'console_URL' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// URL to switch to the role in AWS console.
func (o *AWSInfrastructureAccessRoleGrant) ConsoleURL() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.consoleURL
	}
	return ""
}

// GetConsoleURL returns the value of the 'console_URL' attribute and
// a flag indicating if the attribute has a value.
//
// URL to switch to the role in AWS console.
func (o *AWSInfrastructureAccessRoleGrant) GetConsoleURL() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.consoleURL
	}
	return
}

// Role returns the value of the 'role' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to AWS infrastructure access role.
// Grant must use a 'valid' role.
func (o *AWSInfrastructureAccessRoleGrant) Role() *AWSInfrastructureAccessRole {
	if o != nil && o.bitmap_&16 != 0 {
		return o.role
	}
	return nil
}

// GetRole returns the value of the 'role' attribute and
// a flag indicating if the attribute has a value.
//
// Link to AWS infrastructure access role.
// Grant must use a 'valid' role.
func (o *AWSInfrastructureAccessRoleGrant) GetRole() (value *AWSInfrastructureAccessRole, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.role
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// State of the grant.
func (o *AWSInfrastructureAccessRoleGrant) State() AWSInfrastructureAccessRoleGrantState {
	if o != nil && o.bitmap_&32 != 0 {
		return o.state
	}
	return AWSInfrastructureAccessRoleGrantState("")
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
//
// State of the grant.
func (o *AWSInfrastructureAccessRoleGrant) GetState() (value AWSInfrastructureAccessRoleGrantState, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.state
	}
	return
}

// StateDescription returns the value of the 'state_description' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Description of the state.
// Will be empty unless state is 'Failed'.
func (o *AWSInfrastructureAccessRoleGrant) StateDescription() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.stateDescription
	}
	return ""
}

// GetStateDescription returns the value of the 'state_description' attribute and
// a flag indicating if the attribute has a value.
//
// Description of the state.
// Will be empty unless state is 'Failed'.
func (o *AWSInfrastructureAccessRoleGrant) GetStateDescription() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.stateDescription
	}
	return
}

// UserARN returns the value of the 'user_ARN' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The user AWS IAM ARN we want to grant the role.
func (o *AWSInfrastructureAccessRoleGrant) UserARN() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.userARN
	}
	return ""
}

// GetUserARN returns the value of the 'user_ARN' attribute and
// a flag indicating if the attribute has a value.
//
// The user AWS IAM ARN we want to grant the role.
func (o *AWSInfrastructureAccessRoleGrant) GetUserARN() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.userARN
	}
	return
}

// AWSInfrastructureAccessRoleGrantListKind is the name of the type used to represent list of objects of
// type 'AWS_infrastructure_access_role_grant'.
const AWSInfrastructureAccessRoleGrantListKind = "AWSInfrastructureAccessRoleGrantList"

// AWSInfrastructureAccessRoleGrantListLinkKind is the name of the type used to represent links to list
// of objects of type 'AWS_infrastructure_access_role_grant'.
const AWSInfrastructureAccessRoleGrantListLinkKind = "AWSInfrastructureAccessRoleGrantListLink"

// AWSInfrastructureAccessRoleGrantNilKind is the name of the type used to nil lists of objects of
// type 'AWS_infrastructure_access_role_grant'.
const AWSInfrastructureAccessRoleGrantListNilKind = "AWSInfrastructureAccessRoleGrantListNil"

// AWSInfrastructureAccessRoleGrantList is a list of values of the 'AWS_infrastructure_access_role_grant' type.
type AWSInfrastructureAccessRoleGrantList struct {
	href  string
	link  bool
	items []*AWSInfrastructureAccessRoleGrant
}

// Kind returns the name of the type of the object.
func (l *AWSInfrastructureAccessRoleGrantList) Kind() string {
	if l == nil {
		return AWSInfrastructureAccessRoleGrantListNilKind
	}
	if l.link {
		return AWSInfrastructureAccessRoleGrantListLinkKind
	}
	return AWSInfrastructureAccessRoleGrantListKind
}

// Link returns true iif this is a link.
func (l *AWSInfrastructureAccessRoleGrantList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *AWSInfrastructureAccessRoleGrantList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *AWSInfrastructureAccessRoleGrantList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *AWSInfrastructureAccessRoleGrantList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AWSInfrastructureAccessRoleGrantList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AWSInfrastructureAccessRoleGrantList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AWSInfrastructureAccessRoleGrantList) SetItems(items []*AWSInfrastructureAccessRoleGrant) {
	l.items = items
}

// Items returns the items of the list.
func (l *AWSInfrastructureAccessRoleGrantList) Items() []*AWSInfrastructureAccessRoleGrant {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AWSInfrastructureAccessRoleGrantList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AWSInfrastructureAccessRoleGrantList) Get(i int) *AWSInfrastructureAccessRoleGrant {
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
func (l *AWSInfrastructureAccessRoleGrantList) Slice() []*AWSInfrastructureAccessRoleGrant {
	var slice []*AWSInfrastructureAccessRoleGrant
	if l == nil {
		slice = make([]*AWSInfrastructureAccessRoleGrant, 0)
	} else {
		slice = make([]*AWSInfrastructureAccessRoleGrant, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AWSInfrastructureAccessRoleGrantList) Each(f func(item *AWSInfrastructureAccessRoleGrant) bool) {
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
func (l *AWSInfrastructureAccessRoleGrantList) Range(f func(index int, item *AWSInfrastructureAccessRoleGrant) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
