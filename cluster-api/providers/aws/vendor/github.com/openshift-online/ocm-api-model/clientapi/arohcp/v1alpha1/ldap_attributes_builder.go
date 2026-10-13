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

// LDAP attributes used to configure the LDAP identity provider.
type LDAPAttributesBuilder struct {
	fieldSet_         []bool
	id                []string
	email             []string
	name              []string
	preferredUsername []string
}

// NewLDAPAttributes creates a new builder of 'LDAP_attributes' objects.
func NewLDAPAttributes() *LDAPAttributesBuilder {
	return &LDAPAttributesBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LDAPAttributesBuilder) Empty() bool {
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

// ID sets the value of the 'ID' attribute to the given values.
func (b *LDAPAttributesBuilder) ID(values ...string) *LDAPAttributesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.id = make([]string, len(values))
	copy(b.id, values)
	b.fieldSet_[0] = true
	return b
}

// Email sets the value of the 'email' attribute to the given values.
func (b *LDAPAttributesBuilder) Email(values ...string) *LDAPAttributesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.email = make([]string, len(values))
	copy(b.email, values)
	b.fieldSet_[1] = true
	return b
}

// Name sets the value of the 'name' attribute to the given values.
func (b *LDAPAttributesBuilder) Name(values ...string) *LDAPAttributesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.name = make([]string, len(values))
	copy(b.name, values)
	b.fieldSet_[2] = true
	return b
}

// PreferredUsername sets the value of the 'preferred_username' attribute to the given values.
func (b *LDAPAttributesBuilder) PreferredUsername(values ...string) *LDAPAttributesBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.preferredUsername = make([]string, len(values))
	copy(b.preferredUsername, values)
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LDAPAttributesBuilder) Copy(object *LDAPAttributes) *LDAPAttributesBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.id != nil {
		b.id = make([]string, len(object.id))
		copy(b.id, object.id)
	} else {
		b.id = nil
	}
	if object.email != nil {
		b.email = make([]string, len(object.email))
		copy(b.email, object.email)
	} else {
		b.email = nil
	}
	if object.name != nil {
		b.name = make([]string, len(object.name))
		copy(b.name, object.name)
	} else {
		b.name = nil
	}
	if object.preferredUsername != nil {
		b.preferredUsername = make([]string, len(object.preferredUsername))
		copy(b.preferredUsername, object.preferredUsername)
	} else {
		b.preferredUsername = nil
	}
	return b
}

// Build creates a 'LDAP_attributes' object using the configuration stored in the builder.
func (b *LDAPAttributesBuilder) Build() (object *LDAPAttributes, err error) {
	object = new(LDAPAttributes)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.id != nil {
		object.id = make([]string, len(b.id))
		copy(object.id, b.id)
	}
	if b.email != nil {
		object.email = make([]string, len(b.email))
		copy(object.email, b.email)
	}
	if b.name != nil {
		object.name = make([]string, len(b.name))
		copy(object.name, b.name)
	}
	if b.preferredUsername != nil {
		object.preferredUsername = make([]string, len(b.preferredUsername))
		copy(object.preferredUsername, b.preferredUsername)
	}
	return
}
