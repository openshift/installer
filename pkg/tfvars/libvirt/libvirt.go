package libvirt

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
)

// Libvirt encompasses configuration specific to libvirt.
type Libvirt struct {
	URI             string `json:"libvirt_uri,omitempty"`
	Image           string `json:"os_image,omitempty"`
	Network         `json:",inline"`
	ControlPlaneIPs []string `json:"libvirt_controlplane_ips,omitempty"`
	BootstrapIP     string   `json:"libvirt_bootstrap_ip,omitempty"`
}

// Network describes a libvirt network configuration.
type Network struct {
	IfName string `json:"libvirt_network_if"`
}

// TFVars fills in computed Terraform variables.
func (l *Libvirt) TFVars(machineCIDR *net.IPNet, controlPlaneCount int) error {
	if l.BootstrapIP == "" {
		ip, err := cidr.Host(machineCIDR, 10)
		if err != nil {
			return fmt.Errorf("failed to generate bootstrap IP: %v", err)
		}
		l.BootstrapIP = ip.String()
	}

	if len(l.ControlPlaneIPs) > 0 {
		if len(l.ControlPlaneIPs) != controlPlaneCount {
			return fmt.Errorf("length of ControlPlaneIPs doesn't match control plane count")
		}
	} else {
		if ips, err := generateIPs("controlplane", machineCIDR, controlPlaneCount, 11); err == nil {
			l.ControlPlaneIPs = ips
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
