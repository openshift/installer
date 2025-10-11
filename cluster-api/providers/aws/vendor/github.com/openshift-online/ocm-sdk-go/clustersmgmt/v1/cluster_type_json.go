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
	"io"
	"sort"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
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
	if object.bitmap_&1 != 0 {
		stream.WriteString(ClusterLinkKind)
	} else {
		stream.WriteString(ClusterKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0 && object.api != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("api")
		WriteClusterAPI(object.api, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.aws != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws")
		WriteAWS(object.aws, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.awsInfrastructureAccessRoleGrants != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_infrastructure_access_role_grants")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteAWSInfrastructureAccessRoleGrantList(object.awsInfrastructureAccessRoleGrants.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.ccs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ccs")
		WriteCCS(object.ccs, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0 && object.dns != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("dns")
		WriteDNS(object.dns, stream)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("fips")
		stream.WriteBool(object.fips)
		count++
	}
	present_ = object.bitmap_&512 != 0 && object.gcp != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp")
		WriteGCP(object.gcp, stream)
		count++
	}
	present_ = object.bitmap_&1024 != 0 && object.gcpEncryptionKey != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp_encryption_key")
		WriteGCPEncryptionKey(object.gcpEncryptionKey, stream)
		count++
	}
	present_ = object.bitmap_&2048 != 0 && object.gcpNetwork != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp_network")
		WriteGCPNetwork(object.gcpNetwork, stream)
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("additional_trust_bundle")
		stream.WriteString(object.additionalTrustBundle)
		count++
	}
	present_ = object.bitmap_&8192 != 0 && object.addons != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("addons")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteAddOnInstallationList(object.addons.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = object.bitmap_&16384 != 0 && object.autoscaler != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("autoscaler")
		WriteClusterAutoscaler(object.autoscaler, stream)
		count++
	}
	present_ = object.bitmap_&32768 != 0 && object.azure != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("azure")
		WriteAzure(object.azure, stream)
		count++
	}
	present_ = object.bitmap_&65536 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("billing_model")
		stream.WriteString(string(object.billingModel))
		count++
	}
	present_ = object.bitmap_&131072 != 0 && object.byoOidc != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("byo_oidc")
		WriteByoOidc(object.byoOidc, stream)
		count++
	}
	present_ = object.bitmap_&262144 != 0 && object.capabilities != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("capabilities")
		WriteClusterCapabilities(object.capabilities, stream)
		count++
	}
	present_ = object.bitmap_&524288 != 0 && object.cloudProvider != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		WriteCloudProvider(object.cloudProvider, stream)
		count++
	}
	present_ = object.bitmap_&1048576 != 0 && object.console != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("console")
		WriteClusterConsole(object.console, stream)
		count++
	}
	present_ = object.bitmap_&2097152 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&4194304 != 0 && object.deleteProtection != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("delete_protection")
		WriteDeleteProtection(object.deleteProtection, stream)
		count++
	}
	present_ = object.bitmap_&8388608 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("disable_user_workload_monitoring")
		stream.WriteBool(object.disableUserWorkloadMonitoring)
		count++
	}
	present_ = object.bitmap_&16777216 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("domain_prefix")
		stream.WriteString(object.domainPrefix)
		count++
	}
	present_ = object.bitmap_&33554432 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("etcd_encryption")
		stream.WriteBool(object.etcdEncryption)
		count++
	}
	present_ = object.bitmap_&67108864 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("expiration_timestamp")
		stream.WriteString((object.expirationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&134217728 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_id")
		stream.WriteString(object.externalID)
		count++
	}
	present_ = object.bitmap_&268435456 != 0 && object.externalAuthConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_auth_config")
		WriteExternalAuthConfig(object.externalAuthConfig, stream)
		count++
	}
	present_ = object.bitmap_&536870912 != 0 && object.externalConfiguration != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("external_configuration")
		WriteExternalConfiguration(object.externalConfiguration, stream)
		count++
	}
	present_ = object.bitmap_&1073741824 != 0 && object.flavour != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("flavour")
		WriteFlavour(object.flavour, stream)
		count++
	}
	present_ = object.bitmap_&2147483648 != 0 && object.groups != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("groups")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteGroupList(object.groups.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = object.bitmap_&4294967296 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("health_state")
		stream.WriteString(string(object.healthState))
		count++
	}
	present_ = object.bitmap_&8589934592 != 0 && object.htpasswd != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("htpasswd")
		WriteHTPasswdIdentityProvider(object.htpasswd, stream)
		count++
	}
	present_ = object.bitmap_&17179869184 != 0 && object.hypershift != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hypershift")
		WriteHypershift(object.hypershift, stream)
		count++
	}
	present_ = object.bitmap_&34359738368 != 0 && object.identityProviders != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("identity_providers")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteIdentityProviderList(object.identityProviders.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = object.bitmap_&68719476736 != 0 && object.inflightChecks != nil
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
	present_ = object.bitmap_&137438953472 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("infra_id")
		stream.WriteString(object.infraID)
		count++
	}
	present_ = object.bitmap_&274877906944 != 0 && object.ingresses != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ingresses")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteIngressList(object.ingresses.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = object.bitmap_&549755813888 != 0 && object.kubeletConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kubelet_config")
		WriteKubeletConfig(object.kubeletConfig, stream)
		count++
	}
	present_ = object.bitmap_&1099511627776 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("load_balancer_quota")
		stream.WriteInt(object.loadBalancerQuota)
		count++
	}
	present_ = object.bitmap_&2199023255552 != 0 && object.machinePools != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("machine_pools")
		stream.WriteObjectStart()
		stream.WriteObjectField("items")
		WriteMachinePoolList(object.machinePools.Items(), stream)
		stream.WriteObjectEnd()
		count++
	}
	present_ = object.bitmap_&4398046511104 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed")
		stream.WriteBool(object.managed)
		count++
	}
	present_ = object.bitmap_&8796093022208 != 0 && object.managedService != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("managed_service")
		WriteManagedService(object.managedService, stream)
		count++
	}
	present_ = object.bitmap_&17592186044416 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("multi_az")
		stream.WriteBool(object.multiAZ)
		count++
	}
	present_ = object.bitmap_&35184372088832 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("multi_arch_enabled")
		stream.WriteBool(object.multiArchEnabled)
		count++
	}
	present_ = object.bitmap_&70368744177664 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&140737488355328 != 0 && object.network != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("network")
		WriteNetwork(object.network, stream)
		count++
	}
	present_ = object.bitmap_&281474976710656 != 0 && object.nodeDrainGracePeriod != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("node_drain_grace_period")
		WriteValue(object.nodeDrainGracePeriod, stream)
		count++
	}
	present_ = object.bitmap_&562949953421312 != 0 && object.nodePools != nil
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
	present_ = object.bitmap_&1125899906842624 != 0 && object.nodes != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("nodes")
		WriteClusterNodes(object.nodes, stream)
		count++
	}
	present_ = object.bitmap_&2251799813685248 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("openshift_version")
		stream.WriteString(object.openshiftVersion)
		count++
	}
	present_ = object.bitmap_&4503599627370496 != 0 && object.product != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("product")
		WriteProduct(object.product, stream)
		count++
	}
	present_ = object.bitmap_&9007199254740992 != 0 && object.properties != nil
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
	present_ = object.bitmap_&18014398509481984 != 0 && object.provisionShard != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("provision_shard")
		WriteProvisionShard(object.provisionShard, stream)
		count++
	}
	present_ = object.bitmap_&36028797018963968 != 0 && object.proxy != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("proxy")
		WriteProxy(object.proxy, stream)
		count++
	}
	present_ = object.bitmap_&72057594037927936 != 0 && object.region != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region")
		WriteCloudRegion(object.region, stream)
		count++
	}
	present_ = object.bitmap_&144115188075855872 != 0 && object.registryConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("registry_config")
		WriteClusterRegistryConfig(object.registryConfig, stream)
		count++
	}
	present_ = object.bitmap_&288230376151711744 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(string(object.state))
		count++
	}
	present_ = object.bitmap_&576460752303423488 != 0 && object.status != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		WriteClusterStatus(object.status, stream)
		count++
	}
	present_ = object.bitmap_&1152921504606846976 != 0 && object.storageQuota != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("storage_quota")
		WriteValue(object.storageQuota, stream)
		count++
	}
	present_ = object.bitmap_&2305843009213693952 != 0 && object.subscription != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subscription")
		WriteSubscription(object.subscription, stream)
		count++
	}
	present_ = object.bitmap_&4611686018427387904 != 0 && object.version != nil
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
	object := &Cluster{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ClusterLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "api":
			value := ReadClusterAPI(iterator)
			object.api = value
			object.bitmap_ |= 8
		case "aws":
			value := ReadAWS(iterator)
			object.aws = value
			object.bitmap_ |= 16
		case "aws_infrastructure_access_role_grants":
			value := &AWSInfrastructureAccessRoleGrantList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == AWSInfrastructureAccessRoleGrantListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadAWSInfrastructureAccessRoleGrantList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.awsInfrastructureAccessRoleGrants = value
			object.bitmap_ |= 32
		case "ccs":
			value := ReadCCS(iterator)
			object.ccs = value
			object.bitmap_ |= 64
		case "dns":
			value := ReadDNS(iterator)
			object.dns = value
			object.bitmap_ |= 128
		case "fips":
			value := iterator.ReadBool()
			object.fips = value
			object.bitmap_ |= 256
		case "gcp":
			value := ReadGCP(iterator)
			object.gcp = value
			object.bitmap_ |= 512
		case "gcp_encryption_key":
			value := ReadGCPEncryptionKey(iterator)
			object.gcpEncryptionKey = value
			object.bitmap_ |= 1024
		case "gcp_network":
			value := ReadGCPNetwork(iterator)
			object.gcpNetwork = value
			object.bitmap_ |= 2048
		case "additional_trust_bundle":
			value := iterator.ReadString()
			object.additionalTrustBundle = value
			object.bitmap_ |= 4096
		case "addons":
			value := &AddOnInstallationList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == AddOnInstallationListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadAddOnInstallationList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.addons = value
			object.bitmap_ |= 8192
		case "autoscaler":
			value := ReadClusterAutoscaler(iterator)
			object.autoscaler = value
			object.bitmap_ |= 16384
		case "azure":
			value := ReadAzure(iterator)
			object.azure = value
			object.bitmap_ |= 32768
		case "billing_model":
			text := iterator.ReadString()
			value := BillingModel(text)
			object.billingModel = value
			object.bitmap_ |= 65536
		case "byo_oidc":
			value := ReadByoOidc(iterator)
			object.byoOidc = value
			object.bitmap_ |= 131072
		case "capabilities":
			value := ReadClusterCapabilities(iterator)
			object.capabilities = value
			object.bitmap_ |= 262144
		case "cloud_provider":
			value := ReadCloudProvider(iterator)
			object.cloudProvider = value
			object.bitmap_ |= 524288
		case "console":
			value := ReadClusterConsole(iterator)
			object.console = value
			object.bitmap_ |= 1048576
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.bitmap_ |= 2097152
		case "delete_protection":
			value := ReadDeleteProtection(iterator)
			object.deleteProtection = value
			object.bitmap_ |= 4194304
		case "disable_user_workload_monitoring":
			value := iterator.ReadBool()
			object.disableUserWorkloadMonitoring = value
			object.bitmap_ |= 8388608
		case "domain_prefix":
			value := iterator.ReadString()
			object.domainPrefix = value
			object.bitmap_ |= 16777216
		case "etcd_encryption":
			value := iterator.ReadBool()
			object.etcdEncryption = value
			object.bitmap_ |= 33554432
		case "expiration_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.expirationTimestamp = value
			object.bitmap_ |= 67108864
		case "external_id":
			value := iterator.ReadString()
			object.externalID = value
			object.bitmap_ |= 134217728
		case "external_auth_config":
			value := ReadExternalAuthConfig(iterator)
			object.externalAuthConfig = value
			object.bitmap_ |= 268435456
		case "external_configuration":
			value := ReadExternalConfiguration(iterator)
			object.externalConfiguration = value
			object.bitmap_ |= 536870912
		case "flavour":
			value := ReadFlavour(iterator)
			object.flavour = value
			object.bitmap_ |= 1073741824
		case "groups":
			value := &GroupList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == GroupListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadGroupList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.groups = value
			object.bitmap_ |= 2147483648
		case "health_state":
			text := iterator.ReadString()
			value := ClusterHealthState(text)
			object.healthState = value
			object.bitmap_ |= 4294967296
		case "htpasswd":
			value := ReadHTPasswdIdentityProvider(iterator)
			object.htpasswd = value
			object.bitmap_ |= 8589934592
		case "hypershift":
			value := ReadHypershift(iterator)
			object.hypershift = value
			object.bitmap_ |= 17179869184
		case "identity_providers":
			value := &IdentityProviderList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == IdentityProviderListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadIdentityProviderList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.identityProviders = value
			object.bitmap_ |= 34359738368
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
			object.bitmap_ |= 68719476736
		case "infra_id":
			value := iterator.ReadString()
			object.infraID = value
			object.bitmap_ |= 137438953472
		case "ingresses":
			value := &IngressList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == IngressListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadIngressList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.ingresses = value
			object.bitmap_ |= 274877906944
		case "kubelet_config":
			value := ReadKubeletConfig(iterator)
			object.kubeletConfig = value
			object.bitmap_ |= 549755813888
		case "load_balancer_quota":
			value := iterator.ReadInt()
			object.loadBalancerQuota = value
			object.bitmap_ |= 1099511627776
		case "machine_pools":
			value := &MachinePoolList{}
			for {
				field := iterator.ReadObject()
				if field == "" {
					break
				}
				switch field {
				case "kind":
					text := iterator.ReadString()
					value.SetLink(text == MachinePoolListLinkKind)
				case "href":
					value.SetHREF(iterator.ReadString())
				case "items":
					value.SetItems(ReadMachinePoolList(iterator))
				default:
					iterator.ReadAny()
				}
			}
			object.machinePools = value
			object.bitmap_ |= 2199023255552
		case "managed":
			value := iterator.ReadBool()
			object.managed = value
			object.bitmap_ |= 4398046511104
		case "managed_service":
			value := ReadManagedService(iterator)
			object.managedService = value
			object.bitmap_ |= 8796093022208
		case "multi_az":
			value := iterator.ReadBool()
			object.multiAZ = value
			object.bitmap_ |= 17592186044416
		case "multi_arch_enabled":
			value := iterator.ReadBool()
			object.multiArchEnabled = value
			object.bitmap_ |= 35184372088832
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 70368744177664
		case "network":
			value := ReadNetwork(iterator)
			object.network = value
			object.bitmap_ |= 140737488355328
		case "node_drain_grace_period":
			value := ReadValue(iterator)
			object.nodeDrainGracePeriod = value
			object.bitmap_ |= 281474976710656
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
			object.bitmap_ |= 562949953421312
		case "nodes":
			value := ReadClusterNodes(iterator)
			object.nodes = value
			object.bitmap_ |= 1125899906842624
		case "openshift_version":
			value := iterator.ReadString()
			object.openshiftVersion = value
			object.bitmap_ |= 2251799813685248
		case "product":
			value := ReadProduct(iterator)
			object.product = value
			object.bitmap_ |= 4503599627370496
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
			object.bitmap_ |= 9007199254740992
		case "provision_shard":
			value := ReadProvisionShard(iterator)
			object.provisionShard = value
			object.bitmap_ |= 18014398509481984
		case "proxy":
			value := ReadProxy(iterator)
			object.proxy = value
			object.bitmap_ |= 36028797018963968
		case "region":
			value := ReadCloudRegion(iterator)
			object.region = value
			object.bitmap_ |= 72057594037927936
		case "registry_config":
			value := ReadClusterRegistryConfig(iterator)
			object.registryConfig = value
			object.bitmap_ |= 144115188075855872
		case "state":
			text := iterator.ReadString()
			value := ClusterState(text)
			object.state = value
			object.bitmap_ |= 288230376151711744
		case "status":
			value := ReadClusterStatus(iterator)
			object.status = value
			object.bitmap_ |= 576460752303423488
		case "storage_quota":
			value := ReadValue(iterator)
			object.storageQuota = value
			object.bitmap_ |= 1152921504606846976
		case "subscription":
			value := ReadSubscription(iterator)
			object.subscription = value
			object.bitmap_ |= 2305843009213693952
		case "version":
			value := ReadVersion(iterator)
			object.version = value
			object.bitmap_ |= 4611686018427387904
		default:
			iterator.ReadAny()
		}
	}
	return object
}
