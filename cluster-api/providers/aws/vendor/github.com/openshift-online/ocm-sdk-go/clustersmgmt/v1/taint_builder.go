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

// TaintBuilder contains the data and logic needed to build 'taint' objects.
//
// Representation of a Taint set on a MachinePool in a cluster.
type TaintBuilder struct {
	bitmap_ uint32
	effect  string
	key     string
	value   string
}

// NewTaint creates a new builder of 'taint' objects.
func NewTaint() *TaintBuilder {
	return &TaintBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *TaintBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Effect sets the value of the 'effect' attribute to the given value.
func (b *TaintBuilder) Effect(value string) *TaintBuilder {
	b.effect = value
	b.bitmap_ |= 1
	return b
}

// Key sets the value of the 'key' attribute to the given value.
func (b *TaintBuilder) Key(value string) *TaintBuilder {
	b.key = value
	b.bitmap_ |= 2
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *TaintBuilder) Value(value string) *TaintBuilder {
	b.value = value
	b.bitmap_ |= 4
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *TaintBuilder) Copy(object *Taint) *TaintBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.effect = object.effect
	b.key = object.key
	b.value = object.value
	return b
}

// Build creates a 'taint' object using the configuration stored in the builder.
func (b *TaintBuilder) Build() (object *Taint, err error) {
	object = new(Taint)
	object.bitmap_ = b.bitmap_
	object.effect = b.effect
	object.key = b.key
	object.value = b.value
	return
}
