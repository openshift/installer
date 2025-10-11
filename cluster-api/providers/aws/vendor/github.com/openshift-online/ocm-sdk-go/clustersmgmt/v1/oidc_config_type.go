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

import (
	time "time"
)

// OidcConfig represents the values of the 'oidc_config' type.
//
// Contains the necessary attributes to support oidc configuration hosting under Red Hat or registering a Customer's byo oidc config.
type OidcConfig struct {
	bitmap_             uint32
	href                string
	id                  string
	creationTimestamp   time.Time
	installerRoleArn    string
	issuerUrl           string
	lastUpdateTimestamp time.Time
	lastUsedTimestamp   time.Time
	organizationId      string
	secretArn           string
	managed             bool
	reusable            bool
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *OidcConfig) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// HREF returns the value of the 'HREF' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// HREF for the oidc config, filled in response.
func (o *OidcConfig) HREF() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the value of the 'HREF' attribute and
// a flag indicating if the attribute has a value.
//
// HREF for the oidc config, filled in response.
func (o *OidcConfig) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.href
	}
	return
}

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ID for the oidc config, filled in response.
func (o *OidcConfig) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// ID for the oidc config, filled in response.
func (o *OidcConfig) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Creation timestamp, filled in response.
func (o *OidcConfig) CreationTimestamp() time.Time {
	if o != nil && o.bitmap_&4 != 0 {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Creation timestamp, filled in response.
func (o *OidcConfig) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.creationTimestamp
	}
	return
}

// InstallerRoleArn returns the value of the 'installer_role_arn' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ARN of the AWS role to assume when installing the cluster as to reveal the secret, supplied in request. It is only to be used in Unmanaged Oidc Config.
func (o *OidcConfig) InstallerRoleArn() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.installerRoleArn
	}
	return ""
}

// GetInstallerRoleArn returns the value of the 'installer_role_arn' attribute and
// a flag indicating if the attribute has a value.
//
// ARN of the AWS role to assume when installing the cluster as to reveal the secret, supplied in request. It is only to be used in Unmanaged Oidc Config.
func (o *OidcConfig) GetInstallerRoleArn() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.installerRoleArn
	}
	return
}

// IssuerUrl returns the value of the 'issuer_url' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Issuer URL, filled in response when Managed and supplied in Unmanaged.
func (o *OidcConfig) IssuerUrl() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.issuerUrl
	}
	return ""
}

// GetIssuerUrl returns the value of the 'issuer_url' attribute and
// a flag indicating if the attribute has a value.
//
// Issuer URL, filled in response when Managed and supplied in Unmanaged.
func (o *OidcConfig) GetIssuerUrl() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.issuerUrl
	}
	return
}

// LastUpdateTimestamp returns the value of the 'last_update_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Last update timestamp, filled when patching a valid attribute of this oidc config.
func (o *OidcConfig) LastUpdateTimestamp() time.Time {
	if o != nil && o.bitmap_&32 != 0 {
		return o.lastUpdateTimestamp
	}
	return time.Time{}
}

// GetLastUpdateTimestamp returns the value of the 'last_update_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Last update timestamp, filled when patching a valid attribute of this oidc config.
func (o *OidcConfig) GetLastUpdateTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&32 != 0
	if ok {
		value = o.lastUpdateTimestamp
	}
	return
}

// LastUsedTimestamp returns the value of the 'last_used_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Last used timestamp, filled by the latest cluster that used this oidc config.
func (o *OidcConfig) LastUsedTimestamp() time.Time {
	if o != nil && o.bitmap_&64 != 0 {
		return o.lastUsedTimestamp
	}
	return time.Time{}
}

// GetLastUsedTimestamp returns the value of the 'last_used_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Last used timestamp, filled by the latest cluster that used this oidc config.
func (o *OidcConfig) GetLastUsedTimestamp() (value time.Time, ok bool) {
	ok = o != nil && o.bitmap_&64 != 0
	if ok {
		value = o.lastUsedTimestamp
	}
	return
}

// Managed returns the value of the 'managed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates whether it is Managed or Unmanaged (Customer hosted).
func (o *OidcConfig) Managed() bool {
	if o != nil && o.bitmap_&128 != 0 {
		return o.managed
	}
	return false
}

// GetManaged returns the value of the 'managed' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates whether it is Managed or Unmanaged (Customer hosted).
func (o *OidcConfig) GetManaged() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&128 != 0
	if ok {
		value = o.managed
	}
	return
}

// OrganizationId returns the value of the 'organization_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Organization ID, filled in response respecting token provided.
func (o *OidcConfig) OrganizationId() string {
	if o != nil && o.bitmap_&256 != 0 {
		return o.organizationId
	}
	return ""
}

// GetOrganizationId returns the value of the 'organization_id' attribute and
// a flag indicating if the attribute has a value.
//
// Organization ID, filled in response respecting token provided.
func (o *OidcConfig) GetOrganizationId() (value string, ok bool) {
	ok = o != nil && o.bitmap_&256 != 0
	if ok {
		value = o.organizationId
	}
	return
}

// Reusable returns the value of the 'reusable' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates whether the Oidc Config can be reused.
func (o *OidcConfig) Reusable() bool {
	if o != nil && o.bitmap_&512 != 0 {
		return o.reusable
	}
	return false
}

// GetReusable returns the value of the 'reusable' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates whether the Oidc Config can be reused.
func (o *OidcConfig) GetReusable() (value bool, ok bool) {
	ok = o != nil && o.bitmap_&512 != 0
	if ok {
		value = o.reusable
	}
	return
}

// SecretArn returns the value of the 'secret_arn' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Secrets Manager ARN for the OIDC private key, supplied in request. It is only to be used in Unmanaged Oidc Config.
func (o *OidcConfig) SecretArn() string {
	if o != nil && o.bitmap_&1024 != 0 {
		return o.secretArn
	}
	return ""
}

// GetSecretArn returns the value of the 'secret_arn' attribute and
// a flag indicating if the attribute has a value.
//
// Secrets Manager ARN for the OIDC private key, supplied in request. It is only to be used in Unmanaged Oidc Config.
func (o *OidcConfig) GetSecretArn() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1024 != 0
	if ok {
		value = o.secretArn
	}
	return
}

// OidcConfigListKind is the name of the type used to represent list of objects of
// type 'oidc_config'.
const OidcConfigListKind = "OidcConfigList"

// OidcConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'oidc_config'.
const OidcConfigListLinkKind = "OidcConfigListLink"

// OidcConfigNilKind is the name of the type used to nil lists of objects of
// type 'oidc_config'.
const OidcConfigListNilKind = "OidcConfigListNil"

// OidcConfigList is a list of values of the 'oidc_config' type.
type OidcConfigList struct {
	href  string
	link  bool
	items []*OidcConfig
}

// Len returns the length of the list.
func (l *OidcConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *OidcConfigList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *OidcConfigList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *OidcConfigList) SetItems(items []*OidcConfig) {
	l.items = items
}

// Items returns the items of the list.
func (l *OidcConfigList) Items() []*OidcConfig {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *OidcConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *OidcConfigList) Get(i int) *OidcConfig {
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
func (l *OidcConfigList) Slice() []*OidcConfig {
	var slice []*OidcConfig
	if l == nil {
		slice = make([]*OidcConfig, 0)
	} else {
		slice = make([]*OidcConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *OidcConfigList) Each(f func(item *OidcConfig) bool) {
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
func (l *OidcConfigList) Range(f func(index int, item *OidcConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
