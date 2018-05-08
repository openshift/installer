package libvirt

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
)

// Libvirt-specific configuration
type Libvirt struct {
	URI          string `json:"tectonic_libvirt_uri,omitempty" yaml:"uri"`
	SshKey       string `json:"tectonic_libvirt_ssh_key,omitempty" yaml:"sshKey"`
	QowImagePath string `json:"tectonic_coreos_qow_path,omitempty" yaml:"imagePath"`
	Network      `json:",inline" yaml:"network"`
	MasterIPs    []string `json:"tectonic_libvirt_master_ips,omitempty" yaml:"masterIps"`
}

type Network struct {
	Name      string `json:"tectonic_libvirt_network_name,omitempty" yaml:"name"`
	IfName    string `json:"tectonic_libvirt_network_if,omitempty" yaml"ifName"`
	DnsServer string `json:"tectonic_libvirt_resolver,omitempty" yaml:"dnsServer"`
	IpRange   string `json:"tectonic_libvirt_ip_range,omitempty" yaml:"ipRange"`
}

// Fill in any variables for terraform
func (l *Libvirt) TFVars(masterCount int) error {
	_, network, err := net.ParseCIDR(l.Network.IpRange)
	if err != nil {
		return fmt.Errorf("failed to parse libvirt.network.iprange: %v", err)
	}

	if len(l.MasterIPs) > 0 {
		if len(l.MasterIPs) != masterCount {
			return fmt.Errorf("length of MasterIPs does't match master count")
		} else {
			return nil
		}
	}

	for i := 0; i < masterCount; i++ {
		ip, err := cidr.Host(network, i+10)
		if err != nil {
			return fmt.Errorf("failed to generate masterips: %v", err)
		}
		l.MasterIPs = append(l.MasterIPs, ip.String())
	}

	return nil
}
