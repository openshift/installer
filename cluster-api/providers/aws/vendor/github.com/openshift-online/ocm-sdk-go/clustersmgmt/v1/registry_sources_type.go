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

// RegistrySources represents the values of the 'registry_sources' type.
//
// RegistrySources contains configuration that determines how the container runtime should treat individual
// registries when accessing images for builds and pods. For instance, whether or not to allow insecure access.
// It does not contain configuration for the internal cluster registry.
type RegistrySources struct {
	bitmap_            uint32
	allowedRegistries  []string
	blockedRegistries  []string
	insecureRegistries []string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *RegistrySources) Empty() bool {
	return o == nil || o.bitmap_ == 0
}

// AllowedRegistries returns the value of the 'allowed_registries' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AllowedRegistries: registries for which image pull and push actions are allowed.
// To specify all subdomains, add the asterisk (*) wildcard character as a prefix to the domain name.
// For example, *.example.com. You can specify an individual repository within a registry.
// For example: reg1.io/myrepo/myapp:latest. All other registries are blocked.
// Mutually exclusive with `BlockedRegistries`
func (o *RegistrySources) AllowedRegistries() []string {
	if o != nil && o.bitmap_&1 != 0 {
		return o.allowedRegistries
	}
	return nil
}

// GetAllowedRegistries returns the value of the 'allowed_registries' attribute and
// a flag indicating if the attribute has a value.
//
// AllowedRegistries: registries for which image pull and push actions are allowed.
// To specify all subdomains, add the asterisk (*) wildcard character as a prefix to the domain name.
// For example, *.example.com. You can specify an individual repository within a registry.
// For example: reg1.io/myrepo/myapp:latest. All other registries are blocked.
// Mutually exclusive with `BlockedRegistries`
func (o *RegistrySources) GetAllowedRegistries() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&1 != 0
	if ok {
		value = o.allowedRegistries
	}
	return
}

// BlockedRegistries returns the value of the 'blocked_registries' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// BlockedRegistries: registries for which image pull and push actions are denied.
// To specify all subdomains, add the asterisk (*) wildcard character as a prefix to the domain name.
// For example, *.example.com. You can specify an individual repository within a registry.
// For example: reg1.io/myrepo/myapp:latest. All other registries are allowed.
// Mutually exclusive with `AllowedRegistries`
func (o *RegistrySources) BlockedRegistries() []string {
	if o != nil && o.bitmap_&2 != 0 {
		return o.blockedRegistries
	}
	return nil
}

// GetBlockedRegistries returns the value of the 'blocked_registries' attribute and
// a flag indicating if the attribute has a value.
//
// BlockedRegistries: registries for which image pull and push actions are denied.
// To specify all subdomains, add the asterisk (*) wildcard character as a prefix to the domain name.
// For example, *.example.com. You can specify an individual repository within a registry.
// For example: reg1.io/myrepo/myapp:latest. All other registries are allowed.
// Mutually exclusive with `AllowedRegistries`
func (o *RegistrySources) GetBlockedRegistries() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&2 != 0
	if ok {
		value = o.blockedRegistries
	}
	return
}

// InsecureRegistries returns the value of the 'insecure_registries' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// InsecureRegistries are registries which do not have a valid TLS certificate or only support HTTP connections.
// To specify all subdomains, add the asterisk (*) wildcard character as a prefix to the domain name.
// For example, *.example.com. You can specify an individual repository within a registry.
// For example: reg1.io/myrepo/myapp:latest.
func (o *RegistrySources) InsecureRegistries() []string {
	if o != nil && o.bitmap_&4 != 0 {
		return o.insecureRegistries
	}
	return nil
}

// GetInsecureRegistries returns the value of the 'insecure_registries' attribute and
// a flag indicating if the attribute has a value.
//
// InsecureRegistries are registries which do not have a valid TLS certificate or only support HTTP connections.
// To specify all subdomains, add the asterisk (*) wildcard character as a prefix to the domain name.
// For example, *.example.com. You can specify an individual repository within a registry.
// For example: reg1.io/myrepo/myapp:latest.
func (o *RegistrySources) GetInsecureRegistries() (value []string, ok bool) {
	ok = o != nil && o.bitmap_&4 != 0
	if ok {
		value = o.insecureRegistries
	}
	return
}

// RegistrySourcesListKind is the name of the type used to represent list of objects of
// type 'registry_sources'.
const RegistrySourcesListKind = "RegistrySourcesList"

// RegistrySourcesListLinkKind is the name of the type used to represent links to list
// of objects of type 'registry_sources'.
const RegistrySourcesListLinkKind = "RegistrySourcesListLink"

// RegistrySourcesNilKind is the name of the type used to nil lists of objects of
// type 'registry_sources'.
const RegistrySourcesListNilKind = "RegistrySourcesListNil"

// RegistrySourcesList is a list of values of the 'registry_sources' type.
type RegistrySourcesList struct {
	href  string
	link  bool
	items []*RegistrySources
}

// Len returns the length of the list.
func (l *RegistrySourcesList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Empty returns true if the list is empty.
func (l *RegistrySourcesList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *RegistrySourcesList) Get(i int) *RegistrySources {
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
func (l *RegistrySourcesList) Slice() []*RegistrySources {
	var slice []*RegistrySources
	if l == nil {
		slice = make([]*RegistrySources, 0)
	} else {
		slice = make([]*RegistrySources, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *RegistrySourcesList) Each(f func(item *RegistrySources) bool) {
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
func (l *RegistrySourcesList) Range(f func(index int, item *RegistrySources) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
