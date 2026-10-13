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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1

// Representation of AddonParameters
type AddonParametersBuilder struct {
	fieldSet_ []bool
	items     []*AddonParameterBuilder
}

// NewAddonParameters creates a new builder of 'addon_parameters' objects.
func NewAddonParameters() *AddonParametersBuilder {
	return &AddonParametersBuilder{
		fieldSet_: make([]bool, 1),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonParametersBuilder) Empty() bool {
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

// Items sets the value of the 'items' attribute to the given values.
func (b *AddonParametersBuilder) Items(values ...*AddonParameterBuilder) *AddonParametersBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 1)
	}
	b.items = make([]*AddonParameterBuilder, len(values))
	copy(b.items, values)
	b.fieldSet_[0] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonParametersBuilder) Copy(object *AddonParameters) *AddonParametersBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if object.items != nil {
		b.items = make([]*AddonParameterBuilder, len(object.items))
		for i, v := range object.items {
			b.items[i] = NewAddonParameter().Copy(v)
		}
	} else {
		b.items = nil
	}
	return b
}

// Build creates a 'addon_parameters' object using the configuration stored in the builder.
func (b *AddonParametersBuilder) Build() (object *AddonParameters, err error) {
	object = new(AddonParameters)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.items != nil {
		object.items = make([]*AddonParameter, len(b.items))
		for i, v := range b.items {
			object.items[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
