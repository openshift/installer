/*
Copyright 2019 The Kubernetes Authors.

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

// Package cloudprovider contains API types for the vSphere cloud provider.
package cloudprovider

// Config is the vSphere cloud provider's configuration.
type Config struct {
	// Global is the vSphere cloud provider's global configuration.
	// +optional
	Global GlobalConfig `gcfg:"Global,omitempty" json:"global,omitempty"`

	// VCenter is a list of vCenter configurations.
	// +optional
	VCenter map[string]VCenterConfig `gcfg:"VirtualCenter,omitempty" json:"virtualCenter,omitempty"`

	// Network is the vSphere cloud provider's network configuration.
	// +optional
	Network NetworkConfig `gcfg:"Network,omitempty" json:"network,omitempty"`

	// Disk is the vSphere cloud provider's disk configuration.
	// +optional
	Disk DiskConfig `gcfg:"Disk,omitempty" json:"disk,omitempty"`

	// Workspace is the vSphere cloud provider's workspace configuration.
	// +optional
	Workspace WorkspaceConfig `gcfg:"Workspace,omitempty" json:"workspace,omitempty"`

	// Labels is the vSphere cloud provider's zone and region configuration.
	// +optional
	Labels LabelConfig `gcfg:"Labels,omitempty" json:"labels,omitempty"`

	// ProviderConfig contains extra information used to configure the
	// vSphere cloud provider.
	ProviderConfig ProviderConfig `json:"providerConfig,omitempty"`
}

// ProviderConfig defines any extra information used to configure
// the vSphere external cloud provider
type ProviderConfig struct {
	Cloud   *CloudConfig   `json:"cloud,omitempty"`
	Storage *StorageConfig `json:"storage,omitempty"`
}

type CloudConfig struct {
	ControllerImage string `json:"controllerImage,omitempty"`
}

type StorageConfig struct {
	ControllerImage     string `json:"controllerImage,omitempty"`
	NodeDriverImage     string `json:"nodeDriverImage,omitempty"`
	AttacherImage       string `json:"attacherImage,omitempty"`
	ProvisionerImage    string `json:"provisionerImage,omitempty"`
	MetadataSyncerImage string `json:"metadataSyncerImage,omitempty"`
	LivenessProbeImage  string `json:"livenessProbeImage,omitempty"`
	RegistrarImage      string `json:"registrarImage,omitempty"`
}

// unmarshallableConfig is used to unmarshal the INI data using the gcfg
// package. The package requires fields with map types use *Values. However,
// kubebuilder v2 won't generate CRDs for map types with *Values.
type unmarshallableConfig struct {
	Global    GlobalConfig              `gcfg:"Global,omitempty"`
	VCenter   map[string]*VCenterConfig `gcfg:"VirtualCenter,omitempty"`
	Network   NetworkConfig             `gcfg:"Network,omitempty"`
	Disk      DiskConfig                `gcfg:"Disk,omitempty"`
	Workspace WorkspaceConfig           `gcfg:"Workspace,omitempty"`
	Labels    LabelConfig               `gcfg:"Labels,omitempty"`
}

// GlobalConfig is the vSphere cloud provider's global configuration.
type GlobalConfig struct {
	// Insecure is a flag that disables TLS peer verification.
	// +optional
	Insecure bool `gcfg:"insecure-flag,omitempty" json:"insecure,omitempty"`

	// RoundTripperCount specifies the SOAP round tripper count
	// (retries = RoundTripper - 1)
	// +optional
	RoundTripperCount int32 `gcfg:"soap-roundtrip-count,omitempty" json:"roundTripperCount,omitempty"`

	// Username is the username used to access a vSphere endpoint.
	// +optional
	Username string `gcfg:"user,omitempty" json:"username,omitempty"`

	// Password is the password used to access a vSphere endpoint.
	// +optional
	Password string `gcfg:"password,omitempty" json:"password,omitempty"`

	// SecretName is the name of the Kubernetes secret in which the vSphere
	// credentials are located.
	// +optional
	SecretName string `gcfg:"secret-name,omitempty" json:"secretName,omitempty"`

	// SecretNamespace is the namespace for SecretName.
	// +optional
	SecretNamespace string `gcfg:"secret-namespace,omitempty" json:"secretNamespace,omitempty"`

	// Port is the port on which the vSphere endpoint is listening.
	// Defaults to 443.
	// +optional
	Port string `gcfg:"port,omitempty" json:"port,omitempty"`

	// CAFile Specifies the path to a CA certificate in PEM format.
	// If not configured, the system's CA certificates will be used.
	// +optional
	CAFile string `gcfg:"ca-file,omitempty" json:"caFile,omitempty"`

	// Thumbprint is the cryptographic thumbprint of the vSphere endpoint's
	// certificate.
	// +optional
	Thumbprint string `gcfg:"thumbprint,omitempty" json:"thumbprint,omitempty"`

	// Datacenters is a CSV string of the datacenters in which VMs are located.
	// +optional
	Datacenters string `gcfg:"datacenters,omitempty" json:"datacenters,omitempty"`

	// ServiceAccount is the Kubernetes service account used to launch the cloud
	// controller manager.
	// Defaults to cloud-controller-manager.
	// +optional
	ServiceAccount string `gcfg:"service-account,omitempty" json:"serviceAccount,omitempty"`

	// SecretsDirectory is a directory in which secrets may be found. This
	// may used in the event that:
	// 1. It is not desirable to use the K8s API to watch changes to secrets
	// 2. The cloud controller manager is not running in a K8s environment,
	//    such as DC/OS. For example, the container storage interface (CSI) is
	//    container orcehstrator (CO) agnostic, and should support non-K8s COs.
	// Defaults to /etc/cloud/credentials.
	// +optional
	SecretsDirectory string `gcfg:"secrets-directory,omitempty" json:"secretsDirectory,omitempty"`

	// APIDisable disables the vSphere cloud controller manager API.
	// Defaults to true.
	// +optional
	APIDisable *bool `gcfg:"api-disable,omitempty" json:"apiDisable,omitempty"`

	// APIBindPort configures the vSphere cloud controller manager API port.
	// Defaults to 43001.
	// +optional
	APIBindPort string `gcfg:"api-binding,omitempty" json:"apiBindPort,omitempty"`

	// ClusterID is a unique identifier for a cluster used by the vSphere CSI driver (CNS)
	// NOTE: This field is set internally by CAPV and should not be set by any other consumer of this API
	ClusterID string `gcfg:"cluster-id,omitempty" json:"-"`
}

// VCenterConfig is a vSphere cloud provider's vCenter configuration.
type VCenterConfig struct {
	// Username is the username used to access a vSphere endpoint.
	// +optional
	Username string `gcfg:"user,omitempty" json:"username,omitempty"`

	// Password is the password used to access a vSphere endpoint.
	// +optional
	Password string `gcfg:"password,omitempty" json:"password,omitempty"`

	// Port is the port on which the vSphere endpoint is listening.
	// Defaults to 443.
	// +optional
	Port string `gcfg:"port,omitempty" json:"port,omitempty"`

	// Datacenters is a CSV string of the datacenters in which VMs are located.
	// +optional
	Datacenters string `gcfg:"datacenters,omitempty" json:"datacenters,omitempty"`

	// RoundTripperCount specifies the SOAP round tripper count
	// (retries = RoundTripper - 1)
	// +optional
	RoundTripperCount int32 `gcfg:"soap-roundtrip-count,omitempty" json:"roundTripperCount,omitempty"`

	// Thumbprint is the cryptographic thumbprint of the vSphere endpoint's
	// certificate.
	// +optional
	Thumbprint string `gcfg:"thumbprint,omitempty" json:"thumbprint,omitempty"`
}

// NetworkConfig is the network configuration for the vSphere cloud provider.
type NetworkConfig struct {
	// Name is the name of the network to which VMs are connected.
	// +optional
	Name string `gcfg:"public-network,omitempty" json:"name,omitempty"`
}

// DiskConfig defines the disk configuration for the vSphere cloud provider.
type DiskConfig struct {
	// SCSIControllerType defines SCSI controller to be used.
	// +optional
	SCSIControllerType string `gcfg:"scsicontrollertype,omitempty" json:"scsiControllerType,omitempty"`
}

// WorkspaceConfig defines a workspace configuration for the vSphere cloud
// provider.
type WorkspaceConfig struct {
	// Server is the IP address or FQDN of the vSphere endpoint.
	// +optional
	Server string `gcfg:"server,omitempty" json:"server,omitempty"`

	// Datacenter is the datacenter in which VMs are created/located.
	// +optional
	Datacenter string `gcfg:"datacenter,omitempty" json:"datacenter,omitempty"`

	// Folder is the folder in which VMs are created/located.
	// +optional
	Folder string `gcfg:"folder,omitempty" json:"folder,omitempty"`

	// Datastore is the datastore in which VMs are created/located.
	// +optional
	Datastore string `gcfg:"default-datastore,omitempty" json:"datastore,omitempty"`

	// ResourcePool is the resource pool in which VMs are created/located.
	// +optional
	ResourcePool string `gcfg:"resourcepool-path,omitempty" json:"resourcePool,omitempty"`
}

// LabelConfig defines the categories and tags which correspond to built-in
// node labels, zone and region.
type LabelConfig struct {
	// Zone is the zone in which VMs are created/located.
	// +optional
	Zone string `gcfg:"zone,omitempty" json:"zone,omitempty"`

	// Region is the region in which VMs are created/located.
	// +optional
	Region string `gcfg:"region,omitempty" json:"region,omitempty"`
}
