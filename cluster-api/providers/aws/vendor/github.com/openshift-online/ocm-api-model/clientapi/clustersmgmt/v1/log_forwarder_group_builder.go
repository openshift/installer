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

// Represents a log forwarder group.
type LogForwarderGroupBuilder struct {
	fieldSet_ []bool
	id        string
	version   string
}

// NewLogForwarderGroup creates a new builder of 'log_forwarder_group' objects.
func NewLogForwarderGroup() *LogForwarderGroupBuilder {
	return &LogForwarderGroupBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LogForwarderGroupBuilder) Empty() bool {
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
func (b *LogForwarderGroupBuilder) ID(value string) *LogForwarderGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.id = value
	b.fieldSet_[0] = true
	return b
}

// Version sets the value of the 'version' attribute to the given value.
func (b *LogForwarderGroupBuilder) Version(value string) *LogForwarderGroupBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.version = value
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LogForwarderGroupBuilder) Copy(object *LogForwarderGroup) *LogForwarderGroupBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.version = object.version
	return b
}

// Build creates a 'log_forwarder_group' object using the configuration stored in the builder.
func (b *LogForwarderGroupBuilder) Build() (object *LogForwarderGroup, err error) {
	object = new(LogForwarderGroup)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.id = b.id
	object.version = b.version
	return
}
