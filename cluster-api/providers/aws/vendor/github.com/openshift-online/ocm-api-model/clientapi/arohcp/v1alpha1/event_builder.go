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

// Representation of a trackable event.
type EventBuilder struct {
	fieldSet_ []bool
	body      map[string]string
	key       string
}

// NewEvent creates a new builder of 'event' objects.
func NewEvent() *EventBuilder {
	return &EventBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *EventBuilder) Empty() bool {
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

// Body sets the value of the 'body' attribute to the given value.
func (b *EventBuilder) Body(value map[string]string) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.body = value
	if value != nil {
		b.fieldSet_[0] = true
	} else {
		b.fieldSet_[0] = false
	}
	return b
}

// Key sets the value of the 'key' attribute to the given value.
func (b *EventBuilder) Key(value string) *EventBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.key = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *EventBuilder) Copy(object *Event) *EventBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	if len(object.body) > 0 {
		b.body = map[string]string{}
		for k, v := range object.body {
			b.body[k] = v
		}
	} else {
		b.body = nil
	}
	b.key = object.key
	return b
}

// Build creates a 'event' object using the configuration stored in the builder.
func (b *EventBuilder) Build() (object *Event, err error) {
	object = new(Event)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	if b.body != nil {
		object.body = make(map[string]string)
		for k, v := range b.body {
			object.body[k] = v
		}
	}
	object.key = b.key
	return
}
