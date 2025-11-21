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

import (
	time "time"

	v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
)

// ClusterKind is the name of the type used to represent objects
// of type 'cluster'.
const ClusterKind = "Cluster"

// ClusterLinkKind is the name of the type used to represent links
// to objects of type 'cluster'.
const ClusterLinkKind = "ClusterLink"

// ClusterNilKind is the name of the type used to nil references
// to objects of type 'cluster'.
const ClusterNilKind = "ClusterNil"

// Cluster represents the values of the 'cluster' type.
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
type Cluster struct {
	fieldSet_                         []bool
	id                                string
	href                              string
	api                               *ClusterAPI
	aws                               *AWS
	awsInfrastructureAccessRoleGrants *v1.AWSInfrastructureAccessRoleGrantList
	ccs                               *CCS
	dns                               *DNS
	gcp                               *GCP
	gcpEncryptionKey                  *GCPEncryptionKey
	gcpNetwork                        *GCPNetwork
	additionalTrustBundle             string
	addons                            *v1.AddOnInstallationList
	autoNode                          *ClusterAutoNode
	autoscaler                        *ClusterAutoscaler
	azure                             *Azure
	billingModel                      BillingModel
	byoOidc                           *ByoOidc
	cloudProvider                     *v1.CloudProvider
	console                           *ClusterConsole
	creationTimestamp                 time.Time
	deleteProtection                  *DeleteProtection
	domainPrefix                      string
	expirationTimestamp               time.Time
	externalID                        string
	externalAuthConfig                *ExternalAuthConfig
	externalConfiguration             *ExternalConfiguration
	flavour                           *v1.Flavour
	groups                            *v1.GroupList
	healthState                       ClusterHealthState
	htpasswd                          *HTPasswdIdentityProvider
	hypershift                        *Hypershift
	identityProviders                 *v1.IdentityProviderList
	imageRegistry                     *ClusterImageRegistry
	inflightChecks                    *InflightCheckList
	infraID                           string
	ingresses                         *v1.IngressList
	kubeletConfig                     *KubeletConfig
	loadBalancerQuota                 int
	machinePools                      *v1.MachinePoolList
	managedService                    *ManagedService
	name                              string
	network                           *Network
	nodeDrainGracePeriod              *Value
	nodePools                         *NodePoolList
	nodes                             *ClusterNodes
	openshiftVersion                  string
	product                           *v1.Product
	properties                        map[string]string
	provisionShard                    *ProvisionShard
	proxy                             *Proxy
	region                            *v1.CloudRegion
	registryConfig                    *ClusterRegistryConfig
	state                             ClusterState
	status                            *ClusterStatus
	storageQuota                      *Value
	subscription                      *v1.Subscription
	version                           *Version
	fips                              bool
	disableUserWorkloadMonitoring     bool
	etcdEncryption                    bool
	managed                           bool
	multiAZ                           bool
	multiArchEnabled                  bool
}

// Kind returns the name of the type of the object.
func (o *Cluster) Kind() string {
	if o == nil {
		return ClusterNilKind
	}
	if len(o.fieldSet_) > 0 && o.fieldSet_[0] {
		return ClusterLinkKind
	}
	return ClusterKind
}

// Link returns true if this is a link.
func (o *Cluster) Link() bool {
	return o != nil && len(o.fieldSet_) > 0 && o.fieldSet_[0]
}

// ID returns the identifier of the object.
func (o *Cluster) ID() string {
	if o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1] {
		return o.id
	}
	return ""
}

// GetID returns the identifier of the object and a flag indicating if the
// identifier has a value.
func (o *Cluster) GetID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 1 && o.fieldSet_[1]
	if ok {
		value = o.id
	}
	return
}

// HREF returns the link to the object.
func (o *Cluster) HREF() string {
	if o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2] {
		return o.href
	}
	return ""
}

// GetHREF returns the link of the object and a flag indicating if the
// link has a value.
func (o *Cluster) GetHREF() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 2 && o.fieldSet_[2]
	if ok {
		value = o.href
	}
	return
}

// Empty returns true if the object is empty, i.e. no attribute has a value.
func (o *Cluster) Empty() bool {
	if o == nil || len(o.fieldSet_) == 0 {
		return true
	}

	// Check all fields except the link flag (index 0)
	for i := 1; i < len(o.fieldSet_); i++ {
		if o.fieldSet_[i] {
			return false
		}
	}
	return true
}

// API returns the value of the 'API' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Information about the API of the cluster.
func (o *Cluster) API() *ClusterAPI {
	if o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3] {
		return o.api
	}
	return nil
}

// GetAPI returns the value of the 'API' attribute and
// a flag indicating if the attribute has a value.
//
// Information about the API of the cluster.
func (o *Cluster) GetAPI() (value *ClusterAPI, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 3 && o.fieldSet_[3]
	if ok {
		value = o.api
	}
	return
}

// AWS returns the value of the 'AWS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Amazon Web Services settings of the cluster.
func (o *Cluster) AWS() *AWS {
	if o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4] {
		return o.aws
	}
	return nil
}

// GetAWS returns the value of the 'AWS' attribute and
// a flag indicating if the attribute has a value.
//
// Amazon Web Services settings of the cluster.
func (o *Cluster) GetAWS() (value *AWS, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 4 && o.fieldSet_[4]
	if ok {
		value = o.aws
	}
	return
}

// AWSInfrastructureAccessRoleGrants returns the value of the 'AWS_infrastructure_access_role_grants' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of AWS infrastructure access role grants on this cluster.
func (o *Cluster) AWSInfrastructureAccessRoleGrants() *v1.AWSInfrastructureAccessRoleGrantList {
	if o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5] {
		return o.awsInfrastructureAccessRoleGrants
	}
	return nil
}

