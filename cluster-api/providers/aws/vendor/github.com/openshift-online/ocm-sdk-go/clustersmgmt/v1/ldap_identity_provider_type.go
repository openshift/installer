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

// LDAPIdentityProvider represents the values of the 'LDAP_identity_provider' type.
//
// Details for `ldap` identity providers.
type LDAPIdentityProvider struct {
	bitmap_      uint32
	ca           string
	url          string
	attributes   *LDAPAttributes
	bindDN       string
	bindPassword string
	insecure     bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *LDAPIdentityProvider) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// CA returns the value of the 'CA' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Certificate bundle to use to validate server certificates for the configured URL.
func (o *LDAPIdentityProvider) CA() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.ca
	}
	return ""
}

// GetCA returns the value of the 'CA' attribute and
// a flag indicating if the attribute has a value.
//
// Certificate bundle to use to validate server certificates for the configured URL.
func (o *LDAPIdentityProvider) GetCA() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.ca
	}
	return
}

// URL returns the value of the 'URL' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// An https://tools.ietf.org/html/rfc2255[RFC 2255] URL which specifies the LDAP host and
// search parameters to use.
func (o *LDAPIdentityProvider) URL() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.url
	}
	return ""
}

// GetURL returns the value of the 'URL' attribute and
// a flag indicating if the attribute has a value.
//
// An https://tools.ietf.org/html/rfc2255[RFC 2255] URL which specifies the LDAP host and
// search parameters to use.
func (o *LDAPIdentityProvider) GetURL() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.url
	}
	return
}

// Attributes returns the value of the 'attributes' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// LDAP attributes used to configure the provider.
func (o *LDAPIdentityProvider) Attributes() *LDAPAttributes {
	if o != nil && o.bitmap_&4 != 0 {
		return o.attributes
	}
	return nil
}

// GetAttributes returns the value of the 'attributes' attribute and
// a flag indicating if the attribute has a value.
//
// LDAP attributes used to configure the provider.
func (o *LDAPIdentityProvider) GetAttributes() (value *LDAPAttributes, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.attributes
	}
	return
}

// BindDN returns the value of the 'bind_DN' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional distinguished name to use to bind during the search phase.
func (o *LDAPIdentityProvider) BindDN() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.bindDN
	}
	return ""
}

// GetBindDN returns the value of the 'bind_DN' attribute and
// a flag indicating if the attribute has a value.
//
// Optional distinguished name to use to bind during the search phase.
func (o *LDAPIdentityProvider) GetBindDN() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.bindDN
	}
	return
}

// BindPassword returns the value of the 'bind_password' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional password to use to bind during the search phase.
func (o *LDAPIdentityProvider) BindPassword() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.bindPassword
	}
	return ""
}

// GetBindPassword returns the value of the 'bind_password' attribute and
// a flag indicating if the attribute has a value.
//
// Optional password to use to bind during the search phase.
func (o *LDAPIdentityProvider) GetBindPassword() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.bindPassword
	}
	return
}

// Insecure returns the value of the 'insecure' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// When `true` no TLS connection is made to the server. When `false` `ldaps://...` URLs
// connect using TLS and `ldap://...` are upgraded to TLS.
func (o *LDAPIdentityProvider) Insecure() bool {
	if o != nil && o.bitmap_&32 != 0 {
		return o.insecure
	}
	return false
}

// GetInsecure returns the value of the 'insecure' attribute and
// a flag indicating if the attribute has a value.
//
// When `true` no TLS connection is made to the server. When `false` `ldaps://...` URLs
// connect using TLS and `ldap://...` are upgraded to TLS.
func (o *LDAPIdentityProvider) GetInsecure() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.insecure
	}
	return
}

// LDAPIdentityProviderListKind is the name of the type used to represent list of objects of
// type 'LDAP_identity_provider'.
const LDAPIdentityProviderListKind = "LDAPIdentityProviderList"

// LDAPIdentityProviderListLinkKind is the name of the type used to represent links to list
// of objects of type 'LDAP_identity_provider'.
const LDAPIdentityProviderListLinkKind = "LDAPIdentityProviderListLink"

// LDAPIdentityProviderNilKind is the name of the type used to nil lists of objects of
// type 'LDAP_identity_provider'.
const LDAPIdentityProviderListNilKind = "LDAPIdentityProviderListNil"

// LDAPIdentityProviderList is a list of values of the 'LDAP_identity_provider' type.
type LDAPIdentityProviderList struct {
	href  string
	link  bool
	items []*LDAPIdentityProvider
}

// Len returns the length of the list.
func (l *LDAPIdentityProviderList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *LDAPIdentityProviderList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *LDAPIdentityProviderList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *LDAPIdentityProviderList) SetItems(items []*LDAPIdentityProvider) {
	l.items = items
}

// Items returns the items of the list.
func (l *LDAPIdentityProviderList) Items() []*LDAPIdentityProvider {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *LDAPIdentityProviderList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *LDAPIdentityProviderList) Get(i int) *LDAPIdentityProvider {
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
func (l *LDAPIdentityProviderList) Slice() []*LDAPIdentityProvider {
	var slice []*LDAPIdentityProvider
	if l == nil {
		slice = make([]*LDAPIdentityProvider, 0)
	} else {
		slice = make([]*LDAPIdentityProvider, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *LDAPIdentityProviderList) Each(f func(item *LDAPIdentityProvider) bool) {
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
func (l *LDAPIdentityProviderList) Range(f func(index int, item *LDAPIdentityProvider) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
