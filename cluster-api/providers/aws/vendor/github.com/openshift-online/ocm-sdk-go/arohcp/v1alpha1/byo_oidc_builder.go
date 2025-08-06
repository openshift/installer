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

// ByoOidcBuilder contains the data and logic needed to build 'byo_oidc' objects.
//
// ByoOidc configuration.
type ByoOidcBuilder struct {
	bitmap_ uint32
	enabled bool
}

// NewByoOidc creates a new builder of 'byo_oidc' objects.
func NewByoOidc() *ByoOidcBuilder {
	return &ByoOidcBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ByoOidcBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *ByoOidcBuilder) Enabled(value bool) *ByoOidcBuilder {
	b.enabled = value
	b.bitmap_ |= 1
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ByoOidcBuilder) Copy(object *ByoOidc) *ByoOidcBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.enabled = object.enabled
	return b
}

// Build creates a 'byo_oidc' object using the configuration stored in the builder.
func (b *ByoOidcBuilder) Build() (object *ByoOidc, err error) {
	object = new(ByoOidc)
	object.bitmap_ = b.bitmap_
	object.enabled = b.enabled
	return
}