// GetAWSInfrastructureAccessRoleGrants returns the value of the 'AWS_infrastructure_access_role_grants' attribute and
// a flag indicating if the attribute has a value.
//
// List of AWS infrastructure access role grants on this cluster.
func (o *Cluster) GetAWSInfrastructureAccessRoleGrants() (value *v1.AWSInfrastructureAccessRoleGrantList, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 5 && o.fieldSet_[5]
	if ok {
		value = o.awsInfrastructureAccessRoleGrants
	}
	return
}

// CCS returns the value of the 'CCS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains configuration of a Customer Cloud Subscription cluster.
func (o *Cluster) CCS() *CCS {
	if o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6] {
		return o.ccs
	}
	return nil
}

// GetCCS returns the value of the 'CCS' attribute and
// a flag indicating if the attribute has a value.
//
// Contains configuration of a Customer Cloud Subscription cluster.
func (o *Cluster) GetCCS() (value *CCS, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 6 && o.fieldSet_[6]
	if ok {
		value = o.ccs
	}
	return
}

// DNS returns the value of the 'DNS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// DNS settings of the cluster.
func (o *Cluster) DNS() *DNS {
	if o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7] {
		return o.dns
	}
	return nil
}

// GetDNS returns the value of the 'DNS' attribute and
// a flag indicating if the attribute has a value.
//
// DNS settings of the cluster.
func (o *Cluster) GetDNS() (value *DNS, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 7 && o.fieldSet_[7]
	if ok {
		value = o.dns
	}
	return
}

// FIPS returns the value of the 'FIPS' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Create cluster that uses FIPS Validated / Modules in Process cryptographic libraries.
func (o *Cluster) FIPS() bool {
	if o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8] {
		return o.fips
	}
	return false
}

// GetFIPS returns the value of the 'FIPS' attribute and
// a flag indicating if the attribute has a value.
//
// Create cluster that uses FIPS Validated / Modules in Process cryptographic libraries.
func (o *Cluster) GetFIPS() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 8 && o.fieldSet_[8]
	if ok {
		value = o.fips
	}
	return
}

// GCP returns the value of the 'GCP' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Google cloud platform settings of the cluster.
func (o *Cluster) GCP() *GCP {
	if o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9] {
		return o.gcp
	}
	return nil
}

// GetGCP returns the value of the 'GCP' attribute and
// a flag indicating if the attribute has a value.
//
// Google cloud platform settings of the cluster.
func (o *Cluster) GetGCP() (value *GCP, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 9 && o.fieldSet_[9]
	if ok {
		value = o.gcp
	}
	return
}

// GCPEncryptionKey returns the value of the 'GCP_encryption_key' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Key used for encryption of GCP cluster nodes.
func (o *Cluster) GCPEncryptionKey() *GCPEncryptionKey {
	if o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10] {
		return o.gcpEncryptionKey
	}
	return nil
}

// GetGCPEncryptionKey returns the value of the 'GCP_encryption_key' attribute and
// a flag indicating if the attribute has a value.
//
// Key used for encryption of GCP cluster nodes.
func (o *Cluster) GetGCPEncryptionKey() (value *GCPEncryptionKey, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 10 && o.fieldSet_[10]
	if ok {
		value = o.gcpEncryptionKey
	}
	return
}

// GCPNetwork returns the value of the 'GCP_network' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// GCP Network.
func (o *Cluster) GCPNetwork() *GCPNetwork {
	if o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11] {
		return o.gcpNetwork
	}
	return nil
}

// GetGCPNetwork returns the value of the 'GCP_network' attribute and
// a flag indicating if the attribute has a value.
//
// GCP Network.
func (o *Cluster) GetGCPNetwork() (value *GCPNetwork, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 11 && o.fieldSet_[11]
	if ok {
		value = o.gcpNetwork
	}
	return
}

// AdditionalTrustBundle returns the value of the 'additional_trust_bundle' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Additional trust bundle.
func (o *Cluster) AdditionalTrustBundle() string {
	if o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12] {
		return o.additionalTrustBundle
	}
	return ""
}

// GetAdditionalTrustBundle returns the value of the 'additional_trust_bundle' attribute and
// a flag indicating if the attribute has a value.
//
// Additional trust bundle.
func (o *Cluster) GetAdditionalTrustBundle() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 12 && o.fieldSet_[12]
	if ok {
		value = o.additionalTrustBundle
	}
	return
}

// Addons returns the value of the 'addons' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of add-ons on this cluster.
func (o *Cluster) Addons() *v1.AddOnInstallationList {
	if o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13] {
		return o.addons
	}
	return nil
}

// GetAddons returns the value of the 'addons' attribute and
// a flag indicating if the attribute has a value.
//
// List of add-ons on this cluster.
func (o *Cluster) GetAddons() (value *v1.AddOnInstallationList, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 13 && o.fieldSet_[13]
	if ok {
		value = o.addons
	}
	return
}

// AutoNode returns the value of the 'auto_node' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The AutoNode settings for this cluster.
// This is currently only supported for ROSA HCP
func (o *Cluster) AutoNode() *ClusterAutoNode {
	if o != nil && len(o.fieldSet_) > 14 && o.fieldSet_[14] {
		return o.autoNode
	}
	return nil
}

// GetAutoNode returns the value of the 'auto_node' attribute and
// a flag indicating if the attribute has a value.
//
// The AutoNode settings for this cluster.
// This is currently only supported for ROSA HCP
func (o *Cluster) GetAutoNode() (value *ClusterAutoNode, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 14 && o.fieldSet_[14]
	if ok {
		value = o.autoNode
	}
	return
}

// Autoscaler returns the value of the 'autoscaler' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to an optional _ClusterAutoscaler_ that is coupled with the cluster.
func (o *Cluster) Autoscaler() *ClusterAutoscaler {
	if o != nil && len(o.fieldSet_) > 15 && o.fieldSet_[15] {
		return o.autoscaler
	}
	return nil
}

