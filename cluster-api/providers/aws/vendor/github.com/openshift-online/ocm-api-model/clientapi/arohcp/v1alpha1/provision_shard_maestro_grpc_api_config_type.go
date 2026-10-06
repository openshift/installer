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

// ProvisionShardMaestroGrpcApiConfig represents the values of the 'provision_shard_maestro_grpc_api_config' type.
//
// The Maestro GRPC API configuration of the provision shard.
type ProvisionShardMaestroGrpcApiConfig struct {
	fieldSet_ []bool
	url       string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ProvisionShardMaestroGrpcApiConfig) Empty() bool {
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
// A GRPC endpoint to the Maestro Server GRPC API used by the shard.
// The GRPC endpoint must point to the same Maestro Server as
// the Maestro Server associated to the Maestro Server REST API URL
// specified in `rest_api_config.url`.
// The URL must be a well-formed absolute URL.
// The expected url format naming scheme is: <host>:<port>.
// where both <host> and <port> must be specified and they cannot
// be empty. No whitespace characters are allowed. No URI scheme
// is allowed as part of the URL.
// The <host>:<port> combination cannot be the same as in
// `rest_api_config.url`. This includes different DNS names
// pointing to the same underlying server.
// Example of a URL: maestro.example.com:50051
// Required during creation.
// Immutable.
func (o *ProvisionShardMaestroGrpcApiConfig) Url() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.url
	}
	return ""
}

// GetUrl returns the value of the 'url' attribute and
// a flag indicating if the attribute has a value.
//
// A GRPC endpoint to the Maestro Server GRPC API used by the shard.
// The GRPC endpoint must point to the same Maestro Server as
// the Maestro Server associated to the Maestro Server REST API URL
// specified in `rest_api_config.url`.
// The URL must be a well-formed absolute URL.
// The expected url format naming scheme is: <host>:<port>.
// where both <host> and <port> must be specified and they cannot
// be empty. No whitespace characters are allowed. No URI scheme
// is allowed as part of the URL.
// The <host>:<port> combination cannot be the same as in
// `rest_api_config.url`. This includes different DNS names
// pointing to the same underlying server.
// Example of a URL: maestro.example.com:50051
// Required during creation.
// Immutable.
func (o *ProvisionShardMaestroGrpcApiConfig) GetUrl() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.url
	}
	return
}

// ProvisionShardMaestroGrpcApiConfigListKind is the name of the type used to represent list of objects of
// type 'provision_shard_maestro_grpc_api_config'.
const ProvisionShardMaestroGrpcApiConfigListKind = "ProvisionShardMaestroGrpcApiConfigList"

// ProvisionShardMaestroGrpcApiConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'provision_shard_maestro_grpc_api_config'.
const ProvisionShardMaestroGrpcApiConfigListLinkKind = "ProvisionShardMaestroGrpcApiConfigListLink"

// ProvisionShardMaestroGrpcApiConfigNilKind is the name of the type used to nil lists of objects of
// type 'provision_shard_maestro_grpc_api_config'.
const ProvisionShardMaestroGrpcApiConfigListNilKind = "ProvisionShardMaestroGrpcApiConfigListNil"

// ProvisionShardMaestroGrpcApiConfigList is a list of values of the 'provision_shard_maestro_grpc_api_config' type.
type ProvisionShardMaestroGrpcApiConfigList struct {
	href  string
	link  bool
	items []*ProvisionShardMaestroGrpcApiConfig
}

// Len returns the length of the list.
func (l *ProvisionShardMaestroGrpcApiConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ProvisionShardMaestroGrpcApiConfigList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ProvisionShardMaestroGrpcApiConfigList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ProvisionShardMaestroGrpcApiConfigList) SetItems(items []*ProvisionShardMaestroGrpcApiConfig) {
	l.items = items
}

// Items returns the items of the list.
func (l *ProvisionShardMaestroGrpcApiConfigList) Items() []*ProvisionShardMaestroGrpcApiConfig {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ProvisionShardMaestroGrpcApiConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ProvisionShardMaestroGrpcApiConfigList) Get(i int) *ProvisionShardMaestroGrpcApiConfig {
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
func (l *ProvisionShardMaestroGrpcApiConfigList) Slice() []*ProvisionShardMaestroGrpcApiConfig {
	var slice []*ProvisionShardMaestroGrpcApiConfig
	if l == nil {
		slice = make([]*ProvisionShardMaestroGrpcApiConfig, 0)
	} else {
		slice = make([]*ProvisionShardMaestroGrpcApiConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ProvisionShardMaestroGrpcApiConfigList) Each(f func(item *ProvisionShardMaestroGrpcApiConfig) bool) {
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
func (l *ProvisionShardMaestroGrpcApiConfigList) Range(f func(index int, item *ProvisionShardMaestroGrpcApiConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
