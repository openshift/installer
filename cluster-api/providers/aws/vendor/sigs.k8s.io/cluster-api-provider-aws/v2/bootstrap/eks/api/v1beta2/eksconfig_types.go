/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
)

// EKSConfigSpec defines the desired state of Amazon EKS Bootstrap Configuration.
type EKSConfigSpec struct {
	// KubeletExtraArgs passes the specified kubelet args into the Amazon EKS machine bootstrap script
	// +optional
	KubeletExtraArgs map[string]string `json:"kubeletExtraArgs,omitempty"`
	// ContainerRuntime specify the container runtime to use when bootstrapping EKS.
	// +optional
	ContainerRuntime *string `json:"containerRuntime,omitempty"`
	//  DNSClusterIP overrides the IP address to use for DNS queries within the cluster.
	// +optional
	DNSClusterIP *string `json:"dnsClusterIP,omitempty"`
	// DockerConfigJson is used for the contents of the /etc/docker/daemon.json file. Useful if you want a custom config differing from the default one in the AMI.
	// This is expected to be a json string.
	// +optional
	DockerConfigJSON *string `json:"dockerConfigJson,omitempty"`
	// APIRetryAttempts is the number of retry attempts for AWS API call.
	// +optional
	APIRetryAttempts *int `json:"apiRetryAttempts,omitempty"`
	// PauseContainer allows customization of the pause container to use.
	// +optional
	PauseContainer *PauseContainer `json:"pauseContainer,omitempty"`
	// UseMaxPods  sets --max-pods for the kubelet when true.
	// +optional
	UseMaxPods *bool `json:"useMaxPods,omitempty"`
	// ServiceIPV6Cidr is the ipv6 cidr range of the cluster. If this is specified then
	// the ip family will be set to ipv6.
	// +optional
	ServiceIPV6Cidr *string `json:"serviceIPV6Cidr,omitempty"`
	// PreBootstrapCommands specifies extra commands to run before bootstrapping nodes to the cluster
	// +optional
	PreBootstrapCommands []string `json:"preBootstrapCommands,omitempty"`
	// PostBootstrapCommands specifies extra commands to run after bootstrapping nodes to the cluster
	// +optional
	PostBootstrapCommands []string `json:"postBootstrapCommands,omitempty"`
	// BootstrapCommandOverride allows you to override the bootstrap command to use for EKS nodes.
	// +optional
	BootstrapCommandOverride *string `json:"boostrapCommandOverride,omitempty"`
	// Files specifies extra files to be passed to user_data upon creation.
	// +optional
	Files []File `json:"files,omitempty"`
	// DiskSetup specifies options for the creation of partition tables and file systems on devices.
	// +optional
	DiskSetup *DiskSetup `json:"diskSetup,omitempty"`
	// Mounts specifies a list of mount points to be setup.
	// +optional
	Mounts []MountPoints `json:"mounts,omitempty"`
	// Users specifies extra users to add
	// +optional
	Users []User `json:"users,omitempty"`
	// NTP specifies NTP configuration
	// +optional
	NTP *NTP `json:"ntp,omitempty"`
}

// PauseContainer contains details of pause container.
type PauseContainer struct {
	//  AccountNumber is the AWS account number to pull the pause container from.
	AccountNumber string `json:"accountNumber"`
	// Version is the tag of the pause container to use.
	Version string `json:"version"`
}

// EKSConfigStatus defines the observed state of the Amazon EKS Bootstrap Configuration.
type EKSConfigStatus struct {
	// Ready indicates the BootstrapData secret is ready to be consumed
	Ready bool `json:"ready,omitempty"`

	// DataSecretName is the name of the secret that stores the bootstrap data script.
	// +optional
	DataSecretName *string `json:"dataSecretName,omitempty"`

	// FailureReason will be set on non-retryable errors
	// +optional
	FailureReason string `json:"failureReason,omitempty"`

	// FailureMessage will be set on non-retryable errors
	// +optional
	FailureMessage string `json:"failureMessage,omitempty"`

	// ObservedGeneration is the latest generation observed by the controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// Conditions defines current service state of the EKSConfig.
	// +optional
	Conditions clusterv1beta1.Conditions `json:"conditions,omitempty"`
}

