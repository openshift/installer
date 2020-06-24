// Package libvirt contains libvirt-specific Terraform-variable logic.
package libvirt

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
)

type config struct {
	URI             string   `json:"libvirt_uri,omitempty"`
	Image           string   `json:"os_image,omitempty"`
	IfName          string   `json:"libvirt_network_if"`
	MasterIPs       []string `json:"libvirt_master_ips,omitempty"`
	BootstrapIP     string   `json:"libvirt_bootstrap_ip,omitempty"`
	MasterMemory    string   `json:"libvirt_master_memory,omitempty"`
	MasterVcpu      string   `json:"libvirt_master_vcpu,omitempty"`
	BootstrapMemory int      `json:"libvirt_bootstrap_memory,omitempty"`
	MasterDiskSize  string   `json:"libvirt_master_size,omitempty"`
}

// TFVars generates libvirt-specific Terraform variables.
func TFVars(masterConfig *v1beta1.LibvirtMachineProviderConfig, osImage string, machineCIDR *net.IPNet, bridge string, masterCount int, architecture types.Architecture) ([]byte, error) {
	bootstrapIP, err := cidr.Host(machineCIDR, 10)
	if err != nil {
		return nil, errors.Errorf("failed to generate bootstrap IP: %v", err)
	}

	masterIPs, err := generateIPs("master", machineCIDR, masterCount, 11)
	if err != nil {
		return nil, err
	}

	osImage, err = cache.DownloadImageFile(osImage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached libvirt image")
	}

	cfg := &config{
		URI:          masterConfig.URI,
		Image:        osImage,
		IfName:       bridge,
		BootstrapIP:  bootstrapIP.String(),
		MasterIPs:    masterIPs,
		MasterMemory: strconv.Itoa(masterConfig.DomainMemory),
		MasterVcpu:   strconv.Itoa(masterConfig.DomainVcpu),
	}

	if masterConfig.Volume.VolumeSize != nil {
		cfg.MasterDiskSize = masterConfig.Volume.VolumeSize.String()
	}

	if architecture == types.ArchitecturePPC64LE {
		cfg.BootstrapMemory = 5120
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
