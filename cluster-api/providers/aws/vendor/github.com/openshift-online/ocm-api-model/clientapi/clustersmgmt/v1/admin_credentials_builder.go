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

// Temporary administrator credentials generated during the installation of the
// cluster.
type AdminCredentialsBuilder struct {
	fieldSet_ []bool
	password  string
	user      string
}

// NewAdminCredentials creates a new builder of 'admin_credentials' objects.
func NewAdminCredentials() *AdminCredentialsBuilder {
	return &AdminCredentialsBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AdminCredentialsBuilder) Empty() bool {
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

// Password sets the value of the 'password' attribute to the given value.
func (b *AdminCredentialsBuilder) Password(value string) *AdminCredentialsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.password = value
	b.fieldSet_[0] = true
	return b
}

// User sets the value of the 'user' attribute to the given value.
func (b *AdminCredentialsBuilder) User(value string) *AdminCredentialsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.user = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AdminCredentialsBuilder) Copy(object *AdminCredentials) *AdminCredentialsBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.password = object.password
	b.user = object.user
	return b
}

// Build creates a 'admin_credentials' object using the configuration stored in the builder.
func (b *AdminCredentialsBuilder) Build() (object *AdminCredentials, err error) {
	object = new(AdminCredentials)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.password = b.password
	object.user = b.user
	return
}
