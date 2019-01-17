// Package tfvars converts an InstallConfig to Terraform variables.
package tfvars

import (
	"encoding/json"

	"github.com/openshift/installer/pkg/tfvars/aws"
	"github.com/openshift/installer/pkg/tfvars/libvirt"
	"github.com/openshift/installer/pkg/tfvars/openstack"
	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
)

type config struct {
	ClusterID   string `json:"cluster_id,omitempty"`
	Name        string `json:"cluster_name,omitempty"`
	BaseDomain  string `json:"base_domain,omitempty"`
	MachineCIDR string `json:"machine_cidr"`
	Masters     int    `json:"master_count,omitempty"`

	IgnitionBootstrap string `json:"ignition_bootstrap,omitempty"`
	IgnitionMaster    string `json:"ignition_master,omitempty"`

	aws.AWS             `json:",inline"`
	libvirt.Libvirt     `json:",inline"`
	openstack.OpenStack `json:",inline"`
}

// TFVars converts the InstallConfig and Ignition content to
// terraform.tfvar JSON.
func TFVars(clusterID string, cfg *types.InstallConfig, osImage, bootstrapIgn, masterIgn string) ([]byte, error) {
	config := &config{
		ClusterID:   clusterID,
		Name:        cfg.ObjectMeta.Name,
		BaseDomain:  cfg.BaseDomain,
		MachineCIDR: cfg.Networking.MachineCIDR.String(),

		IgnitionMaster:    masterIgn,
		IgnitionBootstrap: bootstrapIgn,
	}

	for _, m := range cfg.Machines {
		switch m.Name {
		case "master":
			var replicas int
			if m.Replicas == nil {
				replicas = 1
			} else {
				replicas = int(*m.Replicas)
			}

			config.Masters += replicas
			if m.Platform.AWS != nil {
				config.AWS.Master = aws.Master{
					EC2Type:     m.Platform.AWS.InstanceType,
					IAMRoleName: m.Platform.AWS.IAMRoleName,
					MasterRootVolume: aws.MasterRootVolume{
						IOPS: m.Platform.AWS.EC2RootVolume.IOPS,
						Size: m.Platform.AWS.EC2RootVolume.Size,
						Type: m.Platform.AWS.EC2RootVolume.Type,
					},
				}
			}
		case "worker":
			if m.Platform.AWS != nil {
				config.AWS.Worker = aws.Worker{
					IAMRoleName: m.Platform.AWS.IAMRoleName,
				}
			}
		default:
			return nil, errors.Errorf("unrecognized machine pool %q", m.Name)
		}
	}

	if cfg.Platform.AWS != nil {
		config.AWS.Region = cfg.Platform.AWS.Region
		config.AWS.ExtraTags = cfg.Platform.AWS.UserTags
		config.AWS.EC2AMIOverride = osImage
	} else if cfg.Platform.Libvirt != nil {
		masterIPs := make([]string, len(cfg.Platform.Libvirt.MasterIPs))
		for i, ip := range cfg.Platform.Libvirt.MasterIPs {
			masterIPs[i] = ip.String()
		}
		config.Libvirt = libvirt.Libvirt{
			URI: cfg.Platform.Libvirt.URI,
			Network: libvirt.Network{
				IfName: cfg.Platform.Libvirt.Network.IfName,
			},
			Image:     osImage,
			MasterIPs: masterIPs,
		}
		if err := config.Libvirt.TFVars(&cfg.Networking.MachineCIDR.IPNet, config.Masters); err != nil {
			return nil, errors.Wrap(err, "failed to insert libvirt variables")
		}
		if err := config.Libvirt.UseCachedImage(); err != nil {
			return nil, errors.Wrap(err, "failed to use cached libvirt image")
		}
	} else if cfg.Platform.OpenStack != nil {
		config.OpenStack = openstack.OpenStack{
			Region:    cfg.Platform.OpenStack.Region,
			BaseImage: osImage,
		}
		config.OpenStack.Credentials.Cloud = cfg.Platform.OpenStack.Cloud
		config.OpenStack.ExternalNetwork = cfg.Platform.OpenStack.ExternalNetwork
		config.OpenStack.Master.FlavorName = cfg.Platform.OpenStack.FlavorName
		config.OpenStack.TrunkSupport = cfg.Platform.OpenStack.TrunkSupport
	}

	return json.MarshalIndent(config, "", "  ")
}
