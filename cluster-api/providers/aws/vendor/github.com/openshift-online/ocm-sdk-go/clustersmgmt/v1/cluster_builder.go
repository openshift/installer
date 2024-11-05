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
	time "time"
)

// ClusterBuilder contains the data and logic needed to build 'cluster' objects.
//
// Definition of an _OpenShift_ cluster.
//
// The `cloud_provider` attribute is a reference to the cloud provider. When a
// cluster is retrieved it will be a link to the cloud provider, containing only
// the kind, id and href attributes:
//
// ```json
//
//	{
//	  "cloud_provider": {
//	    "kind": "CloudProviderLink",
//	    "id": "123",
//	    "href": "/api/clusters_mgmt/v1/cloud_providers/123"
//	  }
//	}
//
// ```
//
// When a cluster is created this is optional, and if used it should contain the
// identifier of the cloud provider to use:
//
// ```json
//
//	{
//	  "cloud_provider": {
//	    "id": "123",
//	  }
//	}
//
// ```
//
// If not included, then the cluster will be created using the default cloud
// provider, which is currently Amazon Web Services.
//
// The region attribute is mandatory when a cluster is created.
//
// The `aws.access_key_id`, `aws.secret_access_key` and `dns.base_domain`
// attributes are mandatory when creation a cluster with your own Amazon Web
// Services account.
type ClusterBuilder struct {
	bitmap_                           uint64
	id                                string
	href                              string
	api                               *ClusterAPIBuilder
	aws                               *AWSBuilder
	awsInfrastructureAccessRoleGrants *AWSInfrastructureAccessRoleGrantListBuilder
	ccs                               *CCSBuilder
	dns                               *DNSBuilder
	gcp                               *GCPBuilder
	gcpEncryptionKey                  *GCPEncryptionKeyBuilder
	gcpNetwork                        *GCPNetworkBuilder
	additionalTrustBundle             string
	addons                            *AddOnInstallationListBuilder
	autoscaler                        *ClusterAutoscalerBuilder
	azure                             *AzureBuilder
	billingModel                      BillingModel
	byoOidc                           *ByoOidcBuilder
	cloudProvider                     *CloudProviderBuilder
	console                           *ClusterConsoleBuilder
	creationTimestamp                 time.Time
	deleteProtection                  *DeleteProtectionBuilder
	domainPrefix                      string
	expirationTimestamp               time.Time
	externalID                        string
	externalAuthConfig                *ExternalAuthConfigBuilder
	externalConfiguration             *ExternalConfigurationBuilder
	flavour                           *FlavourBuilder
	groups                            *GroupListBuilder
	healthState                       ClusterHealthState
	htpasswd                          *HTPasswdIdentityProviderBuilder
	hypershift                        *HypershiftBuilder
	identityProviders                 *IdentityProviderListBuilder
	inflightChecks                    *InflightCheckListBuilder
	infraID                           string
	ingresses                         *IngressListBuilder
	kubeletConfig                     *KubeletConfigBuilder
	loadBalancerQuota                 int
	machinePools                      *MachinePoolListBuilder
	managedService                    *ManagedServiceBuilder
	name                              string
	network                           *NetworkBuilder
	nodeDrainGracePeriod              *ValueBuilder
	nodePools                         *NodePoolListBuilder
	nodes                             *ClusterNodesBuilder
	openshiftVersion                  string
	product                           *ProductBuilder
	properties                        map[string]string
	provisionShard                    *ProvisionShardBuilder
	proxy                             *ProxyBuilder
	region                            *CloudRegionBuilder
	registryConfig                    *ClusterRegistryConfigBuilder
	state                             ClusterState
	status                            *ClusterStatusBuilder
	storageQuota                      *ValueBuilder
	subscription                      *SubscriptionBuilder
	version                           *VersionBuilder
	fips                              bool
	disableUserWorkloadMonitoring     bool
	etcdEncryption                    bool
	managed                           bool
	multiAZ                           bool
	multiArchEnabled                  bool
}

// NewCluster creates a new builder of 'cluster' objects.
func NewCluster() *ClusterBuilder {
	return &ClusterBuilder{}
}

// Link sets the flag that indicates if this is a link.
func (b *ClusterBuilder) Link(value bool) *ClusterBuilder {
	b.bitmap_ |= 1
	return b
}

// ID sets the identifier of the object.
func (b *ClusterBuilder) ID(value string) *ClusterBuilder {
	b.id = value
	b.bitmap_ |= 2
	return b
}

// HREF sets the link to the object.
func (b *ClusterBuilder) HREF(value string) *ClusterBuilder {
	b.href = value
	b.bitmap_ |= 4
	return b
}

// Empty returns true if the builder is empty, i.e. no attribute has a value.
func (b *ClusterBuilder) Empty() bool {
	return b == nil || b.bitmap_&^1 == 0
}

// API sets the value of the 'API' attribute to the given value.
//
// Information about the API of a cluster.
func (b *ClusterBuilder) API(value *ClusterAPIBuilder) *ClusterBuilder {
	b.api = value
	if value != nil {
		b.bitmap_ |= 8
	} else {
		b.bitmap_ &^= 8
	}
	return b
}

