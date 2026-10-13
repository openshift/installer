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

// RegistrySources contains configuration that determines how the container runtime should treat individual
// registries when accessing images for builds and pods. For instance, whether or not to allow insecure access.
// It does not contain configuration for the internal cluster registry.
type RegistrySourcesBuilder struct {
	fieldSet_          []bool
	allowedRegistries  []string
	blockedRegistries  []string
	insecureRegistries []string
}

// NewRegistrySources creates a new builder of 'registry_sources' objects.
func NewRegistrySources() *RegistrySourcesBuilder {
	return &RegistrySourcesBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RegistrySourcesBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	for _, set := range b.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// AllowedRegistries sets the value of the 'allowed_registries' attribute to the given values.
func (b *RegistrySourcesBuilder) AllowedRegistries(values ...string) *RegistrySourcesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.allowedRegistries = make([]string, len(values))
	copy(b.allowedRegistries, values)
	b.fieldSet_[0] = true
	return b
}

// BlockedRegistries sets the value of the 'blocked_registries' attribute to the given values.
func (b *RegistrySourcesBuilder) BlockedRegistries(values ...string) *RegistrySourcesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.blockedRegistries = make([]string, len(values))
	copy(b.blockedRegistries, values)
	b.fieldSet_[1] = true
	return b
}

// InsecureRegistries sets the value of the 'insecure_registries' attribute to the given values.
func (b *RegistrySourcesBuilder) InsecureRegistries(values ...string) *RegistrySourcesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.insecureRegistries = make([]string, len(values))
	copy(b.insecureRegistries, values)
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RegistrySourcesBuilder) Copy(object *RegistrySources) *RegistrySourcesBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
