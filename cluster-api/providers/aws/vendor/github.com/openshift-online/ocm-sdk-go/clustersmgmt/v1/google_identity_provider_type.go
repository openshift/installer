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

// GoogleIdentityProvider represents the values of the 'google_identity_provider' type.
//
// Details for `google` identity providers.
type GoogleIdentityProvider struct {
	bitmap_      uint32
	clientID     string
	clientSecret string
	hostedDomain string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *GoogleIdentityProvider) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ClientID returns the value of the 'client_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Client identifier of a registered _Google_ project.
func (o *GoogleIdentityProvider) ClientID() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.clientID
	}
	return ""
}

// GetClientID returns the value of the 'client_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Client identifier of a registered _Google_ project.
func (o *GoogleIdentityProvider) GetClientID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.clientID
	}
	return
}

// ClientSecret returns the value of the 'client_secret' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Client secret issued by _Google_.
func (o *GoogleIdentityProvider) ClientSecret() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.clientSecret
	}
	return ""
}

// GetClientSecret returns the value of the 'client_secret' attribute and
// a flag indicating if the attribute has a value.
//
// Client secret issued by _Google_.
func (o *GoogleIdentityProvider) GetClientSecret() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.clientSecret
	}
	return
}

// HostedDomain returns the value of the 'hosted_domain' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional hosted domain to restrict sign-in accounts to.
func (o *GoogleIdentityProvider) HostedDomain() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.hostedDomain
	}
	return ""
}

// GetHostedDomain returns the value of the 'hosted_domain' attribute and
// a flag indicating if the attribute has a value.
//
// Optional hosted domain to restrict sign-in accounts to.
func (o *GoogleIdentityProvider) GetHostedDomain() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.hostedDomain
	}
	return
}

// GoogleIdentityProviderListKind is the name of the type used to represent list of objects of
// type 'google_identity_provider'.
const GoogleIdentityProviderListKind = "GoogleIdentityProviderList"

// GoogleIdentityProviderListLinkKind is the name of the type used to represent links to list
// of objects of type 'google_identity_provider'.
const GoogleIdentityProviderListLinkKind = "GoogleIdentityProviderListLink"

// GoogleIdentityProviderNilKind is the name of the type used to nil lists of objects of
// type 'google_identity_provider'.
const GoogleIdentityProviderListNilKind = "GoogleIdentityProviderListNil"

// GoogleIdentityProviderList is a list of values of the 'google_identity_provider' type.
type GoogleIdentityProviderList struct {
	href  string
	link  bool
	items []*GoogleIdentityProvider
}

// Len returns the length of the list.
func (l *GoogleIdentityProviderList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *GoogleIdentityProviderList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *GoogleIdentityProviderList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *GoogleIdentityProviderList) SetItems(items []*GoogleIdentityProvider) {
	l.items = items
}

// Items returns the items of the list.
func (l *GoogleIdentityProviderList) Items() []*GoogleIdentityProvider {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *GoogleIdentityProviderList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *GoogleIdentityProviderList) Get(i int) *GoogleIdentityProvider {
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
func (l *GoogleIdentityProviderList) Slice() []*GoogleIdentityProvider {
	var slice []*GoogleIdentityProvider
	if l == nil {
		slice = make([]*GoogleIdentityProvider, 0)
	} else {
		slice = make([]*GoogleIdentityProvider, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *GoogleIdentityProviderList) Each(f func(item *GoogleIdentityProvider) bool) {
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
func (l *GoogleIdentityProviderList) Range(f func(index int, item *GoogleIdentityProvider) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
