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

// TokenIssuer represents the values of the 'token_issuer' type.
//
// Representation of a token issuer used in an external authentication.
type TokenIssuer struct {
	bitmap_   uint32
	ca        string
	url       string
	audiences []string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *TokenIssuer) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// CA returns the value of the 'CA' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Certificate bundle to use to validate server certificates for the configured URL.
func (o *TokenIssuer) CA() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.ca
	}
	return ""
}

// GetCA returns the value of the 'CA' attribute and
// a flag indicating if the attribute has a value.
//
// Certificate bundle to use to validate server certificates for the configured URL.
func (o *TokenIssuer) GetCA() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.ca
	}
	return
}

// URL returns the value of the 'URL' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// URL is the serving URL of the token issuer.
func (o *TokenIssuer) URL() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.url
	}
	return ""
}

// GetURL returns the value of the 'URL' attribute and
// a flag indicating if the attribute has a value.
//
// URL is the serving URL of the token issuer.
func (o *TokenIssuer) GetURL() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.url
	}
	return
}

// Audiences returns the value of the 'audiences' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Audiences is an array of audiences that the token was issued for.
// Valid tokens must include at least one of these values in their
// "aud" claim.
// Must be set to exactly one value.
func (o *TokenIssuer) Audiences() []string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.audiences
	}
	return nil
}

// GetAudiences returns the value of the 'audiences' attribute and
// a flag indicating if the attribute has a value.
//
// Audiences is an array of audiences that the token was issued for.
// Valid tokens must include at least one of these values in their
// "aud" claim.
// Must be set to exactly one value.
func (o *TokenIssuer) GetAudiences() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.audiences
	}
	return
}

// TokenIssuerListKind is the name of the type used to represent list of objects of
// type 'token_issuer'.
const TokenIssuerListKind = "TokenIssuerList"

// TokenIssuerListLinkKind is the name of the type used to represent links to list
// of objects of type 'token_issuer'.
const TokenIssuerListLinkKind = "TokenIssuerListLink"

// TokenIssuerNilKind is the name of the type used to nil lists of objects of
// type 'token_issuer'.
const TokenIssuerListNilKind = "TokenIssuerListNil"

// TokenIssuerList is a list of values of the 'token_issuer' type.
type TokenIssuerList struct {
	href  string
	link  bool
	items []*TokenIssuer
}

// Len returns the length of the list.
func (l *TokenIssuerList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *TokenIssuerList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *TokenIssuerList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *TokenIssuerList) SetItems(items []*TokenIssuer) {
	l.items = items
}

// Items returns the items of the list.
func (l *TokenIssuerList) Items() []*TokenIssuer {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *TokenIssuerList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *TokenIssuerList) Get(i int) *TokenIssuer {
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
func (l *TokenIssuerList) Slice() []*TokenIssuer {
	var slice []*TokenIssuer
	if l == nil {
		slice = make([]*TokenIssuer, 0)
	} else {
		slice = make([]*TokenIssuer, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *TokenIssuerList) Each(f func(item *TokenIssuer) bool) {
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
func (l *TokenIssuerList) Range(f func(index int, item *TokenIssuer) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
