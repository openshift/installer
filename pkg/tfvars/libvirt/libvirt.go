// Package libvirt contains libvirt-specific Terraform-variable logic.
package libvirt

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	"github.com/pkg/errors"
)

type config struct {
	URI          string   `json:"libvirt_uri,omitempty"`
	Image        string   `json:"os_image,omitempty"`
	IfName       string   `json:"libvirt_network_if"`
	MasterIPs    []string `json:"libvirt_master_ips,omitempty"`
	BootstrapIP  string   `json:"libvirt_bootstrap_ip,omitempty"`
	MasterMemory string   `json:"libvirt_master_memory,omitempty"`
	MasterVcpu   string   `json:"libvirt_master_vcpu,omitempty"`
}

// TFVars generates libvirt-specific Terraform variables.
func TFVars(masterConfig *v1beta1.LibvirtMachineProviderConfig, osImage string, machineCIDR *net.IPNet, bridge string, masterCount int) ([]byte, error) {
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

	domainMemory := strconv.Itoa(masterConfig.DomainMemory)
	if value := os.Getenv("TF_VAR_libvirt_master_memory"); value != "" {
		domainMemory = value
	}
	domainVcpu := strconv.Itoa(masterConfig.DomainVcpu)
	if value := os.Getenv("TF_VAR_libvirt_master_vcpu"); value != "" {
		domainVcpu = value
	}

	cfg := &config{
		URI:          masterConfig.URI,
		Image:        osImage,
		IfName:       bridge,
		BootstrapIP:  bootstrapIP.String(),
		MasterIPs:    masterIPs,
		MasterMemory: domainMemory,
		MasterVcpu:   domainVcpu,
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
