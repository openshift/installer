package vsphere

import (
	"fmt"
	"strconv"
	"strings"

	ini "gopkg.in/gcfg.v1"
)

// This file contains type definition and conversion method for vsphere-cloud-provider INI config format
// Original code was taken from the vsphere-cloud-provider repository and copied here with as little changes as possible
// List of changes between type definition here and in the upstream:
// 	- Related structures collected into a single file, in cloud-provider-vsphere it split across two different modules
//  - 'TenantRef', `IPFamilyPriority` and `SecretRef` fields was removed,
//     since these fields are not exposed and not intended to come from the config
//  - `createConfig` method was altered for the `cpiConfigINI`
//  - since INI related types are not intended to be used outside of this module, their definitions were made private
//
// Sources:
//  - https://github.com/kubernetes/cloud-provider-vsphere/blob/release-1.25/pkg/cloudprovider/vsphere/config/types_ini_legacy.go
//  - https://github.com/kubernetes/cloud-provider-vsphere/blob/release-1.25/pkg/common/config/types_ini_legacy.go

// globalINI are global values
type globalINI struct {
	// vCenter username.
	User string `gcfg:"user"`
	// vCenter password in clear text.
	Password string `gcfg:"password"`
	// Deprecated. Use VirtualCenter to specify multiple vCenter Servers.
	// vCenter IP.
	VCenterIP string `gcfg:"server"`
	// vCenter port.
	VCenterPort string `gcfg:"port"`
	// True if vCenter uses self-signed cert.
	InsecureFlag bool `gcfg:"insecure-flag"`
	// Datacenter in which VMs are located.
	Datacenters string `gcfg:"datacenters"`
	// Soap round tripper count (retries = RoundTripper - 1)
	RoundTripperCount uint `gcfg:"soap-roundtrip-count"`
	// Specifies the path to a CA certificate in PEM format. Optional; if not
	// configured, the system's CA certificates will be used.
	CAFile string `gcfg:"ca-file"`
	// Thumbprint of the VCenter's certificate thumbprint
	Thumbprint string `gcfg:"thumbprint"`
	// Name of the secret were vCenter credentials are present.
	SecretName string `gcfg:"secret-name"`
	// Secret Namespace where secret will be present that has vCenter credentials.
	SecretNamespace string `gcfg:"secret-namespace"`
	// Secret directory in the event that:
	// 1) we don't want to use the k8s API to listen for changes to secrets
	// 2) we are not in a k8s env, namely DC/OS, since CSI is CO agnostic
	// Default: /etc/cloud/credentials
	SecretsDirectory string `gcfg:"secrets-directory"`
	// Disable the vSphere CCM API
	// Default: true
	APIDisable bool `gcfg:"api-disable"`
	// Configurable vSphere CCM API port
	// Default: 43001
	APIBinding string `gcfg:"api-binding"`
	// IP Family enables the ability to support IPv4 or IPv6
	// Supported values are:
	// ipv4 - IPv4 addresses only (Default)
	// ipv6 - IPv6 addresses only
	IPFamily string `gcfg:"ip-family"`
}

// virtualCenterConfigINI contains information used to access a remote vCenter
// endpoint.
type virtualCenterConfigINI struct {
	// vCenter username.
	User string `gcfg:"user"`
	// vCenter password in clear text.
	Password string `gcfg:"password"`
	// vCenterIP - If this field in the config is set, it is assumed then that value in [VirtualCenter "<value>"]
	// is now the TenantRef above and this field is the actual VCenterIP. Otherwise for backward
	// compatibility, the value by default is the IP or FQDN of the vCenter Server.
	VCenterIP string `gcfg:"server"`
	// vCenter port.
	VCenterPort string `gcfg:"port"`
	// True if vCenter uses self-signed cert.
	InsecureFlag bool `gcfg:"insecure-flag"`
	// Datacenter in which VMs are located.
	Datacenters string `gcfg:"datacenters"`
	// Soap round tripper count (retries = RoundTripper - 1)
	RoundTripperCount uint `gcfg:"soap-roundtrip-count"`
	// Specifies the path to a CA certificate in PEM format. Optional; if not
	// configured, the system's CA certificates will be used.
	CAFile string `gcfg:"ca-file"`
	// Thumbprint of the VCenter's certificate thumbprint
	Thumbprint string `gcfg:"thumbprint"`
	// Name of the secret where vCenter credentials are present.
	SecretName string `gcfg:"secret-name"`
	// Namespace where the secret will be present containing vCenter credentials.
	SecretNamespace string `gcfg:"secret-namespace"`
	// IP Family enables the ability to support IPv4 or IPv6
	// Supported values are:
	// ipv4 - IPv4 addresses only (Default)
	// ipv6 - IPv6 addresses only
	IPFamily string `gcfg:"ip-family"`
}

