/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Annotation constants
const (
	// ClusterIDLabel is the label that a machineset must have to identify the
	// cluster to which it belongs.
	ClusterIDLabel = "machine.openshift.io/cluster-api-cluster"
	// DefaultTenancy creates the instance on a non-dedicated host.
	DefaultTenancy InstanceTenancy = "default"
	// HostTenancy creates the instance on a dedicated host. If you do not specify the DedicatedHostID parameter, Alibaba Cloud automatically selects a dedicated host for the instance.
	HostTenancy InstanceTenancy = "host"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlibabaCloudMachineProviderConfig is the Schema for the alibabacloudmachineproviderconfig API
// +k8s:openapi-gen=true
type AlibabaCloudMachineProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// More detail about alibabacloud ECS
	// https://www.alibabacloud.com/help/doc-detail/25499.htm?spm=a2c63.l28256.b99.727.496d7453jF7Moz
	//The instance type of the instance.
	InstanceType string `json:"instanceType"`

	// The ID of th vpc
	// +optional
	VpcID string `json:"vpcId,omitempty"`

	// The ID of the region in which to create the instance. You can call the DescribeRegions operation to query the most recent region list.
	RegionID string `json:"regionId"`

	// The ID of the zone in which to create the instance. You can call the DescribeZones operation to query the most recent region list.
	ZoneID string `json:"zoneId"`

	// The ID of the image used to create the instance.
	ImageID string `json:"imageId"`

	// DataDisk hold information regarding the extra data disks
	DataDisks []DataDiskProperties `json:"dataDisk,omitempty"`

	// SecurityGroups is an array of references to security groups which to assign the instance. The valid values of N vary based on the
	// maximum number of security groups to which an instance can belong. For more information, see the "Security group limits" section in Limits.
	// https://www.alibabacloud.com/help/doc-detail/101348.htm?spm=a2c63.p38356.879954.48.78f0199aX3dfIE
	SecurityGroups []AlibabaResourceReference `json:"securityGroups,omitempty"`

	//Bandwidth describes the internet bandwidth strategy for the instance
	// +optional
	Bandwidth BandwidthProperties `json:"bandwidth,omitempty"`

	// Information regarding the the system disk for the instance
	// +optional
	SystemDisk SystemDiskProperties `json:"systemDisk,omitempty"`

	// Description of the instance. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
	// This parameter is empty by default.
	Description string `json:"description,omitempty"`

	// VSwitch is a reference to the vswitch to use for this instance
	// This parameter is required when you create an instance of the VPC type.
	// You can call the DescribeVSwitches operation to query the created vSwitches.
	VSwitch AlibabaResourceReference `json:"vSwitch,omitempty"`

	// Subscription holds the billing information for the instance
	// +optional
	Subscription SubscriptionInfo `json:"subscription,omitempty"`

	// RAMRoleName is the name of the instance Resource Access Management (RAM) role. You can call the ListRoles operation provided by RAM to query the instance RAM roles that you have created.
	// +optional
	RAMRoleName string `json:"ramRoleName,omitempty"`

	//S pecifies whether to enable security hardening. Valid values:
	// Active: enables security hardening. This value is applicable only to public images.
	// Deactive: does not enable security hardening. This value is applicable to all images.
	// +optional
	SecurityEnhancementStrategy string `json:"securityEnhancementStrategy,omitempty"`

	// ResourceGroupID is the unique ID of the resource group to which to assign the instance.
	// +optional
	ResourceGroupID string `json:"resourceGroupId,omitempty"`

	// DeletionProtection is release protection property of the instance. It specifies whether you can use the ECS console or call the DeleteInstance operation to manually release the instance. Default value: false. Valid values:
	// true: enables release protection.
	// false: disables release protection.
	// +optional
	DeletionProtection bool `json:"deletionProtection,omitempty"`

	// Affinity specifies whether to associate the instance on a dedicated host with the dedicated host. Valid values:
	// default: does not associate the instance with the dedicated host. When you restart an instance in the No Fees for Stopped Instances (VPC-Connected) state, the instance is automatically deployed to another dedicated host in the automatic deployment resource pool if the available resources of the original dedicated host are insufficient.
	// host: associates the instance with the dedicated host. When you restart an instance in the No Fees for Stopped Instances (VPC-Connected) state, the instance still resides on the original dedicated host. If the available resources of the original dedicated host are insufficient, the instance fails to restart.
	// Default value: default.
	// +optional
	Affinity string `json:"affinity,omitempty"`

	// Tenancy specifies whether to create the instance on a dedicated host. Valid values:
	// default: creates the instance on a non-dedicated host.
	// host: creates the instance on a dedicated host. If you do not specify the DedicatedHostID parameter, Alibaba Cloud automatically selects a dedicated host for the instance.
	// Default value: default.
	// +optional
	Tenancy InstanceTenancy `json:"tenancy,omitempty"`

	// DedicatedHostID is the ID of the dedicated host on which to create the instance.
	// You can call the DescribeDedicatedHosts operation to query the dedicated host list.
	// When the DedicatedHostID parameter is specified, the SpotStrategy and SpotPriceLimit parameters are ignored. This is because preemptible instances cannot be created on dedicated hosts.
	// +optional
	DedicatedHostID string `json:"dedicatedHostId"`

	// UserDataSecret contains a local reference to a secret that contains the
	// UserData to apply to the instance
	// +optional
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	// CredentialsSecret is a reference to the secret with alibabacloud credentials. Otherwise, defaults to permissions
	// provided by attached RAM role where the actuator is running.
	// +optional
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`

	// Tags are the set of metadata to add to an instance.
	// +optional
	Tags []Tag `json:"tags,omitempty"`
}

// ResourceTagReference is a reference to a specific AlibabaCloud resource by ID, or tags.
// Only one of ID or Tags may be specified. Specifying more than one will result in
// a validation error.
type AlibabaResourceReference struct {
	// ID of resource
	// +optional
	ID string `json:"id,omitempty"`

	// Tags is a set of metadata based upon ECS object tags used to identify a resource
	// +optional
	Tags []Tag `json:"tags,omitempty"`
}

// ResourceTagReference is a reference to a specific AlibabaCloud resource by ID, or tags.
// Only one of ID or Tags may be specified. Specifying more than one will result in
// a validation error.
type ResourceTagReference struct {
	// ID of resource
	// +optional
	ID string `json:"id,omitempty"`

	// Tags is a set of tags used to identify a resource
	// +optional
	Tags []Tag `json:"tags,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlibabaCloudMachineProviderConfigList contains a list of AlibabaCloudMachineProviderConfig
type AlibabaCloudMachineProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlibabaCloudMachineProviderConfig `json:"items"`
}

