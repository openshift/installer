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

// Represents a log forwarder group versions configuration.
type LogForwarderGroupVersionsBuilder struct {
	fieldSet_ []bool
	name      string
	versions  []*LogForwarderGroupVersionBuilder
	enabled   bool
}

// NewLogForwarderGroupVersions creates a new builder of 'log_forwarder_group_versions' objects.
func NewLogForwarderGroupVersions() *LogForwarderGroupVersionsBuilder {
	return &LogForwarderGroupVersionsBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *LogForwarderGroupVersionsBuilder) Empty() bool {
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

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *LogForwarderGroupVersionsBuilder) Enabled(value bool) *LogForwarderGroupVersionsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.enabled = value
	b.fieldSet_[0] = true
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *LogForwarderGroupVersionsBuilder) Name(value string) *LogForwarderGroupVersionsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.name = value
	b.fieldSet_[1] = true
	return b
}

// Versions sets the value of the 'versions' attribute to the given values.
func (b *LogForwarderGroupVersionsBuilder) Versions(values ...*LogForwarderGroupVersionBuilder) *LogForwarderGroupVersionsBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.versions = make([]*LogForwarderGroupVersionBuilder, len(values))
	copy(b.versions, values)
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *LogForwarderGroupVersionsBuilder) Copy(object *LogForwarderGroupVersions) *LogForwarderGroupVersionsBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.enabled = object.enabled
	b.name = object.name
	if object.versions != nil {
		b.versions = make([]*LogForwarderGroupVersionBuilder, len(object.versions))
		for i, v := range object.versions {
			b.versions[i] = NewLogForwarderGroupVersion().Copy(v)
		}
	} else {
		b.versions = nil
	}
	return b
}

// Build creates a 'log_forwarder_group_versions' object using the configuration stored in the builder.
func (b *LogForwarderGroupVersionsBuilder) Build() (object *LogForwarderGroupVersions, err error) {
	object = new(LogForwarderGroupVersions)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.enabled = b.enabled
	object.name = b.name
	if b.versions != nil {
		object.versions = make([]*LogForwarderGroupVersion, len(b.versions))
		for i, v := range b.versions {
			object.versions[i], err = v.Build()
			if err != nil {
				return
			}
		}
	}
	return
}