// GetAutoscaler returns the value of the 'autoscaler' attribute and
// a flag indicating if the attribute has a value.
//
// Link to an optional _ClusterAutoscaler_ that is coupled with the cluster.
func (o *Cluster) GetAutoscaler() (value *ClusterAutoscaler, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 15 && o.fieldSet_[15]
	if ok {
		value = o.autoscaler
	}
	return
}

// Azure returns the value of the 'azure' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Microsoft Azure settings of the cluster.
func (o *Cluster) Azure() *Azure {
	if o != nil && len(o.fieldSet_) > 16 && o.fieldSet_[16] {
		return o.azure
	}
	return nil
}

// GetAzure returns the value of the 'azure' attribute and
// a flag indicating if the attribute has a value.
//
// Microsoft Azure settings of the cluster.
func (o *Cluster) GetAzure() (value *Azure, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 16 && o.fieldSet_[16]
	if ok {
		value = o.azure
	}
	return
}

// BillingModel returns the value of the 'billing_model' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Billing model for cluster resources.
func (o *Cluster) BillingModel() BillingModel {
	if o != nil && len(o.fieldSet_) > 17 && o.fieldSet_[17] {
		return o.billingModel
	}
	return BillingModel("")
}

// GetBillingModel returns the value of the 'billing_model' attribute and
// a flag indicating if the attribute has a value.
//
// Billing model for cluster resources.
func (o *Cluster) GetBillingModel() (value BillingModel, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 17 && o.fieldSet_[17]
	if ok {
		value = o.billingModel
	}
	return
}

// ByoOidc returns the value of the 'byo_oidc' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains information about BYO OIDC.
func (o *Cluster) ByoOidc() *ByoOidc {
	if o != nil && len(o.fieldSet_) > 18 && o.fieldSet_[18] {
		return o.byoOidc
	}
	return nil
}

// GetByoOidc returns the value of the 'byo_oidc' attribute and
// a flag indicating if the attribute has a value.
//
// Contains information about BYO OIDC.
func (o *Cluster) GetByoOidc() (value *ByoOidc, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 18 && o.fieldSet_[18]
	if ok {
		value = o.byoOidc
	}
	return
}

// CloudProvider returns the value of the 'cloud_provider' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the cloud provider where the cluster is installed.
func (o *Cluster) CloudProvider() *v1.CloudProvider {
	if o != nil && len(o.fieldSet_) > 19 && o.fieldSet_[19] {
		return o.cloudProvider
	}
	return nil
}

// GetCloudProvider returns the value of the 'cloud_provider' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the cloud provider where the cluster is installed.
func (o *Cluster) GetCloudProvider() (value *v1.CloudProvider, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 19 && o.fieldSet_[19]
	if ok {
		value = o.cloudProvider
	}
	return
}

// Console returns the value of the 'console' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Information about the console of the cluster.
func (o *Cluster) Console() *ClusterConsole {
	if o != nil && len(o.fieldSet_) > 20 && o.fieldSet_[20] {
		return o.console
	}
	return nil
}

// GetConsole returns the value of the 'console' attribute and
// a flag indicating if the attribute has a value.
//
// Information about the console of the cluster.
func (o *Cluster) GetConsole() (value *ClusterConsole, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 20 && o.fieldSet_[20]
	if ok {
		value = o.console
	}
	return
}

// CreationTimestamp returns the value of the 'creation_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the cluster was initially created, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *Cluster) CreationTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 21 && o.fieldSet_[21] {
		return o.creationTimestamp
	}
	return time.Time{}
}

// GetCreationTimestamp returns the value of the 'creation_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the cluster was initially created, using the
// format defined in [RFC3339](https://www.ietf.org/rfc/rfc3339.txt).
func (o *Cluster) GetCreationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 21 && o.fieldSet_[21]
	if ok {
		value = o.creationTimestamp
	}
	return
}

// DeleteProtection returns the value of the 'delete_protection' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Delete protection
func (o *Cluster) DeleteProtection() *DeleteProtection {
	if o != nil && len(o.fieldSet_) > 22 && o.fieldSet_[22] {
		return o.deleteProtection
	}
	return nil
}

// GetDeleteProtection returns the value of the 'delete_protection' attribute and
// a flag indicating if the attribute has a value.
//
// Delete protection
func (o *Cluster) GetDeleteProtection() (value *DeleteProtection, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 22 && o.fieldSet_[22]
	if ok {
		value = o.deleteProtection
	}
	return
}

// DisableUserWorkloadMonitoring returns the value of the 'disable_user_workload_monitoring' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates whether the User workload monitoring is enabled or not
// It is enabled by default
func (o *Cluster) DisableUserWorkloadMonitoring() bool {
	if o != nil && len(o.fieldSet_) > 23 && o.fieldSet_[23] {
		return o.disableUserWorkloadMonitoring
	}
	return false
}

// GetDisableUserWorkloadMonitoring returns the value of the 'disable_user_workload_monitoring' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates whether the User workload monitoring is enabled or not
// It is enabled by default
func (o *Cluster) GetDisableUserWorkloadMonitoring() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 23 && o.fieldSet_[23]
	if ok {
		value = o.disableUserWorkloadMonitoring
	}
	return
}

// DomainPrefix returns the value of the 'domain_prefix' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// DomainPrefix of the cluster. This prefix is optionally assigned by the user when the
// cluster is created. It will appear in the Cluster's domain when the cluster is provisioned.
func (o *Cluster) DomainPrefix() string {
	if o != nil && len(o.fieldSet_) > 24 && o.fieldSet_[24] {
		return o.domainPrefix
	}
	return ""
}

