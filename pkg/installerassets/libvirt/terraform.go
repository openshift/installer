package libvirt

import (
	"context"
	"encoding/json"
	"net"
	"strconv"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/openshift/installer/pkg/assets"
	"github.com/openshift/installer/pkg/installerassets"
	"github.com/pkg/errors"
)

type terraformConfig struct {
	URI         string   `json:"tectonic_libvirt_uri,omitempty"`
	Image       string   `json:"tectonic_os_image,omitempty"`
	IfName      string   `json:"tectonic_libvirt_network_if,omitempty"`
	IPRange     string   `json:"tectonic_libvirt_ip_range,omitempty"`
	BootstrapIP string   `json:"tectonic_libvirt_bootstrap_ip,omitempty"`
	MasterIPs   []string `json:"tectonic_libvirt_master_ips,omitempty"`
}

func terraformRebuilder(ctx context.Context, getByName assets.GetByString) (*assets.Asset, error) {
	asset := &assets.Asset{
		Name:          "terraform/libvirt-terraform.auto.tfvars",
		RebuildHelper: terraformRebuilder,
	}

	parents, err := asset.GetParents(
		ctx,
		getByName,
		"libvirt/image",
		"libvirt/network/interface-name",
		"libvirt/uri",
		"machines/master-count",
		"network/node-cidr",
	)
	if err != nil {
		return nil, err
	}

	masterCount, err := strconv.ParseUint(string(parents["machines/master-count"].Data), 10, 0)
	if err != nil {
		return nil, errors.Wrap(err, "parse master count")
	}

	_, nodeCIDR, err := net.ParseCIDR(string(parents["network/node-cidr"].Data))
	if err != nil {
		return nil, errors.Wrap(err, "parse node CIDR")
	}

	bootstrapIP, err := cidr.Host(nodeCIDR, 10)
	if err != nil {
		return nil, errors.Wrap(err, "generate bootstrap IP")
	}

	masterIPs := make([]string, 0, masterCount)
	for i := 0; i < int(masterCount); i++ {
		masterIP, err := cidr.Host(nodeCIDR, 11+i)
		if err != nil {
			return nil, errors.Wrap(err, "generate master IP")
		}
		masterIPs = append(masterIPs, masterIP.String())
	}

	image := string(parents["libvirt/image"].Data)
	image, err = getCachedImage(image)
	if err != nil {
		return nil, errors.Wrapf(err, "pull %s through the cache", image)
	}

	config := &terraformConfig{
		URI:         string(parents["libvirt/uri"].Data),
		Image:       image,
		IfName:      string(parents["libvirt/network/interface-name"].Data),
		IPRange:     nodeCIDR.String(),
		BootstrapIP: bootstrapIP.String(),
		MasterIPs:   masterIPs,
	}

	asset.Data, err = json.Marshal(config)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func init() {
	installerassets.Rebuilders["terraform/libvirt-terraform.auto.tfvars"] = terraformRebuilder
}
