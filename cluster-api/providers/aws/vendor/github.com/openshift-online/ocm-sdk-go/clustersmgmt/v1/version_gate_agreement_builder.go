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

// VersionGateAgreementBuilder contains the data and logic needed to build 'version_gate_agreement' objects.
//
// VersionGateAgreement represents a version gate that the user agreed to for a specific cluster.
type VersionGateAgreementBuilder struct {
	bitmap_         uint32
	id              string
	href            string
	agreedTimestamp time.Time
	versionGate     *VersionGateBuilder
}

// NewVersionGateAgreement creates a new builder of 'version_gate_agreement' objects.
func NewVersionGateAgreement() *VersionGateAgreementBuilder {
	return &VersionGateAgreementBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *VersionGateAgreementBuilder) Link(value bool) *VersionGateAgreementBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *VersionGateAgreementBuilder) ID(value string) *VersionGateAgreementBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *VersionGateAgreementBuilder) HREF(value string) *VersionGateAgreementBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *VersionGateAgreementBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// AgreedTimestamp sets the value of the 'agreed_timestamp' attribute to the given value.
func (b *VersionGateAgreementBuilder) AgreedTimestamp(value time.Time) *VersionGateAgreementBuilder {
	b.agreedTimestamp = value
	b.bitmap_ |= 8
	return b
}

// VersionGate sets the value of the 'version_gate' attribute to the given value.
//
// Representation of an _OpenShift_ version gate.
func (b *VersionGateAgreementBuilder) VersionGate(value *VersionGateBuilder) *VersionGateAgreementBuilder {
	b.versionGate = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *VersionGateAgreementBuilder) Copy(object *VersionGateAgreement) *VersionGateAgreementBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.agreedTimestamp = object.agreedTimestamp
	if object.versionGate != nil {
		b.versionGate = NewVersionGate().Copy(object.versionGate)
	} else {
		b.versionGate = nil
	}
	return b
}

// Build creates a 'version_gate_agreement' object using the configuration stored in the builder.
func (b *VersionGateAgreementBuilder) Build() (object *VersionGateAgreement, err error) {
	object = new(VersionGateAgreement)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.agreedTimestamp = b.agreedTimestamp
	if b.versionGate != nil {
		object.versionGate, err = b.versionGate.Build()
		if err != nil {
			return
		}
	}
	return
}