// GetDomainPrefix returns the value of the 'domain_prefix' attribute and
// a flag indicating if the attribute has a value.
//
// DomainPrefix of the cluster. This prefix is optionally assigned by the user when the
// cluster is created. It will appear in the Cluster's domain when the cluster is provisioned.
func (o *Cluster) GetDomainPrefix() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 24 && o.fieldSet_[24]
	if ok {
		value = o.domainPrefix
	}
	return
}

// EtcdEncryption returns the value of the 'etcd_encryption' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicates whether that etcd is encrypted or not.
// This is set only during cluster creation.
// For ARO-HCP Clusters, this is a readonly attribute, always set to true.
func (o *Cluster) EtcdEncryption() bool {
	if o != nil && len(o.fieldSet_) > 25 && o.fieldSet_[25] {
		return o.etcdEncryption
	}
	return false
}

// GetEtcdEncryption returns the value of the 'etcd_encryption' attribute and
// a flag indicating if the attribute has a value.
//
// Indicates whether that etcd is encrypted or not.
// This is set only during cluster creation.
// For ARO-HCP Clusters, this is a readonly attribute, always set to true.
func (o *Cluster) GetEtcdEncryption() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 25 && o.fieldSet_[25]
	if ok {
		value = o.etcdEncryption
	}
	return
}

// ExpirationTimestamp returns the value of the 'expiration_timestamp' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Date and time when the cluster will be automatically deleted, using the format defined in
// [RFC3339](https://www.ietf.org/rfc/rfc3339.txt). If no timestamp is provided, the cluster
// will never expire.
//
// This option is unsupported.
func (o *Cluster) ExpirationTimestamp() time.Time {
	if o != nil && len(o.fieldSet_) > 26 && o.fieldSet_[26] {
		return o.expirationTimestamp
	}
	return time.Time{}
}

// GetExpirationTimestamp returns the value of the 'expiration_timestamp' attribute and
// a flag indicating if the attribute has a value.
//
// Date and time when the cluster will be automatically deleted, using the format defined in
// [RFC3339](https://www.ietf.org/rfc/rfc3339.txt). If no timestamp is provided, the cluster
// will never expire.
//
// This option is unsupported.
func (o *Cluster) GetExpirationTimestamp() (value time.Time, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 26 && o.fieldSet_[26]
	if ok {
		value = o.expirationTimestamp
	}
	return
}

// ExternalID returns the value of the 'external_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// External identifier of the cluster, generated by the installer.
func (o *Cluster) ExternalID() string {
	if o != nil && len(o.fieldSet_) > 27 && o.fieldSet_[27] {
		return o.externalID
	}
	return ""
}

// GetExternalID returns the value of the 'external_ID' attribute and
// a flag indicating if the attribute has a value.
//
// External identifier of the cluster, generated by the installer.
func (o *Cluster) GetExternalID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 27 && o.fieldSet_[27]
	if ok {
		value = o.externalID
	}
	return
}

// ExternalAuthConfig returns the value of the 'external_auth_config' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// External authentication configuration.
//
// For ROSA HCP, if this is not specified, external authentication configuration will be disabled by default
// For ARO HCP, if this is not specified, external authentication configuration will be enabled by default
func (o *Cluster) ExternalAuthConfig() *ExternalAuthConfig {
	if o != nil && len(o.fieldSet_) > 28 && o.fieldSet_[28] {
		return o.externalAuthConfig
	}
	return nil
}

// GetExternalAuthConfig returns the value of the 'external_auth_config' attribute and
// a flag indicating if the attribute has a value.
//
// External authentication configuration.
//
// For ROSA HCP, if this is not specified, external authentication configuration will be disabled by default
// For ARO HCP, if this is not specified, external authentication configuration will be enabled by default
func (o *Cluster) GetExternalAuthConfig() (value *ExternalAuthConfig, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 28 && o.fieldSet_[28]
	if ok {
		value = o.externalAuthConfig
	}
	return
}

// ExternalConfiguration returns the value of the 'external_configuration' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ExternalConfiguration shows external configuration on the cluster.
func (o *Cluster) ExternalConfiguration() *ExternalConfiguration {
	if o != nil && len(o.fieldSet_) > 29 && o.fieldSet_[29] {
		return o.externalConfiguration
	}
	return nil
}

// GetExternalConfiguration returns the value of the 'external_configuration' attribute and
// a flag indicating if the attribute has a value.
//
// ExternalConfiguration shows external configuration on the cluster.
func (o *Cluster) GetExternalConfiguration() (value *ExternalConfiguration, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 29 && o.fieldSet_[29]
	if ok {
		value = o.externalConfiguration
	}
	return
}

// Flavour returns the value of the 'flavour' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the _flavour_ that was used to create the cluster.
func (o *Cluster) Flavour() *v1.Flavour {
	if o != nil && len(o.fieldSet_) > 30 && o.fieldSet_[30] {
		return o.flavour
	}
	return nil
}

// GetFlavour returns the value of the 'flavour' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the _flavour_ that was used to create the cluster.
func (o *Cluster) GetFlavour() (value *v1.Flavour, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 30 && o.fieldSet_[30]
	if ok {
		value = o.flavour
	}
	return
}

// Groups returns the value of the 'groups' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the collection of groups of user of the cluster.
func (o *Cluster) Groups() *v1.GroupList {
	if o != nil && len(o.fieldSet_) > 31 && o.fieldSet_[31] {
		return o.groups
	}
	return nil
}

// GetGroups returns the value of the 'groups' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the collection of groups of user of the cluster.
func (o *Cluster) GetGroups() (value *v1.GroupList, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 31 && o.fieldSet_[31]
	if ok {
		value = o.groups
	}
	return
}

// HealthState returns the value of the 'health_state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// HealthState indicates the overall health state of the cluster.
func (o *Cluster) HealthState() ClusterHealthState {
	if o != nil && len(o.fieldSet_) > 32 && o.fieldSet_[32] {
		return o.healthState
	}
	return ClusterHealthState("")
}

