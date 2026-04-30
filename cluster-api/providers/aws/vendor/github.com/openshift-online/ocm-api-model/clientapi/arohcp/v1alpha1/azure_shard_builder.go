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

// The Azure related configuration of the Provision Shard
type AzureShardBuilder struct {
	fieldSet_                                []bool
	aksManagementClusterResourceId           string
	cxManagedIdentitiesKeyVaultUrl           string
	cxSecretsKeyVaultManagedIdentityClientId string
	cxSecretsKeyVaultUrl                     string
	publicDnsZoneResourceId                  string
}

// NewAzureShard creates a new builder of 'azure_shard' objects.
func NewAzureShard() *AzureShardBuilder {
	return &AzureShardBuilder{
		fieldSet_: make([]bool, 5),
	}
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *AzureShardBuilder) Empty() bool {
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

// AksManagementClusterResourceId sets the value of the 'aks_management_cluster_resource_id' attribute to the given value.
func (b *AzureShardBuilder) AksManagementClusterResourceId(value string) *AzureShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.aksManagementClusterResourceId = value
	b.fieldSet_[0] = true
	return b
}

// CxManagedIdentitiesKeyVaultUrl sets the value of the 'cx_managed_identities_key_vault_url' attribute to the given value.
func (b *AzureShardBuilder) CxManagedIdentitiesKeyVaultUrl(value string) *AzureShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.cxManagedIdentitiesKeyVaultUrl = value
	b.fieldSet_[1] = true
	return b
}

// CxSecretsKeyVaultManagedIdentityClientId sets the value of the 'cx_secrets_key_vault_managed_identity_client_id' attribute to the given value.
func (b *AzureShardBuilder) CxSecretsKeyVaultManagedIdentityClientId(value string) *AzureShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.cxSecretsKeyVaultManagedIdentityClientId = value
	b.fieldSet_[2] = true
	return b
}

// CxSecretsKeyVaultUrl sets the value of the 'cx_secrets_key_vault_url' attribute to the given value.
func (b *AzureShardBuilder) CxSecretsKeyVaultUrl(value string) *AzureShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.cxSecretsKeyVaultUrl = value
	b.fieldSet_[3] = true
	return b
}

// PublicDnsZoneResourceId sets the value of the 'public_dns_zone_resource_id' attribute to the given value.
func (b *AzureShardBuilder) PublicDnsZoneResourceId(value string) *AzureShardBuilder {
	if len(b.fieldSet_) == 0 {
		b.fieldSet_ = make([]bool, 5)
	}
	b.publicDnsZoneResourceId = value
	b.fieldSet_[4] = true
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *AzureShardBuilder) Copy(object *AzureShard) *AzureShardBuilder {
	if object == nil {
		return b
	}
	if len(object.fieldSet_) > 0 {
		b.fieldSet_ = make([]bool, len(object.fieldSet_))
		copy(b.fieldSet_, object.fieldSet_)
	}
	b.aksManagementClusterResourceId = object.aksManagementClusterResourceId
	b.cxManagedIdentitiesKeyVaultUrl = object.cxManagedIdentitiesKeyVaultUrl
	b.cxSecretsKeyVaultManagedIdentityClientId = object.cxSecretsKeyVaultManagedIdentityClientId
	b.cxSecretsKeyVaultUrl = object.cxSecretsKeyVaultUrl
	b.publicDnsZoneResourceId = object.publicDnsZoneResourceId
	return b
}

// Build creates a 'azure_shard' object using the configuration stored in the builder.
func (b *AzureShardBuilder) Build() (object *AzureShard, err error) {
	object = new(AzureShard)
	if len(b.fieldSet_) > 0 {
		object.fieldSet_ = make([]bool, len(b.fieldSet_))
		copy(object.fieldSet_, b.fieldSet_)
	}
	object.aksManagementClusterResourceId = b.aksManagementClusterResourceId
	object.cxManagedIdentitiesKeyVaultUrl = b.cxManagedIdentitiesKeyVaultUrl
	object.cxSecretsKeyVaultManagedIdentityClientId = b.cxSecretsKeyVaultManagedIdentityClientId
	object.cxSecretsKeyVaultUrl = b.cxSecretsKeyVaultUrl
	object.publicDnsZoneResourceId = b.publicDnsZoneResourceId
	return
}
