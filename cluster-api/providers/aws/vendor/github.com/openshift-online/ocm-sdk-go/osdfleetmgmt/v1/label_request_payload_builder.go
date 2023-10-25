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

package v1 // github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1

// LabelRequestPayloadBuilder contains the data and logic needed to build 'label_request_payload' objects.
type LabelRequestPayloadBuilder struct {
	bitmap_ uint32
	key     string
	value   string
}

// NewLabelRequestPayload creates a new builder of 'label_request_payload' objects.
func NewLabelRequestPayload() *LabelRequestPayloadBuilder {
	return &LabelRequestPayloadBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LabelRequestPayloadBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// Key sets the value of the 'key' attribute to the given value.
func (b *LabelRequestPayloadBuilder) Key(value string) *LabelRequestPayloadBuilder {
	b.key = value
	b.bitmap_ |= 1
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *LabelRequestPayloadBuilder) Value(value string) *LabelRequestPayloadBuilder {
	b.value = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LabelRequestPayloadBuilder) Copy(object *LabelRequestPayload) *LabelRequestPayloadBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.key = object.key
	b.value = object.value
	return b
}

// Build creates a 'label_request_payload' object using the configuration stored in the builder.
func (b *LabelRequestPayloadBuilder) Build() (object *LabelRequestPayload, err error) {
	object = new(LabelRequestPayload)
	object.bitmap_ = b.bitmap_
	object.key = b.key
	object.value = b.value
	return
}