// Encoding specifies the cloud-init file encoding.
// +kubebuilder:validation:Enum=base64;gzip;gzip+base64
type Encoding string

const (
	// Base64 implies the contents of the file are encoded as base64.
	Base64 Encoding = "base64"
	// Gzip implies the contents of the file are encoded with gzip.
	Gzip Encoding = "gzip"
	// GzipBase64 implies the contents of the file are first base64 encoded and then gzip encoded.
	GzipBase64 Encoding = "gzip+base64"
)

// File defines the input for generating write_files in cloud-init.
type File struct {
	// Path specifies the full path on disk where to store the file.
	Path string `json:"path"`

	// Owner specifies the ownership of the file, e.g. "root:root".
	// +optional
	Owner string `json:"owner,omitempty"`

	// Permissions specifies the permissions to assign to the file, e.g. "0640".
	// +optional
	Permissions string `json:"permissions,omitempty"`

	// Encoding specifies the encoding of the file contents.
	// +optional
	Encoding Encoding `json:"encoding,omitempty"`

	// Append specifies whether to append Content to existing file if Path exists.
	// +optional
	Append bool `json:"append,omitempty"`

	// Content is the actual content of the file.
	// +optional
	Content string `json:"content,omitempty"`

	// ContentFrom is a referenced source of content to populate the file.
	// +optional
	ContentFrom *FileSource `json:"contentFrom,omitempty"`
}

// FileSource is a union of all possible external source types for file data.
// Only one field may be populated in any given instance. Developers adding new
// sources of data for target systems should add them here.
type FileSource struct {
	// Secret represents a secret that should populate this file.
	Secret SecretFileSource `json:"secret"`
}

// SecretFileSource adapts a Secret into a FileSource.
//
// The contents of the target Secret's Data field will be presented
// as files using the keys in the Data field as the file names.
type SecretFileSource struct {
	// Name of the secret in the KubeadmBootstrapConfig's namespace to use.
	Name string `json:"name"`

	// Key is the key in the secret's data map for this value.
	Key string `json:"key"`
}

// PasswdSource is a union of all possible external source types for passwd data.
// Only one field may be populated in any given instance. Developers adding new
// sources of data for target systems should add them here.
type PasswdSource struct {
	// Secret represents a secret that should populate this password.
	Secret SecretPasswdSource `json:"secret"`
}

// SecretPasswdSource adapts a Secret into a PasswdSource.
//
// The contents of the target Secret's Data field will be presented
// as passwd using the keys in the Data field as the file names.
type SecretPasswdSource struct {
	// Name of the secret in the KubeadmBootstrapConfig's namespace to use.
	Name string `json:"name"`

	// Key is the key in the secret's data map for this value.
	Key string `json:"key"`
}

// User defines the input for a generated user in cloud-init.
type User struct {
	// Name specifies the username
	Name string `json:"name"`

	// Gecos specifies the gecos to use for the user
	// +optional
	Gecos *string `json:"gecos,omitempty"`

	// Groups specifies the additional groups for the user
	// +optional
	Groups *string `json:"groups,omitempty"`

	// HomeDir specifies the home directory to use for the user
	// +optional
	HomeDir *string `json:"homeDir,omitempty"`

	// Inactive specifies whether to mark the user as inactive
	// +optional
	Inactive *bool `json:"inactive,omitempty"`

	// Shell specifies the user's shell
	// +optional
	Shell *string `json:"shell,omitempty"`

	// Passwd specifies a hashed password for the user
	// +optional
	Passwd *string `json:"passwd,omitempty"`

	// PasswdFrom is a referenced source of passwd to populate the passwd.
	// +optional
	PasswdFrom *PasswdSource `json:"passwdFrom,omitempty"`

	// PrimaryGroup specifies the primary group for the user
	// +optional
	PrimaryGroup *string `json:"primaryGroup,omitempty"`

	// LockPassword specifies if password login should be disabled
	// +optional
	LockPassword *bool `json:"lockPassword,omitempty"`

	// Sudo specifies a sudo role for the user
	// +optional
	Sudo *string `json:"sudo,omitempty"`

	// SSHAuthorizedKeys specifies a list of ssh authorized keys for the user
	// +optional
	SSHAuthorizedKeys []string `json:"sshAuthorizedKeys,omitempty"`
}

