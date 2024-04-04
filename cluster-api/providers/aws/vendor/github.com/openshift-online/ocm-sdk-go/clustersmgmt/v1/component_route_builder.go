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

// ComponentRouteBuilder contains the data and logic needed to build 'component_route' objects.
//
// Representation of a Component Route.
type ComponentRouteBuilder struct {
	bitmap_      uint32
	id           string
	href         string
	hostname     string
	tlsSecretRef string
}

// NewComponentRoute creates a new builder of 'component_route' objects.
func NewComponentRoute() *ComponentRouteBuilder {
	return &ComponentRouteBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ComponentRouteBuilder) Link(value bool) *ComponentRouteBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ComponentRouteBuilder) ID(value string) *ComponentRouteBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ComponentRouteBuilder) HREF(value string) *ComponentRouteBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ComponentRouteBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// Hostname sets the value of the 'hostname' attribute to the given value.
func (b *ComponentRouteBuilder) Hostname(value string) *ComponentRouteBuilder {
	b.hostname = value
	b.bitmap_ |= 8
	return b
}

// TlsSecretRef sets the value of the 'tls_secret_ref' attribute to the given value.
func (b *ComponentRouteBuilder) TlsSecretRef(value string) *ComponentRouteBuilder {
	b.tlsSecretRef = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ComponentRouteBuilder) Copy(object *ComponentRoute) *ComponentRouteBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	b.hostname = object.hostname
	b.tlsSecretRef = object.tlsSecretRef
	return b
}

// Build creates a 'component_route' object using the configuration stored in the builder.
func (b *ComponentRouteBuilder) Build() (object *ComponentRoute, err error) {
	object = new(ComponentRoute)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	object.hostname = b.hostname
	object.tlsSecretRef = b.tlsSecretRef
	return
}