// AWS sets the value of the 'AWS' attribute to the given value.
//
// _Amazon Web Services_ specific settings of a cluster.
func (b *ClusterBuilder) AWS(value *AWSBuilder) *ClusterBuilder {
	b.aws = value
	if value != nil {
		b.bitmap_ |= 16
	} else {
		b.bitmap_ &^= 16
	}
	return b
}

// AWSInfrastructureAccessRoleGrants sets the value of the 'AWS_infrastructure_access_role_grants' attribute to the given values.
func (b *ClusterBuilder) AWSInfrastructureAccessRoleGrants(value *AWSInfrastructureAccessRoleGrantListBuilder) *ClusterBuilder {
	b.awsInfrastructureAccessRoleGrants = value
	b.bitmap_ |= 32
	return b
}

// CCS sets the value of the 'CCS' attribute to the given value.
func (b *ClusterBuilder) CCS(value *CCSBuilder) *ClusterBuilder {
	b.ccs = value
	if value != nil {
		b.bitmap_ |= 64
	} else {
		b.bitmap_ &^= 64
	}
	return b
}

// DNS sets the value of the 'DNS' attribute to the given value.
//
// DNS settings of the cluster.
func (b *ClusterBuilder) DNS(value *DNSBuilder) *ClusterBuilder {
	b.dns = value
	if value != nil {
		b.bitmap_ |= 128
	} else {
		b.bitmap_ &^= 128
	}
	return b
}

// FIPS sets the value of the 'FIPS' attribute to the given value.
func (b *ClusterBuilder) FIPS(value bool) *ClusterBuilder {
	b.fips = value
	b.bitmap_ |= 256
	return b
}

// GCP sets the value of the 'GCP' attribute to the given value.
//
// Google cloud platform settings of a cluster.
func (b *ClusterBuilder) GCP(value *GCPBuilder) *ClusterBuilder {
	b.gcp = value
	if value != nil {
		b.bitmap_ |= 512
	} else {
		b.bitmap_ &^= 512
	}
	return b
}

// GCPEncryptionKey sets the value of the 'GCP_encryption_key' attribute to the given value.
//
// GCP Encryption Key for CCS clusters.
func (b *ClusterBuilder) GCPEncryptionKey(value *GCPEncryptionKeyBuilder) *ClusterBuilder {
	b.gcpEncryptionKey = value
	if value != nil {
		b.bitmap_ |= 1024
	} else {
		b.bitmap_ &^= 1024
	}
	return b
}

// GCPNetwork sets the value of the 'GCP_network' attribute to the given value.
//
// GCP Network configuration of a cluster.
func (b *ClusterBuilder) GCPNetwork(value *GCPNetworkBuilder) *ClusterBuilder {
	b.gcpNetwork = value
	if value != nil {
		b.bitmap_ |= 2048
	} else {
		b.bitmap_ &^= 2048
	}
	return b
}

// AdditionalTrustBundle sets the value of the 'additional_trust_bundle' attribute to the given value.
func (b *ClusterBuilder) AdditionalTrustBundle(value string) *ClusterBuilder {
	b.additionalTrustBundle = value
	b.bitmap_ |= 4096
	return b
}

// Addons sets the value of the 'addons' attribute to the given values.
func (b *ClusterBuilder) Addons(value *AddOnInstallationListBuilder) *ClusterBuilder {
	b.addons = value
	b.bitmap_ |= 8192
	return b
}

// Autoscaler sets the value of the 'autoscaler' attribute to the given value.
//
// Cluster-wide autoscaling configuration.
func (b *ClusterBuilder) Autoscaler(value *ClusterAutoscalerBuilder) *ClusterBuilder {
	b.autoscaler = value
	if value != nil {
		b.bitmap_ |= 16384
	} else {
		b.bitmap_ &^= 16384
	}
	return b
}

// Azure sets the value of the 'azure' attribute to the given value.
//
// Microsoft Azure settings of a cluster.
func (b *ClusterBuilder) Azure(value *AzureBuilder) *ClusterBuilder {
	b.azure = value
	if value != nil {
		b.bitmap_ |= 32768
	} else {
		b.bitmap_ &^= 32768
	}
	return b
}

// BillingModel sets the value of the 'billing_model' attribute to the given value.
//
// Billing model for cluster resources.
func (b *ClusterBuilder) BillingModel(value BillingModel) *ClusterBuilder {
	b.billingModel = value
	b.bitmap_ |= 65536
	return b
}

// ByoOidc sets the value of the 'byo_oidc' attribute to the given value.
//
// ByoOidc configuration.
func (b *ClusterBuilder) ByoOidc(value *ByoOidcBuilder) *ClusterBuilder {
	b.byoOidc = value
	if value != nil {
		b.bitmap_ |= 131072
	} else {
		b.bitmap_ &^= 131072
	}
	return b
}

