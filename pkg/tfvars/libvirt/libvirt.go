// Package libvirt contains libvirt-specific Terraform-variable logic.
package libvirt

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1alpha1"
	"github.com/pkg/errors"
)

type config struct {
	URI             string   `json:"libvirt_uri,omitempty"`
	Image           string   `json:"os_image,omitempty"`
	IfName          string   `json:"libvirt_network_if"`
	ControlPlaneIPs []string `json:"libvirt_control_plane_ips,omitempty"`
	BootstrapIP     string   `json:"libvirt_bootstrap_ip,omitempty"`
}

// TFVars generates libvirt-specific Terraform variables.
func TFVars(controlPlaneConfig *v1alpha1.LibvirtMachineProviderConfig, osImage string, machineCIDR *net.IPNet, bridge string, controlPlaneCount int) ([]byte, error) {
	bootstrapIP, err := cidr.Host(machineCIDR, 10)
	if err != nil {
		return nil, errors.Errorf("failed to generate bootstrap IP: %v", err)
	}

	controlPlaneIPs, err := generateIPs("control plane", machineCIDR, controlPlaneCount, 11)
	if err != nil {
		return nil, err
	}

	osImage, err = cachedImage(osImage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached libvirt image")
	}

	cfg := &config{
		URI:             controlPlaneConfig.URI,
		Image:           osImage,
		IfName:          bridge,
		BootstrapIP:     bootstrapIP.String(),
		ControlPlaneIPs: controlPlaneIPs,
	}

	return json.MarshalIndent(cfg, "", "  ")
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
