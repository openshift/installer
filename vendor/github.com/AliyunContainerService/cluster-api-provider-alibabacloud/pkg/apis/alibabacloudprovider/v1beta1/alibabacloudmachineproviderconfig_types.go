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
	VpcID string `json:"vpcId,omitempty"`

	// The ID of the region in which to create the instance. You can call the DescribeRegions operation to query the most recent region list.
	RegionID string `json:"regionId"`

	// The ID of the zone in which to create the instance. You can call the DescribeZones operation to query the most recent region list.
	ZoneID string `json:"zoneId"`

	// The ID of the image used to create the instance.
	ImageID string `json:"imageId"`

	// The ID of the security group to which to assign the instance. Instances in the same security group can communicate with each other.
	SecurityGroupID string `json:"securityGroupId"`

	// SecurityGroups is an array of references to security groups which to assign the instance. The valid values of N vary based on the
	// maximum number of security groups to which an instance can belong. For more information, see the "Security group limits" section in Limits.
	// https://www.alibabacloud.com/help/doc-detail/101348.htm?spm=a2c63.p38356.879954.48.78f0199aX3dfIE
	SecurityGroups []ResourceTagReference `json:"securityGroups,omitempty"`

	//The name of the instance. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with
	//http:// or https://. It can contain letters, digits, colons (:), underscores (_), periods (.), and hyphens (-). If you do not specify
	//this parameter, the instance ID is used as the instance name.
	InstanceName string `json:"instanceName,omitempty"`

	// The billing method for network usage. Default value: PayByTraffic. Valid values:
	// PayByBandwidth: pay-by-bandwidth
	// PayByTraffic: pay-by-traffic
	InternetChargeType string `json:"internetChargeType,omitempty"`

	//Specifies whether to enable auto-renewal for the instance. This parameter is valid only when the InstanceChargeType parameter is set to PrePaid. Default value: false. Valid values:
	//true: enables automatic renewal for the instance.
	//false: does not enable auto-renewal for the instance.
	AutoRenew bool `json:"autoRenew,omitempty"`

	//The auto-renewal period of the instance. This parameter is required when AutoRenew is set to true.
	//If PeriodUnit is set to Month, the valid values of the AutoRenewPeriod parameter are 1, 2, 3, 6, and 12.
	AutoRenewPeriod int `json:"autoRenewPeriod,omitempty"`

	//The maximum inbound public bandwidth. Unit: Mbit/s. Valid values:
	//When the purchased outbound public bandwidth is less than or equal to 10 Mbit/s, the valid values of this parameter are 1 to 10, and the default value is 10.
	//When the purchased outbound public bandwidth is greater than 10, the valid values are 1 to the InternetMaxBandwidthOut value, and the default value is the InternetMaxBandwidthOut value.
	InternetMaxBandwidthIn int `json:"internetMaxBandwidthIn,omitempty"`

	//The maximum outbound public bandwidth. Unit: Mbit/s. Valid values: 0 to 100.
	//Default value: 0.
	InternetMaxBandwidthOut int `json:"internetMaxBandwidthOut,omitempty"`

	//The hostname of the instance.
	//The hostname cannot start or end with a period (.) or a hyphen (-).It cannot contain consecutive periods (,) or hyphens (-).
	//For Windows instances, the hostname must be 2 to 15 characters in length and can contain letters, digits, and hyphens (-). It cannot contain periods (.) or contain only digits.
	//For an instance that runs one of other operating systems such as Linux, the hostname must be 2 to 64 characters in length. You can use periods (.) to separate the hostname into multiple segments. Each segment can contain letters, digits, and hyphens (-).
	HostName string `json:"hostName,omitempty"`

	//The password of the instance. The password must be 8 to 30 characters in length and include at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include:
	//( ) ` ~ ! @ # $ % ^ & * - _ + = | { } [ ] : ; ' < > , . ? /
	//Take note of the following items:
	//For security reasons, we recommend that you use HTTPS to send requests if the Password parameter is specified.
	//Passwords of Windows instances cannot start with a forward slash (/).
	//Passwords cannot be set for instances that run some types of operating systems such as Others Linux and Fedora CoreOS. For these instances, only key pairs can be set.
	Password string `json:"password,omitempty"`

	//The category of the system disk. Valid values:
	//cloud_essd: ESSD. When the parameter is set to this value, you can use the SystemDisk.PerformanceLevel parameter to specify the performance level of the disk.
	//cloud_efficiency: ultra disk.
	//cloud_ssd: standard SSD.
	//cloud: basic disk.
	//For non-I/O optimized instances of retired instance types, the default value is cloud. For other instances, the default value is cloud_efficiency.
	SystemDiskCategory string `json:"systemDiskCategory,omitempty"`

	//The performance level of the ESSD used as the system disk. Default value: PL1. Valid values:
	//PL0: A single ESSD can deliver up to 10,000 random read/write IOPS.
	//PL1: A single ESSD can deliver up to 50,000 random read/write IOPS.
	//PL2: A single ESSD can deliver up to 100,000 random read/write IOPS.
	//PL3: A single ESSD can deliver up to 1,000,000 random read/write IOPS.
	//For more information about ESSD performance levels, see ESSDs.
	SystemDiskPerformanceLevel string `json:"systemDiskPerformanceLevel,omitempty"`

	//The name of the system disk. The name must be 2 to 128 characters in length. It must start with a letter and cannot start with http:// or https://. It can contain letters, digits, colons (:), underscores (_), and hyphens (-).
	//This parameter is empty by default.
	SystemDiskDiskName string `json:"systemDiskDiskName,omitempty"`

	//The size of the system disk. Unit: GiB. Valid values: 20 to 500.
	//The value must be at least 20 and greater than or equal to the size of the image.
	//The default value is 40 or the size of the image, depending on whichever is greater.
	SystemDiskSize int `json:"systemDiskSize,omitempty"`

	//The description of the system disk. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
	//This parameter is empty by default.
	SystemDiskDescription string `json:"systemDiskDescription,omitempty"`

	DataDisks []DataDisk `json:"dataDisk,omitempty"`

	//The description of the instance. The description must be 2 to 256 characters in length and cannot start with http:// or https://.
	//This parameter is empty by default.
	Description string `json:"description,omitempty"`

	//The ID of the vSwitch to which to connect the instance. This parameter is required when you create an instance of the VPC type. You can call the DescribeVSwitches operation to query the created vSwitches.
	VSwitchID string `json:"vSwitchId,omitempty"`

	//VSwitch is a reference to the vswitch to use for this instance
	//This parameter is required when you create an instance of the VPC type.
	//You can call the DescribeVSwitches operation to query the created vSwitches.
	VSwitch *ResourceTagReference `json:"vSwitch,omitempty"`

	//Specifies whether the instance is I/O optimized. Valid values:
	//none: The instance is not I/O optimized.
	//optimized: The instance is I/O optimized.
	//For retired instance types, the default value is none. For more information, see Retired instance types.
	//For other instance types, the default value is optimized.
	IoOptimized string `json:"ioOptimized,omitempty"`

	//The billing method of the instance. Default value: PostPaid. Valid values:
	//PrePaid: subscription. If you set this parameter to PrePaid, make sure that you have sufficient balance or credit within your account. Otherwise, an InvalidPayMethod error is returned.
	//PostPaid: pay-as-you-go.
	InstanceChargeType string `json:"instanceChargeType,omitempty"`

	//The subscription period of the instance. The unit is specified by the PeriodUnit parameter. This parameter is valid and required only when InstanceChargeType is set to PrePaid. If the DedicatedHostID parameter is specified, the subscription period of the instance cannot be longer than that of the dedicated host. Valid values:
	//Valid values when PeriodUnit is set to Month: 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 48, and 60.
	Period int `json:"period,omitempty"`

	//The unit of the subscription period. Default value: Month.
	//Set the value to Month.
	PeriodUnit string `json:"periodUnit,omitempty"`

	Tags []Tag `json:"tags,omitempty"`

	//The user data of the instance. The user data must be encoded in Base64. The maximum size of raw data is 16 KB.
	UserData string `json:"userData,omitempty"`

	//The name of the key pair.
	//For Windows instances, this parameter is ignored and is empty by default. The Password parameter takes effect even if the KeyPairName parameter is specified.
	//For Linux instances, the password-based logon method is disabled by default. To make the instance more secure, we recommend that you use key pairs for logons.
	KeyPairName string `json:"keyPairName,omitempty"`

	//The name of the instance Resource Access Management (RAM) role. You can call the ListRoles operation provided by RAM to query the instance RAM roles that you have created.
	RAMRoleName string `json:"ramRoleName,omitempty"`

	//Specifies whether to enable security hardening. Valid values:
	//Active: enables security hardening. This value is applicable only to public images.
	//Deactive: does not enable security hardening. This value is applicable to all images.
	SecurityEnhancementStrategy string `json:"securityEnhancementStrategy,omitempty"`

	//The ID of the resource group to which to assign the instance.
	ResourceGroupID string `json:"resourceGroupId,omitempty"`

	//The release protection property of the instance. It specifies whether you can use the ECS console or call the DeleteInstance operation to manually release the instance. Default value: false. Valid values:
	//true: enables release protection.
	//false: disables release protection.
	DeletionProtection bool `json:"deletionProtection,omitempty"`

	//Specifies whether to associate the instance on a dedicated host with the dedicated host. Valid values:
	//default: does not associate the instance with the dedicated host. When you restart an instance in the No Fees for Stopped Instances (VPC-Connected) state, the instance is automatically deployed to another dedicated host in the automatic deployment resource pool if the available resources of the original dedicated host are insufficient.
	//host: associates the instance with the dedicated host. When you restart an instance in the No Fees for Stopped Instances (VPC-Connected) state, the instance still resides on the original dedicated host. If the available resources of the original dedicated host are insufficient, the instance fails to restart.
	//Default value: default.
	Affinity string `json:"affinity,omitempty"`

	//Specifies whether to create the instance on a dedicated host. Valid values:
	//default: creates the instance on a non-dedicated host.
	//host: creates the instance on a dedicated host. If you do not specify the DedicatedHostID parameter, Alibaba Cloud automatically selects a dedicated host for the instance.
	//Default value: default.
	Tenancy InstanceTenancy `json:"tenancy,omitempty"`

	//The ID of the dedicated host on which to create the instance.
	//You can call the DescribeDedicatedHosts operation to query the dedicated host list.
	//When the DedicatedHostID parameter is specified, the SpotStrategy and SpotPriceLimit parameters are ignored. This is because preemptible instances cannot be created on dedicated hosts.
	DedicatedHostID string `json:"dedicatedHostId"`

	// UserDataSecret contains a local reference to a secret that contains the
	// UserData to apply to the instance
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	// CredentialsSecret is a reference to the secret with alibabacloud credentials. Otherwise, defaults to permissions
	// provided by attached RAM role where the actuator is running.
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`
}

// ResourceTagReference is a reference to a specific AlibabaCloud resource by ID, or tags.
// Only one of ID or Tags may be specified. Specifying more than one will result in
// a validation error.
type ResourceTagReference struct {
	// ID of resource
	// +optional
	ID string `json:"id,omitempty"`

	// Tags is a set of tags used to identify a resource
	Tags []Tag `json:"tags,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AlibabaCloudMachineProviderConfigList contains a list of AlibabaCloudMachineProviderConfig
type AlibabaCloudMachineProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AlibabaCloudMachineProviderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AlibabaCloudMachineProviderConfig{}, &AlibabaCloudMachineProviderConfigList{})
}
