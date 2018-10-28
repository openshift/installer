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
	ClusterID  string `json:"tectonic_cluster_id,omitempty"`
	Name       string `json:"tectonic_cluster_name,omitempty"`
	BaseDomain string `json:"tectonic_base_domain,omitempty"`
	Masters    int    `json:"tectonic_master_count,omitempty"`
	Workers    int    `json:"tectonic_worker_count,omitempty"`

	IgnitionBootstrap string `json:"ignition_bootstrap,omitempty"`
	IgnitionMaster    string `json:"ignition_master,omitempty"`
	IgnitionWorker    string `json:"ignition_worker,omitempty"`

	aws.AWS             `json:",inline"`
	libvirt.Libvirt     `json:",inline"`
	openstack.OpenStack `json:",inline"`
}

// TFVars converts the InstallConfig and Ignition content to
// terraform.tfvar JSON.
func TFVars(cfg *types.InstallConfig, bootstrapIgn, masterIgn, workerIgn string) ([]byte, error) {
	config := &config{
		ClusterID:  cfg.ClusterID,
		Name:       cfg.ObjectMeta.Name,
		BaseDomain: cfg.BaseDomain,

		IgnitionMaster:    masterIgn,
		IgnitionWorker:    workerIgn,
		IgnitionBootstrap: bootstrapIgn,
	}

	for _, m := range cfg.Machines {
		var replicas int
		if m.Replicas == nil {
			replicas = 1
		} else {
			replicas = int(*m.Replicas)
		}

		switch m.Name {
		case "master":
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
			config.Workers += replicas
			if m.Platform.AWS != nil {
				config.AWS.Worker = aws.Worker{
					EC2Type:     m.Platform.AWS.InstanceType,
					IAMRoleName: m.Platform.AWS.IAMRoleName,
					WorkerRootVolume: aws.WorkerRootVolume{
						IOPS: m.Platform.AWS.EC2RootVolume.IOPS,
						Size: m.Platform.AWS.EC2RootVolume.Size,
						Type: m.Platform.AWS.EC2RootVolume.Type,
					},
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
			Endpoints: aws.EndpointsAll, // Default value for endpoints.
			Region:    cfg.Platform.AWS.Region,
			ExtraTags: cfg.Platform.AWS.UserTags,
			External: aws.External{
				VPCID: cfg.Platform.AWS.VPCID,
			},
			VPCCIDRBlock:   cfg.Platform.AWS.VPCCIDRBlock,
			EC2AMIOverride: ami,
		}
	} else if cfg.Platform.Libvirt != nil {
		masterIPs := make([]string, len(cfg.Platform.Libvirt.MasterIPs))
		for i, ip := range cfg.Platform.Libvirt.MasterIPs {
			masterIPs[i] = ip.String()
		}
		config.Libvirt = libvirt.Libvirt{
			URI: cfg.Platform.Libvirt.URI,
			Network: libvirt.Network{
				Name:    cfg.Platform.Libvirt.Network.Name,
				IfName:  cfg.Platform.Libvirt.Network.IfName,
				IPRange: cfg.Platform.Libvirt.Network.IPRange,
			},
			Image:     cfg.Platform.Libvirt.DefaultMachinePlatform.Image,
			MasterIPs: masterIPs,
		}
		if err := config.Libvirt.TFVars(config.Masters, config.Workers); err != nil {
			return nil, errors.Wrap(err, "failed to insert libvirt variables")
		}
		if err := config.Libvirt.UseCachedImage(); err != nil {
			return nil, errors.Wrap(err, "failed to use cached libvirt image")
		}
	} else if cfg.Platform.OpenStack != nil {
		config.OpenStack = openstack.OpenStack{
			Region:           cfg.Platform.OpenStack.Region,
			NetworkCIDRBlock: cfg.Platform.OpenStack.NetworkCIDRBlock,
		}
		config.OpenStack.Credentials.Cloud = cfg.Platform.OpenStack.Cloud
		config.OpenStack.ExternalNetwork = cfg.Platform.OpenStack.ExternalNetwork
	}

	return json.MarshalIndent(config, "", "  ")
}