// CloudProvider sets the value of the 'cloud_provider' attribute to the given value.
//
// Cloud provider.
func (b *ClusterBuilder) CloudProvider(value *CloudProviderBuilder) *ClusterBuilder {
	b.cloudProvider = value
	if value != nil {
		b.bitmap_ |= 262144
	} else {
		b.bitmap_ &^= 262144
	}
	return b
}

// Console sets the value of the 'console' attribute to the given value.
//
// Information about the console of a cluster.
func (b *ClusterBuilder) Console(value *ClusterConsoleBuilder) *ClusterBuilder {
	b.console = value
	if value != nil {
		b.bitmap_ |= 524288
	} else {
		b.bitmap_ &^= 524288
	}
	return b
}

// CreationTimestamp sets the value of the 'creation_timestamp' attribute to the given value.
func (b *ClusterBuilder) CreationTimestamp(value time.Time) *ClusterBuilder {
	b.creationTimestamp = value
	b.bitmap_ |= 1048576
	return b
}

// DeleteProtection sets the value of the 'delete_protection' attribute to the given value.
//
// DeleteProtection configuration.
func (b *ClusterBuilder) DeleteProtection(value *DeleteProtectionBuilder) *ClusterBuilder {
	b.deleteProtection = value
	if value != nil {
		b.bitmap_ |= 2097152
	} else {
		b.bitmap_ &^= 2097152
	}
	return b
}

// DisableUserWorkloadMonitoring sets the value of the 'disable_user_workload_monitoring' attribute to the given value.
func (b *ClusterBuilder) DisableUserWorkloadMonitoring(value bool) *ClusterBuilder {
	b.disableUserWorkloadMonitoring = value
	b.bitmap_ |= 4194304
	return b
}

// DomainPrefix sets the value of the 'domain_prefix' attribute to the given value.
func (b *ClusterBuilder) DomainPrefix(value string) *ClusterBuilder {
	b.domainPrefix = value
	b.bitmap_ |= 8388608
	return b
}

// EtcdEncryption sets the value of the 'etcd_encryption' attribute to the given value.
func (b *ClusterBuilder) EtcdEncryption(value bool) *ClusterBuilder {
	b.etcdEncryption = value
	b.bitmap_ |= 16777216
	return b
}

// ExpirationTimestamp sets the value of the 'expiration_timestamp' attribute to the given value.
func (b *ClusterBuilder) ExpirationTimestamp(value time.Time) *ClusterBuilder {
	b.expirationTimestamp = value
	b.bitmap_ |= 33554432
	return b
}

// ExternalID sets the value of the 'external_ID' attribute to the given value.
func (b *ClusterBuilder) ExternalID(value string) *ClusterBuilder {
	b.externalID = value
	b.bitmap_ |= 67108864
	return b
}

// ExternalAuthConfig sets the value of the 'external_auth_config' attribute to the given value.
//
// ExternalAuthConfig configuration
func (b *ClusterBuilder) ExternalAuthConfig(value *ExternalAuthConfigBuilder) *ClusterBuilder {
	b.externalAuthConfig = value
	if value != nil {
		b.bitmap_ |= 134217728
	} else {
		b.bitmap_ &^= 134217728
	}
	return b
}

// ExternalConfiguration sets the value of the 'external_configuration' attribute to the given value.
//
// Representation of cluster external configuration.
func (b *ClusterBuilder) ExternalConfiguration(value *ExternalConfigurationBuilder) *ClusterBuilder {
	b.externalConfiguration = value
	if value != nil {
		b.bitmap_ |= 268435456
	} else {
		b.bitmap_ &^= 268435456
	}
	return b
}

// Flavour sets the value of the 'flavour' attribute to the given value.
//
// Set of predefined properties of a cluster. For example, a _huge_ flavour can be a cluster
// with 10 infra nodes and 1000 compute nodes.
func (b *ClusterBuilder) Flavour(value *FlavourBuilder) *ClusterBuilder {
	b.flavour = value
	if value != nil {
		b.bitmap_ |= 536870912
	} else {
		b.bitmap_ &^= 536870912
	}
	return b
}

// Groups sets the value of the 'groups' attribute to the given values.
func (b *ClusterBuilder) Groups(value *GroupListBuilder) *ClusterBuilder {
	b.groups = value
	b.bitmap_ |= 1073741824
	return b
}

// HealthState sets the value of the 'health_state' attribute to the given value.
//
// ClusterHealthState indicates the health of a cluster.
func (b *ClusterBuilder) HealthState(value ClusterHealthState) *ClusterBuilder {
	b.healthState = value
	b.bitmap_ |= 2147483648
	return b
}

// Htpasswd sets the value of the 'htpasswd' attribute to the given value.
//
// Details for `htpasswd` identity providers.
func (b *ClusterBuilder) Htpasswd(value *HTPasswdIdentityProviderBuilder) *ClusterBuilder {
	b.htpasswd = value
	if value != nil {
		b.bitmap_ |= 4294967296
	} else {
		b.bitmap_ &^= 4294967296
	}
	return b
}

