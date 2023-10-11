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

// Client is the client of the 'root' resource.
//
// Root of the tree of resources of the clusters management service.
type Client struct {
	transport http.RoundTripper
	path      string
}

// NewClient creates a new client for the 'root'
// resource using the given transport to send the requests and receive the
// responses.
func NewClient(transport http.RoundTripper, path string) *Client {
	return &Client{
		transport: transport,
		path:      path,
	}
}

// Creates a new request for the method that retrieves the metadata.
func (c *Client) Get() *MetadataRequest {
	return &MetadataRequest{
		transport: c.transport,
		path:      c.path,
	}
}

// AWSInfrastructureAccessRoles returns the target 'AWS_infrastructure_access_roles' resource.
//
// Reference to the resource that manages the collection of AWS
// infrastructure access roles.
func (c *Client) AWSInfrastructureAccessRoles() *AWSInfrastructureAccessRolesClient {
	return NewAWSInfrastructureAccessRolesClient(
		c.transport,
		path.Join(c.path, "aws_infrastructure_access_roles"),
	)
}

// AWSInquiries returns the target 'AWS_inquiries' resource.
//
// Reference to the resource that manages the collection of aws inquiries.
func (c *Client) AWSInquiries() *AWSInquiriesClient {
	return NewAWSInquiriesClient(
		c.transport,
		path.Join(c.path, "aws_inquiries"),
	)
}

// DNSDomains returns the target 'DNS_domains' resource.
//
// Reference to the resource that manages dns domains.
func (c *Client) DNSDomains() *DNSDomainsClient {
	return NewDNSDomainsClient(
		c.transport,
		path.Join(c.path, "dns_domains"),
	)
}

// GCPInquiries returns the target 'GCP_inquiries' resource.
//
// Reference to the resource that manages the collection of gcp inquiries.
func (c *Client) GCPInquiries() *GCPInquiriesClient {
	return NewGCPInquiriesClient(
		c.transport,
		path.Join(c.path, "gcp_inquiries"),
	)
}

// Addons returns the target 'add_ons' resource.
//
// Reference to the resource that manages the collection of add-ons.
func (c *Client) Addons() *AddOnsClient {
	return NewAddOnsClient(
		c.transport,
		path.Join(c.path, "addons"),
	)
}

// CloudProviders returns the target 'cloud_providers' resource.
//
// Reference to the resource that manages the collection of cloud providers.
func (c *Client) CloudProviders() *CloudProvidersClient {
	return NewCloudProvidersClient(
		c.transport,
		path.Join(c.path, "cloud_providers"),
	)
}

// Clusters returns the target 'clusters' resource.
//
// Reference to the resource that manages the collection of clusters.
func (c *Client) Clusters() *ClustersClient {
	return NewClustersClient(
		c.transport,
		path.Join(c.path, "clusters"),
	)
}

// Environment returns the target 'environment' resource.
//
// Reference to the resource that manages the environment.
func (c *Client) Environment() *EnvironmentClient {
	return NewEnvironmentClient(
		c.transport,
		path.Join(c.path, "environment"),
	)
}

// Events returns the target 'events' resource.
//
// Reference to the resource that manages the collection of trackable events.
func (c *Client) Events() *EventsClient {
	return NewEventsClient(
		c.transport,
		path.Join(c.path, "events"),
	)
}

// Flavours returns the target 'flavours' resource.
//
// Reference to the service that manages the collection of flavours.
func (c *Client) Flavours() *FlavoursClient {
	return NewFlavoursClient(
		c.transport,
		path.Join(c.path, "flavours"),
	)
}

// LimitedSupportReasonTemplates returns the target 'limited_support_reason_templates' resource.
//
// Reference to limited support reason templates.
func (c *Client) LimitedSupportReasonTemplates() *LimitedSupportReasonTemplatesClient {
	return NewLimitedSupportReasonTemplatesClient(
		c.transport,
		path.Join(c.path, "limited_support_reason_templates"),
	)
}

// MachineTypes returns the target 'machine_types' resource.
//
// Reference to the resource that manage the collection of machine types.
func (c *Client) MachineTypes() *MachineTypesClient {
	return NewMachineTypesClient(
		c.transport,
		path.Join(c.path, "machine_types"),
	)
}

// NetworkVerifications returns the target 'network_verifications' resource.
//
// Reference to the resource that manages network verifications.
func (c *Client) NetworkVerifications() *NetworkVerificationsClient {
	return NewNetworkVerificationsClient(
		c.transport,
		path.Join(c.path, "network_verifications"),
	)
}

// OidcConfigs returns the target 'oidc_configs' resource.
//
// Reference to the resource that manages oidc.
func (c *Client) OidcConfigs() *OidcConfigsClient {
	return NewOidcConfigsClient(
		c.transport,
		path.Join(c.path, "oidc_configs"),
	)
}

// PendingDeleteClusters returns the target 'pending_delete_clusters' resource.
//
// Reference to the resource that manages the collection of pending delete clusters.
func (c *Client) PendingDeleteClusters() *PendingDeleteClustersClient {
	return NewPendingDeleteClustersClient(
		c.transport,
		path.Join(c.path, "pending_delete_clusters"),
	)
}

// Products returns the target 'products' resource.
//
// Reference to the resource that manages the collection of products.
func (c *Client) Products() *ProductsClient {
	return NewProductsClient(
		c.transport,
		path.Join(c.path, "products"),
	)
}

// ProvisionShards returns the target 'provision_shards' resource.
//
// Reference to the resource that manages the collection of provision shards.
func (c *Client) ProvisionShards() *ProvisionShardsClient {
	return NewProvisionShardsClient(
		c.transport,
		path.Join(c.path, "provision_shards"),
	)
}

// TrustedIPAddresses returns the target 'trusted_ips' resource.
//
// Reference to the resource that manages the collection of trusted ip addresses.
func (c *Client) TrustedIPAddresses() *TrustedIpsClient {
	return NewTrustedIpsClient(
		c.transport,
		path.Join(c.path, "trusted_ip_addresses"),
	)
}

// VersionGates returns the target 'version_gates' resource.
//
// Reference to version gates.
func (c *Client) VersionGates() *VersionGatesClient {
	return NewVersionGatesClient(
		c.transport,
		path.Join(c.path, "version_gates"),
	)
}

// Versions returns the target 'versions' resource.
//
// Reference to the resource that manage the collection of versions.
func (c *Client) Versions() *VersionsClient {
	return NewVersionsClient(
		c.transport,
		path.Join(c.path, "versions"),
	)
}
