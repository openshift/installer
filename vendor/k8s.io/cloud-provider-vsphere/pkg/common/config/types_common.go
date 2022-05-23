/*
Copyright 2020 The Kubernetes Authors.

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

package config

/*
	TODO:
	When the INI based cloud-config is deprecated. This file should be deleted and
	the structs in types_yaml.go will be renamed to replace the ones in this file.
*/

// Global struct
type Global struct {
	// vCenter username.
	User string
	// vCenter password in clear text.
	Password string
	// Deprecated. Use VirtualCenter to specify multiple vCenter Servers.
	// vCenter IP.
	VCenterIP string
	// vCenter port.
	VCenterPort string
	// True if vCenter uses self-signed cert.
	InsecureFlag bool
	// Datacenter in which VMs are located.
	Datacenters string
	// Soap round tripper count (retries = RoundTripper - 1)
	RoundTripperCount uint
	// Specifies the path to a CA certificate in PEM format. Optional; if not
	// configured, the system's CA certificates will be used.
	CAFile string
	// Thumbprint of the VCenter's certificate thumbprint
	Thumbprint string
	// Name of the secret were vCenter credentials are present.
	SecretName string
	// Secret Namespace where secret will be present that has vCenter credentials.
	SecretNamespace string
	// Secret directory in the event that:
	// 1) we don't want to use the k8s API to listen for changes to secrets
	// 2) we are not in a k8s env, namely DC/OS, since CSI is CO agnostic
	// Default: /etc/cloud/credentials
	SecretsDirectory string
	// Disable the vSphere CCM API
	// Default: true
	APIDisable bool
	// Configurable vSphere CCM API port
	// Default: 43001
	APIBinding string
}

// VirtualCenterConfig struct
type VirtualCenterConfig struct {
	// vCenter username.
	User string
	// vCenter password in clear text.
	Password string
	// TenantRef (intentionally not exposed via the config) is a unique tenant ref to
	// be used in place of the vcServer as the primary connection key. If one label is set,
	// all virtual center configs must have a unique label.
	TenantRef string
	// vCenterIP - If this field in the config is set, it is assumed then that value in [VirtualCenter "<value>"]
	// is now the TenantRef above and this field is the actual VCenterIP. Otherwise for backward
	// compatibility, the value by default is the IP or FQDN of the vCenter Server.
	VCenterIP string
	// vCenter port.
	VCenterPort string
	// True if vCenter uses self-signed cert.
	InsecureFlag bool
	// Datacenter in which VMs are located.
	Datacenters string
	// Soap round tripper count (retries = RoundTripper - 1)
	RoundTripperCount uint
	// Specifies the path to a CA certificate in PEM format. Optional; if not
	// configured, the system's CA certificates will be used.
	CAFile string
	// Thumbprint of the VCenter's certificate thumbprint
	Thumbprint string
	// SecretRef (intentionally not exposed via the config) is a key to identify which
	// InformerManager holds the secret
	SecretRef string
	// Name of the secret where vCenter credentials are present.
	SecretName string
	// Namespace where the secret will be present containing vCenter credentials.
	SecretNamespace string
	// IP Family enables the ability to support IPv4 or IPv6
	// Supported values are:
	// ipv4 - IPv4 addresses only (Default)
	// ipv6 - IPv6 addresses only
	IPFamilyPriority []string
}

// Labels struct
type Labels struct {
	// Zone describes a zone
	Zone string
	// Region describes a region
	Region string
}

// Config is used to read and store information from the cloud configuration file
type Config struct {
	// Global settings
	Global Global

	// Virtual Center configurations
	VirtualCenter map[string]*VirtualCenterConfig

	// Tag categories and tags which correspond to "built-in node labels: zones and region"
	Labels Labels
}
