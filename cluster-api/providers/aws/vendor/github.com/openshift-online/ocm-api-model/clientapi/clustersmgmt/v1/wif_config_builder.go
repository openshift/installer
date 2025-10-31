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

// Definition of an wif_config resource.
type WifConfigBuilder struct {
	fieldSet_    []bool
	id           string
	href         string
	displayName  string
	gcp          *WifGcpBuilder
	organization *OrganizationLinkBuilder
	wifTemplates []string
}

// NewWifConfig creates a new builder of 'wif_config' objects.
func NewWifConfig() *WifConfigBuilder {
	return &WifConfigBuilder{
		fieldSet_: make([]bool, 7),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *WifConfigBuilder) Link(value bool) *WifConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *WifConfigBuilder) ID(value string) *WifConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *WifConfigBuilder) HREF(value string) *WifConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifConfigBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// DisplayName sets the value of the 'display_name' attribute to the given value.
func (b *WifConfigBuilder) DisplayName(value string) *WifConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.displayName = value
	b.fieldSet_[3] = true
	return b
}

// Gcp sets the value of the 'gcp' attribute to the given value.
func (b *WifConfigBuilder) Gcp(value *WifGcpBuilder) *WifConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.gcp = value
	if value != nil {
		b.fieldSet_[4] = true
	} else {
		b.fieldSet_[4] = false
	}
	return b
}

// Organization sets the value of the 'organization' attribute to the given value.
//
// Definition of an organization link.
func (b *WifConfigBuilder) Organization(value *OrganizationLinkBuilder) *WifConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.organization = value
	if value != nil {
		b.fieldSet_[5] = true
	} else {
		b.fieldSet_[5] = false
	}
	return b
}

// WifTemplates sets the value of the 'wif_templates' attribute to the given values.
func (b *WifConfigBuilder) WifTemplates(values ...string) *WifConfigBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 7)
	}
	b.wifTemplates = make([]string, len(values))
	copy(b.wifTemplates, values)
	b.fieldSet_[6] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifConfigBuilder) Copy(object *WifConfig) *WifConfigBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
