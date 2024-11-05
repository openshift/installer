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

// OidcThumbprintBuilder contains the data and logic needed to build 'oidc_thumbprint' objects.
//
// Contains the necessary attributes to support oidc configuration thumbprint operations such as fetching/creation of a thumbprint
type OidcThumbprintBuilder struct {
	bitmap_      uint32
	href         string
	clusterId    string
	kind         string
	oidcConfigId string
	thumbprint   string
}

// NewOidcThumbprint creates a new builder of 'oidc_thumbprint' objects.
func NewOidcThumbprint() *OidcThumbprintBuilder {
	return &OidcThumbprintBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *OidcThumbprintBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// HREF sets the value of the 'HREF' attribute to the given value.
func (b *OidcThumbprintBuilder) HREF(value string) *OidcThumbprintBuilder {
	b.href = value
	b.bitmap_ |= 1
	return b
}

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *OidcThumbprintBuilder) ClusterId(value string) *OidcThumbprintBuilder {
	b.clusterId = value
	b.bitmap_ |= 2
	return b
}

// Kind sets the value of the 'kind' attribute to the given value.
func (b *OidcThumbprintBuilder) Kind(value string) *OidcThumbprintBuilder {
	b.kind = value
	b.bitmap_ |= 4
	return b
}

// OidcConfigId sets the value of the 'oidc_config_id' attribute to the given value.
func (b *OidcThumbprintBuilder) OidcConfigId(value string) *OidcThumbprintBuilder {
	b.oidcConfigId = value
	b.bitmap_ |= 8
	return b
}

// Thumbprint sets the value of the 'thumbprint' attribute to the given value.
func (b *OidcThumbprintBuilder) Thumbprint(value string) *OidcThumbprintBuilder {
	b.thumbprint = value
	b.bitmap_ |= 16
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *OidcThumbprintBuilder) Copy(object *OidcThumbprint) *OidcThumbprintBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.href = object.href
	b.clusterId = object.clusterId
	b.kind = object.kind
	b.oidcConfigId = object.oidcConfigId
	b.thumbprint = object.thumbprint
	return b
}

// Build creates a 'oidc_thumbprint' object using the configuration stored in the builder.
func (b *OidcThumbprintBuilder) Build() (object *OidcThumbprint, err error) {
	object = new(OidcThumbprint)
	object.bitmap_ = b.bitmap_
	object.href = b.href
	object.clusterId = b.clusterId
	object.kind = b.kind
	object.oidcConfigId = b.oidcConfigId
	object.thumbprint = b.thumbprint
	return
}
