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

// AdminCredentialsBuilder contains the data and logic needed to build 'admin_credentials' objects.
//
// Temporary administrator credentials generated during the installation of the
// cluster.
type AdminCredentialsBuilder struct {
	bitmap_  uint32
	password string
	user     string
}

// NewAdminCredentials creates a new builder of 'admin_credentials' objects.
func NewAdminCredentials() *AdminCredentialsBuilder {
	return &AdminCredentialsBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AdminCredentialsBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Password sets the value of the 'password' attribute to the given value.
func (b *AdminCredentialsBuilder) Password(value string) *AdminCredentialsBuilder {
	b.password = value
	b.bitmap_ |= 1
	return b
}

// User sets the value of the 'user' attribute to the given value.
func (b *AdminCredentialsBuilder) User(value string) *AdminCredentialsBuilder {
	b.user = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AdminCredentialsBuilder) Copy(object *AdminCredentials) *AdminCredentialsBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.password = object.password
	b.user = object.user
	return b
}

// Build creates a 'admin_credentials' object using the configuration stored in the builder.
func (b *AdminCredentialsBuilder) Build() (object *AdminCredentials, err error) {
	object = new(AdminCredentials)
	object.bitmap_ = b.bitmap_
	object.password = b.password
	object.user = b.user
	return
}
