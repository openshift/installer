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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/osdfleetmgmt/v1

type LabelRequestPayloadBuilder struct {
	fieldSet_ []bool
	key       string
	value     string
}

// NewLabelRequestPayload creates a new builder of 'label_request_payload' objects.
func NewLabelRequestPayload() *LabelRequestPayloadBuilder {
	return &LabelRequestPayloadBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LabelRequestPayloadBuilder) Empty() bool {
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
func (b *LabelRequestPayloadBuilder) Key(value string) *LabelRequestPayloadBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.key = value
	b.fieldSet_[0] = true
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *LabelRequestPayloadBuilder) Value(value string) *LabelRequestPayloadBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.value = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LabelRequestPayloadBuilder) Copy(object *LabelRequestPayload) *LabelRequestPayloadBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.key = object.key
	b.value = object.value
	return b
}

// Build creates a 'label_request_payload' object using the configuration stored in the builder.
func (b *LabelRequestPayloadBuilder) Build() (object *LabelRequestPayload, err error) {
	object = new(LabelRequestPayload)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.key = b.key
	object.value = b.value
	return
}
