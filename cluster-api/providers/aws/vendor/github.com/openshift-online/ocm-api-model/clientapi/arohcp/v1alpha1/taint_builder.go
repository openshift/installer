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

// Representation of a Taint set on a MachinePool in a cluster.
type TaintBuilder struct {
	fieldSet_ []bool
	effect    string
	key       string
	value     string
}

// NewTaint creates a new builder of 'taint' objects.
func NewTaint() *TaintBuilder {
	return &TaintBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *TaintBuilder) Empty() bool {
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

// Effect sets the value of the 'effect' attribute to the given value.
func (b *TaintBuilder) Effect(value string) *TaintBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.effect = value
	b.fieldSet_[0] = true
	return b
}

// Key sets the value of the 'key' attribute to the given value.
func (b *TaintBuilder) Key(value string) *TaintBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.key = value
	b.fieldSet_[1] = true
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *TaintBuilder) Value(value string) *TaintBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.value = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *TaintBuilder) Copy(object *Taint) *TaintBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.effect = object.effect
	b.key = object.key
	b.value = object.value
	return b
}

// Build creates a 'taint' object using the configuration stored in the builder.
func (b *TaintBuilder) Build() (object *Taint, err error) {
	object = new(Taint)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.effect = b.effect
	object.key = b.key
	object.value = b.value
	return
}
