package alibabacloud

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AlibabaCloudMachineProviderSpec is the Schema for the alibabacloudmachineproviderconfigs API
// AlibabaCloudMachineProviderSpec: TODO AlibabaCLoud: Future support "cluster-api-provider-alibabacloud"
type AlibabaCloudMachineProviderSpec struct {
	metav1.TypeMeta    `json:",inline"`
	metav1.ObjectMeta  `json:"metadata,omitempty"`
	RegionID           string             `json:"RegionID,omitempty"`
	KeyPairName        string             `json:"KeyPairName,omitempty"`
	ResourceGroupID    string             `json:"ResourceGroupId,omitempty"`
	HostName           string             `json:"HostName,omitempty"`
	Password           string             `json:"Password,omitempty"`
	Tag                []InstanceTag      `json:"Tag,omitempty"`
	VSwitchID          string             `json:"VSwitchId"`
	PrivateIPAddress   string             `json:"PrivateIpAddress,omitempty"`
	InstanceName       string             `json:"InstanceName,omitempty"`
	ZoneID             string             `json:"ZoneId"`
	ImageID            string             `json:"ImageId"`
	SecurityGroupID    string             `json:"SecurityGroupId"`
	SystemDiskCategory string             `json:"SystemDiskCategory"`
	UserData           string             `json:"UserData,omitempty"`
	InstanceType       string             `json:"InstanceType,omitempty"`
	RAMRoleName        string             `json:"RAMRoleName,omitempty"`
	DataDisk           []InstanceDataDisk `json:"DataDisk,omitempty"`
	SystemDiskSize     string             `json:"SystemDiskSize,omitempty"`

	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`
}

// InstanceTag is the set of tags to add to apply to an instance
type InstanceTag struct {
	Value string `name:"Value"`
	Key   string `name:"Key"`
}

// InstanceDataDisk describes a data disk.
type InstanceDataDisk struct {
	DiskName           string `name:"DiskName"`
	SnapshotID         string `name:"SnapshotId"`
	Size               string `name:"Size"`
	Encrypted          string `name:"Encrypted"`
	PerformanceLevel   string `name:"PerformanceLevel"`
	EncryptAlgorithm   string `name:"EncryptAlgorithm"`
	Description        string `name:"Description"`
	Category           string `name:"Category"`
	KMSKeyID           string `name:"KMSKeyId"`
	Device             string `name:"Device"`
	DeleteWithInstance string `name:"DeleteWithInstance"`
}
