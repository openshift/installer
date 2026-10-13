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

// Representation of an add-on sub operator. A sub operator is an operator
// who's life cycle is controlled by the add-on umbrella operator.
type AddOnSubOperatorBuilder struct {
	fieldSet_         []bool
	operatorName      string
	operatorNamespace string
	enabled           bool
}

// NewAddOnSubOperator creates a new builder of 'add_on_sub_operator' objects.
func NewAddOnSubOperator() *AddOnSubOperatorBuilder {
	return &AddOnSubOperatorBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AddOnSubOperatorBuilder) Empty() bool {
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
func (b *AddOnSubOperatorBuilder) Enabled(value bool) *AddOnSubOperatorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.enabled = value
	b.fieldSet_[0] = true
	return b
}

// OperatorName sets the value of the 'operator_name' attribute to the given value.
func (b *AddOnSubOperatorBuilder) OperatorName(value string) *AddOnSubOperatorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.operatorName = value
	b.fieldSet_[1] = true
	return b
}

// OperatorNamespace sets the value of the 'operator_namespace' attribute to the given value.
func (b *AddOnSubOperatorBuilder) OperatorNamespace(value string) *AddOnSubOperatorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.operatorNamespace = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AddOnSubOperatorBuilder) Copy(object *AddOnSubOperator) *AddOnSubOperatorBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.enabled = object.enabled
	b.operatorName = object.operatorName
	b.operatorNamespace = object.operatorNamespace
	return b
}

// Build creates a 'add_on_sub_operator' object using the configuration stored in the builder.
func (b *AddOnSubOperatorBuilder) Build() (object *AddOnSubOperator, err error) {
	object = new(AddOnSubOperator)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.enabled = b.enabled
	object.operatorName = b.operatorName
	object.operatorNamespace = b.operatorNamespace
	return
}
