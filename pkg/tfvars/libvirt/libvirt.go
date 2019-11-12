// Package libvirt contains libvirt-specific Terraform-variable logic.
package libvirt

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
	"strconv"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	"github.com/openshift/installer/pkg/tfvars/internal/cache"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

	cfg := &config{
		URI:          masterConfig.URI,
		IfName:       bridge,
		BootstrapIP:  bootstrapIP.String(),
		MasterIPs:    masterIPs,
		MasterMemory: strconv.Itoa(masterConfig.DomainMemory),
		MasterVcpu:   strconv.Itoa(masterConfig.DomainVcpu),
	}

	baseImageURL, err := url.ParseRequestURI(osImage)
	osImage, err = cachedImage(osImage)
	if strings.HasSuffix(baseImageURL.Path, ".gz") {
		osImageUncompressed := osImage + ".uncompressed"
		// Do nothing if we already have the uncompressed file in cache, otherwise decompress the data
		_, err = os.Stat(strings.TrimPrefix(osImageUncompressed, "file://"))
		if err != nil {
			if os.IsNotExist(err) {
				logrus.Infof("Decompress image data from %v to %v", osImage, osImageUncompressed)
				err = cache.DecompressFile(strings.TrimPrefix(osImage, "file://"), strings.TrimPrefix(osImageUncompressed, "file://"))
				if err != nil {
					logrus.Infof("Decompress %v %v failed: %v", osImage, osImageUncompressed, err)
					return nil, err
				}
			} else {
				logrus.Infof("Misc failure looking for %v: %v", osImageUncompressed, err)
				return nil, err
			}
		}
		fmt.Printf("Using uncompressed image %s\n", osImageUncompressed)
		cfg.Image = osImageUncompressed
	} else {
		cfg.Image = osImage
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to use cached libvirt image")
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