// Hypershift sets the value of the 'hypershift' attribute to the given value.
//
// Hypershift configuration.
func (b *ClusterBuilder) Hypershift(value *HypershiftBuilder) *ClusterBuilder {
	b.hypershift = value
	if value != nil {
		b.bitmap_ |= 8589934592
	} else {
		b.bitmap_ &^= 8589934592
	}
	return b
}

// IdentityProviders sets the value of the 'identity_providers' attribute to the given values.
func (b *ClusterBuilder) IdentityProviders(value *IdentityProviderListBuilder) *ClusterBuilder {
	b.identityProviders = value
	b.bitmap_ |= 17179869184
	return b
}

// InflightChecks sets the value of the 'inflight_checks' attribute to the given values.
func (b *ClusterBuilder) InflightChecks(value *InflightCheckListBuilder) *ClusterBuilder {
	b.inflightChecks = value
	b.bitmap_ |= 34359738368
	return b
}

// InfraID sets the value of the 'infra_ID' attribute to the given value.
func (b *ClusterBuilder) InfraID(value string) *ClusterBuilder {
	b.infraID = value
	b.bitmap_ |= 68719476736
	return b
}

// Ingresses sets the value of the 'ingresses' attribute to the given values.
func (b *ClusterBuilder) Ingresses(value *IngressListBuilder) *ClusterBuilder {
	b.ingresses = value
	b.bitmap_ |= 137438953472
	return b
}

// KubeletConfig sets the value of the 'kubelet_config' attribute to the given value.
//
// OCM representation of KubeletConfig, exposing the fields of Kubernetes
// KubeletConfig that can be managed by users
func (b *ClusterBuilder) KubeletConfig(value *KubeletConfigBuilder) *ClusterBuilder {
	b.kubeletConfig = value
	if value != nil {
		b.bitmap_ |= 274877906944
	} else {
		b.bitmap_ &^= 274877906944
	}
	return b
}

// LoadBalancerQuota sets the value of the 'load_balancer_quota' attribute to the given value.
func (b *ClusterBuilder) LoadBalancerQuota(value int) *ClusterBuilder {
	b.loadBalancerQuota = value
	b.bitmap_ |= 549755813888
	return b
}

// MachinePools sets the value of the 'machine_pools' attribute to the given values.
func (b *ClusterBuilder) MachinePools(value *MachinePoolListBuilder) *ClusterBuilder {
	b.machinePools = value
	b.bitmap_ |= 1099511627776
	return b
}

// Managed sets the value of the 'managed' attribute to the given value.
func (b *ClusterBuilder) Managed(value bool) *ClusterBuilder {
	b.managed = value
	b.bitmap_ |= 2199023255552
	return b
}

// ManagedService sets the value of the 'managed_service' attribute to the given value.
//
// Contains the necessary attributes to support role-based authentication on AWS.
func (b *ClusterBuilder) ManagedService(value *ManagedServiceBuilder) *ClusterBuilder {
	b.managedService = value
	if value != nil {
		b.bitmap_ |= 4398046511104
	} else {
		b.bitmap_ &^= 4398046511104
	}
	return b
}

// MultiAZ sets the value of the 'multi_AZ' attribute to the given value.
func (b *ClusterBuilder) MultiAZ(value bool) *ClusterBuilder {
	b.multiAZ = value
	b.bitmap_ |= 8796093022208
	return b
}

// MultiArchEnabled sets the value of the 'multi_arch_enabled' attribute to the given value.
func (b *ClusterBuilder) MultiArchEnabled(value bool) *ClusterBuilder {
	b.multiArchEnabled = value
	b.bitmap_ |= 17592186044416
	return b
}

// Name sets the value of the 'name' attribute to the given value.
func (b *ClusterBuilder) Name(value string) *ClusterBuilder {
	b.name = value
	b.bitmap_ |= 35184372088832
	return b
}

// Network sets the value of the 'network' attribute to the given value.
//
// Network configuration of a cluster.
func (b *ClusterBuilder) Network(value *NetworkBuilder) *ClusterBuilder {
	b.network = value
	if value != nil {
		b.bitmap_ |= 70368744177664
	} else {
		b.bitmap_ &^= 70368744177664
	}
	return b
}

// NodeDrainGracePeriod sets the value of the 'node_drain_grace_period' attribute to the given value.
//
// Numeric value and the unit used to measure it.
//
// Units are not mandatory, and they're not specified for some resources. For
// resources that use bytes, the accepted units are:
//
// - 1 B = 1 byte
// - 1 KB = 10^3 bytes
// - 1 MB = 10^6 bytes
// - 1 GB = 10^9 bytes
// - 1 TB = 10^12 bytes
// - 1 PB = 10^15 bytes
//
// - 1 B = 1 byte
// - 1 KiB = 2^10 bytes
// - 1 MiB = 2^20 bytes
// - 1 GiB = 2^30 bytes
// - 1 TiB = 2^40 bytes
// - 1 PiB = 2^50 bytes
func (b *ClusterBuilder) NodeDrainGracePeriod(value *ValueBuilder) *ClusterBuilder {
	b.nodeDrainGracePeriod = value
	if value != nil {
		b.bitmap_ |= 140737488355328
	} else {
		b.bitmap_ &^= 140737488355328
	}
	return b
}