// GetHealthState returns the value of the 'health_state' attribute and
// a flag indicating if the attribute has a value.
//
// HealthState indicates the overall health state of the cluster.
func (o *Cluster) GetHealthState() (value ClusterHealthState, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 32 && o.fieldSet_[32]
	if ok {
		value = o.healthState
	}
	return
}

// Htpasswd returns the value of the 'htpasswd' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Details for `htpasswd` identity provider.
func (o *Cluster) Htpasswd() *HTPasswdIdentityProvider {
	if o != nil && len(o.fieldSet_) > 33 && o.fieldSet_[33] {
		return o.htpasswd
	}
	return nil
}

// GetHtpasswd returns the value of the 'htpasswd' attribute and
// a flag indicating if the attribute has a value.
//
// Details for `htpasswd` identity provider.
func (o *Cluster) GetHtpasswd() (value *HTPasswdIdentityProvider, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 33 && o.fieldSet_[33]
	if ok {
		value = o.htpasswd
	}
	return
}

// Hypershift returns the value of the 'hypershift' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Hypershift configuration.
func (o *Cluster) Hypershift() *Hypershift {
	if o != nil && len(o.fieldSet_) > 34 && o.fieldSet_[34] {
		return o.hypershift
	}
	return nil
}

// GetHypershift returns the value of the 'hypershift' attribute and
// a flag indicating if the attribute has a value.
//
// Hypershift configuration.
func (o *Cluster) GetHypershift() (value *Hypershift, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 34 && o.fieldSet_[34]
	if ok {
		value = o.hypershift
	}
	return
}

// IdentityProviders returns the value of the 'identity_providers' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the collection of identity providers of the cluster.
func (o *Cluster) IdentityProviders() *v1.IdentityProviderList {
	if o != nil && len(o.fieldSet_) > 35 && o.fieldSet_[35] {
		return o.identityProviders
	}
	return nil
}

// GetIdentityProviders returns the value of the 'identity_providers' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the collection of identity providers of the cluster.
func (o *Cluster) GetIdentityProviders() (value *v1.IdentityProviderList, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 35 && o.fieldSet_[35]
	if ok {
		value = o.identityProviders
	}
	return
}

// ImageRegistry returns the value of the 'image_registry' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// The OpenShift Image Registry configuration
// It provides an internal, integrated container image registry to locally manage images.
// For non ARO-HCP clusters, it is readonly and always enabled
func (o *Cluster) ImageRegistry() *ClusterImageRegistry {
	if o != nil && len(o.fieldSet_) > 36 && o.fieldSet_[36] {
		return o.imageRegistry
	}
	return nil
}

// GetImageRegistry returns the value of the 'image_registry' attribute and
// a flag indicating if the attribute has a value.
//
// The OpenShift Image Registry configuration
// It provides an internal, integrated container image registry to locally manage images.
// For non ARO-HCP clusters, it is readonly and always enabled
func (o *Cluster) GetImageRegistry() (value *ClusterImageRegistry, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 36 && o.fieldSet_[36]
	if ok {
		value = o.imageRegistry
	}
	return
}

// InflightChecks returns the value of the 'inflight_checks' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of inflight checks on this cluster.
func (o *Cluster) InflightChecks() *InflightCheckList {
	if o != nil && len(o.fieldSet_) > 37 && o.fieldSet_[37] {
		return o.inflightChecks
	}
	return nil
}

// GetInflightChecks returns the value of the 'inflight_checks' attribute and
// a flag indicating if the attribute has a value.
//
// List of inflight checks on this cluster.
func (o *Cluster) GetInflightChecks() (value *InflightCheckList, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 37 && o.fieldSet_[37]
	if ok {
		value = o.inflightChecks
	}
	return
}

// InfraID returns the value of the 'infra_ID' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// InfraID is used for example to name the VPCs.
func (o *Cluster) InfraID() string {
	if o != nil && len(o.fieldSet_) > 38 && o.fieldSet_[38] {
		return o.infraID
	}
	return ""
}

// GetInfraID returns the value of the 'infra_ID' attribute and
// a flag indicating if the attribute has a value.
//
// InfraID is used for example to name the VPCs.
func (o *Cluster) GetInfraID() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 38 && o.fieldSet_[38]
	if ok {
		value = o.infraID
	}
	return
}

// Ingresses returns the value of the 'ingresses' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of ingresses on this cluster.
func (o *Cluster) Ingresses() *v1.IngressList {
	if o != nil && len(o.fieldSet_) > 39 && o.fieldSet_[39] {
		return o.ingresses
	}
	return nil
}

// GetIngresses returns the value of the 'ingresses' attribute and
// a flag indicating if the attribute has a value.
//
// List of ingresses on this cluster.
func (o *Cluster) GetIngresses() (value *v1.IngressList, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 39 && o.fieldSet_[39]
	if ok {
		value = o.ingresses
	}
	return
}

// KubeletConfig returns the value of the 'kubelet_config' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Details of cluster-wide KubeletConfig
func (o *Cluster) KubeletConfig() *KubeletConfig {
	if o != nil && len(o.fieldSet_) > 40 && o.fieldSet_[40] {
		return o.kubeletConfig
	}
	return nil
}

// GetKubeletConfig returns the value of the 'kubelet_config' attribute and
// a flag indicating if the attribute has a value.
//
// Details of cluster-wide KubeletConfig
func (o *Cluster) GetKubeletConfig() (value *KubeletConfig, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 40 && o.fieldSet_[40]
	if ok {
		value = o.kubeletConfig
	}
	return
}

// LoadBalancerQuota returns the value of the 'load_balancer_quota' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Load Balancer quota to be assigned to the cluster.
func (o *Cluster) LoadBalancerQuota() int {
	if o != nil && len(o.fieldSet_) > 41 && o.fieldSet_[41] {
		return o.loadBalancerQuota
	}
	return 0
}

