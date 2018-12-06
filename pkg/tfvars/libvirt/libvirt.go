package libvirt

import (
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
)

// Libvirt encompasses configuration specific to libvirt.
type Libvirt struct {
	URI         string `json:"libvirt_uri,omitempty"`
	Image       string `json:"os_image,omitempty"`
	Network     `json:",inline"`
	MasterIPs   []string `json:"libvirt_master_ips,omitempty"`
	BootstrapIP string   `json:"libvirt_bootstrap_ip,omitempty"`
	MasterIPs   []string `json:"tectonic_libvirt_master_ips,omitempty"`
	BootstrapIP string   `json:"tectonic_libvirt_bootstrap_ip,omitempty"`
	SSHKey string `json:"ssh_key,omitempty"`
}

// Network describes a libvirt network configuration.
type Network struct {
	IfName  string `json:"libvirt_network_if,omitempty"`
	IPRange string `json:"libvirt_ip_range,omitempty"`
}

// TFVars fills in computed Terraform variables.
func (l *Libvirt) TFVars(masterCount int) error {
	_, network, err := net.ParseCIDR(l.Network.IPRange)
	if err != nil {
		return fmt.Errorf("failed to parse libvirt network ipRange: %v", err)
	}
	l.SSHKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDL5/TfQpp2kEUmhpPWRCuHDPOMxuQLdwTlK9TiiLWzHSCRVSET6vwhg8Vrph5SxgsGR9rC+KQkAfNLyS800RmG+svxaCJSNFRzweK5kC9kj4KalxbiYKANfQbKMk8iehQ76GBa6WCYc+GAaOBnVxFaqQ0AuvIeyyRGTPsjaqslHujEijCfyhWWOtMa0v5i4GyY0T0Rvi+sF0f7p6IEkH8rgdMwFju7wVisrSGKRBN3KR/cr7lI/oKPfN6ApSNPuSo2TIWyUfy3P+y2njwG+3DuaQTUM2W2M7dXeosZ2dCYtzV3EeeXk6aZhIfrTgi+ckGC73onhvb5D6ob+3Iu5U8T4eG3HewwlcohMvPS8Syi6GM+Nh7EjrO8TXJACKUfaDDQTlHoPc5Ws52dt1UBd067XEN6IQxKlM2qTqy9B4tFag5veMK70sSQqrMI4RGKRwwdN+KWPqcupX2F81uR12zNCc7RjuGi4ofI6hI8yw89KYMVsKqrRz3xTz9LyXlaxuXFfvIK+I9Q13OC/6AahX7bdvLpK7KvIuaxFt7jUExzpHtcztJ1KDEUYpuMhRFsor0PH+KlkLsTqhvKRA5W0/pCUqmjmsRuenG6dZwkKHZs54KEGmq0KOhShuBwd/2dOeIZH9JENgUdy5x9hEzJoEhOpWXQVjc8UC6xq0H3VGO59Q== mgugino@redhat.com"
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
