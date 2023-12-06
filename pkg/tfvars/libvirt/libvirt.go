// Package libvirt contains libvirt-specific Terraform-variable logic.
package libvirt

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/pkg/errors"

	"github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	"github.com/openshift/installer/pkg/rhcos/cache"
	"github.com/openshift/installer/pkg/types"
)

type config struct {
	URI             string              `json:"libvirt_uri,omitempty"`
	Image           string              `json:"os_image,omitempty"`
	IfName          string              `json:"libvirt_network_if"`
	MasterIPs       []string            `json:"libvirt_master_ips,omitempty"`
	BootstrapIP     string              `json:"libvirt_bootstrap_ip,omitempty"`
	MasterMemory    string              `json:"libvirt_master_memory,omitempty"`
	MasterVcpu      string              `json:"libvirt_master_vcpu,omitempty"`
	BootstrapMemory int                 `json:"libvirt_bootstrap_memory,omitempty"`
	MasterDiskSize  string              `json:"libvirt_master_size,omitempty"`
	DnsmasqOptions  []map[string]string `json:"libvirt_dnsmasq_options,omitempty"`
}

// TFVarsSources contains the parameters to be converted into Terraform variables
type TFVarsSources struct {
	MasterConfig   *v1beta1.LibvirtMachineProviderConfig
	OsImage        string
	MachineCIDR    *net.IPNet
	Bridge         string
	MasterCount    int
	Architecture   types.Architecture
	DnsmasqOptions []map[string]string
}

// TFVars generates libvirt-specific Terraform variables.
func TFVars(sources TFVarsSources) ([]byte, error) {
	bootstrapIP, err := cidr.Host(sources.MachineCIDR, 10)
	if err != nil {
		return nil, errors.Errorf("failed to generate bootstrap IP: %v", err)
	}

	masterIPs, err := generateIPs("master", sources.MachineCIDR, sources.MasterCount, 11)
	if err != nil {
		return nil, err
	}

	osImage := sources.OsImage
	url, err := url.Parse(osImage)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse image url")
	}

	if url.Scheme == "file" {
		osImage = filepath.FromSlash(url.Path)
		if _, err = os.Stat(osImage); err != nil {
			return nil, errors.Wrap(err, "failed to access file or directory")
		}
	} else {
		osImage, err = cache.DownloadImageFile(osImage, cache.InstallerApplicationName)
		if err != nil {
			return nil, errors.Wrap(err, "failed to use cached libvirt image")
		}
	}

	cfg := &config{
		URI:            sources.MasterConfig.URI,
		Image:          osImage,
		IfName:         sources.Bridge,
		BootstrapIP:    bootstrapIP.String(),
		MasterIPs:      masterIPs,
		MasterMemory:   strconv.Itoa(sources.MasterConfig.DomainMemory),
		MasterVcpu:     strconv.Itoa(sources.MasterConfig.DomainVcpu),
		DnsmasqOptions: sources.DnsmasqOptions,
	}

	if sources.MasterConfig.Volume.VolumeSize != nil {
		// As per https://github.com/hashicorp/terraform/issues/3287 the
		// master disk size needs to be specified as the number of bytes.
		diskSizeInBytes, converted := sources.MasterConfig.Volume.VolumeSize.AsInt64()
		if !converted {
			msgTemplate := "failed to convert master disk size to bytes: %s"
			return nil, fmt.Errorf(msgTemplate, sources.MasterConfig.Volume.VolumeSize)
		}
		cfg.MasterDiskSize = fmt.Sprintf("%d", diskSizeInBytes)
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
