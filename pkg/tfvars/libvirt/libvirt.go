package libvirt

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
)

// Libvirt encompasses configuration specific to libvirt.
type Libvirt struct {
	URI         string `json:"tectonic_libvirt_uri,omitempty"`
	Image       string `json:"tectonic_os_image,omitempty"`
	Network     `json:",inline"`
	MasterIPs   []string `json:"tectonic_libvirt_master_ips,omitempty"`
	BootstrapIP string   `json:"tectonic_libvirt_bootstrap_ip,omitempty"`
}

// Network describes a libvirt network configuration.
type Network struct {
	IfName  string `json:"tectonic_libvirt_network_if,omitempty"`
	IPRange string `json:"tectonic_libvirt_ip_range,omitempty"`
}

// TFVars fills in computed Terraform variables.
func (l *Libvirt) TFVars(masterCount int) error {
	_, network, err := net.ParseCIDR(l.Network.IPRange)
	if err != nil {
		return fmt.Errorf("failed to parse libvirt network ipRange: %v", err)
	}

	if l.BootstrapIP == "" {
		ip, err := cidr.Host(network, 10)
		if err != nil {
			return fmt.Errorf("failed to generate bootstrap IP: %v", err)
		}
		l.BootstrapIP = ip.String()
	}

	if len(l.MasterIPs) > 0 {
		if len(l.MasterIPs) != masterCount {
			return fmt.Errorf("length of MasterIPs doesn't match master count")
		}
	} else {
		if ips, err := generateIPs("master", network, masterCount, 11); err == nil {
			l.MasterIPs = ips
		} else {
			return err
		}
	}

	return nil
}

func generateIPs(name string, network *net.IPNet, count int, offset int) ([]string, error) {
	var ips []string
	for i := 0; i < count; i++ {
		ip, err := cidr.Host(network, offset+i)
		if err != nil {
			return nil, fmt.Errorf("failed to generate %s IPs: %v", name, err)
		}
		ips = append(ips, ip.String())
	}

	return ips, nil
}
