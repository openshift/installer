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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

// AzureShard represents the values of the 'azure_shard' type.
//
// The Azure related configuration of the Provision Shard
type AzureShard struct {
	fieldSet_                                []bool
	aksManagementClusterResourceId           string
	cxManagedIdentitiesKeyVaultUrl           string
	cxSecretsKeyVaultManagedIdentityClientId string
	cxSecretsKeyVaultUrl                     string
	publicDnsZoneResourceId                  string
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *AzureShard) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}
	for _, set := range o.fieldSet_ {
		if set {
			return false
		}
	}
	return true
}

// AksManagementClusterResourceId returns the value of the 'aks_management_cluster_resource_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Resource ID of the AKS Management Cluster associated
// to the provision shard.
// It is the Azure Resource ID of a pre-existing Azure Kubernetes
// Service (AKS) Managed Cluster.
// `aks_management_cluster_resource_id` must be located in the same
// Azure Tenant as the Clusters Service's Azure infrastructure.
// The Azure Resource Group Name specified as part of `aks_management_cluster_resource_id`
// must be in the same Azure location to where Clusters Service is deployed.
// It must be unique across Azure based provision shards.
// Required during creation.
// Immutable.
func (o *AzureShard) AksManagementClusterResourceId() string {
	if o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return o.aksManagementClusterResourceId
	}
	return ""
}

// GetAksManagementClusterResourceId returns the value of the 'aks_management_cluster_resource_id' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure Resource ID of the AKS Management Cluster associated
// to the provision shard.
// It is the Azure Resource ID of a pre-existing Azure Kubernetes
// Service (AKS) Managed Cluster.
// `aks_management_cluster_resource_id` must be located in the same
// Azure Tenant as the Clusters Service's Azure infrastructure.
// The Azure Resource Group Name specified as part of `aks_management_cluster_resource_id`
// must be in the same Azure location to where Clusters Service is deployed.
// It must be unique across Azure based provision shards.
// Required during creation.
// Immutable.
func (o *AzureShard) GetAksManagementClusterResourceId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
	if ok {
		value = o.aksManagementClusterResourceId
	}
	return
}

// CxManagedIdentitiesKeyVaultUrl returns the value of the 'cx_managed_identities_key_vault_url' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Key Vault URL used to store the certificates associated to the
// managed identities of the aro-hcp clusters control plane operators.
// The Azure Key Vault associated to the URL must pre-exist and be in the same Azure location
// to where Clusters Service is deployed.
// The URL must be a well-formed absolute URL.
// The URL must have an HTTPS scheme.
// The URL cannot contain HTTP query parameters.
// The URL path must be '/'.
// The expected url format naming scheme is: https://<name>.<known-host-domain>/
// where both <name> and <known-host-domain> must be specified. The <name> part
// cannot contain '.' characters. <known-host-domain> must be one of the
// known domain names for Key Vault. See https://aka.ms/azsdk/blog/vault-uri#recommended-actions
// for more information.
// Example of a URL: https://mykeyvault.vault.azure.net/
// Required during creation.
// Immutable.
func (o *AzureShard) CxManagedIdentitiesKeyVaultUrl() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.cxManagedIdentitiesKeyVaultUrl
	}
	return ""
}

// GetCxManagedIdentitiesKeyVaultUrl returns the value of the 'cx_managed_identities_key_vault_url' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure Key Vault URL used to store the certificates associated to the
// managed identities of the aro-hcp clusters control plane operators.
// The Azure Key Vault associated to the URL must pre-exist and be in the same Azure location
// to where Clusters Service is deployed.
// The URL must be a well-formed absolute URL.
// The URL must have an HTTPS scheme.
// The URL cannot contain HTTP query parameters.
// The URL path must be '/'.
// The expected url format naming scheme is: https://<name>.<known-host-domain>/
// where both <name> and <known-host-domain> must be specified. The <name> part
// cannot contain '.' characters. <known-host-domain> must be one of the
// known domain names for Key Vault. See https://aka.ms/azsdk/blog/vault-uri#recommended-actions
// for more information.
// Example of a URL: https://mykeyvault.vault.azure.net/
// Required during creation.
// Immutable.
func (o *AzureShard) GetCxManagedIdentitiesKeyVaultUrl() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.cxManagedIdentitiesKeyVaultUrl
	}
	return
}

// CxSecretsKeyVaultManagedIdentityClientId returns the value of the 'cx_secrets_key_vault_managed_identity_client_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Client ID of an Azure User-Assigned Managed Identity
// that would give access to the CX Secrets Key Vault.
// The Azure User-Assigned Managed Identity associated to the Client ID
// must be pre-exist and be in the same Azure location to where Clusters
// Service is deployed.
// The Azure Resource Group of the Azure User-Assigned Managed Identity
// must be in the same Azure location to where Clusters Service is
// deployed.
// Required during creation.
// Immutable.
func (o *AzureShard) CxSecretsKeyVaultManagedIdentityClientId() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.cxSecretsKeyVaultManagedIdentityClientId
	}
	return ""
}

// GetCxSecretsKeyVaultManagedIdentityClientId returns the value of the 'cx_secrets_key_vault_managed_identity_client_id' attribute and
// a flag indicating if the attribute has a value.
//
// The Client ID of an Azure User-Assigned Managed Identity
// that would give access to the CX Secrets Key Vault.
// The Azure User-Assigned Managed Identity associated to the Client ID
// must be pre-exist and be in the same Azure location to where Clusters
// Service is deployed.
// The Azure Resource Group of the Azure User-Assigned Managed Identity
// must be in the same Azure location to where Clusters Service is
// deployed.
// Required during creation.
// Immutable.
func (o *AzureShard) GetCxSecretsKeyVaultManagedIdentityClientId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.cxSecretsKeyVaultManagedIdentityClientId
	}
	return
}

