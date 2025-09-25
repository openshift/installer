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

// RegistryLocationBuilder contains the data and logic needed to build 'registry_location' objects.
//
// RegistryLocation contains a location of the registry specified by the registry domain
// name. The domain name might include wildcards, like '*' or '??'.
type RegistryLocationBuilder struct {
	bitmap_    uint32
	domainName string
	insecure   bool
}

// NewRegistryLocation creates a new builder of 'registry_location' objects.
func NewRegistryLocation() *RegistryLocationBuilder {
	return &RegistryLocationBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *RegistryLocationBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// DomainName sets the value of the 'domain_name' attribute to the given value.
func (b *RegistryLocationBuilder) DomainName(value string) *RegistryLocationBuilder {
	b.domainName = value
	b.bitmap_ |= 1
	return b
}

// Insecure sets the value of the 'insecure' attribute to the given value.
func (b *RegistryLocationBuilder) Insecure(value bool) *RegistryLocationBuilder {
	b.insecure = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *RegistryLocationBuilder) Copy(object *RegistryLocation) *RegistryLocationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.domainName = object.domainName
	b.insecure = object.insecure
	return b
}

// Build creates a 'registry_location' object using the configuration stored in the builder.
func (b *RegistryLocationBuilder) Build() (object *RegistryLocation, err error) {
	object = new(RegistryLocation)
	object.bitmap_ = b.bitmap_
	object.domainName = b.domainName
	object.insecure = b.insecure
	return
}
