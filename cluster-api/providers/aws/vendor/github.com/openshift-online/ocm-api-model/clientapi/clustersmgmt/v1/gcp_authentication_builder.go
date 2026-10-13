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

// Google cloud platform authentication method of a cluster.
type GcpAuthenticationBuilder struct {
	fieldSet_ []bool
	href      string
	id        string
	kind      string
}

// NewGcpAuthentication creates a new builder of 'gcp_authentication' objects.
func NewGcpAuthentication() *GcpAuthenticationBuilder {
	return &GcpAuthenticationBuilder{
		fieldSet_: make([]bool, 3),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *GcpAuthenticationBuilder) Empty() bool {
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

// Href sets the value of the 'href' attribute to the given value.
func (b *GcpAuthenticationBuilder) Href(value string) *GcpAuthenticationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.href = value
	b.fieldSet_[0] = true
	return b
}

// Id sets the value of the 'id' attribute to the given value.
func (b *GcpAuthenticationBuilder) Id(value string) *GcpAuthenticationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.id = value
	b.fieldSet_[1] = true
	return b
}

// Kind sets the value of the 'kind' attribute to the given value.
func (b *GcpAuthenticationBuilder) Kind(value string) *GcpAuthenticationBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 3)
	}
	b.kind = value
	b.fieldSet_[2] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *GcpAuthenticationBuilder) Copy(object *GcpAuthentication) *GcpAuthenticationBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.href = object.href
	b.id = object.id
	b.kind = object.kind
	return b
}

// Build creates a 'gcp_authentication' object using the configuration stored in the builder.
func (b *GcpAuthenticationBuilder) Build() (object *GcpAuthentication, err error) {
	object = new(GcpAuthentication)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.href = b.href
	object.id = b.id
	object.kind = b.kind
	return
}
