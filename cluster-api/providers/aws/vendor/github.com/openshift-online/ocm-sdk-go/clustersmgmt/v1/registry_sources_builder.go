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

// RegistrySourcesBuilder contains the data and logic needed to build 'registry_sources' objects.
//
// RegistrySources contains configuration that determines how the container runtime should treat individual
// registries when accessing images for builds and pods. For instance, whether or not to allow insecure access.
// It does not contain configuration for the internal cluster registry.
type RegistrySourcesBuilder struct {
	bitmap_            uint32
	allowedRegistries  []string
	blockedRegistries  []string
	insecureRegistries []string
}

// NewRegistrySources creates a new builder of 'registry_sources' objects.
func NewRegistrySources() *RegistrySourcesBuilder {
	return &RegistrySourcesBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RegistrySourcesBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// AllowedRegistries sets the value of the 'allowed_registries' attribute to the given values.
func (b *RegistrySourcesBuilder) AllowedRegistries(values ...string) *RegistrySourcesBuilder {
	b.allowedRegistries = make([]string, len(values))
	copy(b.allowedRegistries, values)
	b.bitmap_ |= 1
	return b
}

// BlockedRegistries sets the value of the 'blocked_registries' attribute to the given values.
func (b *RegistrySourcesBuilder) BlockedRegistries(values ...string) *RegistrySourcesBuilder {
	b.blockedRegistries = make([]string, len(values))
	copy(b.blockedRegistries, values)
	b.bitmap_ |= 2
	return b
}

// InsecureRegistries sets the value of the 'insecure_registries' attribute to the given values.
func (b *RegistrySourcesBuilder) InsecureRegistries(values ...string) *RegistrySourcesBuilder {
	b.insecureRegistries = make([]string, len(values))
	copy(b.insecureRegistries, values)
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RegistrySourcesBuilder) Copy(object *RegistrySources) *RegistrySourcesBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.allowedRegistries != nil {
		b.allowedRegistries = make([]string, len(object.allowedRegistries))
		copy(b.allowedRegistries, object.allowedRegistries)
	} else {
		b.allowedRegistries = nil
	}
	if object.blockedRegistries != nil {
		b.blockedRegistries = make([]string, len(object.blockedRegistries))
		copy(b.blockedRegistries, object.blockedRegistries)
	} else {
		b.blockedRegistries = nil
	}
	if object.insecureRegistries != nil {
		b.insecureRegistries = make([]string, len(object.insecureRegistries))
		copy(b.insecureRegistries, object.insecureRegistries)
	} else {
		b.insecureRegistries = nil
	}
	return b
}

// Build creates a 'registry_sources' object using the configuration stored in the builder.
func (b *RegistrySourcesBuilder) Build() (object *RegistrySources, err error) {
	object = new(RegistrySources)
	object.bitmap_ = b.bitmap_
	if b.allowedRegistries != nil {
		object.allowedRegistries = make([]string, len(b.allowedRegistries))
		copy(object.allowedRegistries, b.allowedRegistries)
	}
	if b.blockedRegistries != nil {
		object.blockedRegistries = make([]string, len(b.blockedRegistries))
		copy(object.blockedRegistries, b.blockedRegistries)
	}
	if b.insecureRegistries != nil {
		object.insecureRegistries = make([]string, len(b.insecureRegistries))
		copy(object.insecureRegistries, b.insecureRegistries)
	}
	return
}
