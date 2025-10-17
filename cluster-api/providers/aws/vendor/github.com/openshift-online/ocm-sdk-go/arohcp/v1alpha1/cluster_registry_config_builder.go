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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

// ClusterRegistryConfigBuilder contains the data and logic needed to build 'cluster_registry_config' objects.
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
type ClusterRegistryConfigBuilder struct {
	bitmap_                    uint32
	additionalTrustedCa        map[string]string
	allowedRegistriesForImport []*RegistryLocationBuilder
	platformAllowlist          *RegistryAllowlistBuilder
	registrySources            *RegistrySourcesBuilder
}

// NewClusterRegistryConfig creates a new builder of 'cluster_registry_config' objects.
func NewClusterRegistryConfig() *ClusterRegistryConfigBuilder {
	return &ClusterRegistryConfigBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterRegistryConfigBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AdditionalTrustedCa sets the value of the 'additional_trusted_ca' attribute to the given value.
func (b *ClusterRegistryConfigBuilder) AdditionalTrustedCa(value map[string]string) *ClusterRegistryConfigBuilder {
	b.additionalTrustedCa = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// AllowedRegistriesForImport sets the value of the 'allowed_registries_for_import' attribute to the given values.
func (b *ClusterRegistryConfigBuilder) AllowedRegistriesForImport(values ...*RegistryLocationBuilder) *ClusterRegistryConfigBuilder {
	b.allowedRegistriesForImport = make([]*RegistryLocationBuilder, len(values))
	copy(b.allowedRegistriesForImport, values)
	b.bitmap_ |= 2
	return b
}

// PlatformAllowlist sets the value of the 'platform_allowlist' attribute to the given value.
//
// RegistryAllowlist represents a single registry allowlist.
func (b *ClusterRegistryConfigBuilder) PlatformAllowlist(value *RegistryAllowlistBuilder) *ClusterRegistryConfigBuilder {
	b.platformAllowlist = value
	if value != nil {
		b.bitmap_ |= 4
	} else {
		b.bitmap_ &^= 4
	}
	return b
}

// RegistrySources sets the value of the 'registry_sources' attribute to the given value.
//
// RegistrySources contains configuration that determines how the container runtime should treat individual
// registries when accessing images for builds and pods. For instance, whether or not to allow insecure access.
// It does not contain configuration for the internal cluster registry.
func (b *ClusterRegistryConfigBuilder) RegistrySources(value *RegistrySourcesBuilder) *ClusterRegistryConfigBuilder {
	b.registrySources = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterRegistryConfigBuilder) Copy(object *ClusterRegistryConfig) *ClusterRegistryConfigBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if len(object.additionalTrustedCa) > 0 {
		b.additionalTrustedCa = map[string]string{}
		for k, v := range object.additionalTrustedCa {
			b.additionalTrustedCa[k] = v
		}
	} else {
		b.additionalTrustedCa = nil
	}
	if object.allowedRegistriesForImport != nil {
		b.allowedRegistriesForImport = make([]*RegistryLocationBuilder, len(object.allowedRegistriesForImport))
		for i, v := range object.allowedRegistriesForImport {
			b.allowedRegistriesForImport[i] = NewRegistryLocation().Copy(v)
		}
	} else {
		b.allowedRegistriesForImport = nil
	}
	if object.platformAllowlist != nil {
		b.platformAllowlist = NewRegistryAllowlist().Copy(object.platformAllowlist)
	} else {
		b.platformAllowlist = nil
	}
	if object.registrySources != nil {
		b.registrySources = NewRegistrySources().Copy(object.registrySources)
	} else {
		b.registrySources = nil
	}
	return b
}

// Build creates a 'cluster_registry_config' object using the configuration stored in the builder.
func (b *ClusterRegistryConfigBuilder) Build() (object *ClusterRegistryConfig, err error) {
	object = new(ClusterRegistryConfig)
	object.bitmap_ = b.bitmap_
	if b.additionalTrustedCa != nil {
		object.additionalTrustedCa = make(map[string]string)
		for k, v := range b.additionalTrustedCa {
			object.additionalTrustedCa[k] = v
		}
	}
	if b.allowedRegistriesForImport != nil {
		object.allowedRegistriesForImport = make([]*RegistryLocation, len(b.allowedRegistriesForImport))
		for i, v := range b.allowedRegistriesForImport {
			object.allowedRegistriesForImport[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	if b.platformAllowlist != nil {
		object.platformAllowlist, err = b.platformAllowlist.Build()
		if err != nil {
			return
		}
	}
	if b.registrySources != nil {
		object.registrySources, err = b.registrySources.Build()
		if err != nil {
			return
		}
	}
	return
}
