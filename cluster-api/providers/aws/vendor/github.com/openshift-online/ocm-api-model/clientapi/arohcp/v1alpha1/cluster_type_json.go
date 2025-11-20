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
	"io"
	"sort"
	"time"

	jsoniter "github.com/json-iterator/go"
	v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalCluster writes a value of the 'cluster' type to the given writer.
func MarshalCluster(object *Cluster, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteCluster(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteCluster writes a value of the 'cluster' type to the given stream.
func WriteCluster(object *Cluster, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(ClusterLinkKind)
	} else {
		stream.WriteString(ClusterKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.api != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("api")
		WriteClusterAPI(object.api, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.aws != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws")
		WriteAWS(object.aws, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.awsInfrastructureAccessRoleGrants != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_infrastructure_access_role_grants")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		v1.WriteAWSInfrastructureAccessRoleGrantList(object.awsInfrastructureAccessRoleGrants.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.ccs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ccs")
		WriteCCS(object.ccs, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.dns != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("dns")
		WriteDNS(object.dns, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("fips")
		stream.WriteBool(object.fips)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9] && object.gcp != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp")
		WriteGCP(object.gcp, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10] && object.gcpEncryptionKey != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp_encryption_key")
		WriteGCPEncryptionKey(object.gcpEncryptionKey, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11] && object.gcpNetwork != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp_network")
		WriteGCPNetwork(object.gcpNetwork, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_trust_bundle")
		stream.WriteString(object.additionalTrustBundle)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13] && object.addons != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("addons")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		v1.WriteAddOnInstallationList(object.addons.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14] && object.autoNode != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("auto_node")
		WriteClusterAutoNode(object.autoNode, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15] && object.autoscaler != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("autoscaler")
		WriteClusterAutoscaler(object.autoscaler, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 16 && object.fieldSet_[16] && object.azure != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("azure")
		WriteAzure(object.azure, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 17 && object.fieldSet_[17]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("billing_model")
		stream.WriteString(string(object.billingModel))
		count++
	}
	present_ = len(object.fieldSet_) > 18 && object.fieldSet_[18] && object.byoOidc != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("byo_oidc")
		WriteByoOidc(object.byoOidc, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 19 && object.fieldSet_[19] && object.cloudProvider != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		v1.WriteCloudProvider(object.cloudProvider, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 20 && object.fieldSet_[20] && object.console != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("console")
		WriteClusterConsole(object.console, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 21 && object.fieldSet_[21]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 22 && object.fieldSet_[22] && object.deleteProtection != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("delete_protection")
		WriteDeleteProtection(object.deleteProtection, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 23 && object.fieldSet_[23]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("disable_user_workload_monitoring")
		stream.WriteBool(object.disableUserWorkloadMonitoring)
		count++
	}
	present_ = len(object.fieldSet_) > 24 && object.fieldSet_[24]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("domain_prefix")
		stream.WriteString(object.domainPrefix)
		count++
	}
	present_ = len(object.fieldSet_) > 25 && object.fieldSet_[25]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("etcd_encryption")
		stream.WriteBool(object.etcdEncryption)
		count++
	}
	present_ = len(object.fieldSet_) > 26 && object.fieldSet_[26]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("expiration_timestamp")
		stream.WriteString((object.expirationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 27 && object.fieldSet_[27]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_id")
		stream.WriteString(object.externalID)
		count++
	}
	present_ = len(object.fieldSet_) > 28 && object.fieldSet_[28] && object.externalAuthConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_auth_config")
		WriteExternalAuthConfig(object.externalAuthConfig, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 29 && object.fieldSet_[29] && object.externalConfiguration != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_configuration")
		WriteExternalConfiguration(object.externalConfiguration, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 30 && object.fieldSet_[30] && object.flavour != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("flavour")
		v1.WriteFlavour(object.flavour, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 31 && object.fieldSet_[31] && object.groups != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("groups")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		v1.WriteGroupList(object.groups.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = len(object.fieldSet_) > 32 && object.fieldSet_[32]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("health_state")
		stream.WriteString(string(object.healthState))
		count++
	}
	present_ = len(object.fieldSet_) > 33 && object.fieldSet_[33] && object.htpasswd != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("htpasswd")
		WriteHTPasswdIdentityProvider(object.htpasswd, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 34 && object.fieldSet_[34] && object.hypershift != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hypershift")
		WriteHypershift(object.hypershift, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 35 && object.fieldSet_[35] && object.identityProviders != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("identity_providers")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		v1.WriteIdentityProviderList(object.identityProviders.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = len(object.fieldSet_) > 36 && object.fieldSet_[36] && object.imageRegistry != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("image_registry")
		WriteClusterImageRegistry(object.imageRegistry, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 37 && object.fieldSet_[37] && object.inflightChecks != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("inflight_checks")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteInflightCheckList(object.inflightChecks.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = len(object.fieldSet_) > 38 && object.fieldSet_[38]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("infra_id")
		stream.WriteString(object.infraID)
		count++
	}
	present_ = len(object.fieldSet_) > 39 && object.fieldSet_[39] && object.ingresses != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ingresses")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		v1.WriteIngressList(object.ingresses.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = len(object.fieldSet_) > 40 && object.fieldSet_[40] && object.kubeletConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kubelet_config")
		WriteKubeletConfig(object.kubeletConfig, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 41 && object.fieldSet_[41]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("load_balancer_quota")
		stream.WriteInt(object.loadBalancerQuota)
		count++
	}
	present_ = len(object.fieldSet_) > 42 && object.fieldSet_[42] && object.machinePools != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("machine_pools")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		v1.WriteMachinePoolList(object.machinePools.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = len(object.fieldSet_) > 43 && object.fieldSet_[43]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed")
		stream.WriteBool(object.managed)
		count++
	}
	present_ = len(object.fieldSet_) > 44 && object.fieldSet_[44] && object.managedService != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed_service")
		WriteManagedService(object.managedService, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 45 && object.fieldSet_[45]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("multi_az")
		stream.WriteBool(object.multiAZ)
		count++
	}
	present_ = len(object.fieldSet_) > 46 && object.fieldSet_[46]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("multi_arch_enabled")
		stream.WriteBool(object.multiArchEnabled)
		count++
	}
	present_ = len(object.fieldSet_) > 47 && object.fieldSet_[47]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 48 && object.fieldSet_[48] && object.network != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("network")
		WriteNetwork(object.network, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 49 && object.fieldSet_[49] && object.nodeDrainGracePeriod != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("node_drain_grace_period")
		WriteValue(object.nodeDrainGracePeriod, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 50 && object.fieldSet_[50] && object.nodePools != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("node_pools")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteNodePoolList(object.nodePools.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = len(object.fieldSet_) > 51 && object.fieldSet_[51] && object.nodes != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("nodes")
		WriteClusterNodes(object.nodes, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 52 && object.fieldSet_[52]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("openshift_version")
		stream.WriteString(object.openshiftVersion)
		count++
	}
	present_ = len(object.fieldSet_) > 53 && object.fieldSet_[53] && object.product != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product")
		v1.WriteProduct(object.product, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 54 && object.fieldSet_[54] && object.properties != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("properties")
		if object.properties != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.properties))
			i := 0
			for key := range object.properties {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.properties[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 55 && object.fieldSet_[55] && object.provisionShard != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("provision_shard")
		WriteProvisionShard(object.provisionShard, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 56 && object.fieldSet_[56] && object.proxy != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("proxy")
		WriteProxy(object.proxy, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 57 && object.fieldSet_[57] && object.region != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region")
		v1.WriteCloudRegion(object.region, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 58 && object.fieldSet_[58] && object.registryConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("registry_config")
		WriteClusterRegistryConfig(object.registryConfig, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 59 && object.fieldSet_[59]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(string(object.state))
		count++
	}
	present_ = len(object.fieldSet_) > 60 && object.fieldSet_[60] && object.status != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		WriteClusterStatus(object.status, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 61 && object.fieldSet_[61] && object.storageQuota != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("storage_quota")
		WriteValue(object.storageQuota, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 62 && object.fieldSet_[62] && object.subscription != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription")
		v1.WriteSubscription(object.subscription, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 63 && object.fieldSet_[63] && object.version != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		WriteVersion(object.version, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalCluster reads a value of the 'cluster' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalCluster(source interface{}) (object *Cluster, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadCluster(iterator)
	err = iterator.Error
	return
}

// ReadCluster reads a value of the 'cluster' type from the given iterator.
func ReadCluster(iterator *jsoniter.Iterator) *Cluster {
	object := &Cluster{
		fieldSet_: make([]bool, 64),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ClusterLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "api":
			value := ReadClusterAPI(iterator)
			object.api = value
			object.fieldSet_[3] = true
		case "aws":
			value := ReadAWS(iterator)
			object.aws = value
			object.fieldSet_[4] = true
		case "aws_infrastructure_access_role_grants":
			value := &v1.AWSInfrastructureAccessRoleGrantList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == v1.AWSInfrastructureAccessRoleGrantListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(v1.ReadAWSInfrastructureAccessRoleGrantList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.awsInfrastructureAccessRoleGrants = value
			object.fieldSet_[5] = true
		case "ccs":
			value := ReadCCS(iterator)
			object.ccs = value
			object.fieldSet_[6] = true
		case "dns":
			value := ReadDNS(iterator)
			object.dns = value
			object.fieldSet_[7] = true
		case "fips":
			value := iterator.ReadBool()
			object.fips = value
			object.fieldSet_[8] = true
		case "gcp":
			value := ReadGCP(iterator)
			object.gcp = value
			object.fieldSet_[9] = true
		case "gcp_encryption_key":
			value := ReadGCPEncryptionKey(iterator)
			object.gcpEncryptionKey = value
			object.fieldSet_[10] = true
		case "gcp_network":
			value := ReadGCPNetwork(iterator)
			object.gcpNetwork = value
			object.fieldSet_[11] = true
		case "additional_trust_bundle":
			value := iterator.ReadString()
			object.additionalTrustBundle = value
			object.fieldSet_[12] = true
		case "addons":
			value := &v1.AddOnInstallationList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == v1.AddOnInstallationListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(v1.ReadAddOnInstallationList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.addons = value
			object.fieldSet_[13] = true
		case "auto_node":
			value := ReadClusterAutoNode(iterator)
			object.autoNode = value
			object.fieldSet_[14] = true
		case "autoscaler":
			value := ReadClusterAutoscaler(iterator)
			object.autoscaler = value
			object.fieldSet_[15] = true
		case "azure":
			value := ReadAzure(iterator)
			object.azure = value
			object.fieldSet_[16] = true
		case "billing_model":
			text := iterator.ReadString()
			value := BillingModel(text)
			object.billingModel = value
			object.fieldSet_[17] = true
		case "byo_oidc":
			value := ReadByoOidc(iterator)
			object.byoOidc = value
			object.fieldSet_[18] = true
		case "cloud_provider":
			value := v1.ReadCloudProvider(iterator)
			object.cloudProvider = value
			object.fieldSet_[19] = true
		case "console":
			value := ReadClusterConsole(iterator)
			object.console = value
			object.fieldSet_[20] = true
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.fieldSet_[21] = true
		case "delete_protection":
			value := ReadDeleteProtection(iterator)
			object.deleteProtection = value
			object.fieldSet_[22] = true
		case "disable_user_workload_monitoring":
			value := iterator.ReadBool()
			object.disableUserWorkloadMonitoring = value
			object.fieldSet_[23] = true
		case "domain_prefix":
			value := iterator.ReadString()
			object.domainPrefix = value
			object.fieldSet_[24] = true
		case "etcd_encryption":
			value := iterator.ReadBool()
			object.etcdEncryption = value
			object.fieldSet_[25] = true
		case "expiration_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.expirationTimestamp = value
			object.fieldSet_[26] = true
		case "external_id":
			value := iterator.ReadString()
			object.externalID = value
			object.fieldSet_[27] = true
		case "external_auth_config":
			value := ReadExternalAuthConfig(iterator)
			object.externalAuthConfig = value
			object.fieldSet_[28] = true
		case "external_configuration":
			value := ReadExternalConfiguration(iterator)
			object.externalConfiguration = value
			object.fieldSet_[29] = true
		case "flavour":
			value := v1.ReadFlavour(iterator)
			object.flavour = value
			object.fieldSet_[30] = true
		case "groups":
			value := &v1.GroupList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == v1.GroupListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(v1.ReadGroupList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.groups = value
			object.fieldSet_[31] = true
		case "health_state":
			text := iterator.ReadString()
			value := ClusterHealthState(text)
			object.healthState = value
			object.fieldSet_[32] = true
		case "htpasswd":
			value := ReadHTPasswdIdentityProvider(iterator)
			object.htpasswd = value
			object.fieldSet_[33] = true
		case "hypershift":
			value := ReadHypershift(iterator)
			object.hypershift = value
			object.fieldSet_[34] = true
		case "identity_providers":
			value := &v1.IdentityProviderList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == v1.IdentityProviderListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(v1.ReadIdentityProviderList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.identityProviders = value
			object.fieldSet_[35] = true
		case "image_registry":
			value := ReadClusterImageRegistry(iterator)
			object.imageRegistry = value
			object.fieldSet_[36] = true
		case "inflight_checks":
			value := &InflightCheckList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == InflightCheckListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadInflightCheckList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.inflightChecks = value
			object.fieldSet_[37] = true
		case "infra_id":
			value := iterator.ReadString()
			object.infraID = value
			object.fieldSet_[38] = true
		case "ingresses":
			value := &v1.IngressList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == v1.IngressListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(v1.ReadIngressList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.ingresses = value
			object.fieldSet_[39] = true
		case "kubelet_config":
			value := ReadKubeletConfig(iterator)
			object.kubeletConfig = value
			object.fieldSet_[40] = true
		case "load_balancer_quota":
			value := iterator.ReadInt()
			object.loadBalancerQuota = value
			object.fieldSet_[41] = true
		case "machine_pools":
			value := &v1.MachinePoolList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == v1.MachinePoolListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(v1.ReadMachinePoolList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.machinePools = value
			object.fieldSet_[42] = true
		case "managed":
			value := iterator.ReadBool()
			object.managed = value
			object.fieldSet_[43] = true
		case "managed_service":
			value := ReadManagedService(iterator)
			object.managedService = value
			object.fieldSet_[44] = true
		case "multi_az":
			value := iterator.ReadBool()
			object.multiAZ = value
			object.fieldSet_[45] = true
		case "multi_arch_enabled":
			value := iterator.ReadBool()
			object.multiArchEnabled = value
			object.fieldSet_[46] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[47] = true
		case "network":
			value := ReadNetwork(iterator)
			object.network = value
			object.fieldSet_[48] = true
		case "node_drain_grace_period":
			value := ReadValue(iterator)
			object.nodeDrainGracePeriod = value
			object.fieldSet_[49] = true
		case "node_pools":
			value := &NodePoolList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == NodePoolListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadNodePoolList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.nodePools = value
			object.fieldSet_[50] = true
		case "nodes":
			value := ReadClusterNodes(iterator)
			object.nodes = value
			object.fieldSet_[51] = true
		case "openshift_version":
			value := iterator.ReadString()
			object.openshiftVersion = value
			object.fieldSet_[52] = true
		case "product":
			value := v1.ReadProduct(iterator)
			object.product = value
			object.fieldSet_[53] = true
		case "properties":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.properties = value
			object.fieldSet_[54] = true
		case "provision_shard":
			value := ReadProvisionShard(iterator)
			object.provisionShard = value
			object.fieldSet_[55] = true
		case "proxy":
			value := ReadProxy(iterator)
			object.proxy = value
			object.fieldSet_[56] = true
		case "region":
			value := v1.ReadCloudRegion(iterator)
			object.region = value
			object.fieldSet_[57] = true
		case "registry_config":
			value := ReadClusterRegistryConfig(iterator)
			object.registryConfig = value
			object.fieldSet_[58] = true
		case "state":
			text := iterator.ReadString()
			value := ClusterState(text)
			object.state = value
			object.fieldSet_[59] = true
		case "status":
			value := ReadClusterStatus(iterator)
			object.status = value
			object.fieldSet_[60] = true
		case "storage_quota":
			value := ReadValue(iterator)
			object.storageQuota = value
			object.fieldSet_[61] = true
		case "subscription":
			value := v1.ReadSubscription(iterator)
			object.subscription = value
			object.fieldSet_[62] = true
		case "version":
			value := ReadVersion(iterator)
			object.version = value
			object.fieldSet_[63] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
