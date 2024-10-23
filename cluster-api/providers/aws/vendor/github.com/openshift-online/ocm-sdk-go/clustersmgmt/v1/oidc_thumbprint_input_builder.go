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

// OidcThumbprintInputBuilder contains the data and logic needed to build 'oidc_thumbprint_input' objects.
//
// Contains the necessary attributes to fetch an OIDC Configuration thumbprint
type OidcThumbprintInputBuilder struct {
	bitmap_      uint32
	clusterId    string
	oidcConfigId string
}

// NewOidcThumbprintInput creates a new builder of 'oidc_thumbprint_input' objects.
func NewOidcThumbprintInput() *OidcThumbprintInputBuilder {
	return &OidcThumbprintInputBuilder{}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *OidcThumbprintInputBuilder) Empty() bool {
	return b == nil || b.bitmap_ == 0
}

// ClusterId sets the value of the 'cluster_id' attribute to the given value.
func (b *OidcThumbprintInputBuilder) ClusterId(value string) *OidcThumbprintInputBuilder {
	b.clusterId = value
	b.bitmap_ |= 1
	return b
}

// OidcConfigId sets the value of the 'oidc_config_id' attribute to the given value.
func (b *OidcThumbprintInputBuilder) OidcConfigId(value string) *OidcThumbprintInputBuilder {
	b.oidcConfigId = value
	b.bitmap_ |= 2
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *OidcThumbprintInputBuilder) Copy(object *OidcThumbprintInput) *OidcThumbprintInputBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.clusterId = object.clusterId
	b.oidcConfigId = object.oidcConfigId
	return b
}

// Build creates a 'oidc_thumbprint_input' object using the configuration stored in the builder.
func (b *OidcThumbprintInputBuilder) Build() (object *OidcThumbprintInput, err error) {
	object = new(OidcThumbprintInput)
	object.bitmap_ = b.bitmap_
	object.clusterId = b.clusterId
	object.oidcConfigId = b.oidcConfigId
	return
}