// NodePools sets the value of the 'node_pools' attribute to the given values.
func (b *ClusterBuilder) NodePools(value *NodePoolListBuilder) *ClusterBuilder {
	b.nodePools = value
	b.bitmap_ |= 281474976710656
	return b
}

// Nodes sets the value of the 'nodes' attribute to the given value.
//
// Counts of different classes of nodes inside a cluster.
func (b *ClusterBuilder) Nodes(value *ClusterNodesBuilder) *ClusterBuilder {
	b.nodes = value
	if value != nil {
		b.bitmap_ |= 562949953421312
	} else {
		b.bitmap_ &^= 562949953421312
	}
	return b
}

// OpenshiftVersion sets the value of the 'openshift_version' attribute to the given value.
func (b *ClusterBuilder) OpenshiftVersion(value string) *ClusterBuilder {
	b.openshiftVersion = value
	b.bitmap_ |= 1125899906842624
	return b
}

// Product sets the value of the 'product' attribute to the given value.
//
// Representation of an product that can be selected as a cluster type.
func (b *ClusterBuilder) Product(value *ProductBuilder) *ClusterBuilder {
	b.product = value
	if value != nil {
		b.bitmap_ |= 2251799813685248
	} else {
		b.bitmap_ &^= 2251799813685248
	}
	return b
}

// Properties sets the value of the 'properties' attribute to the given value.
func (b *ClusterBuilder) Properties(value map[string]string) *ClusterBuilder {
	b.properties = value
	if value != nil {
		b.bitmap_ |= 4503599627370496
	} else {
		b.bitmap_ &^= 4503599627370496
	}
	return b
}

// ProvisionShard sets the value of the 'provision_shard' attribute to the given value.
//
// Contains the properties of the provision shard, including AWS and GCP related configurations
func (b *ClusterBuilder) ProvisionShard(value *ProvisionShardBuilder) *ClusterBuilder {
	b.provisionShard = value
	if value != nil {
		b.bitmap_ |= 9007199254740992
	} else {
		b.bitmap_ &^= 9007199254740992
	}
	return b
}

// Proxy sets the value of the 'proxy' attribute to the given value.
//
// Proxy configuration of a cluster.
func (b *ClusterBuilder) Proxy(value *ProxyBuilder) *ClusterBuilder {
	b.proxy = value
	if value != nil {
		b.bitmap_ |= 18014398509481984
	} else {
		b.bitmap_ &^= 18014398509481984
	}
	return b
}

// Region sets the value of the 'region' attribute to the given value.
//
// Description of a region of a cloud provider.
func (b *ClusterBuilder) Region(value *CloudRegionBuilder) *ClusterBuilder {
	b.region = value
	if value != nil {
		b.bitmap_ |= 36028797018963968
	} else {
		b.bitmap_ &^= 36028797018963968
	}
	return b
}

// RegistryConfig sets the value of the 'registry_config' attribute to the given value.
//
// ClusterRegistryConfig describes the configuration of registries for the cluster.
// Its format reflects the OpenShift Image Configuration, for which docs are available on
// [docs.openshift.com](https://docs.openshift.com/container-platform/4.16/openshift_images/image-configuration.html)
// ```json
//
//	{
//	   "registry_config": {
//	     "registry_sources": {
//	       "blocked_registries": [
//	         "badregistry.io",
//	         "badregistry8.io"
//	       ]
//	     }
//	   }
//	}
//
// ```
func (b *ClusterBuilder) RegistryConfig(value *ClusterRegistryConfigBuilder) *ClusterBuilder {
	b.registryConfig = value
	if value != nil {
		b.bitmap_ |= 72057594037927936
	} else {
		b.bitmap_ &^= 72057594037927936
	}
	return b
}

// State sets the value of the 'state' attribute to the given value.
//
// Overall state of a cluster.
func (b *ClusterBuilder) State(value ClusterState) *ClusterBuilder {
	b.state = value
	b.bitmap_ |= 144115188075855872
	return b
}

// Status sets the value of the 'status' attribute to the given value.
//
// Detailed status of a cluster.
func (b *ClusterBuilder) Status(value *ClusterStatusBuilder) *ClusterBuilder {
	b.status = value
	if value != nil {
		b.bitmap_ |= 288230376151711744
	} else {
		b.bitmap_ &^= 288230376151711744
	}
	return b
}

// StorageQuota sets the value of the 'storage_quota' attribute to the given value.
//
// Numeric value and the unit used to measure it.
//
// Units are not mandatory, and they're not specified for some resources. For
// resources that use bytes, the accepted units are:
//
// - 1 B = 1 byte
// - 1 KB = 10^3 bytes
// - 1 MB = 10^6 bytes
// - 1 GB = 10^9 bytes
// - 1 TB = 10^12 bytes
// - 1 PB = 10^15 bytes
//
// - 1 B = 1 byte
// - 1 KiB = 2^10 bytes
// - 1 MiB = 2^20 bytes
// - 1 GiB = 2^30 bytes
// - 1 TiB = 2^40 bytes
// - 1 PiB = 2^50 bytes
func (b *ClusterBuilder) StorageQuota(value *ValueBuilder) *ClusterBuilder {
	b.storageQuota = value
	if value != nil {
		b.bitmap_ |= 576460752303423488
	} else {
		b.bitmap_ &^= 576460752303423488
	}
	return b
}

