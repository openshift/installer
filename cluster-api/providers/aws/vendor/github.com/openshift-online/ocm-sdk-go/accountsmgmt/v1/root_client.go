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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

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

// AccessToken returns the target 'access_token' resource.
//
// Reference to the resource that manages generates access tokens.
func (c *Client) AccessToken() *AccessTokenClient {
	return NewAccessTokenClient(
		c.transport,
		path.Join(c.path, "access_token"),
	)
}

// Accounts returns the target 'accounts' resource.
//
// Reference to the resource that manages the collection of accounts.
func (c *Client) Accounts() *AccountsClient {
	return NewAccountsClient(
		c.transport,
		path.Join(c.path, "accounts"),
	)
}

// BillingModels returns the target 'billing_models' resource.
//
// Reference to the resource that manages billing models.
func (c *Client) BillingModels() *BillingModelsClient {
	return NewBillingModelsClient(
		c.transport,
		path.Join(c.path, "billing_models"),
	)
}

// Capabilities returns the target 'capabilities' resource.
//
// Reference to the resource that manages the collection of capabilities.
func (c *Client) Capabilities() *CapabilitiesClient {
	return NewCapabilitiesClient(
		c.transport,
		path.Join(c.path, "capabilities"),
	)
}

// CloudResources returns the target 'cloud_resources' resource.
//
// Reference to the resource that manages the collection of cloud resources.
func (c *Client) CloudResources() *CloudResourcesClient {
	return NewCloudResourcesClient(
		c.transport,
		path.Join(c.path, "cloud_resources"),
	)
}

// ClusterAuthorizations returns the target 'cluster_authorizations' resource.
//
// Reference to the resource that manages cluster authorizations.
func (c *Client) ClusterAuthorizations() *ClusterAuthorizationsClient {
	return NewClusterAuthorizationsClient(
		c.transport,
		path.Join(c.path, "cluster_authorizations"),
	)
}

// ClusterRegistrations returns the target 'cluster_registrations' resource.
//
// Reference to the resource that manages cluster registrations.
func (c *Client) ClusterRegistrations() *ClusterRegistrationsClient {
	return NewClusterRegistrationsClient(
		c.transport,
		path.Join(c.path, "cluster_registrations"),
	)
}

// CurrentAccess returns the target 'roles' resource.
//
// Reference to the resource that manages the current authenticated
// account.
func (c *Client) CurrentAccess() *RolesClient {
	return NewRolesClient(
		c.transport,
		path.Join(c.path, "current_access"),
	)
}

// CurrentAccount returns the target 'current_account' resource.
//
// Reference to the resource that manages the current authenticated
// account.
func (c *Client) CurrentAccount() *CurrentAccountClient {
	return NewCurrentAccountClient(
		c.transport,
		path.Join(c.path, "current_account"),
	)
}

// DeletedSubscriptions returns the target 'deleted_subscriptions' resource.
//
// Reference to the resource that manages the collection of deleted subscriptions.
func (c *Client) DeletedSubscriptions() *DeletedSubscriptionsClient {
	return NewDeletedSubscriptionsClient(
		c.transport,
		path.Join(c.path, "deleted_subscriptions"),
	)
}

// FeatureToggles returns the target 'feature_toggles' resource.
//
// Reference to the resource that manages feature toggles.
func (c *Client) FeatureToggles() *FeatureTogglesClient {
	return NewFeatureTogglesClient(
		c.transport,
		path.Join(c.path, "feature_toggles"),
	)
}

// Labels returns the target 'labels' resource.
//
// Reference to the resource that manages the collection of labels.
func (c *Client) Labels() *LabelsClient {
	return NewLabelsClient(
		c.transport,
		path.Join(c.path, "labels"),
	)
}

// Notify returns the target 'notify' resource.
//
// Reference to the resource that manages the notifications.
func (c *Client) Notify() *NotifyClient {
	return NewNotifyClient(
		c.transport,
		path.Join(c.path, "notify"),
	)
}

