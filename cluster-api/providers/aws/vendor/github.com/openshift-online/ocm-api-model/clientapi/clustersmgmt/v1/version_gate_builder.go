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

import (
	time "time"
)

// Representation of an _OpenShift_ version gate.
type VersionGateBuilder struct {
	fieldSet_          []bool
	id                 string
	href               string
	clusterCondition   string
	creationTimestamp  time.Time
	description        string
	documentationURL   string
	label              string
	value              string
	versionRawIDPrefix string
	warningMessage     string
	stsOnly            bool
}

// NewVersionGate creates a new builder of 'version_gate' objects.
func NewVersionGate() *VersionGateBuilder {
	return &VersionGateBuilder{
		fieldSet_: make([]bool, 12),
	}
}

// Link sets the flag that indicates if this is a link.
func (b *VersionGateBuilder) Link(value bool) *VersionGateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.fieldSet_[0] = true
	return b
}

// ID sets the identifier of the object.
func (b *VersionGateBuilder) ID(value string) *VersionGateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// HREF sets the link to the object.
func (b *VersionGateBuilder) HREF(value string) *VersionGateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.href = value
	b.fieldSet_[2] = true
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *VersionGateBuilder) Empty() bool {
	if b == nil || len(b.fieldSet_) == 0 {
		return true
	}
	// Check all fields except the link flag (index 0)
	for i := 1; i < len(b.fieldSet_); i++ {
		if b.fieldSet_[i] {
			return false
		}
	}
	return true
}

// STSOnly sets the value of the 'STS_only' attribute to the given value.
func (b *VersionGateBuilder) STSOnly(value bool) *VersionGateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.stsOnly = value
	b.fieldSet_[3] = true
	return b
}

// ClusterCondition sets the value of the 'cluster_condition' attribute to the given value.
func (b *VersionGateBuilder) ClusterCondition(value string) *VersionGateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.clusterCondition = value
	b.fieldSet_[4] = true
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *VersionGateBuilder) CreationTimestamp(value time.Time) *VersionGateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.creationTimestamp = value
	b.fieldSet_[5] = true
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *VersionGateBuilder) Description(value string) *VersionGateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.description = value
	b.fieldSet_[6] = true
	return b
}

// DocumentationURL sets the value of the 'documentation_URL' attribute to the given value.
func (b *VersionGateBuilder) DocumentationURL(value string) *VersionGateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.documentationURL = value
	b.fieldSet_[7] = true
	return b
}

// Label sets the value of the 'label' attribute to the given value.
func (b *VersionGateBuilder) Label(value string) *VersionGateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.label = value
	b.fieldSet_[8] = true
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *VersionGateBuilder) Value(value string) *VersionGateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.value = value
	b.fieldSet_[9] = true
	return b
}

// VersionRawIDPrefix sets the value of the 'version_raw_ID_prefix' attribute to the given value.
func (b *VersionGateBuilder) VersionRawIDPrefix(value string) *VersionGateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.versionRawIDPrefix = value
	b.fieldSet_[10] = true
	return b
}

// WarningMessage sets the value of the 'warning_message' attribute to the given value.
func (b *VersionGateBuilder) WarningMessage(value string) *VersionGateBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 12)
	}
	b.warningMessage = value
	b.fieldSet_[11] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *VersionGateBuilder) Copy(object *VersionGate) *VersionGateBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.id = object.id
	b.href = object.href
	b.stsOnly = object.stsOnly
	b.clusterCondition = object.clusterCondition
	b.creationTimestamp = object.creationTimestamp
	b.description = object.description
	b.documentationURL = object.documentationURL
	b.label = object.label
	b.value = object.value
	b.versionRawIDPrefix = object.versionRawIDPrefix
	b.warningMessage = object.warningMessage
	return b
}

// Build creates a 'version_gate' object using the configuration stored in the builder.
func (b *VersionGateBuilder) Build() (object *VersionGate, err error) {
	object = new(VersionGate)
	object.id = b.id
	object.href = b.href
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.stsOnly = b.stsOnly
	object.clusterCondition = b.clusterCondition
	object.creationTimestamp = b.creationTimestamp
	object.description = b.description
	object.documentationURL = b.documentationURL
	object.label = b.label
	object.value = b.value
	object.versionRawIDPrefix = b.versionRawIDPrefix
	object.warningMessage = b.warningMessage
	return
}
