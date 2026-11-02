package vsphere

import (
	"fmt"

	"sigs.k8s.io/yaml"
)

// This file contains type definition fir vsphere-cloud-provider YAML config format
// Original code was taken from the vsphere-cloud-provider repository and copied here with as little changes as possible.
// Type definition below uses for serializing vsphere-cloud-provider into a yaml document.
// List of changes between type definition here and in the upstream:
// 	- Related structs collected into a single file, in cloud-provider-vsphere it split across two different modules
//  - 'TenantRef' and `SecretRef` fields was removed,
//     since these fields are not exposed and not intended to come from the config
//  - 'YAML' suffix was removed from struct definition names
//  - yaml related tags in struct definitions were removed,
//    since sigs.k8s.io/yaml falls back to json tags for YAML serialization
//
// Sources:
//  - https://github.com/kubernetes/cloud-provider-vsphere/blob/release-1.25/pkg/cloudprovider/vsphere/config/types_yaml.go
//  - https://github.com/kubernetes/cloud-provider-vsphere/blob/release-1.25/pkg/common/config/types_yaml.go

// Global are global values
type Global struct {
	// vCenter username.
	User string `json:"user,omitempty"`
	// vCenter password in clear text.
	Password string `json:"password,omitempty"`
	// Deprecated. Use VirtualCenter to specify multiple vCenter Servers.
	// vCenter IP.
	VCenterIP string `json:"server,omitempty"`
	// vCenter port.
	VCenterPort uint `json:"port,omitempty"`
	// True if vCenter uses self-signed cert.
	InsecureFlag bool `json:"insecureFlag,omitempty"`
	// Datacenter in which VMs are located.
	Datacenters []string `json:"datacenters,omitempty"`
	// Soap round tripper count (retries = RoundTripper - 1)
	RoundTripperCount uint `json:"soapRoundtripCount,omitempty"`
	// Specifies the path to a CA certificate in PEM format. Optional; if not
	// configured, the system's CA certificates will be used.
	CAFile string `json:"caFile,omitempty"`
	// Thumbprint of the VCenter's certificate thumbprint
	Thumbprint string `json:"thumbprint,omitempty"`
	// Name of the secret were vCenter credentials are present.
	SecretName string `json:"secretName,omitempty"`
	// Secret Namespace where secret will be present that has vCenter credentials.
	SecretNamespace string `json:"secretNamespace,omitempty"`
	// Secret directory in the event that:
	// 1) we don't want to use the k8s API to listen for changes to secrets
	// 2) we are not in a k8s env, namely DC/OS, since CSI is CO agnostic
	// Default: /etc/cloud/credentials
	SecretsDirectory string `json:"secretsDirectory,omitempty"`
	// Disable the vSphere CCM API
	// Default: true
	APIDisable bool `json:"apiDisable,omitempty"`
	// Configurable vSphere CCM API port
	// Default: 43001
	APIBinding string `json:"apiBinding,omitempty"`
	// IP Family enables the ability to support IPv4 or IPv6
	// Supported values are:
	// ipv4 - IPv4 addresses only (Default)
	// ipv6 - IPv6 addresses only
	IPFamilyPriority []string `json:"ipFamily,omitempty"`
}

// VirtualCenterConfig contains information used to access a remote vCenter
// endpoint.
type VirtualCenterConfig struct {
	// vCenter username.
	User string `json:"user,omitempty"`
	// vCenter password in clear text.
	Password string `json:"password,omitempty"`
	// vCenterIP - If this field in the config is set, it is assumed then that value in [VirtualCenter "<value>"]
	// is now the TenantRef above and this field is the actual VCenterIP. Otherwise for backward
	// compatibility, the value by default is the IP or FQDN of the vCenter Server.
	VCenterIP string `json:"server,omitempty"`
	// vCenter port.
	VCenterPort uint `json:"port,omitempty"`
	// True if vCenter uses self-signed cert.
	InsecureFlag bool `json:"insecureFlag,omitempty"`
	// Datacenter in which VMs are located.
	Datacenters []string `json:"datacenters,omitempty"`
	// Soap round tripper count (retries = RoundTripper - 1)
	RoundTripperCount uint `json:"soapRoundtripCount,omitempty"`
	// Specifies the path to a CA certificate in PEM format. Optional; if not
	// configured, the system's CA certificates will be used.
	CAFile string `json:"caFile,omitempty"`
	// Thumbprint of the VCenter's certificate thumbprint
	Thumbprint string `json:"thumbprint,omitempty"`
	// Name of the secret where vCenter credentials are present.
	SecretName string `json:"secretName,omitempty"`
	// Namespace where the secret will be present containing vCenter credentials.
	SecretNamespace string `json:"secretNamespace,omitempty"`
	// IP Family enables the ability to support IPv4 or IPv6
	// Supported values are:
	// ipv4 - IPv4 addresses only (Default)
	// ipv6 - IPv6 addresses only
	IPFamilyPriority []string `json:"ipFamily,omitempty"`
}

// Labels tags categories and tags which correspond to "built-in node labels: zones and region"
type Labels struct {
	Zone   string `json:"zone,omitempty"`
	Region string `json:"region,omitempty"`
}

// CommonConfig is used to read and store information from the cloud configuration file
type CommonConfig struct {
	// Global values...
	Global Global `json:"global,omitempty"`

	// Virtual Center configurations
	Vcenter map[string]*VirtualCenterConfig `json:"vcenter,omitempty"`

	// Tag categories and tags which correspond to "built-in node labels: zones and region"
	Labels *Labels `json:"labels,omitempty"`
}

// Nodes captures internal/external networks
type Nodes struct {
	// IP address on VirtualMachine's network interfaces included in the fields' CIDRs
	// that will be used in respective status.addresses fields.
	InternalNetworkSubnetCIDR string `json:"internalNetworkSubnetCidr,omitempty"`
	ExternalNetworkSubnetCIDR string `json:"externalNetworkSubnetCidr,omitempty"`
	// IP address on VirtualMachine's VM Network names that will be used to when searching
	// for status.addresses fields. Note that if InternalNetworkSubnetCIDR and
	// ExternalNetworkSubnetCIDR are not set, then the vNIC associated to this network must
	// only have a single IP address assigned to it.
	InternalVMNetworkName string `json:"internalVmNetworkName,omitempty"`
	ExternalVMNetworkName string `json:"externalVmNetworkName,omitempty"`
	// IP addresses in these subnet ranges will be excluded when selecting
	// the IP address from the VirtualMachine's VM for use in the
	// status.addresses fields.
	ExcludeInternalNetworkSubnetCIDR string `json:"excludeInternalNetworkSubnetCidr,omitempty"`
	ExcludeExternalNetworkSubnetCIDR string `json:"excludeExternalNetworkSubnetCidr,omitempty"`
}

// CPIConfig is the YAML representation of vsphere-cloud-provider config
type CPIConfig struct {
	CommonConfig `json:",inline"`
	Nodes        *Nodes `json:"nodes,omitempty"`
}

// readCPIConfigYAML parses vSphere cloud config file and stores it into CPIConfig
func readCPIConfigYAML(byConfig []byte) (*CPIConfig, error) {
	if len(byConfig) == 0 {
		return nil, fmt.Errorf("empty YAML file")
	}

	cfg := &CPIConfig{}

	if err := yaml.Unmarshal(byConfig, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
