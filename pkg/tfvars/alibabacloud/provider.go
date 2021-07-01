package alibabacloud

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MachineProviderSpec is the type that will be embedded in a Machine.Spec.ProviderSpec field
// for an Alibaba Cloud virtual machine.
type MachineProviderSpec struct {
	// MachineProviderSpec: TODO AlibabaCLoud: Future support "cluster-api-provider-alibabacloud"
	metav1.TypeMeta    `json:",inline"`
	metav1.ObjectMeta  `json:"metadata,omitempty"`
	RegionID           string              `json:"RegionID,omitempty"`
	KeyPairName        string              `json:"KeyPairName,omitempty"`
	ResourceGroupID    string              `json:"ResourceGroupId,omitempty"`
	HostName           string              `json:"HostName,omitempty"`
	Password           string              `json:"Password,omitempty"`
	Tag                []*InstanceTag      `json:"Tag,omitempty"`
	VSwitchID          string              `json:"VSwitchId"`
	PrivateIPAddress   string              `json:"PrivateIpAddress,omitempty"`
	InstanceName       string              `json:"InstanceName,omitempty"`
	ZoneID             string              `json:"ZoneId"`
	ImageID            string              `json:"ImageId"`
	SecurityGroupID    string              `json:"SecurityGroupId"`
	SystemDiskCategory string              `json:"SystemDiskCategory"`
	UserData           string              `json:"UserData,omitempty"`
	InstanceType       string              `json:"InstanceType,omitempty"`
	RAMRoleName        string              `json:"RAMRoleName,omitempty"`
	DataDisk           []*InstanceDataDisk `json:"DataDisk,omitempty"`
	SystemDiskSize     string              `json:"SystemDiskSize,omitempty"`

	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`
}

// InstanceTag is the set of tags to add to apply to an instance
type InstanceTag struct {
	Value string `name:"Value,omitempty"`
	Key   string `name:"Key,omitempty"`
}

// InstanceDataDisk describes a data disk.
type InstanceDataDisk struct {
	DiskName           string `name:"DiskName,omitempty"`
	SnapshotID         string `name:"SnapshotId,omitempty"`
	Size               string `name:"Size,omitempty"`
	Encrypted          string `name:"Encrypted,omitempty"`
	PerformanceLevel   string `name:"PerformanceLevel,omitempty"`
	EncryptAlgorithm   string `name:"EncryptAlgorithm,omitempty"`
	Description        string `name:"Description,omitempty"`
	Category           string `name:"Category,omitempty"`
	KMSKeyID           string `name:"KMSKeyId,omitempty"`
	Device             string `name:"Device,omitempty"`
	DeleteWithInstance string `name:"DeleteWithInstance,omitempty"`
}
