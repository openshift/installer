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

// OpenIDIdentityProvider represents the values of the 'open_ID_identity_provider' type.
//
// Details for `openid` identity providers.
type OpenIDIdentityProvider struct {
	bitmap_                  uint32
	ca                       string
	claims                   *OpenIDClaims
	clientID                 string
	clientSecret             string
	extraAuthorizeParameters map[string]string
	extraScopes              []string
	issuer                   string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *OpenIDIdentityProvider) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// CA returns the value of the 'CA' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Certificate bunde to use to validate server certificates for the configured URL.
func (o *OpenIDIdentityProvider) CA() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.ca
	}
	return ""
}

// GetCA returns the value of the 'CA' attribute and
// a flag indicating if the attribute has a value.
//
// Certificate bunde to use to validate server certificates for the configured URL.
func (o *OpenIDIdentityProvider) GetCA() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.ca
	}
	return
}

// Claims returns the value of the 'claims' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Claims used to configure the provider.
func (o *OpenIDIdentityProvider) Claims() *OpenIDClaims {
	if o != nil && o.bitmap_&2 != 0 {
		return o.claims
	}
	return nil
}

// GetClaims returns the value of the 'claims' attribute and
// a flag indicating if the attribute has a value.
//
// Claims used to configure the provider.
func (o *OpenIDIdentityProvider) GetClaims() (value *OpenIDClaims, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.claims
	}
	return
}

// ClientID returns the value of the 'client_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Identifier of a client registered with the _OpenID_ provider.
func (o *OpenIDIdentityProvider) ClientID() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.clientID
	}
	return ""
}

// GetClientID returns the value of the 'client_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Identifier of a client registered with the _OpenID_ provider.
func (o *OpenIDIdentityProvider) GetClientID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.clientID
	}
	return
}

// ClientSecret returns the value of the 'client_secret' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Client secret.
func (o *OpenIDIdentityProvider) ClientSecret() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.clientSecret
	}
	return ""
}

// GetClientSecret returns the value of the 'client_secret' attribute and
// a flag indicating if the attribute has a value.
//
// Client secret.
func (o *OpenIDIdentityProvider) GetClientSecret() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.clientSecret
	}
	return
}

// ExtraAuthorizeParameters returns the value of the 'extra_authorize_parameters' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional map of extra parameters to add to the authorization token request.
func (o *OpenIDIdentityProvider) ExtraAuthorizeParameters() map[string]string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.extraAuthorizeParameters
	}
	return nil
}

// GetExtraAuthorizeParameters returns the value of the 'extra_authorize_parameters' attribute and
// a flag indicating if the attribute has a value.
//
// Optional map of extra parameters to add to the authorization token request.
func (o *OpenIDIdentityProvider) GetExtraAuthorizeParameters() (value map[string]string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.extraAuthorizeParameters
	}
	return
}

// ExtraScopes returns the value of the 'extra_scopes' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional list of scopes to request, in addition to the `openid` scope, during the
// authorization token request.
func (o *OpenIDIdentityProvider) ExtraScopes() []string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.extraScopes
	}
	return nil
}

// GetExtraScopes returns the value of the 'extra_scopes' attribute and
// a flag indicating if the attribute has a value.
//
// Optional list of scopes to request, in addition to the `openid` scope, during the
// authorization token request.
func (o *OpenIDIdentityProvider) GetExtraScopes() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.extraScopes
	}
	return
}

// Issuer returns the value of the 'issuer' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The URL that the OpenID Provider asserts as the Issuer Identifier
func (o *OpenIDIdentityProvider) Issuer() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.issuer
	}
	return ""
}

// GetIssuer returns the value of the 'issuer' attribute and
// a flag indicating if the attribute has a value.
//
// The URL that the OpenID Provider asserts as the Issuer Identifier
func (o *OpenIDIdentityProvider) GetIssuer() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.issuer
	}
	return
}

// OpenIDIdentityProviderListKind is the name of the type used to represent list of objects of
// type 'open_ID_identity_provider'.
const OpenIDIdentityProviderListKind = "OpenIDIdentityProviderList"

// OpenIDIdentityProviderListLinkKind is the name of the type used to represent links to list
// of objects of type 'open_ID_identity_provider'.
const OpenIDIdentityProviderListLinkKind = "OpenIDIdentityProviderListLink"

// OpenIDIdentityProviderNilKind is the name of the type used to nil lists of objects of
// type 'open_ID_identity_provider'.
const OpenIDIdentityProviderListNilKind = "OpenIDIdentityProviderListNil"

// OpenIDIdentityProviderList is a list of values of the 'open_ID_identity_provider' type.
type OpenIDIdentityProviderList struct {
	href  string
	link  bool
	items []*OpenIDIdentityProvider
}

// Len returns the length of the list.
func (l *OpenIDIdentityProviderList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *OpenIDIdentityProviderList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *OpenIDIdentityProviderList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *OpenIDIdentityProviderList) SetItems(items []*OpenIDIdentityProvider) {
	l.items = items
}

// Items returns the items of the list.
func (l *OpenIDIdentityProviderList) Items() []*OpenIDIdentityProvider {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *OpenIDIdentityProviderList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *OpenIDIdentityProviderList) Get(i int) *OpenIDIdentityProvider {
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
func (l *OpenIDIdentityProviderList) Slice() []*OpenIDIdentityProvider {
	var slice []*OpenIDIdentityProvider
	if l == nil {
		slice = make([]*OpenIDIdentityProvider, 0)
	} else {
		slice = make([]*OpenIDIdentityProvider, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *OpenIDIdentityProviderList) Each(f func(item *OpenIDIdentityProvider) bool) {
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
func (l *OpenIDIdentityProviderList) Range(f func(index int, item *OpenIDIdentityProvider) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
