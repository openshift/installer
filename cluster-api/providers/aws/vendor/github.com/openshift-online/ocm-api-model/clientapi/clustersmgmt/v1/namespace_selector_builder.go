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

// Representation of a NamespaceSelector
type NamespaceSelectorBuilder struct {
	fieldSet_ []bool
	key       string
	values    []string
}

// NewNamespaceSelector creates a new builder of 'namespace_selector' objects.
func NewNamespaceSelector() *NamespaceSelectorBuilder {
	return &NamespaceSelectorBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *NamespaceSelectorBuilder) Empty() bool {
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

// Key sets the value of the 'key' attribute to the given value.
func (b *NamespaceSelectorBuilder) Key(value string) *NamespaceSelectorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.key = value
	b.fieldSet_[0] = true
	return b
}

// Values sets the value of the 'values' attribute to the given values.
func (b *NamespaceSelectorBuilder) Values(values ...string) *NamespaceSelectorBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.values = make([]string, len(values))
	copy(b.values, values)
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *NamespaceSelectorBuilder) Copy(object *NamespaceSelector) *NamespaceSelectorBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.key = object.key
	if object.values != nil {
		b.values = make([]string, len(object.values))
		copy(b.values, object.values)
	} else {
		b.values = nil
	}
	return b
}

// Build creates a 'namespace_selector' object using the configuration stored in the builder.
func (b *NamespaceSelectorBuilder) Build() (object *NamespaceSelector, err error) {
	object = new(NamespaceSelector)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.key = b.key
	if b.values != nil {
		object.values = make([]string, len(b.values))
		copy(object.values, b.values)
	}
	return
}
