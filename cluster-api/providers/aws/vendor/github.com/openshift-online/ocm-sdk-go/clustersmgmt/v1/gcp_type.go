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

// GCP represents the values of the 'GCP' type.
//
// Google cloud platform settings of a cluster.
type GCP struct {
	bitmap_                 uint32
	authURI                 string
	authProviderX509CertURL string
	authentication          *GcpAuthentication
	clientID                string
	clientX509CertURL       string
	clientEmail             string
	privateKey              string
	privateKeyID            string
	privateServiceConnect   *GcpPrivateServiceConnect
	projectID               string
	security                *GcpSecurity
	tokenURI                string
	type_                   string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *GCP) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AuthURI returns the value of the 'auth_URI' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP authentication uri
func (o *GCP) AuthURI() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.authURI
	}
	return ""
}

// GetAuthURI returns the value of the 'auth_URI' attribute and
// a flag indicating if the attribute has a value.
//
// GCP authentication uri
func (o *GCP) GetAuthURI() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.authURI
	}
	return
}

// AuthProviderX509CertURL returns the value of the 'auth_provider_X509_cert_URL' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP Authentication provider x509 certificate url
func (o *GCP) AuthProviderX509CertURL() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.authProviderX509CertURL
	}
	return ""
}

// GetAuthProviderX509CertURL returns the value of the 'auth_provider_X509_cert_URL' attribute and
// a flag indicating if the attribute has a value.
//
// GCP Authentication provider x509 certificate url
func (o *GCP) GetAuthProviderX509CertURL() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.authProviderX509CertURL
	}
	return
}

// Authentication returns the value of the 'authentication' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP Authentication Method
func (o *GCP) Authentication() *GcpAuthentication {
	if o != nil && o.bitmap_&4 != 0 {
		return o.authentication
	}
	return nil
}

// GetAuthentication returns the value of the 'authentication' attribute and
// a flag indicating if the attribute has a value.
//
// GCP Authentication Method
func (o *GCP) GetAuthentication() (value *GcpAuthentication, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.authentication
	}
	return
}

// ClientID returns the value of the 'client_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP client identifier
func (o *GCP) ClientID() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.clientID
	}
	return ""
}

// GetClientID returns the value of the 'client_ID' attribute and
// a flag indicating if the attribute has a value.
//
// GCP client identifier
func (o *GCP) GetClientID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.clientID
	}
	return
}

// ClientX509CertURL returns the value of the 'client_X509_cert_URL' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP client x509 certificate url
func (o *GCP) ClientX509CertURL() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.clientX509CertURL
	}
	return ""
}

// GetClientX509CertURL returns the value of the 'client_X509_cert_URL' attribute and
// a flag indicating if the attribute has a value.
//
// GCP client x509 certificate url
func (o *GCP) GetClientX509CertURL() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.clientX509CertURL
	}
	return
}

// ClientEmail returns the value of the 'client_email' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP client email
func (o *GCP) ClientEmail() string {
	if o != nil && o.bitmap_&32 != 0 {
		return o.clientEmail
	}
	return ""
}

// GetClientEmail returns the value of the 'client_email' attribute and
// a flag indicating if the attribute has a value.
//
// GCP client email
func (o *GCP) GetClientEmail() (value string, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.clientEmail
	}
	return
}

// PrivateKey returns the value of the 'private_key' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP private key
func (o *GCP) PrivateKey() string {
	if o != nil && o.bitmap_&64 != 0 {
		return o.privateKey
	}
	return ""
}

// GetPrivateKey returns the value of the 'private_key' attribute and
// a flag indicating if the attribute has a value.
//
// GCP private key
func (o *GCP) GetPrivateKey() (value string, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.privateKey
	}
	return
}

// PrivateKeyID returns the value of the 'private_key_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP private key identifier
func (o *GCP) PrivateKeyID() string {
	if o != nil && o.bitmap_&128 != 0 {
		return o.privateKeyID
	}
	return ""
}

// GetPrivateKeyID returns the value of the 'private_key_ID' attribute and
// a flag indicating if the attribute has a value.
//
// GCP private key identifier
func (o *GCP) GetPrivateKeyID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.privateKeyID
	}
	return
}

// PrivateServiceConnect returns the value of the 'private_service_connect' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP PrivateServiceConnect configuration
func (o *GCP) PrivateServiceConnect() *GcpPrivateServiceConnect {
	if o != nil && o.bitmap_&256 != 0 {
		return o.privateServiceConnect
	}
	return nil
}

// GetPrivateServiceConnect returns the value of the 'private_service_connect' attribute and
// a flag indicating if the attribute has a value.
//
// GCP PrivateServiceConnect configuration
func (o *GCP) GetPrivateServiceConnect() (value *GcpPrivateServiceConnect, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.privateServiceConnect
	}
	return
}

// ProjectID returns the value of the 'project_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP project identifier.
func (o *GCP) ProjectID() string {
	if o != nil && o.bitmap_&512 != 0 {
		return o.projectID
	}
	return ""
}

// GetProjectID returns the value of the 'project_ID' attribute and
// a flag indicating if the attribute has a value.
//
// GCP project identifier.
func (o *GCP) GetProjectID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.projectID
	}
	return
}

// Security returns the value of the 'security' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP Security Settings
func (o *GCP) Security() *GcpSecurity {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.security
	}
	return nil
}

// GetSecurity returns the value of the 'security' attribute and
// a flag indicating if the attribute has a value.
//
// GCP Security Settings
func (o *GCP) GetSecurity() (value *GcpSecurity, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.security
	}
	return
}

// TokenURI returns the value of the 'token_URI' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP token uri
func (o *GCP) TokenURI() string {
	if o != nil && o.bitmap_&2048 != 0 {
		return o.tokenURI
	}
	return ""
}

// GetTokenURI returns the value of the 'token_URI' attribute and
// a flag indicating if the attribute has a value.
//
// GCP token uri
func (o *GCP) GetTokenURI() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2048 != 0
	if ok {
		value = o.tokenURI
	}
	return
}

// Type returns the value of the 'type' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP the type of the service the key belongs to
func (o *GCP) Type() string {
	if o != nil && o.bitmap_&4096 != 0 {
		return o.type_
	}
	return ""
}

// GetType returns the value of the 'type' attribute and
// a flag indicating if the attribute has a value.
//
// GCP the type of the service the key belongs to
func (o *GCP) GetType() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4096 != 0
	if ok {
		value = o.type_
	}
	return
}

// GCPListKind is the name of the type used to represent list of objects of
// type 'GCP'.
const GCPListKind = "GCPList"

// GCPListLinkKind is the name of the type used to represent links to list
// of objects of type 'GCP'.
const GCPListLinkKind = "GCPListLink"

// GCPNilKind is the name of the type used to nil lists of objects of
// type 'GCP'.
const GCPListNilKind = "GCPListNil"

// GCPList is a list of values of the 'GCP' type.
type GCPList struct {
	href  string
	link  bool
	items []*GCP
}

// Len returns the length of the list.
func (l *GCPList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *GCPList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *GCPList) Get(i int) *GCP {
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
func (l *GCPList) Slice() []*GCP {
	var slice []*GCP
	if l == nil {
		slice = make([]*GCP, 0)
	} else {
		slice = make([]*GCP, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *GCPList) Each(f func(item *GCP) bool) {
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
func (l *GCPList) Range(f func(index int, item *GCP) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