// CxSecretsKeyVaultUrl returns the value of the 'cx_secrets_key_vault_url' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Key Vault URL used to store different types of secrets
// like k8s api certificates, image pull secrets, ...
// The Azure Key Vault associated to the URL must pre-exist and be in the same Azure location
// to where Clusters Service is deployed.
// The URL must be a well-formed absolute URL.
// The URL must have an HTTPS scheme.
// The URL cannot contain HTTP query parameters.
// The URL path must be '/'.
// The expected url format naming scheme is: https://<name>.<known-host-domain>/
// where both <name> and <known-host-domain> must be specified. The <name> part
// cannot contain '.' characters. <known-host-domain> must be one of the
// known domain names for Key Vault. See https://aka.ms/azsdk/blog/vault-uri#recommended-actions
// for more information.
// Example of a URL: https://mykeyvault.vault.azure.net/
// Required during creation.
// Immutable.
func (o *AzureShard) CxSecretsKeyVaultUrl() string {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.cxSecretsKeyVaultUrl
	}
	return ""
}

// GetCxSecretsKeyVaultUrl returns the value of the 'cx_secrets_key_vault_url' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure Key Vault URL used to store different types of secrets
// like k8s api certificates, image pull secrets, ...
// The Azure Key Vault associated to the URL must pre-exist and be in the same Azure location
// to where Clusters Service is deployed.
// The URL must be a well-formed absolute URL.
// The URL must have an HTTPS scheme.
// The URL cannot contain HTTP query parameters.
// The URL path must be '/'.
// The expected url format naming scheme is: https://<name>.<known-host-domain>/
// where both <name> and <known-host-domain> must be specified. The <name> part
// cannot contain '.' characters. <known-host-domain> must be one of the
// known domain names for Key Vault. See https://aka.ms/azsdk/blog/vault-uri#recommended-actions
// for more information.
// Example of a URL: https://mykeyvault.vault.azure.net/
// Required during creation.
// Immutable.
func (o *AzureShard) GetCxSecretsKeyVaultUrl() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.cxSecretsKeyVaultUrl
	}
	return
}

// PublicDnsZoneResourceId returns the value of the 'public_dns_zone_resource_id' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The Azure Public DNS Zone associated to the provision shard.
// It is the Azure Resource ID of a pre-existing Public DNS Zone.
// `public_dns_zone_resource_id` must be located in the same
// Azure Tenant as the Clusters Service's Azure infrastructure
// The Azure Resource Group Name specified as part of `public_dns_zone_resource_id`
// must be in the same Azure location to where Clusters Service is deployed.
// Required during creation.
// Immutable.
func (o *AzureShard) PublicDnsZoneResourceId() string {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.publicDnsZoneResourceId
	}
	return ""
}

// GetPublicDnsZoneResourceId returns the value of the 'public_dns_zone_resource_id' attribute and
// a flag indicating if the attribute has a value.
//
// The Azure Public DNS Zone associated to the provision shard.
// It is the Azure Resource ID of a pre-existing Public DNS Zone.
// `public_dns_zone_resource_id` must be located in the same
// Azure Tenant as the Clusters Service's Azure infrastructure
// The Azure Resource Group Name specified as part of `public_dns_zone_resource_id`
// must be in the same Azure location to where Clusters Service is deployed.
// Required during creation.
// Immutable.
func (o *AzureShard) GetPublicDnsZoneResourceId() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.publicDnsZoneResourceId
	}
	return
}

// AzureShardListKind is the name of the type used to represent list of objects of
// type 'azure_shard'.
const AzureShardListKind = "AzureShardList"

// AzureShardListLinkKind is the name of the type used to represent links to list
// of objects of type 'azure_shard'.
const AzureShardListLinkKind = "AzureShardListLink"

// AzureShardNilKind is the name of the type used to nil lists of objects of
// type 'azure_shard'.
const AzureShardListNilKind = "AzureShardListNil"

// AzureShardList is a list of values of the 'azure_shard' type.
type AzureShardList struct {
	href  string
	link  bool
	items []*AzureShard
}

// Len returns the length of the list.
func (l *AzureShardList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *AzureShardList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *AzureShardList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *AzureShardList) SetItems(items []*AzureShard) {
	l.items = items
}

// Items returns the items of the list.
func (l *AzureShardList) Items() []*AzureShard {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *AzureShardList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *AzureShardList) Get(i int) *AzureShard {
	if l == nil || i < 0 || i >= len(l.items) {
		return nil
	}
	return l.items[i]
}

// Slice returns an slice containing the items of the list. The returned slice is a
// copy of the one used internally, so it can be modified without affecting the
// internal representation.
//
// If you don't need to modify the returned slice consider using the Each or Range
// functions, as they don't need to allocate a new slice.
func (l *AzureShardList) Slice() []*AzureShard {
	var slice []*AzureShard
	if l == nil {
		slice = make([]*AzureShard, 0)
	} else {
		slice = make([]*AzureShard, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *AzureShardList) Each(f func(item *AzureShard) bool) {
	if l == nil {
		return
	}
	for _, item := range l.items {
		if !f(item) {
			break
		}
	}
}

// Range runs the given function for each index and item of the list, in order. If
// the function returns false the iteration stops, otherwise it continues till all
// the elements of the list have been processed.
func (l *AzureShardList) Range(f func(index int, item *AzureShard) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
