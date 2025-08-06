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

// WifIdentityProvider represents the values of the 'wif_identity_provider' type.
type WifIdentityProvider struct {
	bitmap_            uint32
	allowedAudiences   []string
	identityProviderId string
	issuerUrl          string
	jwks               string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *WifIdentityProvider) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AllowedAudiences returns the value of the 'allowed_audiences' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifIdentityProvider) AllowedAudiences() []string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.allowedAudiences
	}
	return nil
}

// GetAllowedAudiences returns the value of the 'allowed_audiences' attribute and
// a flag indicating if the attribute has a value.
func (o *WifIdentityProvider) GetAllowedAudiences() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.allowedAudiences
	}
	return
}

// IdentityProviderId returns the value of the 'identity_provider_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifIdentityProvider) IdentityProviderId() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.identityProviderId
	}
	return ""
}

// GetIdentityProviderId returns the value of the 'identity_provider_id' attribute and
// a flag indicating if the attribute has a value.
func (o *WifIdentityProvider) GetIdentityProviderId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.identityProviderId
	}
	return
}

// IssuerUrl returns the value of the 'issuer_url' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifIdentityProvider) IssuerUrl() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.issuerUrl
	}
	return ""
}

// GetIssuerUrl returns the value of the 'issuer_url' attribute and
// a flag indicating if the attribute has a value.
func (o *WifIdentityProvider) GetIssuerUrl() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.issuerUrl
	}
	return
}

// Jwks returns the value of the 'jwks' attribute, or
// the zero value of the type if the attribute doesn't have a value.
func (o *WifIdentityProvider) Jwks() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.jwks
	}
	return ""
}

// GetJwks returns the value of the 'jwks' attribute and
// a flag indicating if the attribute has a value.
func (o *WifIdentityProvider) GetJwks() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.jwks
	}
	return
}

// WifIdentityProviderListKind is the name of the type used to represent list of objects of
// type 'wif_identity_provider'.
const WifIdentityProviderListKind = "WifIdentityProviderList"

// WifIdentityProviderListLinkKind is the name of the type used to represent links to list
// of objects of type 'wif_identity_provider'.
const WifIdentityProviderListLinkKind = "WifIdentityProviderListLink"

// WifIdentityProviderNilKind is the name of the type used to nil lists of objects of
// type 'wif_identity_provider'.
const WifIdentityProviderListNilKind = "WifIdentityProviderListNil"

// WifIdentityProviderList is a list of values of the 'wif_identity_provider' type.
type WifIdentityProviderList struct {
	href  string
	link  bool
	items []*WifIdentityProvider
}

// Len returns the length of the list.
func (l *WifIdentityProviderList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *WifIdentityProviderList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *WifIdentityProviderList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *WifIdentityProviderList) SetItems(items []*WifIdentityProvider) {
	l.items = items
}

// Items returns the items of the list.
func (l *WifIdentityProviderList) Items() []*WifIdentityProvider {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *WifIdentityProviderList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *WifIdentityProviderList) Get(i int) *WifIdentityProvider {
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
func (l *WifIdentityProviderList) Slice() []*WifIdentityProvider {
	var slice []*WifIdentityProvider
	if l == nil {
		slice = make([]*WifIdentityProvider, 0)
	} else {
		slice = make([]*WifIdentityProvider, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *WifIdentityProviderList) Each(f func(item *WifIdentityProvider) bool) {
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
func (l *WifIdentityProviderList) Range(f func(index int, item *WifIdentityProvider) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