// Subscription sets the value of the 'subscription' attribute to the given value.
//
// Definition of a subscription.
func (b *ClusterBuilder) Subscription(value *SubscriptionBuilder) *ClusterBuilder {
	b.subscription = value
	if value != nil {
		b.bitmap_ |= 1152921504606846976
	} else {
		b.bitmap_ &^= 1152921504606846976
	}
	return b
}

// Version sets the value of the 'version' attribute to the given value.
//
// Representation of an _OpenShift_ version.
func (b *ClusterBuilder) Version(value *VersionBuilder) *ClusterBuilder {
	b.version = value
	if value != nil {
		b.bitmap_ |= 2305843009213693952
	} else {
		b.bitmap_ &^= 2305843009213693952
	}
	return b
}

// Copy copies the attributes of the given object into this builder, discarding any previous values.
func (b *ClusterBuilder) Copy(object *Cluster) *ClusterBuilder {
	if object == nil {
		return b
	}
	b.bitmap_ = object.bitmap_
	b.id = object.id
	b.href = object.href
	if object.api != nil {
		b.api = NewClusterAPI().Copy(object.api)
	} else {
		b.api = nil
	}
	if object.aws != nil {
		b.aws = NewAWS().Copy(object.aws)
	} else {
		b.aws = nil
	}
	if object.awsInfrastructureAccessRoleGrants != nil {
		b.awsInfrastructureAccessRoleGrants = NewAWSInfrastructureAccessRoleGrantList().Copy(object.awsInfrastructureAccessRoleGrants)
	} else {
		b.awsInfrastructureAccessRoleGrants = nil
	}
	if object.ccs != nil {
		b.ccs = NewCCS().Copy(object.ccs)
	} else {
		b.ccs = nil
	}
	if object.dns != nil {
		b.dns = NewDNS().Copy(object.dns)
	} else {
		b.dns = nil
	}
	b.fips = object.fips
	if object.gcp != nil {
		b.gcp = NewGCP().Copy(object.gcp)
	} else {
		b.gcp = nil
	}
	if object.gcpEncryptionKey != nil {
		b.gcpEncryptionKey = NewGCPEncryptionKey().Copy(object.gcpEncryptionKey)
	} else {
		b.gcpEncryptionKey = nil
	}
	if object.gcpNetwork != nil {
		b.gcpNetwork = NewGCPNetwork().Copy(object.gcpNetwork)
	} else {
		b.gcpNetwork = nil
	}
	b.additionalTrustBundle = object.additionalTrustBundle
	if object.addons != nil {
		b.addons = NewAddOnInstallationList().Copy(object.addons)
	} else {
		b.addons = nil
	}
	if object.autoscaler != nil {
		b.autoscaler = NewClusterAutoscaler().Copy(object.autoscaler)
	} else {
		b.autoscaler = nil
	}
	if object.azure != nil {
		b.azure = NewAzure().Copy(object.azure)
	} else {
		b.azure = nil
	}
	b.billingModel = object.billingModel
	if object.byoOidc != nil {
		b.byoOidc = NewByoOidc().Copy(object.byoOidc)
	} else {
		b.byoOidc = nil
	}
	if object.cloudProvider != nil {
		b.cloudProvider = NewCloudProvider().Copy(object.cloudProvider)
	} else {
		b.cloudProvider = nil
	}
	if object.console != nil {
		b.console = NewClusterConsole().Copy(object.console)
	} else {
		b.console = nil
	}
	b.creationTimestamp = object.creationTimestamp
	if object.deleteProtection != nil {
		b.deleteProtection = NewDeleteProtection().Copy(object.deleteProtection)
	} else {
		b.deleteProtection = nil
	}
	b.disableUserWorkloadMonitoring = object.disableUserWorkloadMonitoring
	b.domainPrefix = object.domainPrefix
	b.etcdEncryption = object.etcdEncryption
	b.expirationTimestamp = object.expirationTimestamp
	b.externalID = object.externalID
	if object.externalAuthConfig != nil {
		b.externalAuthConfig = NewExternalAuthConfig().Copy(object.externalAuthConfig)
	} else {
		b.externalAuthConfig = nil
	}
	if object.externalConfiguration != nil {
		b.externalConfiguration = NewExternalConfiguration().Copy(object.externalConfiguration)
	} else {
		b.externalConfiguration = nil
	}
	if object.flavour != nil {
		b.flavour = NewFlavour().Copy(object.flavour)
	} else {
		b.flavour = nil
	}
	if object.groups != nil {
		b.groups = NewGroupList().Copy(object.groups)
	} else {
		b.groups = nil
	}
	b.healthState = object.healthState
	if object.htpasswd != nil {
		b.htpasswd = NewHTPasswdIdentityProvider().Copy(object.htpasswd)
	} else {
		b.htpasswd = nil
	}
	if object.hypershift != nil {
		b.hypershift = NewHypershift().Copy(object.hypershift)
	} else {
		b.hypershift = nil
	}
	if object.identityProviders != nil {
		b.identityProviders = NewIdentityProviderList().Copy(object.identityProviders)
	} else {
		b.identityProviders = nil
	}
	if object.inflightChecks != nil {
		b.inflightChecks = NewInflightCheckList().Copy(object.inflightChecks)
	} else {
		b.inflightChecks = nil
	}
	b.infraID = object.infraID
	if object.ingresses != nil {
		b.ingresses = NewIngressList().Copy(object.ingresses)
	} else {
		b.ingresses = nil
	}
	if object.kubeletConfig != nil {
		b.kubeletConfig = NewKubeletConfig().Copy(object.kubeletConfig)
	} else {
		b.kubeletConfig = nil
	}
	b.loadBalancerQuota = object.loadBalancerQuota
	if object.machinePools != nil {
		b.machinePools = NewMachinePoolList().Copy(object.machinePools)
	} else {
		b.machinePools = nil
	}
	b.managed = object.managed
	if object.managedService != nil {
		b.managedService = NewManagedService().Copy(object.managedService)
	} else {
		b.managedService = nil
	}
	b.multiAZ = object.multiAZ
	b.multiArchEnabled = object.multiArchEnabled
	b.name = object.name
	if object.network != nil {
		b.network = NewNetwork().Copy(object.network)
	} else {
		b.network = nil
	}
	if object.nodeDrainGracePeriod != nil {
		b.nodeDrainGracePeriod = NewValue().Copy(object.nodeDrainGracePeriod)
	} else {
		b.nodeDrainGracePeriod = nil
	}
	if object.nodePools != nil {
		b.nodePools = NewNodePoolList().Copy(object.nodePools)
	} else {
		b.nodePools = nil
	}
	if object.nodes != nil {
		b.nodes = NewClusterNodes().Copy(object.nodes)
	} else {
		b.nodes = nil
	}
	b.openshiftVersion = object.openshiftVersion
	if object.product != nil {
		b.product = NewProduct().Copy(object.product)
	} else {
		b.product = nil
	}
	if len(object.properties) > 0 {
		b.properties = map[string]string{}
		for k, v := range object.properties {
			b.properties[k] = v
		}
	} else {
		b.properties = nil
	}
	if object.provisionShard != nil {
		b.provisionShard = NewProvisionShard().Copy(object.provisionShard)
	} else {
		b.provisionShard = nil
	}
	if object.proxy != nil {
		b.proxy = NewProxy().Copy(object.proxy)
	} else {
		b.proxy = nil
	}
	if object.region != nil {
		b.region = NewCloudRegion().Copy(object.region)
	} else {
		b.region = nil
	}
	if object.registryConfig != nil {
		b.registryConfig = NewClusterRegistryConfig().Copy(object.registryConfig)
	} else {
		b.registryConfig = nil
	}
	b.state = object.state
	if object.status != nil {
		b.status = NewClusterStatus().Copy(object.status)
	} else {
		b.status = nil
	}
	if object.storageQuota != nil {
		b.storageQuota = NewValue().Copy(object.storageQuota)
	} else {
		b.storageQuota = nil
	}
	if object.subscription != nil {
		b.subscription = NewSubscription().Copy(object.subscription)
	} else {
		b.subscription = nil
	}
	if object.version != nil {
		b.version = NewVersion().Copy(object.version)
	} else {
		b.version = nil
	}
	return b
}