// labelsINI tags categories and tags which correspond to "built-in node labels: zones and region"
type labelsINI struct {
	Zone   string `gcfg:"zone"`
	Region string `gcfg:"region"`
}

// commonConfigINI is used to read and store information from the cloud configuration file
type commonConfigINI struct {
	// Global values...
	Global globalINI

	// Virtual Center configurations
	VirtualCenter map[string]*virtualCenterConfigINI

	// Tag categories and tags which correspond to "built-in node labels: zones and region"
	Labels labelsINI
}

// nodesINI captures internal/external networks
type nodesINI struct {
	// IP address on VirtualMachine's network interfaces included in the fields' CIDRs
	// that will be used in respective status.addresses fields.
	InternalNetworkSubnetCIDR string `gcfg:"internal-network-subnet-cidr"`
	ExternalNetworkSubnetCIDR string `gcfg:"external-network-subnet-cidr"`
	// IP address on VirtualMachine's VM Network names that will be used to when searching
	// for status.addresses fields. Note that if InternalNetworkSubnetCIDR and
	// ExternalNetworkSubnetCIDR are not set, then the vNIC associated to this network must
	// only have a single IP address assigned to it.
	InternalVMNetworkName string `gcfg:"internal-vm-network-name"`
	ExternalVMNetworkName string `gcfg:"external-vm-network-name"`
	// IP addresses in these subnet ranges will be excluded when selecting
	// the IP address from the VirtualMachine's VM for use in the
	// status.addresses fields.
	ExcludeInternalNetworkSubnetCIDR string `gcfg:"exclude-internal-network-subnet-cidr"`
	ExcludeExternalNetworkSubnetCIDR string `gcfg:"exclude-external-network-subnet-cidr"`
}

// cpiConfigINI is the INI representation
type cpiConfigINI struct {
	commonConfigINI
	Nodes nodesINI
}

// parseUIntOrZero parses string to uint, returns error for negative numbers
func parseUIntOrZero(s string) (uint, error) {
	var parsedInt int
	var err error
	if s == "" {
		parsedInt = 0
	} else {
		parsedInt, err = strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("can not parse vCenter port from ini config: %w", err)
		}
		if parsedInt < 0 {
			return 0, fmt.Errorf("parsed int bigger than zero")
		}
	}
	return uint(parsedInt), nil
}

// parseUIntOrZero parses datacenters string into a slice of strings
// example "DC0,DC1" -> []string{"DC0", "DC1"}
func splitDatacenters(datacentersString string) []string {
	splitted := strings.Split(datacentersString, ",")
	result := make([]string, 0)
	for _, dc := range splitted {
		if dc != "" {
			result = append(result, strings.Trim(dc, " "))
		}
	}
	return result
}

