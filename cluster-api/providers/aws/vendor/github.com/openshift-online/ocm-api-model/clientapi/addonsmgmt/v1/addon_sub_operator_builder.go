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

// Representation of an addon sub operator. A sub operator is an operator
// who's life cycle is controlled by the addon umbrella operator.
type AddonSubOperatorBuilder struct {
	fieldSet_         []bool
	addon             *AddonBuilder
	operatorName      string
	operatorNamespace string
	enabled           bool
}

// NewAddonSubOperator creates a new builder of 'addon_sub_operator' objects.
func NewAddonSubOperator() *AddonSubOperatorBuilder {
	return &AddonSubOperatorBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddonSubOperatorBuilder) Empty() bool {
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

// Addon sets the value of the 'addon' attribute to the given value.
//
// Representation of an addon that can be installed in a cluster.
func (b *AddonSubOperatorBuilder) Addon(value *AddonBuilder) *AddonSubOperatorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.addon = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddonSubOperatorBuilder) Enabled(value bool) *AddonSubOperatorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.enabled = value
	b.fieldSet_[1] = true
	return b
}

// OperatorName sets the value of the 'operator_name' attribute to the given value.
func (b *AddonSubOperatorBuilder) OperatorName(value string) *AddonSubOperatorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.operatorName = value
	b.fieldSet_[2] = true
	return b
}

// OperatorNamespace sets the value of the 'operator_namespace' attribute to the given value.
func (b *AddonSubOperatorBuilder) OperatorNamespace(value string) *AddonSubOperatorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.operatorNamespace = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddonSubOperatorBuilder) Copy(object *AddonSubOperator) *AddonSubOperatorBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
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
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
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
