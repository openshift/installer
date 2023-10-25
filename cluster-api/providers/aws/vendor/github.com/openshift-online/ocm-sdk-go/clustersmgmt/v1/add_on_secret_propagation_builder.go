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

// AddOnSecretPropagationBuilder contains the data and logic needed to build 'add_on_secret_propagation' objects.
//
// Representation of an addon secret propagation
type AddOnSecretPropagationBuilder struct {
	bitmap_           uint32
	id                string
	destinationSecret string
	sourceSecret      string
	enabled           bool
}

// NewAddOnSecretPropagation creates a new builder of 'add_on_secret_propagation' objects.
func NewAddOnSecretPropagation() *AddOnSecretPropagationBuilder {
	return &AddOnSecretPropagationBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnSecretPropagationBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ID sets the value of the 'ID' attribute to the given value.
func (b *AddOnSecretPropagationBuilder) ID(value string) *AddOnSecretPropagationBuilder {
	b.id = value
	b.bitmap_ |= 1
	return b
}

// DestinationSecret sets the value of the 'destination_secret' attribute to the given value.
func (b *AddOnSecretPropagationBuilder) DestinationSecret(value string) *AddOnSecretPropagationBuilder {
	b.destinationSecret = value
	b.bitmap_ |= 2
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddOnSecretPropagationBuilder) Enabled(value bool) *AddOnSecretPropagationBuilder {
	b.enabled = value
	b.bitmap_ |= 4
	return b
}

// SourceSecret sets the value of the 'source_secret' attribute to the given value.
func (b *AddOnSecretPropagationBuilder) SourceSecret(value string) *AddOnSecretPropagationBuilder {
	b.sourceSecret = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnSecretPropagationBuilder) Copy(object *AddOnSecretPropagation) *AddOnSecretPropagationBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.destinationSecret = object.destinationSecret
	b.enabled = object.enabled
	b.sourceSecret = object.sourceSecret
	return b
}

// Build creates a 'add_on_secret_propagation' object using the configuration stored in the builder.
func (b *AddOnSecretPropagationBuilder) Build() (object *AddOnSecretPropagation, err error) {
	object = new(AddOnSecretPropagation)
	object.bitmap_ = b.bitmap_
	object.id = b.id
	object.destinationSecret = b.destinationSecret
	object.enabled = b.enabled
	object.sourceSecret = b.sourceSecret
	return
}
