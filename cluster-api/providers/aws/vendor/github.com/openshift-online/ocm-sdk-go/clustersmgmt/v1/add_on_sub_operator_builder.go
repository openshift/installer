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

// AddOnSubOperatorBuilder contains the data and logic needed to build 'add_on_sub_operator' objects.
//
// Representation of an add-on sub operator. A sub operator is an operator
// who's life cycle is controlled by the add-on umbrella operator.
type AddOnSubOperatorBuilder struct {
	bitmap_           uint32
	operatorName      string
	operatorNamespace string
	enabled           bool
}

// NewAddOnSubOperator creates a new builder of 'add_on_sub_operator' objects.
func NewAddOnSubOperator() *AddOnSubOperatorBuilder {
	return &AddOnSubOperatorBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnSubOperatorBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *AddOnSubOperatorBuilder) Enabled(value bool) *AddOnSubOperatorBuilder {
	b.enabled = value
	b.bitmap_ |= 1
	return b
}

// OperatorName sets the value of the 'operator_name' attribute to the given value.
func (b *AddOnSubOperatorBuilder) OperatorName(value string) *AddOnSubOperatorBuilder {
	b.operatorName = value
	b.bitmap_ |= 2
	return b
}

// OperatorNamespace sets the value of the 'operator_namespace' attribute to the given value.
func (b *AddOnSubOperatorBuilder) OperatorNamespace(value string) *AddOnSubOperatorBuilder {
	b.operatorNamespace = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnSubOperatorBuilder) Copy(object *AddOnSubOperator) *AddOnSubOperatorBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.enabled = object.enabled
	b.operatorName = object.operatorName
	b.operatorNamespace = object.operatorNamespace
	return b
}

// Build creates a 'add_on_sub_operator' object using the configuration stored in the builder.
func (b *AddOnSubOperatorBuilder) Build() (object *AddOnSubOperator, err error) {
	object = new(AddOnSubOperator)
	object.bitmap_ = b.bitmap_
	object.enabled = b.enabled
	object.operatorName = b.operatorName
	object.operatorNamespace = b.operatorNamespace
	return
}
