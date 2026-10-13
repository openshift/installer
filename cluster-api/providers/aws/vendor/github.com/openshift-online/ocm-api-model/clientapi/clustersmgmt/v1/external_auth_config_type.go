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

// ExternalAuthConfigKind is the name of the type used to represent objects
// of type 'external_auth_config'.
const ExternalAuthConfigKind = "ExternalAuthConfig"

// ExternalAuthConfigLinkKind is the name of the type used to represent links
// to objects of type 'external_auth_config'.
const ExternalAuthConfigLinkKind = "ExternalAuthConfigLink"

// ExternalAuthConfigNilKind is the name of the type used to nil references
// to objects of type 'external_auth_config'.
const ExternalAuthConfigNilKind = "ExternalAuthConfigNil"

// ExternalAuthConfig represents the values of the 'external_auth_config' type.
//
// Represents an external authentication configuration
type ExternalAuthConfig struct {
	fieldSet_     []bool
	id            string
	href          string
	externalAuths *ExternalAuthList
	state         ExternalAuthConfigState
	enabled       bool
}

// Kind returns the name of the type of the object.
func (o *ExternalAuthConfig) Kind() string {
	if o == nil {
		return ExternalAuthConfigNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ExternalAuthConfigLinkKind
	}
	return ExternalAuthConfigKind
}

// Link returns true if this is a link.
func (o *ExternalAuthConfig) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *ExternalAuthConfig) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ExternalAuthConfig) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ExternalAuthConfig) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ExternalAuthConfig) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ExternalAuthConfig) Empty() bool {
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

// Enabled returns the value of the 'enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Boolean flag indicating if the cluster should use an external authentication configuration for ROSA HCP clusters.
//
// By default this is false.
//
// To enable it the cluster needs to be ROSA HCP cluster and the organization of the user needs
// to have the `external-authentication` feature toggle enabled.
//
// For ARO HCP clusters, use the "State" property to enable/disable this feature instead.
func (o *ExternalAuthConfig) Enabled() bool {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.enabled
	}
	return false
}

// GetEnabled returns the value of the 'enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Boolean flag indicating if the cluster should use an external authentication configuration for ROSA HCP clusters.
//
// By default this is false.
//
// To enable it the cluster needs to be ROSA HCP cluster and the organization of the user needs
// to have the `external-authentication` feature toggle enabled.
//
// For ARO HCP clusters, use the "State" property to enable/disable this feature instead.
func (o *ExternalAuthConfig) GetEnabled() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.enabled
	}
	return
}

// ExternalAuths returns the value of the 'external_auths' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// A list of external authentication providers configured for the cluster.
//
// Only one external authentication provider can be configured.
func (o *ExternalAuthConfig) ExternalAuths() *ExternalAuthList {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.externalAuths
	}
	return nil
}

// GetExternalAuths returns the value of the 'external_auths' attribute and
// a flag indicating if the attribute has a value.
//
// A list of external authentication providers configured for the cluster.
//
// Only one external authentication provider can be configured.
func (o *ExternalAuthConfig) GetExternalAuths() (value *ExternalAuthList, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.externalAuths
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Controls whether the cluster uses an external authentication configuration for ARO HCP clusters.
//
// For ARO HCP clusters, this will be "enabled" by default and cannot be set to "disabled".
//
// FOR ROSA HCP clusters, use the "Enabled" boolean flag to enable/disable this feature instead.
func (o *ExternalAuthConfig) State() ExternalAuthConfigState {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.state
	}
	return ExternalAuthConfigState("")
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
//
// Controls whether the cluster uses an external authentication configuration for ARO HCP clusters.
//
// For ARO HCP clusters, this will be "enabled" by default and cannot be set to "disabled".
//
// FOR ROSA HCP clusters, use the "Enabled" boolean flag to enable/disable this feature instead.
func (o *ExternalAuthConfig) GetState() (value ExternalAuthConfigState, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.state
	}
	return
}

// ExternalAuthConfigListKind is the name of the type used to represent list of objects of
// type 'external_auth_config'.
const ExternalAuthConfigListKind = "ExternalAuthConfigList"

// ExternalAuthConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'external_auth_config'.
const ExternalAuthConfigListLinkKind = "ExternalAuthConfigListLink"

// ExternalAuthConfigNilKind is the name of the type used to nil lists of objects of
// type 'external_auth_config'.
const ExternalAuthConfigListNilKind = "ExternalAuthConfigListNil"

// ExternalAuthConfigList is a list of values of the 'external_auth_config' type.
type ExternalAuthConfigList struct {
	href  string
	link  bool
	items []*ExternalAuthConfig
}

// Kind returns the name of the type of the object.
func (l *ExternalAuthConfigList) Kind() string {
	if l == nil {
		return ExternalAuthConfigListNilKind
	}
	if l.link {
		return ExternalAuthConfigListLinkKind
	}
	return ExternalAuthConfigListKind
}

// Link returns true iif this is a link.
func (l *ExternalAuthConfigList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ExternalAuthConfigList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ExternalAuthConfigList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ExternalAuthConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ExternalAuthConfigList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ExternalAuthConfigList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ExternalAuthConfigList) SetItems(items []*ExternalAuthConfig) {
	l.items = items
}

// Items returns the items of the list.
func (l *ExternalAuthConfigList) Items() []*ExternalAuthConfig {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ExternalAuthConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ExternalAuthConfigList) Get(i int) *ExternalAuthConfig {
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
func (l *ExternalAuthConfigList) Slice() []*ExternalAuthConfig {
	var slice []*ExternalAuthConfig
	if l == nil {
		slice = make([]*ExternalAuthConfig, 0)
	} else {
		slice = make([]*ExternalAuthConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ExternalAuthConfigList) Each(f func(item *ExternalAuthConfig) bool) {
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
func (l *ExternalAuthConfigList) Range(f func(index int, item *ExternalAuthConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