// Build creates a 'cluster' object using the configuration stored in the builder.
func (b *ClusterBuilder) Build() (object *Cluster, err error) {
	object = new(Cluster)
	object.id = b.id
	object.href = b.href
	object.bitmap_ = b.bitmap_
	if b.api != nil {
		object.api, err = b.api.Build()
		if err != nil {
			return
		}
	}
	if b.aws != nil {
		object.aws, err = b.aws.Build()
		if err != nil {
			return
		}
	}
	if b.awsInfrastructureAccessRoleGrants != nil {
		object.awsInfrastructureAccessRoleGrants, err = b.awsInfrastructureAccessRoleGrants.Build()
		if err != nil {
			return
		}
	}
	if b.ccs != nil {
		object.ccs, err = b.ccs.Build()
		if err != nil {
			return
		}
	}
	if b.dns != nil {
		object.dns, err = b.dns.Build()
		if err != nil {
			return
		}
	}
	object.fips = b.fips
	if b.gcp != nil {
		object.gcp, err = b.gcp.Build()
		if err != nil {
			return
		}
	}
	if b.gcpEncryptionKey != nil {
		object.gcpEncryptionKey, err = b.gcpEncryptionKey.Build()
		if err != nil {
			return
		}
	}
	if b.gcpNetwork != nil {
		object.gcpNetwork, err = b.gcpNetwork.Build()
		if err != nil {
			return
		}
	}
	object.additionalTrustBundle = b.additionalTrustBundle
	if b.addons != nil {
		object.addons, err = b.addons.Build()
		if err != nil {
			return
		}
	}
	if b.autoscaler != nil {
		object.autoscaler, err = b.autoscaler.Build()
		if err != nil {
			return
		}
	}
	if b.azure != nil {
		object.azure, err = b.azure.Build()
		if err != nil {
			return
		}
	}
	object.billingModel = b.billingModel
	if b.byoOidc != nil {
		object.byoOidc, err = b.byoOidc.Build()
		if err != nil {
			return
		}
	}
	if b.cloudProvider != nil {
		object.cloudProvider, err = b.cloudProvider.Build()
		if err != nil {
			return
		}
	}
	if b.console != nil {
		object.console, err = b.console.Build()
		if err != nil {
			return
		}
	}
	object.creationTimestamp = b.creationTimestamp
	if b.deleteProtection != nil {
		object.deleteProtection, err = b.deleteProtection.Build()
		if err != nil {
			return
		}
	}
	object.disableUserWorkloadMonitoring = b.disableUserWorkloadMonitoring
	object.domainPrefix = b.domainPrefix
	object.etcdEncryption = b.etcdEncryption
	object.expirationTimestamp = b.expirationTimestamp
	object.externalID = b.externalID
	if b.externalAuthConfig != nil {
		object.externalAuthConfig, err = b.externalAuthConfig.Build()
		if err != nil {
			return
		}
	}
	if b.externalConfiguration != nil {
		object.externalConfiguration, err = b.externalConfiguration.Build()
		if err != nil {
			return
		}
	}
	if b.flavour != nil {
		object.flavour, err = b.flavour.Build()
		if err != nil {
			return
		}
	}
	if b.groups != nil {
		object.groups, err = b.groups.Build()
		if err != nil {
			return
		}
	}
	object.healthState = b.healthState
	if b.htpasswd != nil {
		object.htpasswd, err = b.htpasswd.Build()
		if err != nil {
			return
		}
	}
	if b.hypershift != nil {
		object.hypershift, err = b.hypershift.Build()
		if err != nil {
			return
		}
	}
	if b.identityProviders != nil {
		object.identityProviders, err = b.identityProviders.Build()
		if err != nil {
			return
		}
	}
	if b.inflightChecks != nil {
		object.inflightChecks, err = b.inflightChecks.Build()
		if err != nil {
			return
		}
	}
	object.infraID = b.infraID
	if b.ingresses != nil {
		object.ingresses, err = b.ingresses.Build()
		if err != nil {
			return
		}
	}
	if b.kubeletConfig != nil {
		object.kubeletConfig, err = b.kubeletConfig.Build()
		if err != nil {
			return
		}
	}
	object.loadBalancerQuota = b.loadBalancerQuota
	if b.machinePools != nil {
		object.machinePools, err = b.machinePools.Build()
		if err != nil {
			return
		}
	}
	object.managed = b.managed
	if b.managedService != nil {
		object.managedService, err = b.managedService.Build()
		if err != nil {
			return
		}
	}
	object.multiAZ = b.multiAZ
	object.multiArchEnabled = b.multiArchEnabled
	object.name = b.name
	if b.network != nil {
		object.network, err = b.network.Build()
		if err != nil {
			return
		}
	}
	if b.nodeDrainGracePeriod != nil {
		object.nodeDrainGracePeriod, err = b.nodeDrainGracePeriod.Build()
		if err != nil {
			return
		}
	}
	if b.nodePools != nil {
		object.nodePools, err = b.nodePools.Build()
		if err != nil {
			return
		}
	}
	if b.nodes != nil {
		object.nodes, err = b.nodes.Build()
		if err != nil {
			return
		}
	}
	object.openshiftVersion = b.openshiftVersion
	if b.product != nil {
		object.product, err = b.product.Build()
		if err != nil {
			return
		}
	}
	if b.properties != nil {
		object.properties = make(map[string]string)
		for k, v := range b.properties {
			object.properties[k] = v
		}
	}
	if b.provisionShard != nil {
		object.provisionShard, err = b.provisionShard.Build()
		if err != nil {
			return
		}
	}
	if b.proxy != nil {
		object.proxy, err = b.proxy.Build()
		if err != nil {
			return
		}
	}
	if b.region != nil {
		object.region, err = b.region.Build()
		if err != nil {
			return
		}
	}
	if b.registryConfig != nil {
		object.registryConfig, err = b.registryConfig.Build()
		if err != nil {
			return
		}
	}
	object.state = b.state
	if b.status != nil {
		object.status, err = b.status.Build()
		if err != nil {
			return
		}
	}
	if b.storageQuota != nil {
		object.storageQuota, err = b.storageQuota.Build()
		if err != nil {
			return
		}
	}
	if b.subscription != nil {
		object.subscription, err = b.subscription.Build()
		if err != nil {
			return
		}
	}
	if b.version != nil {
		object.version, err = b.version.Build()
		if err != nil {
			return
		}
	}
	return
}