// GetLoadBalancerQuota returns the value of the 'load_balancer_quota' attribute and
// a flag indicating if the attribute has a value.
//
// Load Balancer quota to be assigned to the cluster.
func (o *Cluster) GetLoadBalancerQuota() (value int, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 41 && o.fieldSet_[41]
	if ok {
		value = o.loadBalancerQuota
	}
	return
}

// MachinePools returns the value of the 'machine_pools' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of machine pools on this cluster.
func (o *Cluster) MachinePools() *v1.MachinePoolList {
	if o != nil && len(o.fieldSet_) > 42 && o.fieldSet_[42] {
		return o.machinePools
	}
	return nil
}

// GetMachinePools returns the value of the 'machine_pools' attribute and
// a flag indicating if the attribute has a value.
//
// List of machine pools on this cluster.
func (o *Cluster) GetMachinePools() (value *v1.MachinePoolList, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 42 && o.fieldSet_[42]
	if ok {
		value = o.machinePools
	}
	return
}

// Managed returns the value of the 'managed' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Flag indicating if the cluster is managed (by Red Hat) or
// self-managed by the user.
func (o *Cluster) Managed() bool {
	if o != nil && len(o.fieldSet_) > 43 && o.fieldSet_[43] {
		return o.managed
	}
	return false
}

// GetManaged returns the value of the 'managed' attribute and
// a flag indicating if the attribute has a value.
//
// Flag indicating if the cluster is managed (by Red Hat) or
// self-managed by the user.
func (o *Cluster) GetManaged() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 43 && o.fieldSet_[43]
	if ok {
		value = o.managed
	}
	return
}

// ManagedService returns the value of the 'managed_service' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Contains information about Managed Service
func (o *Cluster) ManagedService() *ManagedService {
	if o != nil && len(o.fieldSet_) > 44 && o.fieldSet_[44] {
		return o.managedService
	}
	return nil
}

// GetManagedService returns the value of the 'managed_service' attribute and
// a flag indicating if the attribute has a value.
//
// Contains information about Managed Service
func (o *Cluster) GetManagedService() (value *ManagedService, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 44 && o.fieldSet_[44]
	if ok {
		value = o.managedService
	}
	return
}

// MultiAZ returns the value of the 'multi_AZ' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Flag indicating if the cluster should be created with nodes in
// different availability zones or all the nodes in a single one
// randomly selected.
// For ARO-HCP Clusters, this attribute is unused, and the control plane
// is deployed in multiple availability zones when the Azure region where
// it is deployed supports multiple availability zones.
func (o *Cluster) MultiAZ() bool {
	if o != nil && len(o.fieldSet_) > 45 && o.fieldSet_[45] {
		return o.multiAZ
	}
	return false
}

// GetMultiAZ returns the value of the 'multi_AZ' attribute and
// a flag indicating if the attribute has a value.
//
// Flag indicating if the cluster should be created with nodes in
// different availability zones or all the nodes in a single one
// randomly selected.
// For ARO-HCP Clusters, this attribute is unused, and the control plane
// is deployed in multiple availability zones when the Azure region where
// it is deployed supports multiple availability zones.
func (o *Cluster) GetMultiAZ() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 45 && o.fieldSet_[45]
	if ok {
		value = o.multiAZ
	}
	return
}

// MultiArchEnabled returns the value of the 'multi_arch_enabled' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Indicate whether the cluster is enabled for multi arch workers
func (o *Cluster) MultiArchEnabled() bool {
	if o != nil && len(o.fieldSet_) > 46 && o.fieldSet_[46] {
		return o.multiArchEnabled
	}
	return false
}

// GetMultiArchEnabled returns the value of the 'multi_arch_enabled' attribute and
// a flag indicating if the attribute has a value.
//
// Indicate whether the cluster is enabled for multi arch workers
func (o *Cluster) GetMultiArchEnabled() (value bool, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 46 && o.fieldSet_[46]
	if ok {
		value = o.multiArchEnabled
	}
	return
}

// Name returns the value of the 'name' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Name of the cluster. This name is assigned by the user when the
// cluster is created. This is used to uniquely identify the cluster
func (o *Cluster) Name() string {
	if o != nil && len(o.fieldSet_) > 47 && o.fieldSet_[47] {
		return o.name
	}
	return ""
}

// GetName returns the value of the 'name' attribute and
// a flag indicating if the attribute has a value.
//
// Name of the cluster. This name is assigned by the user when the
// cluster is created. This is used to uniquely identify the cluster
func (o *Cluster) GetName() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 47 && o.fieldSet_[47]
	if ok {
		value = o.name
	}
	return
}

// Network returns the value of the 'network' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Network settings of the cluster.
func (o *Cluster) Network() *Network {
	if o != nil && len(o.fieldSet_) > 48 && o.fieldSet_[48] {
		return o.network
	}
	return nil
}

// GetNetwork returns the value of the 'network' attribute and
// a flag indicating if the attribute has a value.
//
// Network settings of the cluster.
func (o *Cluster) GetNetwork() (value *Network, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 48 && o.fieldSet_[48]
	if ok {
		value = o.network
	}
	return
}

// NodeDrainGracePeriod returns the value of the 'node_drain_grace_period' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Node drain grace period.
func (o *Cluster) NodeDrainGracePeriod() *Value {
	if o != nil && len(o.fieldSet_) > 49 && o.fieldSet_[49] {
		return o.nodeDrainGracePeriod
	}
	return nil
}

// GetNodeDrainGracePeriod returns the value of the 'node_drain_grace_period' attribute and
// a flag indicating if the attribute has a value.
//
// Node drain grace period.
func (o *Cluster) GetNodeDrainGracePeriod() (value *Value, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 49 && o.fieldSet_[49]
	if ok {
		value = o.nodeDrainGracePeriod
	}
	return
}

