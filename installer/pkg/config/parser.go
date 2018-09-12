package config

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/coreos/tectonic-config/config/tectonic-network"
	"github.com/openshift/installer/installer/pkg/config/aws"
	"github.com/openshift/installer/installer/pkg/config/libvirt"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	"gopkg.in/yaml.v2"
)

// ParseConfig parses a yaml string and returns, if successful, a Cluster.
func ParseConfig(data []byte) (*Cluster, error) {
	cluster := defaultCluster

	err := parseInstallConfig(data, &cluster)
	if err != nil {
		err2 := parseLegacyConfig(data, &cluster)
		if err2 != nil {
			return nil, err
		}
	}

	if cluster.EC2AMIOverride == "" {
		ami, err := rhcos.AMI(DefaultChannel, cluster.AWS.Region)
		if err != nil {
			return nil, fmt.Errorf("failed to determine default AMI: %v", err)
		}
		cluster.EC2AMIOverride = ami
	}

	return &cluster, nil
}

func parseInstallConfig(data []byte, cluster *Cluster) (err error) {
	installConfig := &types.InstallConfig{}
	err = yaml.Unmarshal(data, &installConfig)
	if err != nil {
		return err
	}

	cluster.Name = installConfig.Name
	cluster.Internal.ClusterID = installConfig.ClusterID
	cluster.Admin = Admin{
		Email:    installConfig.Admin.Email,
		Password: installConfig.Admin.Password,
		SSHKey:   installConfig.Admin.SSHKey,
	}
	cluster.BaseDomain = installConfig.BaseDomain
	cluster.Networking = Networking{
		Type:        tectonicnetwork.NetworkType(installConfig.Networking.Type),
		ServiceCIDR: installConfig.Networking.ServiceCIDR.String(),
		PodCIDR:     installConfig.Networking.PodCIDR.String(),
	}

	for _, machinePool := range installConfig.Machines {
		nodePool := NodePool{
			Name: machinePool.Name,
		}
		if machinePool.Replicas == nil {
			nodePool.Count = 1
		} else {
			nodePool.Count = int(*machinePool.Replicas)
		}
		cluster.NodePools = append(cluster.NodePools, nodePool)
		switch machinePool.Name {
		case "master":
			if machinePool.Platform.AWS != nil {
				cluster.AWS.Master = aws.Master{
					EC2Type:     machinePool.Platform.AWS.InstanceType,
					IAMRoleName: machinePool.Platform.AWS.IAMRoleName,
					MasterRootVolume: aws.MasterRootVolume{
						IOPS: machinePool.Platform.AWS.EC2RootVolume.IOPS,
						Size: machinePool.Platform.AWS.EC2RootVolume.Size,
						Type: machinePool.Platform.AWS.EC2RootVolume.Type,
					},
				}
			}
		case "worker":
			if machinePool.Platform.AWS != nil {
				cluster.AWS.Worker = aws.Worker{
					EC2Type:     machinePool.Platform.AWS.InstanceType,
					IAMRoleName: machinePool.Platform.AWS.IAMRoleName,
					WorkerRootVolume: aws.WorkerRootVolume{
						IOPS: machinePool.Platform.AWS.EC2RootVolume.IOPS,
						Size: machinePool.Platform.AWS.EC2RootVolume.Size,
						Type: machinePool.Platform.AWS.EC2RootVolume.Type,
					},
				}
			}
		default:
			return fmt.Errorf("unrecognized machine pool %q", machinePool.Name)
		}

		if machinePool.Platform.Libvirt != nil {
			if cluster.Libvirt.QCOWImagePath != "" && cluster.Libvirt.QCOWImagePath != machinePool.Platform.Libvirt.QCOWImagePath {
				return fmt.Errorf("per-machine-pool images are not yet supported")
			}
			cluster.Libvirt.QCOWImagePath = machinePool.Platform.Libvirt.QCOWImagePath
		}
	}

	if installConfig.Platform.AWS != nil {
		cluster.AWS = aws.AWS{
			Region:    installConfig.Platform.AWS.Region,
			ExtraTags: installConfig.Platform.AWS.UserTags,
			External: aws.External{
				VPCID: installConfig.Platform.AWS.VPCID,
			},
			VPCCIDRBlock: installConfig.Platform.AWS.VPCCIDRBlock,
		}
	}

	if installConfig.Platform.Libvirt != nil {
		masterIPs := make([]string, len(installConfig.Platform.Libvirt.MasterIPs))
		for i, ip := range installConfig.Platform.Libvirt.MasterIPs {
			masterIPs[i] = ip.String()
		}
		cluster.Libvirt = libvirt.Libvirt{
			URI: installConfig.Platform.Libvirt.URI,
			Network: libvirt.Network{
				Name:    installConfig.Platform.Libvirt.Network.Name,
				IfName:  installConfig.Platform.Libvirt.Network.IfName,
				IPRange: installConfig.Platform.Libvirt.Network.IPRange,
			},
			MasterIPs: masterIPs,
		}
	}

	cluster.PullSecret = installConfig.PullSecret

	return nil
}

func parseLegacyConfig(data []byte, cluster *Cluster) (err error) {
	if err := yaml.Unmarshal(data, cluster); err != nil {
		return err
	}

	// Deprecated: remove after openshift/release is ported to pullSecret
	if cluster.PullSecretPath != "" {
		if cluster.PullSecret != "" {
			return errors.New("pullSecretPath is deprecated; just set pullSecret")
		}

		data, err := ioutil.ReadFile(cluster.PullSecretPath)
		if err != nil {
			return err
		}
		cluster.PullSecret = string(data)
	}

	return nil
}

// ParseConfigFile parses a yaml file and returns, if successful, a Cluster.
func ParseConfigFile(path string) (*Cluster, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseConfig(data)
}

// ParseInternal parses a yaml string and returns, if successful, an internal.
func ParseInternal(data []byte) (*Internal, error) {
	internal := &Internal{}

	if err := yaml.Unmarshal(data, internal); err != nil {
		return nil, err
	}

	return internal, nil
}

// ParseInternalFile parses a yaml file and returns, if successful, an internal.
func ParseInternalFile(path string) (*Internal, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return ParseInternal(data)
}
