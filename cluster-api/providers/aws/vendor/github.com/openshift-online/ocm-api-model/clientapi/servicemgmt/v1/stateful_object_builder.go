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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/servicemgmt/v1

type StatefulObjectBuilder struct {
	fieldSet_ []bool
	id        string
	href      string
	kind      string
	state     string
}

// NewStatefulObject creates a new builder of 'stateful_object' objects.
func NewStatefulObject() *StatefulObjectBuilder {
	return &StatefulObjectBuilder{
		fieldSet_: make([]bool, 4),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *StatefulObjectBuilder) Empty() bool {
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

// ID sets the value of the 'ID' attribute to the given value.
func (b *StatefulObjectBuilder) ID(value string) *StatefulObjectBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.id = value
	b.fieldSet_[0] = true
	return b
}

// Href sets the value of the 'href' attribute to the given value.
func (b *StatefulObjectBuilder) Href(value string) *StatefulObjectBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.href = value
	b.fieldSet_[1] = true
	return b
}

// Kind sets the value of the 'kind' attribute to the given value.
func (b *StatefulObjectBuilder) Kind(value string) *StatefulObjectBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.kind = value
	b.fieldSet_[2] = true
	return b
}

// State sets the value of the 'state' attribute to the given value.
func (b *StatefulObjectBuilder) State(value string) *StatefulObjectBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 4)
	}
	b.state = value
	b.fieldSet_[3] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *StatefulObjectBuilder) Copy(object *StatefulObject) *StatefulObjectBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.kind = object.kind
	b.state = object.state
	return b
}

// Build creates a 'stateful_object' object using the configuration stored in the builder.
func (b *StatefulObjectBuilder) Build() (object *StatefulObject, err error) {
	object = new(StatefulObject)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.id = b.id
	object.href = b.href
	object.kind = b.kind
	object.state = b.state
	return
}
