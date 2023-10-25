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

// EventBuilder contains the data and logic needed to build 'event' objects.
//
// Representation of a trackable event.
type EventBuilder struct {
	bitmap_ uint32
	body    map[string]string
	key     string
}

// NewEvent creates a new builder of 'event' objects.
func NewEvent() *EventBuilder {
	return &EventBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *EventBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Body sets the value of the 'body' attribute to the given value.
func (b *EventBuilder) Body(value map[string]string) *EventBuilder {
	b.body = value
	if value != nil {
		b.bitmap_ |= 1
	} else {
		b.bitmap_ &^= 1
	}
	return b
}

// Key sets the value of the 'key' attribute to the given value.
func (b *EventBuilder) Key(value string) *EventBuilder {
	b.key = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *EventBuilder) Copy(object *Event) *EventBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
	if b.body != nil {
		object.body = make(map[string]string)
		for k, v := range b.body {
			object.body[k] = v
		}
	}
	object.key = b.key
	return
}
