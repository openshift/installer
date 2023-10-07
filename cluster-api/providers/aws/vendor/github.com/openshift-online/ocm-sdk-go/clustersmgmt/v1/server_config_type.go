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

// ServerConfigKind is the name of the type used to represent objects
// of type 'server_config'.
const ServerConfigKind = "ServerConfig"

// ServerConfigLinkKind is the name of the type used to represent links
// to objects of type 'server_config'.
const ServerConfigLinkKind = "ServerConfigLink"

// ServerConfigNilKind is the name of the type used to nil references
// to objects of type 'server_config'.
const ServerConfigNilKind = "ServerConfigNil"

// ServerConfig represents the values of the 'server_config' type.
//
// Representation of a server config
type ServerConfig struct {
	bitmap_    uint32
	id         string
	href       string
	kubeconfig string
	server     string
}

// Kind returns the name of the type of the object.
func (o *ServerConfig) Kind() string {
	if o == nil {
		return ServerConfigNilKind
	}
	if o.bitmap_&1 != 0 {
		return ServerConfigLinkKind
	}
	return ServerConfigKind
}

// Link returns true iif this is a link.
func (o *ServerConfig) Link() bool {
	return o != nil && o.bitmap_&1 != 0
}

// ID returns the identifier of the object.
func (o *ServerConfig) ID() string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *ServerConfig) GetID() (value string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *ServerConfig) HREF() string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *ServerConfig) GetHREF() (value string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ServerConfig) Empty() bool {
	return o == nil || o.bitmap_&^1 == 0
}

// Kubeconfig returns the value of the 'kubeconfig' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The kubeconfig of the server
func (o *ServerConfig) Kubeconfig() string {
	if o != nil && o.bitmap_&8 != 0 {
		return o.kubeconfig
	}
	return ""
}

// GetKubeconfig returns the value of the 'kubeconfig' attribute and
// a flag indicating if the attribute has a value.
//
// The kubeconfig of the server
func (o *ServerConfig) GetKubeconfig() (value string, ok bool) {
	ok = o != nil && o.bitmap_&8 != 0
	if ok {
		value = o.kubeconfig
	}
	return
}

// Server returns the value of the 'server' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The URL of the server
func (o *ServerConfig) Server() string {
	if o != nil && o.bitmap_&16 != 0 {
		return o.server
	}
	return ""
}

// GetServer returns the value of the 'server' attribute and
// a flag indicating if the attribute has a value.
//
// The URL of the server
func (o *ServerConfig) GetServer() (value string, ok bool) {
	ok = o != nil && o.bitmap_&16 != 0
	if ok {
		value = o.server
	}
	return
}

// ServerConfigListKind is the name of the type used to represent list of objects of
// type 'server_config'.
const ServerConfigListKind = "ServerConfigList"

// ServerConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'server_config'.
const ServerConfigListLinkKind = "ServerConfigListLink"

// ServerConfigNilKind is the name of the type used to nil lists of objects of
// type 'server_config'.
const ServerConfigListNilKind = "ServerConfigListNil"

// ServerConfigList is a list of values of the 'server_config' type.
type ServerConfigList struct {
	href  string
	link  bool
	items []*ServerConfig
}

// Kind returns the name of the type of the object.
func (l *ServerConfigList) Kind() string {
	if l == nil {
		return ServerConfigListNilKind
	}
	if l.link {
		return ServerConfigListLinkKind
	}
	return ServerConfigListKind
}

// Link returns true iif this is a link.
func (l *ServerConfigList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ServerConfigList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ServerConfigList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ServerConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *ServerConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ServerConfigList) Get(i int) *ServerConfig {
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
func (l *ServerConfigList) Slice() []*ServerConfig {
	var slice []*ServerConfig
	if l == nil {
		slice = make([]*ServerConfig, 0)
	} else {
		slice = make([]*ServerConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ServerConfigList) Each(f func(item *ServerConfig) bool) {
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
func (l *ServerConfigList) Range(f func(index int, item *ServerConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
