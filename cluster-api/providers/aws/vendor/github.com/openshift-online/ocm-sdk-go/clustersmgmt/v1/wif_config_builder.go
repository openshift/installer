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

// WifConfigBuilder contains the data and logic needed to build 'wif_config' objects.
//
// Definition of an wif_config resource.
type WifConfigBuilder struct {
	bitmap_      uint32
	id           string
	href         string
	displayName  string
	gcp          *WifGcpBuilder
	organization *OrganizationLinkBuilder
	wifTemplates []string
}

// NewWifConfig creates a new builder of 'wif_config' objects.
func NewWifConfig() *WifConfigBuilder {
	return &WifConfigBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *WifConfigBuilder) Link(value bool) *WifConfigBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *WifConfigBuilder) ID(value string) *WifConfigBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *WifConfigBuilder) HREF(value string) *WifConfigBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifConfigBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *WifConfigBuilder) DisplayName(value string) *WifConfigBuilder {
	b.displayName = value
	b.bitmap_ |= 8
	return b
}

// Gcp sets the value of the 'gcp' attribute to the given value.
func (b *WifConfigBuilder) Gcp(value *WifGcpBuilder) *WifConfigBuilder {
	b.gcp = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// Organization sets the value of the 'organization' attribute to the given value.
//
// Definition of an organization link.
func (b *WifConfigBuilder) Organization(value *OrganizationLinkBuilder) *WifConfigBuilder {
	b.organization = value
	if value != nil {
		b.bitmap_ |= 32
	} else {
		b.bitmap_ &^= 32
	}
	return b
}

// WifTemplates sets the value of the 'wif_templates' attribute to the given values.
func (b *WifConfigBuilder) WifTemplates(values ...string) *WifConfigBuilder {
	b.wifTemplates = make([]string, len(values))
	copy(b.wifTemplates, values)
	b.bitmap_ |= 64
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifConfigBuilder) Copy(object *WifConfig) *WifConfigBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.displayName = object.displayName
	if object.gcp != nil {
		b.gcp = NewWifGcp().Copy(object.gcp)
	} else {
		b.gcp = nil
	}
	if object.organization != nil {
		b.organization = NewOrganizationLink().Copy(object.organization)
	} else {
		b.organization = nil
	}
	if object.wifTemplates != nil {
		b.wifTemplates = make([]string, len(object.wifTemplates))
		copy(b.wifTemplates, object.wifTemplates)
	} else {
		b.wifTemplates = nil
	}
	return b
}

// Build creates a 'wif_config' object using the configuration stored in the builder.
func (b *WifConfigBuilder) Build() (object *WifConfig, err error) {
	object = new(WifConfig)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.displayName = b.displayName
	if b.gcp != nil {
		object.gcp, err = b.gcp.Build()
		if err != nil {
			return
		}
	}
	if b.organization != nil {
		object.organization, err = b.organization.Build()
		if err != nil {
			return
		}
	}
	if b.wifTemplates != nil {
		object.wifTemplates = make([]string, len(b.wifTemplates))
		copy(object.wifTemplates, b.wifTemplates)
	}
	return
}
