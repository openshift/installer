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

// ProvisionShardMaestroConfig represents the values of the 'provision_shard_maestro_config' type.
//
// The Maestro related configuration of the Provision Shard.
// The combination of `consumer_name` and `rest_api_config.url`
// must be unique across shards.
// The combination of `consumer_name` and `grpc_api_config.url`
// must be unique across shards.
type ProvisionShardMaestroConfig struct {
	fieldSet_     []bool
	consumerName  string
	grpcApiConfig *ProvisionShardMaestroGrpcApiConfig
	restApiConfig *ProvisionShardMaestroRestApiConfig
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ProvisionShardMaestroConfig) Empty() bool {
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

// ConsumerName returns the value of the 'consumer_name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The name of the Maestro Consumer used by the shard.
// A Maestro Consumer with the provided name must pre-exist
// and be pre-registered in the Maestro Server associated
// to the REST and GRPC Maestro URLs specified in `rest_api_config.url`
// and `grpc_api_config.url`.
// Required during creation.
// Immutable.
func (o *ProvisionShardMaestroConfig) ConsumerName() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.consumerName
	}
	return ""
}

// GetConsumerName returns the value of the 'consumer_name' attribute and
// a flag indicating if the attribute has a value.
//
// The name of the Maestro Consumer used by the shard.
// A Maestro Consumer with the provided name must pre-exist
// and be pre-registered in the Maestro Server associated
// to the REST and GRPC Maestro URLs specified in `rest_api_config.url`
// and `grpc_api_config.url`.
// Required during creation.
// Immutable.
func (o *ProvisionShardMaestroConfig) GetConsumerName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.consumerName
	}
	return
}

// GrpcApiConfig returns the value of the 'grpc_api_config' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Maestro Server GRPC API configuration used by the shard.
// Required during creation.
func (o *ProvisionShardMaestroConfig) GrpcApiConfig() *ProvisionShardMaestroGrpcApiConfig {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.grpcApiConfig
	}
	return nil
}

// GetGrpcApiConfig returns the value of the 'grpc_api_config' attribute and
// a flag indicating if the attribute has a value.
//
// The Maestro Server GRPC API configuration used by the shard.
// Required during creation.
func (o *ProvisionShardMaestroConfig) GetGrpcApiConfig() (value *ProvisionShardMaestroGrpcApiConfig, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.grpcApiConfig
	}
	return
}

// RestApiConfig returns the value of the 'rest_api_config' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Maestro Server REST API configuration used by the shard.
// Required during creation.
func (o *ProvisionShardMaestroConfig) RestApiConfig() *ProvisionShardMaestroRestApiConfig {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.restApiConfig
	}
	return nil
}

// GetRestApiConfig returns the value of the 'rest_api_config' attribute and
// a flag indicating if the attribute has a value.
//
// The Maestro Server REST API configuration used by the shard.
// Required during creation.
func (o *ProvisionShardMaestroConfig) GetRestApiConfig() (value *ProvisionShardMaestroRestApiConfig, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.restApiConfig
	}
	return
}

// ProvisionShardMaestroConfigListKind is the name of the type used to represent list of objects of
// type 'provision_shard_maestro_config'.
const ProvisionShardMaestroConfigListKind = "ProvisionShardMaestroConfigList"

// ProvisionShardMaestroConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'provision_shard_maestro_config'.
const ProvisionShardMaestroConfigListLinkKind = "ProvisionShardMaestroConfigListLink"

// ProvisionShardMaestroConfigNilKind is the name of the type used to nil lists of objects of
// type 'provision_shard_maestro_config'.
const ProvisionShardMaestroConfigListNilKind = "ProvisionShardMaestroConfigListNil"

// ProvisionShardMaestroConfigList is a list of values of the 'provision_shard_maestro_config' type.
type ProvisionShardMaestroConfigList struct {
	href  string
	link  bool
	items []*ProvisionShardMaestroConfig
}

// Len returns the length of the list.
func (l *ProvisionShardMaestroConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ProvisionShardMaestroConfigList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ProvisionShardMaestroConfigList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ProvisionShardMaestroConfigList) SetItems(items []*ProvisionShardMaestroConfig) {
	l.items = items
}

// Items returns the items of the list.
func (l *ProvisionShardMaestroConfigList) Items() []*ProvisionShardMaestroConfig {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ProvisionShardMaestroConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ProvisionShardMaestroConfigList) Get(i int) *ProvisionShardMaestroConfig {
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
func (l *ProvisionShardMaestroConfigList) Slice() []*ProvisionShardMaestroConfig {
	var slice []*ProvisionShardMaestroConfig
	if l == nil {
		slice = make([]*ProvisionShardMaestroConfig, 0)
	} else {
		slice = make([]*ProvisionShardMaestroConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ProvisionShardMaestroConfigList) Each(f func(item *ProvisionShardMaestroConfig) bool) {
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
func (l *ProvisionShardMaestroConfigList) Range(f func(index int, item *ProvisionShardMaestroConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