// createConfig creates CPIConfig instance which is ready to further YAML serialization
// from the intermediate type cpiConfigINI.
// Due to differences between CPIConfig and cpiConfigINI number of extra checks and conversions are happening here.
func (iniConfig *cpiConfigINI) createConfig() (*CPIConfig, error) {

	globalVcenterPort, err := parseUIntOrZero(iniConfig.Global.VCenterPort)
	if err != nil {
		return nil, fmt.Errorf("can not create CPIConfig, invalid global port parameter: %w", err)
	}

	cfg := &CPIConfig{
		CommonConfig: CommonConfig{
			Global: Global{
				User:              iniConfig.Global.User,
				Password:          iniConfig.Global.Password,
				VCenterIP:         iniConfig.Global.VCenterIP,
				VCenterPort:       globalVcenterPort,
				InsecureFlag:      iniConfig.Global.InsecureFlag,
				Datacenters:       splitDatacenters(iniConfig.Global.Datacenters),
				RoundTripperCount: iniConfig.Global.RoundTripperCount,
				CAFile:            iniConfig.Global.CAFile,
				Thumbprint:        iniConfig.Global.Thumbprint,
				SecretName:        iniConfig.Global.SecretName,
				SecretNamespace:   iniConfig.Global.SecretNamespace,
				SecretsDirectory:  iniConfig.Global.SecretsDirectory,
			},
			Vcenter: make(map[string]*VirtualCenterConfig),
		},
	}

	// Only create Labels if either Zone or Region is set
	if iniConfig.Labels.Zone != "" || iniConfig.Labels.Region != "" {
		cfg.Labels = &Labels{
			Zone:   iniConfig.Labels.Zone,
			Region: iniConfig.Labels.Region,
		}
	}

	// Only create Nodes if any field is set
	if iniConfig.Nodes.InternalNetworkSubnetCIDR != "" ||
		iniConfig.Nodes.ExternalNetworkSubnetCIDR != "" ||
		iniConfig.Nodes.InternalVMNetworkName != "" ||
		iniConfig.Nodes.ExternalVMNetworkName != "" ||
		iniConfig.Nodes.ExcludeInternalNetworkSubnetCIDR != "" ||
		iniConfig.Nodes.ExcludeExternalNetworkSubnetCIDR != "" {
		cfg.Nodes = &Nodes{
			InternalNetworkSubnetCIDR:        iniConfig.Nodes.InternalNetworkSubnetCIDR,
			ExternalNetworkSubnetCIDR:        iniConfig.Nodes.ExternalNetworkSubnetCIDR,
			InternalVMNetworkName:            iniConfig.Nodes.InternalVMNetworkName,
			ExternalVMNetworkName:            iniConfig.Nodes.ExternalVMNetworkName,
			ExcludeInternalNetworkSubnetCIDR: iniConfig.Nodes.ExcludeInternalNetworkSubnetCIDR,
			ExcludeExternalNetworkSubnetCIDR: iniConfig.Nodes.ExcludeExternalNetworkSubnetCIDR,
		}
	}

	for keyVcConfig, valVcConfig := range iniConfig.VirtualCenter {
		vcenterPort, err := parseUIntOrZero(valVcConfig.VCenterPort)
		if err != nil {
			return nil, fmt.Errorf("invalid port parameter for vc %s: %w", keyVcConfig, err)
		}

		// For YAML based config format VCenterIP is mandatory
		// If this field in the config in INI config is not set, it is assumed then that value in [VirtualCenter "<value>"]
		// section header
		vcenterIP := valVcConfig.VCenterIP
		if vcenterIP == "" {
			vcenterIP = keyVcConfig
		}

		ipFamilyPriority := []string{}
		if valVcConfig.IPFamily != "" {
			ipFamilyPriority = append(ipFamilyPriority, valVcConfig.IPFamily)
		}

		cfg.Vcenter[keyVcConfig] = &VirtualCenterConfig{
			User:              valVcConfig.User,
			Password:          valVcConfig.Password,
			VCenterIP:         vcenterIP,
			VCenterPort:       vcenterPort,
			InsecureFlag:      valVcConfig.InsecureFlag,
			Datacenters:       splitDatacenters(valVcConfig.Datacenters),
			RoundTripperCount: valVcConfig.RoundTripperCount,
			CAFile:            valVcConfig.CAFile,
			Thumbprint:        valVcConfig.Thumbprint,
			SecretName:        valVcConfig.SecretName,
			SecretNamespace:   valVcConfig.SecretNamespace,
			IPFamilyPriority:  ipFamilyPriority,
		}
	}

	return cfg, nil
}

// readCPIConfigINI parses vSphere cloud config file, stores it into cpiConfigINI immediately, and converts
// it into CPIConfig with the further return.
func readCPIConfigINI(byConfig []byte) (*CPIConfig, error) {
	if len(byConfig) == 0 {
		return nil, fmt.Errorf("empty INI file")
	}

	strConfig := string(byConfig[:])

	cfg := &cpiConfigINI{}

	if err := ini.FatalOnly(ini.ReadStringInto(cfg, strConfig)); err != nil {
		return nil, err
	}

	return cfg.createConfig()
}
