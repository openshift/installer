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

// GithubIdentityProvider represents the values of the 'github_identity_provider' type.
//
// Details for `github` identity providers.
type GithubIdentityProvider struct {
	fieldSet_     []bool
	ca            string
	clientID      string
	clientSecret  string
	hostname      string
	organizations []string
	teams         []string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *GithubIdentityProvider) Empty() bool {
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

// CA returns the value of the 'CA' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional trusted certificate authority bundle to use when making requests tot he server.
func (o *GithubIdentityProvider) CA() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.ca
	}
	return ""
}

// GetCA returns the value of the 'CA' attribute and
// a flag indicating if the attribute has a value.
//
// Optional trusted certificate authority bundle to use when making requests tot he server.
func (o *GithubIdentityProvider) GetCA() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.ca
	}
	return
}

// ClientID returns the value of the 'client_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Client identifier of a registered _GitHub_ OAuth application.
func (o *GithubIdentityProvider) ClientID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.clientID
	}
	return ""
}

// GetClientID returns the value of the 'client_ID' attribute and
// a flag indicating if the attribute has a value.
//
// Client identifier of a registered _GitHub_ OAuth application.
func (o *GithubIdentityProvider) GetClientID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.clientID
	}
	return
}

// ClientSecret returns the value of the 'client_secret' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Client secret of a registered _GitHub_ OAuth application.
func (o *GithubIdentityProvider) ClientSecret() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.clientSecret
	}
	return ""
}

// GetClientSecret returns the value of the 'client_secret' attribute and
// a flag indicating if the attribute has a value.
//
// Client secret of a registered _GitHub_ OAuth application.
func (o *GithubIdentityProvider) GetClientSecret() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.clientSecret
	}
	return
}

// Hostname returns the value of the 'hostname' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// For _GitHub Enterprise_ you must provide the host name of your instance, such as
// `example.com`. This value must match the _GitHub Enterprise_ host name value in the
// `/setup/settings` file and cannot include a port number.
//
// For plain _GitHub_ omit this parameter.
func (o *GithubIdentityProvider) Hostname() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.hostname
	}
	return ""
}

// GetHostname returns the value of the 'hostname' attribute and
// a flag indicating if the attribute has a value.
//
// For _GitHub Enterprise_ you must provide the host name of your instance, such as
// `example.com`. This value must match the _GitHub Enterprise_ host name value in the
// `/setup/settings` file and cannot include a port number.
//
// For plain _GitHub_ omit this parameter.
func (o *GithubIdentityProvider) GetHostname() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.hostname
	}
	return
}

// Organizations returns the value of the 'organizations' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional list of organizations. Cannot be used in combination with the Teams field.
func (o *GithubIdentityProvider) Organizations() []string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.organizations
	}
	return nil
}

// GetOrganizations returns the value of the 'organizations' attribute and
// a flag indicating if the attribute has a value.
//
// Optional list of organizations. Cannot be used in combination with the Teams field.
func (o *GithubIdentityProvider) GetOrganizations() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.organizations
	}
	return
}

// Teams returns the value of the 'teams' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Optional list of teams. Cannot be used in combination with the Organizations field.
func (o *GithubIdentityProvider) Teams() []string {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.teams
	}
	return nil
}

// GetTeams returns the value of the 'teams' attribute and
// a flag indicating if the attribute has a value.
//
// Optional list of teams. Cannot be used in combination with the Organizations field.
func (o *GithubIdentityProvider) GetTeams() (value []string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.teams
	}
	return
}

// GithubIdentityProviderListKind is the name of the type used to represent list of objects of
// type 'github_identity_provider'.
const GithubIdentityProviderListKind = "GithubIdentityProviderList"

// GithubIdentityProviderListLinkKind is the name of the type used to represent links to list
// of objects of type 'github_identity_provider'.
const GithubIdentityProviderListLinkKind = "GithubIdentityProviderListLink"

// GithubIdentityProviderNilKind is the name of the type used to nil lists of objects of
// type 'github_identity_provider'.
const GithubIdentityProviderListNilKind = "GithubIdentityProviderListNil"

// GithubIdentityProviderList is a list of values of the 'github_identity_provider' type.
type GithubIdentityProviderList struct {
	href  string
	link  bool
	items []*GithubIdentityProvider
}

// Len returns the length of the list.
func (l *GithubIdentityProviderList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *GithubIdentityProviderList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *GithubIdentityProviderList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *GithubIdentityProviderList) SetItems(items []*GithubIdentityProvider) {
	l.items = items
}

// Items returns the items of the list.
func (l *GithubIdentityProviderList) Items() []*GithubIdentityProvider {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *GithubIdentityProviderList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *GithubIdentityProviderList) Get(i int) *GithubIdentityProvider {
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
func (l *GithubIdentityProviderList) Slice() []*GithubIdentityProvider {
	var slice []*GithubIdentityProvider
	if l == nil {
		slice = make([]*GithubIdentityProvider, 0)
	} else {
		slice = make([]*GithubIdentityProvider, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *GithubIdentityProviderList) Each(f func(item *GithubIdentityProvider) bool) {
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
func (l *GithubIdentityProviderList) Range(f func(index int, item *GithubIdentityProvider) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
