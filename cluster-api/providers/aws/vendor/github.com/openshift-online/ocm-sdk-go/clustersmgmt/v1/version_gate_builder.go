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

import (
	time "time"
)

// VersionGateBuilder contains the data and logic needed to build 'version_gate' objects.
//
// Representation of an _OpenShift_ version gate.
type VersionGateBuilder struct {
	bitmap_            uint32
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
	return &VersionGateBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *VersionGateBuilder) Link(value bool) *VersionGateBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *VersionGateBuilder) ID(value string) *VersionGateBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *VersionGateBuilder) HREF(value string) *VersionGateBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *VersionGateBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// STSOnly sets the value of the 'STS_only' attribute to the given value.
func (b *VersionGateBuilder) STSOnly(value bool) *VersionGateBuilder {
	b.stsOnly = value
	b.bitmap_ |= 8
	return b
}

// ClusterCondition sets the value of the 'cluster_condition' attribute to the given value.
func (b *VersionGateBuilder) ClusterCondition(value string) *VersionGateBuilder {
	b.clusterCondition = value
	b.bitmap_ |= 16
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *VersionGateBuilder) CreationTimestamp(value time.Time) *VersionGateBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 32
	return b
}

// Description sets the value of the 'description' attribute to the given value.
func (b *VersionGateBuilder) Description(value string) *VersionGateBuilder {
	b.description = value
	b.bitmap_ |= 64
	return b
}

// DocumentationURL sets the value of the 'documentation_URL' attribute to the given value.
func (b *VersionGateBuilder) DocumentationURL(value string) *VersionGateBuilder {
	b.documentationURL = value
	b.bitmap_ |= 128
	return b
}

// Label sets the value of the 'label' attribute to the given value.
func (b *VersionGateBuilder) Label(value string) *VersionGateBuilder {
	b.label = value
	b.bitmap_ |= 256
	return b
}

// Value sets the value of the 'value' attribute to the given value.
func (b *VersionGateBuilder) Value(value string) *VersionGateBuilder {
	b.value = value
	b.bitmap_ |= 512
	return b
}

// VersionRawIDPrefix sets the value of the 'version_raw_ID_prefix' attribute to the given value.
func (b *VersionGateBuilder) VersionRawIDPrefix(value string) *VersionGateBuilder {
	b.versionRawIDPrefix = value
	b.bitmap_ |= 1024
	return b
}

// WarningMessage sets the value of the 'warning_message' attribute to the given value.
func (b *VersionGateBuilder) WarningMessage(value string) *VersionGateBuilder {
	b.warningMessage = value
	b.bitmap_ |= 2048
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *VersionGateBuilder) Copy(object *VersionGate) *VersionGateBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
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
	object.bitmap_ = b.bitmap_
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
