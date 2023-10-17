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

// cloudprovider_types contains API types for the vSphere cloud provider.
// The configuration may be unmarshalled from an INI-style configuration using
// the "gopkg.in/gcfg.v1" package.
//
// The configuration may be marshalled to an INI-style configuration using a Go
// template.
//
// The "gopkg.in/go-ini/ini.v1" package was investigated, but it does not
// support reflecting a struct with a field of type "map[string]TYPE" to INI.
//nolint:gocritic,godot
package v1alpha3

// CPIConfig is the vSphere cloud provider's configuration.
// DEPRECATED: will be removed in v1alpha4
type CPIConfig struct {
	// Global is the vSphere cloud provider's global configuration.
	// +optional
	Global CPIGlobalConfig `gcfg:"Global,omitempty" json:"global,omitempty"`

	// VCenter is a list of vCenter configurations.
	// +optional
	VCenter map[string]CPIVCenterConfig `gcfg:"VirtualCenter,omitempty" json:"virtualCenter,omitempty"`

	// Network is the vSphere cloud provider's network configuration.
	// +optional
	Network CPINetworkConfig `gcfg:"Network,omitempty" json:"network,omitempty"`

	// Disk is the vSphere cloud provider's disk configuration.
	// +optional
	Disk CPIDiskConfig `gcfg:"Disk,omitempty" json:"disk,omitempty"`

	// Workspace is the vSphere cloud provider's workspace configuration.
	// +optional
	Workspace CPIWorkspaceConfig `gcfg:"Workspace,omitempty" json:"workspace,omitempty"`

	// Labels is the vSphere cloud provider's zone and region configuration.
	// +optional
	Labels CPILabelConfig `gcfg:"Labels,omitempty" json:"labels,omitempty"`

	// CPIProviderConfig contains extra information used to configure the
	// vSphere cloud provider.
	ProviderConfig CPIProviderConfig `json:"providerConfig,omitempty"`
}

// CPIProviderConfig defines any extra information used to configure
// the vSphere external cloud provider
type CPIProviderConfig struct {
	Cloud   *CPICloudConfig   `json:"cloud,omitempty"`
	Storage *CPIStorageConfig `json:"storage,omitempty"`
}

type CPICloudConfig struct {
	ControllerImage string `json:"controllerImage,omitempty"`
	// ExtraArgs passes through extra arguments to the cloud provider.
	// The arguments here are passed to the cloud provider daemonset specification
	// +optional
	ExtraArgs map[string]string `json:"extraArgs,omitempty"`
}

type CPIStorageConfig struct {
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
	Global    CPIGlobalConfig              `gcfg:"Global,omitempty"`
	VCenter   map[string]*CPIVCenterConfig `gcfg:"VirtualCenter,omitempty"`
	Network   CPINetworkConfig             `gcfg:"Network,omitempty"`
	Disk      CPIDiskConfig                `gcfg:"Disk,omitempty"`
	Workspace CPIWorkspaceConfig           `gcfg:"Workspace,omitempty"`
	Labels    CPILabelConfig               `gcfg:"Labels,omitempty"`
}

// CPIGlobalConfig is the vSphere cloud provider's global configuration.
type CPIGlobalConfig struct {
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

// CPIVCenterConfig is a vSphere cloud provider's vCenter configuration.
type CPIVCenterConfig struct {
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

// CPINetworkConfig is the network configuration for the vSphere cloud provider.
type CPINetworkConfig struct {
	// Name is the name of the network to which VMs are connected.
	// +optional
	Name string `gcfg:"public-network,omitempty" json:"name,omitempty"`
}

// CPIDiskConfig defines the disk configuration for the vSphere cloud provider.
type CPIDiskConfig struct {
	// SCSIControllerType defines SCSI controller to be used.
	// +optional
	SCSIControllerType string `gcfg:"scsicontrollertype,omitempty" json:"scsiControllerType,omitempty"`
}

// CPIWorkspaceConfig defines a workspace configuration for the vSphere cloud
// provider.
type CPIWorkspaceConfig struct {
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

// CPILabelConfig defines the categories and tags which correspond to built-in
// node labels, zone and region.
type CPILabelConfig struct {
	// Zone is the zone in which VMs are created/located.
	// +optional
	Zone string `gcfg:"zone,omitempty" json:"zone,omitempty"`

	// Region is the region in which VMs are created/located.
	// +optional
	Region string `gcfg:"region,omitempty" json:"region,omitempty"`
}