// NodePools returns the value of the 'node_pools' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// List of node pools on this cluster.
// NodePool is a scalable set of worker nodes attached to a hosted cluster.
func (o *Cluster) NodePools() *NodePoolList {
	if o != nil && len(o.fieldSet_) > 50 && o.fieldSet_[50] {
		return o.nodePools
	}
	return nil
}

// GetNodePools returns the value of the 'node_pools' attribute and
// a flag indicating if the attribute has a value.
//
// List of node pools on this cluster.
// NodePool is a scalable set of worker nodes attached to a hosted cluster.
func (o *Cluster) GetNodePools() (value *NodePoolList, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 50 && o.fieldSet_[50]
	if ok {
		value = o.nodePools
	}
	return
}

// Nodes returns the value of the 'nodes' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Information about the nodes of the cluster.
func (o *Cluster) Nodes() *ClusterNodes {
	if o != nil && len(o.fieldSet_) > 51 && o.fieldSet_[51] {
		return o.nodes
	}
	return nil
}

// GetNodes returns the value of the 'nodes' attribute and
// a flag indicating if the attribute has a value.
//
// Information about the nodes of the cluster.
func (o *Cluster) GetNodes() (value *ClusterNodes, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 51 && o.fieldSet_[51]
	if ok {
		value = o.nodes
	}
	return
}

// OpenshiftVersion returns the value of the 'openshift_version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Version of _OpenShift_ installed in the cluster, for example `4.0.0-0.2`.
//
// When retrieving a cluster this will always be reported.
//
// When provisioning a cluster this will be ignored, as the version to
// deploy will be determined internally.
func (o *Cluster) OpenshiftVersion() string {
	if o != nil && len(o.fieldSet_) > 52 && o.fieldSet_[52] {
		return o.openshiftVersion
	}
	return ""
}

// GetOpenshiftVersion returns the value of the 'openshift_version' attribute and
// a flag indicating if the attribute has a value.
//
// Version of _OpenShift_ installed in the cluster, for example `4.0.0-0.2`.
//
// When retrieving a cluster this will always be reported.
//
// When provisioning a cluster this will be ignored, as the version to
// deploy will be determined internally.
func (o *Cluster) GetOpenshiftVersion() (value string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 52 && o.fieldSet_[52]
	if ok {
		value = o.openshiftVersion
	}
	return
}

// Product returns the value of the 'product' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the product type of this cluster.
func (o *Cluster) Product() *v1.Product {
	if o != nil && len(o.fieldSet_) > 53 && o.fieldSet_[53] {
		return o.product
	}
	return nil
}

// GetProduct returns the value of the 'product' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the product type of this cluster.
func (o *Cluster) GetProduct() (value *v1.Product, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 53 && o.fieldSet_[53]
	if ok {
		value = o.product
	}
	return
}

// Properties returns the value of the 'properties' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// User defined properties for tagging and querying.
func (o *Cluster) Properties() map[string]string {
	if o != nil && len(o.fieldSet_) > 54 && o.fieldSet_[54] {
		return o.properties
	}
	return nil
}

// GetProperties returns the value of the 'properties' attribute and
// a flag indicating if the attribute has a value.
//
// User defined properties for tagging and querying.
func (o *Cluster) GetProperties() (value map[string]string, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 54 && o.fieldSet_[54]
	if ok {
		value = o.properties
	}
	return
}

// ProvisionShard returns the value of the 'provision_shard' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// ProvisionShard contains the properties of the provision shard, including AWS and GCP related configurations
func (o *Cluster) ProvisionShard() *ProvisionShard {
	if o != nil && len(o.fieldSet_) > 55 && o.fieldSet_[55] {
		return o.provisionShard
	}
	return nil
}

// GetProvisionShard returns the value of the 'provision_shard' attribute and
// a flag indicating if the attribute has a value.
//
// ProvisionShard contains the properties of the provision shard, including AWS and GCP related configurations
func (o *Cluster) GetProvisionShard() (value *ProvisionShard, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 55 && o.fieldSet_[55]
	if ok {
		value = o.provisionShard
	}
	return
}

// Proxy returns the value of the 'proxy' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Proxy.
func (o *Cluster) Proxy() *Proxy {
	if o != nil && len(o.fieldSet_) > 56 && o.fieldSet_[56] {
		return o.proxy
	}
	return nil
}

// GetProxy returns the value of the 'proxy' attribute and
// a flag indicating if the attribute has a value.
//
// Proxy.
func (o *Cluster) GetProxy() (value *Proxy, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 56 && o.fieldSet_[56]
	if ok {
		value = o.proxy
	}
	return
}

// Region returns the value of the 'region' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the cloud provider region where the cluster is installed.
func (o *Cluster) Region() *v1.CloudRegion {
	if o != nil && len(o.fieldSet_) > 57 && o.fieldSet_[57] {
		return o.region
	}
	return nil
}

// GetRegion returns the value of the 'region' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the cloud provider region where the cluster is installed.
func (o *Cluster) GetRegion() (value *v1.CloudRegion, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 57 && o.fieldSet_[57]
	if ok {
		value = o.region
	}
	return
}

// RegistryConfig returns the value of the 'registry_config' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// External registry configuration for the cluster
func (o *Cluster) RegistryConfig() *ClusterRegistryConfig {
	if o != nil && len(o.fieldSet_) > 58 && o.fieldSet_[58] {
		return o.registryConfig
	}
	return nil
}

// GetRegistryConfig returns the value of the 'registry_config' attribute and
// a flag indicating if the attribute has a value.
//
// External registry configuration for the cluster
func (o *Cluster) GetRegistryConfig() (value *ClusterRegistryConfig, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 58 && o.fieldSet_[58]
	if ok {
		value = o.registryConfig
	}
	return
}

