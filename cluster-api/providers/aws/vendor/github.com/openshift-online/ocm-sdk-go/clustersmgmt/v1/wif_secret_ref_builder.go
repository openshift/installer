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

// WifSecretRefBuilder contains the data and logic needed to build 'wif_secret_ref' objects.
type WifSecretRefBuilder struct {
	bitmap_   uint32
	name      string
	namespace string
}

// NewWifSecretRef creates a new builder of 'wif_secret_ref' objects.
func NewWifSecretRef() *WifSecretRefBuilder {
	return &WifSecretRefBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *WifSecretRefBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Name sets the value of the 'name' attribute to the given value.
func (b *WifSecretRefBuilder) Name(value string) *WifSecretRefBuilder {
	b.name = value
	b.bitmap_ |= 1
	return b
}

// Namespace sets the value of the 'namespace' attribute to the given value.
func (b *WifSecretRefBuilder) Namespace(value string) *WifSecretRefBuilder {
	b.namespace = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *WifSecretRefBuilder) Copy(object *WifSecretRef) *WifSecretRefBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.name = object.name
	b.namespace = object.namespace
	return b
}

// Build creates a 'wif_secret_ref' object using the configuration stored in the builder.
func (b *WifSecretRefBuilder) Build() (object *WifSecretRef, err error) {
	object = new(WifSecretRef)
	object.bitmap_ = b.bitmap_
	object.name = b.name
	object.namespace = b.namespace
	return
}