// InstanceTenancy Specifies whether to create the instance on a dedicated host
type InstanceTenancy string

type SystemDiskProperties struct {
	//The category of the system disk. Valid values:
	//cloud_essd: ESSD. When the parameter is set to this value, you can use the SystemDisk.PerformanceLevel parameter to specify the performance level of the disk.
	//cloud_efficiency: ultra disk.
	//cloud_ssd: standard SSD.
	//cloud: basic disk.
	//For non-I/O optimized instances of retired instance types, the default value is cloud. For other instances, the default value is cloud_efficiency.
	// +optional
	Category string `json:"category,omitempty"`

	//The performance level of the ESSD used as the system disk. Default value: PL1. Valid values:
	//PL0: A single ESSD can deliver up to 10,000 random read/write IOPS.
	//PL1: A single ESSD can deliver up to 50,000 random read/write IOPS.
	//PL2: A single ESSD can deliver up to 100,000 random read/write IOPS.
	//PL3: A single ESSD can deliver up to 1,000,000 random read/write IOPS.
	//For more information about ESSD performance levels, see ESSDs.
	// +optional
	PerformanceLevel string `json:"performanceLevel,omitempty"`

	//The name of the system disk. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with http:// or https://. It can contain letters, digits, colons (:), underscores (_), and hyphens (-).
	//This parameter is empty by default.
	// +optional
	Name string `json:"name,omitempty"`

	//The size of the system disk. Unit: GiB. Valid values: 20 to 500.
	//The value must be at least 20 and greater than or equal to the size of the image.
	//The default value is 40 or the size of the image, depending on whichever is greater.
	// +optional
	Size int `json:"size,omitempty"`

	//The description of the system disk. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
	//This parameter is empty by default.
	// +optional
	Description string `json:"description,omitempty"`
}