// State returns the value of the 'state' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Overall state of the cluster.
func (o *Cluster) State() ClusterState {
	if o != nil && len(o.fieldSet_) > 59 && o.fieldSet_[59] {
		return o.state
	}
	return ClusterState("")
}

// GetState returns the value of the 'state' attribute and
// a flag indicating if the attribute has a value.
//
// Overall state of the cluster.
func (o *Cluster) GetState() (value ClusterState, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 59 && o.fieldSet_[59]
	if ok {
		value = o.state
	}
	return
}

// Status returns the value of the 'status' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Status of cluster
func (o *Cluster) Status() *ClusterStatus {
	if o != nil && len(o.fieldSet_) > 60 && o.fieldSet_[60] {
		return o.status
	}
	return nil
}

// GetStatus returns the value of the 'status' attribute and
// a flag indicating if the attribute has a value.
//
// Status of cluster
func (o *Cluster) GetStatus() (value *ClusterStatus, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 60 && o.fieldSet_[60]
	if ok {
		value = o.status
	}
	return
}

// StorageQuota returns the value of the 'storage_quota' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Storage quota to be assigned to the cluster.
func (o *Cluster) StorageQuota() *Value {
	if o != nil && len(o.fieldSet_) > 61 && o.fieldSet_[61] {
		return o.storageQuota
	}
	return nil
}

// GetStorageQuota returns the value of the 'storage_quota' attribute and
// a flag indicating if the attribute has a value.
//
// Storage quota to be assigned to the cluster.
func (o *Cluster) GetStorageQuota() (value *Value, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 61 && o.fieldSet_[61]
	if ok {
		value = o.storageQuota
	}
	return
}

// Subscription returns the value of the 'subscription' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the subscription that comes from the account management service when the cluster
// is registered.
func (o *Cluster) Subscription() *v1.Subscription {
	if o != nil && len(o.fieldSet_) > 62 && o.fieldSet_[62] {
		return o.subscription
	}
	return nil
}

// GetSubscription returns the value of the 'subscription' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the subscription that comes from the account management service when the cluster
// is registered.
func (o *Cluster) GetSubscription() (value *v1.Subscription, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 62 && o.fieldSet_[62]
	if ok {
		value = o.subscription
	}
	return
}

// Version returns the value of the 'version' attribute, or
// the zero value of the type if the attribute doesn't have a value.
//
// Link to the version of _OpenShift_ that will be used to install the cluster.
func (o *Cluster) Version() *Version {
	if o != nil && len(o.fieldSet_) > 63 && o.fieldSet_[63] {
		return o.version
	}
	return nil
}

// GetVersion returns the value of the 'version' attribute and
// a flag indicating if the attribute has a value.
//
// Link to the version of _OpenShift_ that will be used to install the cluster.
func (o *Cluster) GetVersion() (value *Version, ok bool) {
	ok = o != nil && len(o.fieldSet_) > 63 && o.fieldSet_[63]
	if ok {
		value = o.version
	}
	return
}

// ClusterListKind is the name of the type used to represent list of objects of
// type 'cluster'.
const ClusterListKind = "ClusterList"

// ClusterListLinkKind is the name of the type used to represent links to list
// of objects of type 'cluster'.
const ClusterListLinkKind = "ClusterListLink"

// ClusterNilKind is the name of the type used to nil lists of objects of
// type 'cluster'.
const ClusterListNilKind = "ClusterListNil"

// ClusterList is a list of values of the 'cluster' type.
type ClusterList struct {
	href  string
	link  bool
	items []*Cluster
}

// Kind returns the name of the type of the object.
func (l *ClusterList) Kind() string {
	if l == nil {
		return ClusterListNilKind
	}
	if l.link {
		return ClusterListLinkKind
	}
	return ClusterListKind
}

// Link returns true iif this is a link.
func (l *ClusterList) Link() bool {
	return l != nil && l.link
}

// HREF returns the link to the list.
func (l *ClusterList) HREF() string {
	if l != nil {
		return l.href
	}
	return ""
}

// GetHREF returns the link of the list and a flag indicating if the
// link has a value.
func (l *ClusterList) GetHREF() (value string, ok bool) {
	ok = l != nil && l.href != ""
	if ok {
		value = l.href
	}
	return
}

// Len returns the length of the list.
func (l *ClusterList) Len() int {
	if l == nil {
		return 0
	}
	return len(l.items)
}

// Items sets the items of the list.
func (l *ClusterList) SetLink(link bool) {
	l.link = link
}

// Items sets the items of the list.
func (l *ClusterList) SetHREF(href string) {
	l.href = href
}

// Items sets the items of the list.
func (l *ClusterList) SetItems(items []*Cluster) {
	l.items = items
}

// Items returns the items of the list.
func (l *ClusterList) Items() []*Cluster {
	if l == nil {
		return nil
	}
	return l.items
}

// Empty returns true if the list is empty.
func (l *ClusterList) Empty() bool {
	return l == nil || len(l.items) == 0
}

// Get returns the item of the list with the given index. If there is no item with
// that index it returns nil.
func (l *ClusterList) Get(i int) *Cluster {
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
func (l *ClusterList) Slice() []*Cluster {
	var slice []*Cluster
	if l == nil {
		slice = make([]*Cluster, 0)
	} else {
		slice = make([]*Cluster, len(l.items))
		copy(slice, l.items)
	}
	return slice
}

// Each runs the given function for each item of the list, in order. If the function
// returns false the iteration stops, otherwise it continues till all the elements
// of the list have been processed.
func (l *ClusterList) Each(f func(item *Cluster) bool) {
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
func (l *ClusterList) Range(f func(index int, item *Cluster) bool) {
	if l == nil {
		return
	}
	for index, item := range l.items {
		if !f(index, item) {
			break
		}
	}
}
