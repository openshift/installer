// Package tfvars converts an InstallConfig to Terraform variables.
package tfvars

import (
	"context"
	"encoding/json"
	"time"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/tfvars/aws"
	"github.com/openshift/installer/pkg/tfvars/libvirt"
	"github.com/openshift/installer/pkg/tfvars/openstack"
	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
)

type config struct {
	ClusterID    string `json:"cluster_id,omitempty"`
	Name         string `json:"cluster_name,omitempty"`
	BaseDomain   string `json:"base_domain,omitempty"`
	MachineCIDR  string `json:"machine_cidr"`
	ControlPlane int    `json:"controlplane_count,omitempty"`

	IgnitionBootstrap    string `json:"ignition_bootstrap,omitempty"`
	IgnitionControlPlane string `json:"ignition_controlplane,omitempty"`

	aws.AWS             `json:",inline"`
	libvirt.Libvirt     `json:",inline"`
	openstack.OpenStack `json:",inline"`
}

// TFVars converts the InstallConfig and Ignition content to
// terraform.tfvar JSON.
func TFVars(cfg *types.InstallConfig, bootstrapIgn, controlPlaneIgn string) ([]byte, error) {
	config := &config{
		ClusterID:   cfg.ClusterID,
		Name:        cfg.ObjectMeta.Name,
		BaseDomain:  cfg.BaseDomain,
		MachineCIDR: cfg.Networking.MachineCIDR.String(),

		IgnitionControlPlane: controlPlaneIgn,
		IgnitionBootstrap:    bootstrapIgn,
	}

	for _, m := range cfg.Machines {
		switch m.Name {
		case "controlplane":
			var replicas int
			if m.Replicas == nil {
				replicas = 1
			} else {
				replicas = int(*m.Replicas)
			}

			config.ControlPlane += replicas
			if m.Platform.AWS != nil {
				config.AWS.ControlPlane = aws.ControlPlane{
					EC2Type:     m.Platform.AWS.InstanceType,
					IAMRoleName: m.Platform.AWS.IAMRoleName,
					ControlPlaneRootVolume: aws.ControlPlaneRootVolume{
						IOPS: m.Platform.AWS.EC2RootVolume.IOPS,
						Size: m.Platform.AWS.EC2RootVolume.Size,
						Type: m.Platform.AWS.EC2RootVolume.Type,
					},
				}
			}
		case "compute":
			if m.Platform.AWS != nil {
				config.AWS.Compute = aws.Compute{
					IAMRoleName: m.Platform.AWS.IAMRoleName,
				}
			}
		default:
			return nil, errors.Errorf("unrecognized machine pool %q", m.Name)
		}
	}

	if cfg.Platform.AWS != nil {
		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancel()
		ami, err := rhcos.AMI(ctx, rhcos.DefaultChannel, cfg.Platform.AWS.Region)
		if err != nil {
			return nil, errors.Wrap(err, "failed to determine default AMI")
		}

		config.AWS = aws.AWS{
			Region:         cfg.Platform.AWS.Region,
			ExtraTags:      cfg.Platform.AWS.UserTags,
			EC2AMIOverride: ami,
		}
	} else if cfg.Platform.Libvirt != nil {
		controlPlaneIPs := make([]string, len(cfg.Platform.Libvirt.ControlPlaneIPs))
		for i, ip := range cfg.Platform.Libvirt.ControlPlaneIPs {
			controlPlaneIPs[i] = ip.String()
		}
		config.Libvirt = libvirt.Libvirt{
			URI: cfg.Platform.Libvirt.URI,
			Network: libvirt.Network{
				IfName: cfg.Platform.Libvirt.Network.IfName,
			},
			Image:           cfg.Platform.Libvirt.DefaultMachinePlatform.Image,
			ControlPlaneIPs: controlPlaneIPs,
		}
		if err := config.Libvirt.TFVars(&cfg.Networking.MachineCIDR.IPNet, config.ControlPlane); err != nil {
			return nil, errors.Wrap(err, "failed to insert libvirt variables")
		}
		if err := config.Libvirt.UseCachedImage(); err != nil {
			return nil, errors.Wrap(err, "failed to use cached libvirt image")
		}
	} else if cfg.Platform.OpenStack != nil {
		config.OpenStack = openstack.OpenStack{
			Region:    cfg.Platform.OpenStack.Region,
			BaseImage: cfg.Platform.OpenStack.BaseImage,
		}
		config.OpenStack.Credentials.Cloud = cfg.Platform.OpenStack.Cloud
		config.OpenStack.ExternalNetwork = cfg.Platform.OpenStack.ExternalNetwork
		config.OpenStack.ControlPlane.FlavorName = cfg.Platform.OpenStack.FlavorName
		config.OpenStack.TrunkSupport = cfg.Platform.OpenStack.TrunkSupport
	}

	return json.MarshalIndent(config, "", "  ")
}
