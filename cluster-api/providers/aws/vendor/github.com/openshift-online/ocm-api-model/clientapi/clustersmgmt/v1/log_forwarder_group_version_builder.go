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

// Represents a version of a log forwarder group.
type LogForwarderGroupVersionBuilder struct {
	fieldSet_    []bool
	id           string
	applications []string
}

// NewLogForwarderGroupVersion creates a new builder of 'log_forwarder_group_version' objects.
func NewLogForwarderGroupVersion() *LogForwarderGroupVersionBuilder {
	return &LogForwarderGroupVersionBuilder{
		fieldSet_: make([]bool, 2),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LogForwarderGroupVersionBuilder) Empty() bool {
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
func (b *LogForwarderGroupVersionBuilder) ID(value string) *LogForwarderGroupVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.id = value
	b.fieldSet_[0] = true
	return b
}

// Applications sets the value of the 'applications' attribute to the given values.
func (b *LogForwarderGroupVersionBuilder) Applications(values ...string) *LogForwarderGroupVersionBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 2)
	}
	b.applications = make([]string, len(values))
	copy(b.applications, values)
	b.fieldSet_[1] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LogForwarderGroupVersionBuilder) Copy(object *LogForwarderGroupVersion) *LogForwarderGroupVersionBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	if object.applications != nil {
		b.applications = make([]string, len(object.applications))
		copy(b.applications, object.applications)
	} else {
		b.applications = nil
	}
	return b
}

// Build creates a 'log_forwarder_group_version' object using the configuration stored in the builder.
func (b *LogForwarderGroupVersionBuilder) Build() (object *LogForwarderGroupVersion, err error) {
	object = new(LogForwarderGroupVersion)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.id = b.id
	if b.applications != nil {
		object.applications = make([]string, len(b.applications))
		copy(object.applications, b.applications)
	}
	return
}
