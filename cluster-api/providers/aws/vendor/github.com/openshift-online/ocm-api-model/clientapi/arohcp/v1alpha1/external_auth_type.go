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

// ExternalAuthKind is the name of the type used to represent objects
// of type 'external_auth'.
const ExternalAuthKind = "ExternalAuth"

// ExternalAuthLinkKind is the name of the type used to represent links
// to objects of type 'external_auth'.
const ExternalAuthLinkKind = "ExternalAuthLink"

// ExternalAuthNilKind is the name of the type used to nil references
// to objects of type 'external_auth'.
const ExternalAuthNilKind = "ExternalAuthNil"

// ExternalAuth represents the values of the 'external_auth' type.
//
// Representation of an external authentication provider.
type ExternalAuth struct {
	fieldSet_ []bool
	id        string
	href      string
	claim     *ExternalAuthClaim
	clients   []*ExternalAuthClientConfig
	issuer    *TokenIssuer
	status    *ExternalAuthStatus
}

// Kind returns the name of the type of the object.
func (o *ExternalAuth) Kind() string {
	if o == nil {
		return ExternalAuthNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ExternalAuthLinkKind
	}
	return ExternalAuthKind
}

// Link returns true if this is a link.
func (o *ExternalAuth) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *ExternalAuth) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ExternalAuth) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ExternalAuth) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ExternalAuth) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ExternalAuth) Empty() bool {
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

// Claim returns the value of the 'claim' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The rules on how to transform information from an ID token into a cluster identity.
func (o *ExternalAuth) Claim() *ExternalAuthClaim {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.claim
	}
	return nil
}

// GetClaim returns the value of the 'claim' attribute and
// a flag indicating if the attribute has a value.
//
// The rules on how to transform information from an ID token into a cluster identity.
func (o *ExternalAuth) GetClaim() (value *ExternalAuthClaim, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.claim
	}
	return
}

// Clients returns the value of the 'clients' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The list of the platform's clients that need to request tokens from the issuer.
func (o *ExternalAuth) Clients() []*ExternalAuthClientConfig {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.clients
	}
	return nil
}

// GetClients returns the value of the 'clients' attribute and
// a flag indicating if the attribute has a value.
//
// The list of the platform's clients that need to request tokens from the issuer.
func (o *ExternalAuth) GetClients() (value []*ExternalAuthClientConfig, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.clients
	}
	return
}

// Issuer returns the value of the 'issuer' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The issuer describes the attributes of the OIDC token issuer.
func (o *ExternalAuth) Issuer() *TokenIssuer {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.issuer
	}
	return nil
}

// GetIssuer returns the value of the 'issuer' attribute and
// a flag indicating if the attribute has a value.
//
// The issuer describes the attributes of the OIDC token issuer.
func (o *ExternalAuth) GetIssuer() (value *TokenIssuer, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.issuer
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The status describes the current state of the external authentication provider.
// This is read-only.
func (o *ExternalAuth) Status() *ExternalAuthStatus {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.status
	}
	return nil
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// The status describes the current state of the external authentication provider.
// This is read-only.
func (o *ExternalAuth) GetStatus() (value *ExternalAuthStatus, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.status
	}
	return
}

// ExternalAuthListKind is the name of the type used to represent list of objects of
// type 'external_auth'.
const ExternalAuthListKind = "ExternalAuthList"

// ExternalAuthListLinkKind is the name of the type used to represent links to list
// of objects of type 'external_auth'.
const ExternalAuthListLinkKind = "ExternalAuthListLink"

// ExternalAuthNilKind is the name of the type used to nil lists of objects of
// type 'external_auth'.
const ExternalAuthListNilKind = "ExternalAuthListNil"

// ExternalAuthList is a list of values of the 'external_auth' type.
type ExternalAuthList struct {
	href  string
	link  bool
	items []*ExternalAuth
}

// Kind returns the name of the type of the object.
func (l *ExternalAuthList) Kind() string {
	if l == nil {
		return ExternalAuthListNilKind
	}
	if l.link {
		return ExternalAuthListLinkKind
	}
	return ExternalAuthListKind
}

// Link returns true iif this is a link.
func (l *ExternalAuthList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ExternalAuthList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ExternalAuthList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ExternalAuthList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ExternalAuthList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ExternalAuthList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ExternalAuthList) SetItems(items []*ExternalAuth) {
	l.items = items
}

// Items returns the items of the list.
func (l *ExternalAuthList) Items() []*ExternalAuth {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ExternalAuthList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ExternalAuthList) Get(i int) *ExternalAuth {
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
func (l *ExternalAuthList) Slice() []*ExternalAuth {
	var slice []*ExternalAuth
	if l == nil {
		slice = make([]*ExternalAuth, 0)
	} else {
		slice = make([]*ExternalAuth, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ExternalAuthList) Each(f func(item *ExternalAuth) bool) {
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
func (l *ExternalAuthList) Range(f func(index int, item *ExternalAuth) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
