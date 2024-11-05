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

// AWSInquiriesClient is the client of the 'AWS_inquiries' resource.
//
// Manages the collection of aws inquiries.
type AWSInquiriesClient struct {
	transport http.RoundTripper
	path      string
}

// NewAWSInquiriesClient creates a new client for the 'AWS_inquiries'
// resource using the given transport to send the requests and receive the
// responses.
func NewAWSInquiriesClient(transport http.RoundTripper, path string) *AWSInquiriesClient {
	return &AWSInquiriesClient{
		transport: transport,
		path:      path,
	}
}

// STSAccountRoles returns the target 'AWSSTS_account_roles_inquiry' resource.
//
// Reference to the resource that manages aws sts roles.
func (c *AWSInquiriesClient) STSAccountRoles() *AWSSTSAccountRolesInquiryClient {
	return NewAWSSTSAccountRolesInquiryClient(
		c.transport,
		path.Join(c.path, "sts_account_roles"),
	)
}

// STSCredentialRequests returns the target 'STS_credential_requests_inquiry' resource.
//
// Reference to the resource that manages sts cred request.
func (c *AWSInquiriesClient) STSCredentialRequests() *STSCredentialRequestsInquiryClient {
	return NewSTSCredentialRequestsInquiryClient(
		c.transport,
		path.Join(c.path, "sts_credential_requests"),
	)
}

// STSPolicies returns the target 'AWSSTS_policies_inquiry' resource.
//
// Reference to the resource that manages aws sts policies.
func (c *AWSInquiriesClient) STSPolicies() *AWSSTSPoliciesInquiryClient {
	return NewAWSSTSPoliciesInquiryClient(
		c.transport,
		path.Join(c.path, "sts_policies"),
	)
}

// MachineTypes returns the target 'AWS_region_machine_types_inquiry' resource.
//
// Reference to the resource that manages aws machine types by regions.
func (c *AWSInquiriesClient) MachineTypes() *AWSRegionMachineTypesInquiryClient {
	return NewAWSRegionMachineTypesInquiryClient(
		c.transport,
		path.Join(c.path, "machine_types"),
	)
}

// OidcThumbprint returns the target 'oidc_thumbprint' resource.
//
// Reference to the resource that manages OIDC Config Thumbprint fetching.
func (c *AWSInquiriesClient) OidcThumbprint() *OidcThumbprintClient {
	return NewOidcThumbprintClient(
		c.transport,
		path.Join(c.path, "oidc_thumbprint"),
	)
}

// Regions returns the target 'available_regions_inquiry' resource.
//
// Reference to the resource that manages a collection of regions.
func (c *AWSInquiriesClient) Regions() *AvailableRegionsInquiryClient {
	return NewAvailableRegionsInquiryClient(
		c.transport,
		path.Join(c.path, "regions"),
	)
}

// ValidateCredentials returns the target 'aws_validate_credentials' resource.
//
// Reference to the resource that manages creds validation.
func (c *AWSInquiriesClient) ValidateCredentials() *AwsValidateCredentialsClient {
	return NewAwsValidateCredentialsClient(
		c.transport,
		path.Join(c.path, "validate_credentials"),
	)
}

// Vpcs returns the target 'vpcs_inquiry' resource.
//
// Reference to the resource that manages a collection of vpcs.
func (c *AWSInquiriesClient) Vpcs() *VpcsInquiryClient {
	return NewVpcsInquiryClient(
		c.transport,
		path.Join(c.path, "vpcs"),
	)
}
