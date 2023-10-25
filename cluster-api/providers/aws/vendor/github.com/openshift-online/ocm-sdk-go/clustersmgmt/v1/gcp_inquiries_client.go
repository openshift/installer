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
	"net/http"
	"path"
)

// GCPInquiriesClient is the client of the 'GCP_inquiries' resource.
//
// Manages the collection of gcp inquiries.
type GCPInquiriesClient struct {
	transport http.RoundTripper
	path      string
}

// NewGCPInquiriesClient creates a new client for the 'GCP_inquiries'
// resource using the given transport to send the requests and receive the
// responses.
func NewGCPInquiriesClient(transport http.RoundTripper, path string) *GCPInquiriesClient {
	return &GCPInquiriesClient{
		transport: transport,
		path:      path,
	}
}

// EncryptionKeys returns the target 'encryption_keys_inquiry' resource.
//
// Reference to the resource that manages a collection of encryption keys.
func (c *GCPInquiriesClient) EncryptionKeys() *EncryptionKeysInquiryClient {
	return NewEncryptionKeysInquiryClient(
		c.transport,
		path.Join(c.path, "encryption_keys"),
	)
}

// KeyRings returns the target 'key_rings_inquiry' resource.
//
// Reference to the resource that manages a collection of key rings.
func (c *GCPInquiriesClient) KeyRings() *KeyRingsInquiryClient {
	return NewKeyRingsInquiryClient(
		c.transport,
		path.Join(c.path, "key_rings"),
	)
}

// Regions returns the target 'available_regions_inquiry' resource.
//
// Reference to the resource that manages a collection of regions.
func (c *GCPInquiriesClient) Regions() *AvailableRegionsInquiryClient {
	return NewAvailableRegionsInquiryClient(
		c.transport,
		path.Join(c.path, "regions"),
	)
}

// Vpcs returns the target 'vpcs_inquiry' resource.
//
// Reference to the resource that manages a collection of vpcs.
func (c *GCPInquiriesClient) Vpcs() *VpcsInquiryClient {
	return NewVpcsInquiryClient(
		c.transport,
		path.Join(c.path, "vpcs"),
	)
}
