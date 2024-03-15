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

// ExternalAuthClientConfig represents the values of the 'external_auth_client_config' type.
//
// ExternalAuthClientConfig contains configuration for the platform's clients that
// need to request tokens from the issuer.
type ExternalAuthClientConfig struct {
	bitmap_     uint32
	id          string
	component   *ClientComponent
	extraScopes []string
	secret      string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ExternalAuthClientConfig) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// ID returns the value of the 'ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The identifier of the OIDC client from the OIDC provider.
func (o *ExternalAuthClientConfig) ID() string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the value of the 'ID' attribute and
// a flag indicating if the attribute has a value.
//
// The identifier of the OIDC client from the OIDC provider.
func (o *ExternalAuthClientConfig) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.id
	}
	return
}

// Component returns the value of the 'component' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The component that is supposed to consume this client configuration.
func (o *ExternalAuthClientConfig) Component() *ClientComponent {
	if o != nil && o.bitmap_&2 != 0 {
		return o.component
	}
	return nil
}

// GetComponent returns the value of the 'component' attribute and
// a flag indicating if the attribute has a value.
//
// The component that is supposed to consume this client configuration.
func (o *ExternalAuthClientConfig) GetComponent() (value *ClientComponent, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.component
	}
	return
}

// ExtraScopes returns the value of the 'extra_scopes' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ExtraScopes is an optional set of scopes to request tokens with.
func (o *ExternalAuthClientConfig) ExtraScopes() []string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.extraScopes
	}
	return nil
}

// GetExtraScopes returns the value of the 'extra_scopes' attribute and
// a flag indicating if the attribute has a value.
//
// ExtraScopes is an optional set of scopes to request tokens with.
func (o *ExternalAuthClientConfig) GetExtraScopes() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.extraScopes
	}
	return
}

// Secret returns the value of the 'secret' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The secret of the OIDC client from the OIDC provider.
func (o *ExternalAuthClientConfig) Secret() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.secret
	}
	return ""
}

// GetSecret returns the value of the 'secret' attribute and
// a flag indicating if the attribute has a value.
//
// The secret of the OIDC client from the OIDC provider.
func (o *ExternalAuthClientConfig) GetSecret() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.secret
	}
	return
}

// ExternalAuthClientConfigListKind is the name of the type used to represent list of objects of
// type 'external_auth_client_config'.
const ExternalAuthClientConfigListKind = "ExternalAuthClientConfigList"

// ExternalAuthClientConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'external_auth_client_config'.
const ExternalAuthClientConfigListLinkKind = "ExternalAuthClientConfigListLink"

// ExternalAuthClientConfigNilKind is the name of the type used to nil lists of objects of
// type 'external_auth_client_config'.
const ExternalAuthClientConfigListNilKind = "ExternalAuthClientConfigListNil"

// ExternalAuthClientConfigList is a list of values of the 'external_auth_client_config' type.
type ExternalAuthClientConfigList struct {
	href  string
	link  bool
	items []*ExternalAuthClientConfig
}

// Len returns the length of the list.
func (l *ExternalAuthClientConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ExternalAuthClientConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ExternalAuthClientConfigList) Get(i int) *ExternalAuthClientConfig {
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
func (l *ExternalAuthClientConfigList) Slice() []*ExternalAuthClientConfig {
	var slice []*ExternalAuthClientConfig
	if l == nil {
		slice = make([]*ExternalAuthClientConfig, 0)
	} else {
		slice = make([]*ExternalAuthClientConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ExternalAuthClientConfigList) Each(f func(item *ExternalAuthClientConfig) bool) {
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
func (l *ExternalAuthClientConfigList) Range(f func(index int, item *ExternalAuthClientConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