// Organizations returns the target 'organizations' resource.
//
// Reference to the resource that manages the collection of
// organizations.
func (c *Client) Organizations() *OrganizationsClient {
	return NewOrganizationsClient(
		c.transport,
		path.Join(c.path, "organizations"),
	)
}

// Permissions returns the target 'permissions' resource.
//
// Reference to the resource that manages the collection of permissions.
func (c *Client) Permissions() *PermissionsClient {
	return NewPermissionsClient(
		c.transport,
		path.Join(c.path, "permissions"),
	)
}

// PullSecrets returns the target 'pull_secrets' resource.
//
// Reference to the resource that manages generates access tokens.
func (c *Client) PullSecrets() *PullSecretsClient {
	return NewPullSecretsClient(
		c.transport,
		path.Join(c.path, "pull_secrets"),
	)
}

// QuotaAuthorizations returns the target 'quota_authorizations' resource.
//
// Reference to the resource that manages quota authorizations.
func (c *Client) QuotaAuthorizations() *QuotaAuthorizationsClient {
	return NewQuotaAuthorizationsClient(
		c.transport,
		path.Join(c.path, "quota_authorizations"),
	)
}

// Registries returns the target 'registries' resource.
//
// Reference to the resource that manages the collection of registries.
func (c *Client) Registries() *RegistriesClient {
	return NewRegistriesClient(
		c.transport,
		path.Join(c.path, "registries"),
	)
}

// RegistryCredentials returns the target 'registry_credentials' resource.
//
// Reference to the resource that manages the collection of registry
// credentials.
func (c *Client) RegistryCredentials() *RegistryCredentialsClient {
	return NewRegistryCredentialsClient(
		c.transport,
		path.Join(c.path, "registry_credentials"),
	)
}

// ResourceQuota returns the target 'resource_quotas' resource.
//
// Reference to the resource that manages the collection of resource
// quota.
func (c *Client) ResourceQuota() *ResourceQuotasClient {
	return NewResourceQuotasClient(
		c.transport,
		path.Join(c.path, "resource_quota"),
	)
}

// RoleBindings returns the target 'role_bindings' resource.
//
// Reference to the resource that manages the collection of role
// bindings.
func (c *Client) RoleBindings() *RoleBindingsClient {
	return NewRoleBindingsClient(
		c.transport,
		path.Join(c.path, "role_bindings"),
	)
}

// Roles returns the target 'roles' resource.
//
// Reference to the resource that manages the collection of roles.
func (c *Client) Roles() *RolesClient {
	return NewRolesClient(
		c.transport,
		path.Join(c.path, "roles"),
	)
}

// SkuRules returns the target 'sku_rules' resource.
//
// Reference to the resource that manages the collection of
// Sku Rules
func (c *Client) SkuRules() *SkuRulesClient {
	return NewSkuRulesClient(
		c.transport,
		path.Join(c.path, "sku_rules"),
	)
}

// Subscriptions returns the target 'subscriptions' resource.
//
// Reference to the resource that manages the collection of
// subscriptions.
func (c *Client) Subscriptions() *SubscriptionsClient {
	return NewSubscriptionsClient(
		c.transport,
		path.Join(c.path, "subscriptions"),
	)
}

// SupportCases returns the target 'support_cases' resource.
//
// Reference to the resource that manages the support cases.
func (c *Client) SupportCases() *SupportCasesClient {
	return NewSupportCasesClient(
		c.transport,
		path.Join(c.path, "support_cases"),
	)
}

// TokenAuthorization returns the target 'token_authorization' resource.
//
// Reference to the resource that manages token authorization.
func (c *Client) TokenAuthorization() *TokenAuthorizationClient {
	return NewTokenAuthorizationClient(
		c.transport,
		path.Join(c.path, "token_authorization"),
	)
}
