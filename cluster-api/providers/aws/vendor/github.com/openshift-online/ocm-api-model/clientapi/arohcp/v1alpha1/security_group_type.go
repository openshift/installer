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

// SecurityGroup represents the values of the 'security_group' type.
//
// AWS security group object
type SecurityGroup struct {
	fieldSet_     []bool
	id            string
	name          string
	redHatManaged bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *SecurityGroup) Empty() bool {
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

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The security group ID.
func (o *SecurityGroup) ID() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// The security group ID.
func (o *SecurityGroup) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.id
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the security group according to its `Name` tag on AWS.
func (o *SecurityGroup) Name() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the security group according to its `Name` tag on AWS.
func (o *SecurityGroup) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.name
	}
	return
}

// RedHatManaged returns the value of the 'red_hat_managed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// If the resource is RH managed.
func (o *SecurityGroup) RedHatManaged() bool {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.redHatManaged
	}
	return false
}

// GetRedHatManaged returns the value of the 'red_hat_managed' attribute and
// a flag indicating if the attribute has a value.
//
// If the resource is RH managed.
func (o *SecurityGroup) GetRedHatManaged() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.redHatManaged
	}
	return
}

// SecurityGroupListKind is the name of the type used to represent list of objects of
// type 'security_group'.
const SecurityGroupListKind = "SecurityGroupList"

// SecurityGroupListLinkKind is the name of the type used to represent links to list
// of objects of type 'security_group'.
const SecurityGroupListLinkKind = "SecurityGroupListLink"

// SecurityGroupNilKind is the name of the type used to nil lists of objects of
// type 'security_group'.
const SecurityGroupListNilKind = "SecurityGroupListNil"

// SecurityGroupList is a list of values of the 'security_group' type.
type SecurityGroupList struct {
	href  string
	link  bool
	items []*SecurityGroup
}

// Len returns the length of the list.
func (l *SecurityGroupList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *SecurityGroupList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *SecurityGroupList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *SecurityGroupList) SetItems(items []*SecurityGroup) {
	l.items = items
}

// Items returns the items of the list.
func (l *SecurityGroupList) Items() []*SecurityGroup {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *SecurityGroupList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *SecurityGroupList) Get(i int) *SecurityGroup {
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
func (l *SecurityGroupList) Slice() []*SecurityGroup {
	var slice []*SecurityGroup
	if l == nil {
		slice = make([]*SecurityGroup, 0)
	} else {
		slice = make([]*SecurityGroup, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *SecurityGroupList) Each(f func(item *SecurityGroup) bool) {
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
func (l *SecurityGroupList) Range(f func(index int, item *SecurityGroup) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
