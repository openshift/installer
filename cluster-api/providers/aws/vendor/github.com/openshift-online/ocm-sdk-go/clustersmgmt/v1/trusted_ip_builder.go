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

// TrustedIpBuilder contains the data and logic needed to build 'trusted_ip' objects.
//
// Representation of a trusted ip address in clusterdeployment.
type TrustedIpBuilder struct {
	bitmap_ uint32
	id      string
	href    string
	enabled bool
}

// NewTrustedIp creates a new builder of 'trusted_ip' objects.
func NewTrustedIp() *TrustedIpBuilder {
	return &TrustedIpBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *TrustedIpBuilder) Link(value bool) *TrustedIpBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *TrustedIpBuilder) ID(value string) *TrustedIpBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *TrustedIpBuilder) HREF(value string) *TrustedIpBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *TrustedIpBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Enabled sets the value of the 'enabled' attribute to the given value.
func (b *TrustedIpBuilder) Enabled(value bool) *TrustedIpBuilder {
	b.enabled = value
	b.bitmap_ |= 8
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *TrustedIpBuilder) Copy(object *TrustedIp) *TrustedIpBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.enabled = object.enabled
	return b
}

// Build creates a 'trusted_ip' object using the configuration stored in the builder.
func (b *TrustedIpBuilder) Build() (object *TrustedIp, err error) {
	object = new(TrustedIp)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.enabled = b.enabled
	return
}