// NTP defines input for generated ntp in cloud-init.
type NTP struct {
	// Servers specifies which NTP servers to use
	// +optional
	Servers []string `json:"servers,omitempty"`

	// Enabled specifies whether NTP should be enabled
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
}

// DiskSetup defines input for generated disk_setup and fs_setup in cloud-init.
type DiskSetup struct {
	// Partitions specifies the list of the partitions to setup.
	// +optional
	Partitions []Partition `json:"partitions,omitempty"`

	// Filesystems specifies the list of file systems to setup.
	// +optional
	Filesystems []Filesystem `json:"filesystems,omitempty"`
}

// Partition defines how to create and layout a partition.
type Partition struct {
	// Device is the name of the device.
	Device string `json:"device"`
	// Layout specifies the device layout.
	// If it is true, a single partition will be created for the entire device.
	// When layout is false, it means don't partition or ignore existing partitioning.
	Layout bool `json:"layout"`
	// Overwrite describes whether to skip checks and create the partition if a partition or filesystem is found on the device.
	// Use with caution. Default is 'false'.
	// +optional
	Overwrite *bool `json:"overwrite,omitempty"`
	// TableType specifies the tupe of partition table. The following are supported:
	// 'mbr': default and setups a MS-DOS partition table
	// 'gpt': setups a GPT partition table
	// +optional
	TableType *string `json:"tableType,omitempty"`
}

// Filesystem defines the file systems to be created.
type Filesystem struct {
	// Device specifies the device name
	Device string `json:"device"`
	// Filesystem specifies the file system type.
	Filesystem string `json:"filesystem"`
	// Label specifies the file system label to be used. If set to None, no label is used.
	Label string `json:"label"`
	// Partition specifies the partition to use. The valid options are: "auto|any", "auto", "any", "none", and <NUM>, where NUM is the actual partition number.
	// +optional
	Partition *string `json:"partition,omitempty"`
	// Overwrite defines whether or not to overwrite any existing filesystem.
	// If true, any pre-existing file system will be destroyed. Use with Caution.
	// +optional
	Overwrite *bool `json:"overwrite,omitempty"`
	// ExtraOpts defined extra options to add to the command for creating the file system.
	// +optional
	ExtraOpts []string `json:"extraOpts,omitempty"`
}

// MountPoints defines input for generated mounts in cloud-init.
type MountPoints []string

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=eksconfigs,scope=Namespaced,categories=cluster-api,shortName=eksc
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Bootstrap configuration is ready"
// +kubebuilder:printcolumn:name="DataSecretName",type="string",JSONPath=".status.dataSecretName",description="Name of Secret containing bootstrap data"

// EKSConfig is the schema for the Amazon EKS Machine Bootstrap Configuration API.
type EKSConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EKSConfigSpec   `json:"spec,omitempty"`
	Status EKSConfigStatus `json:"status,omitempty"`
}

// GetConditions returns the observations of the operational state of the EKSConfig resource.
func (r *EKSConfig) GetConditions() clusterv1beta1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the EKSConfig to the predescribed clusterv1beta1.Conditions.
func (r *EKSConfig) SetConditions(conditions clusterv1beta1.Conditions) {
	r.Status.Conditions = conditions
}

// +kubebuilder:object:root=true

// EKSConfigList contains a list of EKSConfig.
type EKSConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EKSConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EKSConfig{}, &EKSConfigList{})
}