// DataDisk The datadisk of Instance
type DataDiskProperties struct {
	//The name of data disk N. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with http:// or https://. It can contain letters, digits, colons (:), underscores (_), and hyphens (-).
	//
	//This parameter is empty by default.
	// +optional
	Name string `name:"name,omitempty"`

	//The ID of the snapshot used to create data disk N. Valid values of N: 1 to 16.
	//
	//When the DataDisk.N.SnapshotID parameter is specified, the DataDisk.N.Size parameter is ignored. The data disk is created based on the size of the specified snapshot.
	//Use snapshots created after July 15, 2013. Otherwise, an error is returned and your request is rejected.
	// +optional
	SnapshotID string `name:"snapshotId,omitempty"`

	//The size of data disk N. Valid values of N: 1 to 16. Unit: GiB. Valid values:
	//
	//Valid values when DataDisk.N.Category is set to cloud_efficiency: 20 to 32768
	//Valid values when DataDisk.N.Category is set to cloud_ssd: 20 to 32768
	//Valid values when DataDisk.N.Category is set to cloud_essd: 20 to 32768
	//Valid values when DataDisk.N.Category is set to cloud: 5 to 2000
	//The value of this parameter must be greater than or equal to the size of the snapshot specified by the SnapshotID parameter.
	// +optional
	Size int `name:"size,omitempty"`

	//Specifies whether to encrypt data disk N.
	//
	//Default value: false.
	// +optional
	Encrypted bool `name:"encrypted,omitempty"`

	//
	//The performance level of the ESSD used as data disk N. The N value must be the same as that in DataDisk.N.Category when DataDisk.N.Category is set to cloud_essd. Default value: PL1. Valid values:
	//
	//PL0: A single ESSD can deliver up to 10,000 random read/write IOPS.
	//PL1: A single ESSD can deliver up to 50,000 random read/write IOPS.
	//PL2: A single ESSD can deliver up to 100,000 random read/write IOPS.
	//PL3: A single ESSD can deliver up to 1,000,000 random read/write IOPS.
	//For more information about ESSD performance levels, see ESSDs.
	// +optional
	PerformanceLevel string `name:"performanceLevel,omitempty"`

	//TODO
	//EncryptAlgorithm string `name:"EncryptAlgorithm"`

	//The description of data disk N. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
	//
	//This parameter is empty by default.
	// +optional
	Description string `name:"description,omitempty"`

	//The category of data disk N. Valid values:
	//
	//cloud_efficiency: ultra disk
	//cloud_ssd: standard SSD
	//cloud_essd: ESSD
	//cloud: basic disk
	//For I/O optimized instances, the default value is cloud_efficiency. For non-I/O optimized instances, the default value is cloud.
	// +optional
	Category string `name:"category,omitempty"`

	//The ID of the Key Management Service (KMS) key to be used by data disk N.
	// +optional
	KMSKeyID string `name:"kmsKeyId,omitempty"`

	//The mount point of data disk N.
	// +optional
	Device string `name:"device,omitempty"`

	//Specifies whether to release data disk N along with the instance.
	//
	//Default value: true.
	// +optional
	DeleteWithInstance *bool `name:"deleteWithInstance,omitempty"`
}

// Subscription information for the instance
type SubscriptionInfo struct {
	//The billing method of the instance. Default value: PostPaid. Valid values:
	//PrePaid: subscription. If you set this parameter to PrePaid, make sure that you have sufficient balance or credit within your account. Otherwise, an InvalidPayMethod error is returned.
	//PostPaid: pay-as-you-go.
	// +optional
	InstanceChargeType string `json:"instanceChargeType,omitempty"`

	//The subscription period of the instance. The unit is specified by the PeriodUnit parameter. This parameter is valid and required only when InstanceChargeType is set to PrePaid. If the DedicatedHostID parameter is specified, the subscription period of the instance cannot be longer than that of the dedicated host. Valid values:
	//Valid values when PeriodUnit is set to Month: 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 48, and 60.
	// +optional
	Period int `json:"period,omitempty"`

	//The unit of the subscription period. Default value: Month.
	//Set the value to Month.
	// +optional
	PeriodUnit string `json:"periodUnit,omitempty"`

	// The billing method for network usage. Default value: PayByTraffic. Valid values:
	// PayByBandwidth: pay-by-bandwidth
	// PayByTraffic: pay-by-traffic
	// +optional
	InternetChargeType string `json:"internetChargeType,omitempty"`

	//Specifies whether to enable auto-renewal for the instance. This parameter is valid only when the InstanceChargeType parameter is set to PrePaid. Default value: false. Valid values:
	//true: enables automatic renewal for the instance.
	//false: does not enable auto-renewal for the instance.
	// +optional
	AutoRenew bool `json:"autoRenew,omitempty"`

	//The auto-renewal period of the instance. This parameter is required when AutoRenew is set to true.
	//If PeriodUnit is set to Month, the valid values of the AutoRenewPeriod parameter are 1, 2, 3, 6, and 12.
	// +optional
	AutoRenewPeriod int `json:"autoRenewPeriod,omitempty"`
}

// Bandwidth describes the bandwidth strategy for the network of the instance
type BandwidthProperties struct {
	//The maximum inbound public bandwidth. Unit: Mbit/s. Valid values:
	//When the purchased outbound public bandwidth is less than or equal to 10 Mbit/s, the valid values of this parameter are 1 to 10, and the default value is 10.
	//When the purchased outbound public bandwidth is greater than 10, the valid values are 1 to the InternetMaxBandwidthOut value, and the default value is the InternetMaxBandwidthOut value.
	// +optional
	InternetMaxBandwidthIn int `json:"internetMaxBandwidthIn,omitempty"`

	//The maximum outbound public bandwidth. Unit: Mbit/s. Valid values: 0 to 100.
	//Default value: 0.
	// +optional
	InternetMaxBandwidthOut int `json:"internetMaxBandwidthOut,omitempty"`
}

// Tag  The tags of ECS Instance
type Tag struct {
	Value string `name:"value"`
	Key   string `name:"key"`
}

func init() {
	SchemeBuilder.Register(&AlibabaCloudMachineProviderConfig{}, &AlibabaCloudMachineProviderConfigList{})
}
