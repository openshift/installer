package libvirt

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
)

const (
	// DefaultDNSServer is the default DNS server for libvirt.
	DefaultDNSServer = "8.8.8.8"
	// DefaultIfName is the default interface name for libvirt.
	DefaultIfName = "osbr0"
)

// Libvirt-specific configuration
type Libvirt struct {
	URI           string `json:"tectonic_libvirt_uri,omitempty" yaml:"uri"`
	SSHKey        string `json:"tectonic_libvirt_ssh_key,omitempty" yaml:"sshKey"`
	QCOWImagePath string `json:"tectonic_coreos_qcow_path,omitempty" yaml:"imagePath"`
	Network       `json:",inline" yaml:"network"`
	MasterIPs     []string `json:"tectonic_libvirt_master_ips,omitempty" yaml:"masterIPs"`
}

type Network struct {
	Name      string `json:"tectonic_libvirt_network_name,omitempty" yaml:"name"`
	IfName    string `json:"tectonic_libvirt_network_if,omitempty" yaml"ifName"`
	DNSServer string `json:"tectonic_libvirt_resolver,omitempty" yaml:"dnsServer"`
	IPRange   string `json:"tectonic_libvirt_ip_range,omitempty" yaml:"ipRange"`
}

// Fill in any variables for terraform
func (l *Libvirt) TFVars(masterCount int) error {
	_, network, err := net.ParseCIDR(l.Network.IPRange)
	if err != nil {
		return fmt.Errorf("failed to parse libvirt network ipRange: %v", err)
	}

	if len(l.MasterIPs) > 0 {
		if len(l.MasterIPs) != masterCount {
			return fmt.Errorf("length of MasterIPs doesn't match master count")
		} else {
			return nil
		}
	}

	for i := 0; i < masterCount; i++ {
		ip, err := cidr.Host(network, i+10)
		if err != nil {
			return fmt.Errorf("failed to generate master IPs: %v", err)
		}
		l.MasterIPs = append(l.MasterIPs, ip.String())
	}

	return nil
}
