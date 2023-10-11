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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

// AddonSubOperatorBuilder contains the data and logic needed to build 'addon_sub_operator' objects.
//
// Representation of an addon sub operator. A sub operator is an operator
// who's life cycle is controlled by the addon umbrella operator.
type AddonSubOperatorBuilder struct {
	bitmap_           uint32
	addon             *AddonBuilder
	operatorName      string
	operatorNamespace string
	enabled           bool
}

// NewAddonSubOperator creates a new builder of 'addon_sub_operator' objects.
func NewAddonSubOperator() *AddonSubOperatorBuilder {
	return &AddonSubOperatorBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonSubOperatorBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Addon sets the value of the 'addon' attribute to the given value.
//
// Representation of an addon that can be installed in a cluster.
func (b *AddonSubOperatorBuilder) Addon(value *AddonBuilder) *AddonSubOperatorBuilder {
	b.addon = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddonSubOperatorBuilder) Enabled(value bool) *AddonSubOperatorBuilder {
	b.enabled = value
	b.bitmap_ |= 2
	return b
}

// OperatorName sets the value of the 'operator_name' attribute to the given value.
func (b *AddonSubOperatorBuilder) OperatorName(value string) *AddonSubOperatorBuilder {
	b.operatorName = value
	b.bitmap_ |= 4
	return b
}

// OperatorNamespace sets the value of the 'operator_namespace' attribute to the given value.
func (b *AddonSubOperatorBuilder) OperatorNamespace(value string) *AddonSubOperatorBuilder {
	b.operatorNamespace = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonSubOperatorBuilder) Copy(object *AddonSubOperator) *AddonSubOperatorBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	if object.addon != nil {
		b.addon = NewAddon().Copy(object.addon)
	} else {
		b.addon = nil
	}
	b.enabled = object.enabled
	b.operatorName = object.operatorName
	b.operatorNamespace = object.operatorNamespace
	return b
}

// Build creates a 'addon_sub_operator' object using the configuration stored in the builder.
func (b *AddonSubOperatorBuilder) Build() (object *AddonSubOperator, err error) {
	object = new(AddonSubOperator)
	object.bitmap_ = b.bitmap_
	if b.addon != nil {
		object.addon, err = b.addon.Build()
		if err != nil {
			return
		}
	}
	object.enabled = b.enabled
	object.operatorName = b.operatorName
	object.operatorNamespace = b.operatorNamespace
	return
}
