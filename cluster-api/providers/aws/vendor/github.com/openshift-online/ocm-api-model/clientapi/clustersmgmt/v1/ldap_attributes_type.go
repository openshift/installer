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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

// LDAPAttributes represents the values of the 'LDAP_attributes' type.
//
// LDAP attributes used to configure the LDAP identity provider.
type LDAPAttributes struct {
	fieldSet_         []bool
	id                []string
	email             []string
	name              []string
	preferredUsername []string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *LDAPAttributes) Empty() bool {
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
// List of attributes to use as the identity.
func (o *LDAPAttributes) ID() []string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.id
	}
	return nil
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// List of attributes to use as the identity.
func (o *LDAPAttributes) GetID() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.id
	}
	return
}

// Email returns the value of the 'email' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of attributes to use as the mail address.
func (o *LDAPAttributes) Email() []string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.email
	}
	return nil
}

// GetEmail returns the value of the 'email' attribute and
// a flag indicating if the attribute has a value.
//
// List of attributes to use as the mail address.
func (o *LDAPAttributes) GetEmail() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.email
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of attributes to use as the display name.
func (o *LDAPAttributes) Name() []string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.name
	}
	return nil
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// List of attributes to use as the display name.
func (o *LDAPAttributes) GetName() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.name
	}
	return
}

// PreferredUsername returns the value of the 'preferred_username' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of attributes to use as the preferred user name when provisioning a user.
func (o *LDAPAttributes) PreferredUsername() []string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.preferredUsername
	}
	return nil
}

// GetPreferredUsername returns the value of the 'preferred_username' attribute and
// a flag indicating if the attribute has a value.
//
// List of attributes to use as the preferred user name when provisioning a user.
func (o *LDAPAttributes) GetPreferredUsername() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.preferredUsername
	}
	return
}

// LDAPAttributesListKind is the name of the type used to represent list of objects of
// type 'LDAP_attributes'.
const LDAPAttributesListKind = "LDAPAttributesList"

// LDAPAttributesListLinkKind is the name of the type used to represent links to list
// of objects of type 'LDAP_attributes'.
const LDAPAttributesListLinkKind = "LDAPAttributesListLink"

// LDAPAttributesNilKind is the name of the type used to nil lists of objects of
// type 'LDAP_attributes'.
const LDAPAttributesListNilKind = "LDAPAttributesListNil"

// LDAPAttributesList is a list of values of the 'LDAP_attributes' type.
type LDAPAttributesList struct {
	href  string
	link  bool
	items []*LDAPAttributes
}

// Len returns the length of the list.
func (l *LDAPAttributesList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *LDAPAttributesList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *LDAPAttributesList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *LDAPAttributesList) SetItems(items []*LDAPAttributes) {
	l.items = items
}

// Items returns the items of the list.
func (l *LDAPAttributesList) Items() []*LDAPAttributes {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *LDAPAttributesList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *LDAPAttributesList) Get(i int) *LDAPAttributes {
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
func (l *LDAPAttributesList) Slice() []*LDAPAttributes {
	var slice []*LDAPAttributes
	if l == nil {
		slice = make([]*LDAPAttributes, 0)
	} else {
		slice = make([]*LDAPAttributes, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *LDAPAttributesList) Each(f func(item *LDAPAttributes) bool) {
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
func (l *LDAPAttributesList) Range(f func(index int, item *LDAPAttributes) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
