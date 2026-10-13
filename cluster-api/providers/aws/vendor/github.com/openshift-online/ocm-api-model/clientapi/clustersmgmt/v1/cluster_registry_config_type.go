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

// ClusterRegistryConfig represents the values of the 'cluster_registry_config' type.
//
// ClusterRegistryConfig describes the configuration of registries for the cluster.
// Its format reflects the OpenShift Image Configuration, for which docs are available on
// [docs.openshift.com](https://docs.openshift.com/container-platform/4.16/openshift_images/image-configuration.html)
// ```json
//
//	{
//	   "registry_config": {
//	     "registry_sources": {
//	       "blocked_registries": [
//	         "badregistry.io",
//	         "badregistry8.io"
//	       ]
//	     }
//	   }
//	}
//
// ```
type ClusterRegistryConfig struct {
	fieldSet_                  []bool
	additionalTrustedCa        map[string]string
	allowedRegistriesForImport []*RegistryLocation
	platformAllowlist          *RegistryAllowlist
	registrySources            *RegistrySources
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *ClusterRegistryConfig) Empty() bool {
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

// AdditionalTrustedCa returns the value of the 'additional_trusted_ca' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// A map containing the registry hostname as the key, and the PEM-encoded certificate as the value,
// for each additional registry CA to trust.
func (o *ClusterRegistryConfig) AdditionalTrustedCa() map[string]string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.additionalTrustedCa
	}
	return nil
}

// GetAdditionalTrustedCa returns the value of the 'additional_trusted_ca' attribute and
// a flag indicating if the attribute has a value.
//
// A map containing the registry hostname as the key, and the PEM-encoded certificate as the value,
// for each additional registry CA to trust.
func (o *ClusterRegistryConfig) GetAdditionalTrustedCa() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.additionalTrustedCa
	}
	return
}

// AllowedRegistriesForImport returns the value of the 'allowed_registries_for_import' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// AllowedRegistriesForImport limits the container image registries that normal users may import
// images from. Set this list to the registries that you trust to contain valid Docker
// images and that you want applications to be able to import from. Users with
// permission to create Images or ImageStreamMappings via the API are not affected by
// this policy - typically only administrators or system integrations will have those
// permissions.
func (o *ClusterRegistryConfig) AllowedRegistriesForImport() []*RegistryLocation {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.allowedRegistriesForImport
	}
	return nil
}

// GetAllowedRegistriesForImport returns the value of the 'allowed_registries_for_import' attribute and
// a flag indicating if the attribute has a value.
//
// AllowedRegistriesForImport limits the container image registries that normal users may import
// images from. Set this list to the registries that you trust to contain valid Docker
// images and that you want applications to be able to import from. Users with
// permission to create Images or ImageStreamMappings via the API are not affected by
// this policy - typically only administrators or system integrations will have those
// permissions.
func (o *ClusterRegistryConfig) GetAllowedRegistriesForImport() (value []*RegistryLocation, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.allowedRegistriesForImport
	}
	return
}

// PlatformAllowlist returns the value of the 'platform_allowlist' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// PlatformAllowlist contains a reference to a RegistryAllowlist which is a list of internal registries
// which needs to be whitelisted for the platform to work. It can be omitted at creation and
// updating and its lifecycle can be managed separately if needed.
func (o *ClusterRegistryConfig) PlatformAllowlist() *RegistryAllowlist {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.platformAllowlist
	}
	return nil
}

// GetPlatformAllowlist returns the value of the 'platform_allowlist' attribute and
// a flag indicating if the attribute has a value.
//
// PlatformAllowlist contains a reference to a RegistryAllowlist which is a list of internal registries
// which needs to be whitelisted for the platform to work. It can be omitted at creation and
// updating and its lifecycle can be managed separately if needed.
func (o *ClusterRegistryConfig) GetPlatformAllowlist() (value *RegistryAllowlist, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.platformAllowlist
	}
	return
}

// RegistrySources returns the value of the 'registry_sources' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// RegistrySources contains configuration that determines how the container runtime
// should treat individual registries when accessing images for builds+pods. (e.g.
// whether or not to allow insecure access). It does not contain configuration for the
// internal cluster registry.
func (o *ClusterRegistryConfig) RegistrySources() *RegistrySources {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.registrySources
	}
	return nil
}

// GetRegistrySources returns the value of the 'registry_sources' attribute and
// a flag indicating if the attribute has a value.
//
// RegistrySources contains configuration that determines how the container runtime
// should treat individual registries when accessing images for builds+pods. (e.g.
// whether or not to allow insecure access). It does not contain configuration for the
// internal cluster registry.
func (o *ClusterRegistryConfig) GetRegistrySources() (value *RegistrySources, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.registrySources
	}
	return
}

// ClusterRegistryConfigListKind is the name of the type used to represent list of objects of
// type 'cluster_registry_config'.
const ClusterRegistryConfigListKind = "ClusterRegistryConfigList"

// ClusterRegistryConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster_registry_config'.
const ClusterRegistryConfigListLinkKind = "ClusterRegistryConfigListLink"

// ClusterRegistryConfigNilKind is the name of the type used to nil lists of objects of
// type 'cluster_registry_config'.
const ClusterRegistryConfigListNilKind = "ClusterRegistryConfigListNil"

// ClusterRegistryConfigList is a list of values of the 'cluster_registry_config' type.
type ClusterRegistryConfigList struct {
	href  string
	link  bool
	items []*ClusterRegistryConfig
}

// Len returns the length of the list.
func (l *ClusterRegistryConfigList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ClusterRegistryConfigList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ClusterRegistryConfigList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ClusterRegistryConfigList) SetItems(items []*ClusterRegistryConfig) {
	l.items = items
}

// Items returns the items of the list.
func (l *ClusterRegistryConfigList) Items() []*ClusterRegistryConfig {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ClusterRegistryConfigList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterRegistryConfigList) Get(i int) *ClusterRegistryConfig {
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
func (l *ClusterRegistryConfigList) Slice() []*ClusterRegistryConfig {
	var slice []*ClusterRegistryConfig
	if l == nil {
		slice = make([]*ClusterRegistryConfig, 0)
	} else {
		slice = make([]*ClusterRegistryConfig, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterRegistryConfigList) Each(f func(item *ClusterRegistryConfig) bool) {
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
func (l *ClusterRegistryConfigList) Range(f func(index int, item *ClusterRegistryConfig) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
