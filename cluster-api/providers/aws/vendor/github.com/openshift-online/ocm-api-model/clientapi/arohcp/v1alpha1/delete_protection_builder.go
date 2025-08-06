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

// DeleteProtection configuration.
type DeleteProtectionBuilder struct {
	fieldSet_ []bool
	enabled   bool
}

// NewDeleteProtection creates a new builder of 'delete_protection' objects.
func NewDeleteProtection() *DeleteProtectionBuilder {
	return &DeleteProtectionBuilder{
		fieldSet_: make([]bool, 1),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *DeleteProtectionBuilder) Empty() bool {
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

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *DeleteProtectionBuilder) Enabled(value bool) *DeleteProtectionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 1)
	}
	b.enabled = value
	b.fieldSet_[0] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *DeleteProtectionBuilder) Copy(object *DeleteProtection) *DeleteProtectionBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.enabled = object.enabled
	return b
}

// Build creates a 'delete_protection' object using the configuration stored in the builder.
func (b *DeleteProtectionBuilder) Build() (object *DeleteProtection, err error) {
	object = new(DeleteProtection)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.enabled = b.enabled
	return
}
