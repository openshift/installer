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

// ProvisionShardMaestroRestApiConfig represents the values of the 'provision_shard_maestro_rest_api_config' type.
//
// The Maestro REST API configuration of the provision shard.
type ProvisionShardMaestroRestApiConfig struct {
	fieldSet_ []bool
	url       string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ProvisionShardMaestroRestApiConfig) Empty() bool {
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

// Url returns the value of the 'url' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// A URL to the Maestro Server REST API used by the shard.
// The URL must point to the same Maestro Server as the
// Maestro Server associated to the Maestro Server GRPC API
// endpoint specified in `grpc_api_config.url`
// The URL must be a well-formed absolute URL.
// The expected url format naming scheme is: <scheme>://<host>:<port>
// where both <scheme>, <host> and <port> must be specified.
// The <host>:<port> combination cannot be the same as in
// `grpc_api_config.url`. This includes different DNS names
// pointing to the same underlying server.
// Example of URL: https://maestro.example.com:50052
// Required during creation.
// Immutable.
func (o *ProvisionShardMaestroRestApiConfig) Url() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.url
	}
	return ""
}

// GetUrl returns the value of the 'url' attribute and
// a flag indicating if the attribute has a value.
//
// A URL to the Maestro Server REST API used by the shard.
// The URL must point to the same Maestro Server as the
// Maestro Server associated to the Maestro Server GRPC API
// endpoint specified in `grpc_api_config.url`
// The URL must be a well-formed absolute URL.
// The expected url format naming scheme is: <scheme>://<host>:<port>
// where both <scheme>, <host> and <port> must be specified.
// The <host>:<port> combination cannot be the same as in
// `grpc_api_config.url`. This includes different DNS names
// pointing to the same underlying server.
// Example of URL: https://maestro.example.com:50052
// Required during creation.
// Immutable.
func (o *ProvisionShardMaestroRestApiConfig) GetUrl() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.url
	}
	return
}

// ProvisionShardMaestroRestApiConfigListKind is the name of the type used to represent list of objects of
// type 'provision_shard_maestro_rest_api_config'.
const ProvisionShardMaestroRestApiConfigListKind = "ProvisionShardMaestroRestApiConfigList"

// ProvisionShardMaestroRestApiConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'provision_shard_maestro_rest_api_config'.
const ProvisionShardMaestroRestApiConfigListLinkKind = "ProvisionShardMaestroRestApiConfigListLink"

// ProvisionShardMaestroRestApiConfigNilKind is the name of the type used to nil lists of objects of
// type 'provision_shard_maestro_rest_api_config'.
const ProvisionShardMaestroRestApiConfigListNilKind = "ProvisionShardMaestroRestApiConfigListNil"

// ProvisionShardMaestroRestApiConfigList is a list of values of the 'provision_shard_maestro_rest_api_config' type.
type ProvisionShardMaestroRestApiConfigList struct {
	href  string
	link  bool
	items []*ProvisionShardMaestroRestApiConfig
}

// Len returns the length of the list.
func (l *ProvisionShardMaestroRestApiConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ProvisionShardMaestroRestApiConfigList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ProvisionShardMaestroRestApiConfigList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ProvisionShardMaestroRestApiConfigList) SetItems(items []*ProvisionShardMaestroRestApiConfig) {
	l.items = items
}

// Items returns the items of the list.
func (l *ProvisionShardMaestroRestApiConfigList) Items() []*ProvisionShardMaestroRestApiConfig {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ProvisionShardMaestroRestApiConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ProvisionShardMaestroRestApiConfigList) Get(i int) *ProvisionShardMaestroRestApiConfig {
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
func (l *ProvisionShardMaestroRestApiConfigList) Slice() []*ProvisionShardMaestroRestApiConfig {
	var slice []*ProvisionShardMaestroRestApiConfig
	if l == nil {
		slice = make([]*ProvisionShardMaestroRestApiConfig, 0)
	} else {
		slice = make([]*ProvisionShardMaestroRestApiConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ProvisionShardMaestroRestApiConfigList) Each(f func(item *ProvisionShardMaestroRestApiConfig) bool) {
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
func (l *ProvisionShardMaestroRestApiConfigList) Range(f func(index int, item *ProvisionShardMaestroRestApiConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
