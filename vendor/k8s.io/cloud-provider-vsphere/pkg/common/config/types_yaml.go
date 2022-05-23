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
	When the INI based cloud-config is deprecated, this file should be renamed
	from types_yaml.go to types.go and the structs within this file should be named:

	GlobalYAML -> Global
	VirtualCenterConfigYAML -> VirtualCenterConfig
	LabelsYAML -> Labels
	ConfigYAML -> Config
*/

// GlobalYAML are global values
type GlobalYAML struct {
	// vCenter username.
	User string `yaml:"user"`
	// vCenter password in clear text.
	Password string `yaml:"password"`
	// Deprecated. Use VirtualCenter to specify multiple vCenter Servers.
	// vCenter IP.
	VCenterIP string `yaml:"server"`
	// vCenter port.
	VCenterPort uint `yaml:"port"`
	// True if vCenter uses self-signed cert.
	InsecureFlag bool `yaml:"insecureFlag"`
	// Datacenter in which VMs are located.
	Datacenters []string `yaml:"datacenters"`
	// Soap round tripper count (retries = RoundTripper - 1)
	RoundTripperCount uint `yaml:"soapRoundtripCount"`
	// Specifies the path to a CA certificate in PEM format. Optional; if not
	// configured, the system's CA certificates will be used.
	CAFile string `yaml:"caFile"`
	// Thumbprint of the VCenter's certificate thumbprint
	Thumbprint string `yaml:"thumbprint"`
	// Name of the secret were vCenter credentials are present.
	SecretName string `yaml:"secretName"`
	// Secret Namespace where secret will be present that has vCenter credentials.
	SecretNamespace string `yaml:"secretNamespace"`
	// Secret directory in the event that:
	// 1) we don't want to use the k8s API to listen for changes to secrets
	// 2) we are not in a k8s env, namely DC/OS, since CSI is CO agnostic
	// Default: /etc/cloud/credentials
	SecretsDirectory string `yaml:"secretsDirectory"`
	// Disable the vSphere CCM API
	// Default: true
	APIDisable bool `yaml:"apiDisable"`
	// Configurable vSphere CCM API port
	// Default: 43001
	APIBinding string `yaml:"apiBinding"`
	// IP Family enables the ability to support IPv4 or IPv6
	// Supported values are:
	// ipv4 - IPv4 addresses only (Default)
	// ipv6 - IPv6 addresses only
	IPFamilyPriority []string `yaml:"ipFamily"`
}

// VirtualCenterConfigYAML contains information used to access a remote vCenter
// endpoint.
type VirtualCenterConfigYAML struct {
	// vCenter username.
	User string `yaml:"user"`
	// vCenter password in clear text.
	Password string `yaml:"password"`
	// TenantRef (intentionally not exposed via the config) is a unique tenant ref to
	// be used in place of the vcServer as the primary connection key. If one label is set,
	// all virtual center configs must have a unique label.
	TenantRef string
	// vCenterIP - If this field in the config is set, it is assumed then that value in [VirtualCenter "<value>"]
	// is now the TenantRef above and this field is the actual VCenterIP. Otherwise for backward
	// compatibility, the value by default is the IP or FQDN of the vCenter Server.
	VCenterIP string `yaml:"server"`
	// vCenter port.
	VCenterPort uint `yaml:"port"`
	// True if vCenter uses self-signed cert.
	InsecureFlag bool `yaml:"insecureFlag"`
	// Datacenter in which VMs are located.
	Datacenters []string `yaml:"datacenters"`
	// Soap round tripper count (retries = RoundTripper - 1)
	RoundTripperCount uint `yaml:"soapRoundtripCount"`
	// Specifies the path to a CA certificate in PEM format. Optional; if not
	// configured, the system's CA certificates will be used.
	CAFile string `yaml:"caFile"`
	// Thumbprint of the VCenter's certificate thumbprint
	Thumbprint string `yaml:"thumbprint"`
	// SecretRef (intentionally not exposed via the config) is a key to identify which
	// InformerManager holds the secret
	SecretRef string
	// Name of the secret where vCenter credentials are present.
	SecretName string `yaml:"secretName"`
	// Namespace where the secret will be present containing vCenter credentials.
	SecretNamespace string `yaml:"secretNamespace"`
	// IP Family enables the ability to support IPv4 or IPv6
	// Supported values are:
	// ipv4 - IPv4 addresses only (Default)
	// ipv6 - IPv6 addresses only
	IPFamilyPriority []string `yaml:"ipFamily"`
}

// LabelsYAML tags categories and tags which correspond to "built-in node labels: zones and region"
type LabelsYAML struct {
	Zone   string `yaml:"zone"`
	Region string `yaml:"region"`
}

// CommonConfigYAML is used to read and store information from the cloud configuration file
type CommonConfigYAML struct {
	// Global values...
	Global GlobalYAML

	// Virtual Center configurations
	Vcenter map[string]*VirtualCenterConfigYAML

	// Tag categories and tags which correspond to "built-in node labels: zones and region"
	Labels LabelsYAML
}
