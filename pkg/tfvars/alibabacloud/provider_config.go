package alibabacloud

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: AlibabaCLoud: Future support "cluster-api-provider-alibabacloud"
type AlibabacloudMachineProviderSpec struct {
	metav1.TypeMeta    `json:",inline"`
	metav1.ObjectMeta  `json:"metadata,omitempty"`
	RegionId           string             `json:"RegionId,omitempty"`
	KeyPairName        string             `json:"KeyPairName,omitempty"`
	ResourceGroupId    string             `json:"ResourceGroupId,omitempty"`
	HostName           string             `json:"HostName,omitempty"`
	Password           string             `json:"Password,omitempty"`
	Tag                []InstanceTag      `json:"Tag,omitempty"`
	VSwitchId          string             `json:"VSwitchId"`
	PrivateIpAddress   string             `json:"PrivateIpAddress,omitempty"`
	InstanceName       string             `json:"InstanceName,omitempty"`
	ZoneId             string             `json:"ZoneId"`
	ImageId            string             `json:"ImageId"`
	SecurityGroupId    string             `json:"SecurityGroupId"`
	SystemDiskCategory string             `json:"SystemDiskCategory"`
	UserData           string             `json:"UserData,omitempty"`
	InstanceType       string             `json:"InstanceType,omitempty"`
	RamRoleName        string             `json:"RamRoleName,omitempty"`
	DataDisk           []InstanceDataDisk `json:"DataDisk,omitempty"`
	SystemDiskSize     string             `json:"SystemDiskSize,omitempty"`

	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret,omitempty"`

	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret,omitempty"`
}

type InstanceTag struct {
	Value string `name:"value"`
	Key   string `name:"Key"`
}

type InstanceDataDisk struct {
	DiskName           string `name:"DiskName"`
	SnapshotId         string `name:"SnapshotId"`
	Size               string `name:"Size"`
	Encrypted          string `name:"Encrypted"`
	PerformanceLevel   string `name:"PerformanceLevel"`
	EncryptAlgorithm   string `name:"EncryptAlgorithm"`
	Description        string `name:"Description"`
	Category           string `name:"Category"`
	KMSKeyId           string `name:"KMSKeyId"`
	Device             string `name:"Device"`
	DeleteWithInstance string `name:"DeleteWithInstance"`
}
