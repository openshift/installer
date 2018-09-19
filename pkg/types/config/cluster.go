package config

import (
	"encoding/json"
	"fmt"

	"github.com/coreos/tectonic-config/config/tectonic-network"
	"gopkg.in/yaml.v2"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/config/aws"
	"github.com/openshift/installer/pkg/types/config/libvirt"
)

const (
	// IgnitionPathMaster is the relative path to the ign master cfg from the tf working directory
	// This is a format string so that the index can be populated later
	IgnitionPathMaster = "master-%d.ign"
	// IgnitionPathWorker is the relative path to the ign worker cfg from the tf working directory
	IgnitionPathWorker = "worker.ign"
	// PlatformAWS is the platform for a cluster launched on AWS.
	PlatformAWS Platform = "aws"
	// PlatformLibvirt is the platform for a cluster launched on libvirt.
	PlatformLibvirt Platform = "libvirt"
)

// Platform indicates the target platform of the cluster.
type Platform string

// UnmarshalYAML unmarshals and verifies the platform.
func (p *Platform) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var data string
	if err := unmarshal(&data); err != nil {
		return err
	}

	platform := Platform(data)
	switch platform {
	case PlatformAWS, PlatformLibvirt:
	default:
		return fmt.Errorf("invalid platform specified (%s); must be one of %s", platform, []Platform{PlatformAWS, PlatformLibvirt})
	}

	*p = platform
	return nil
}

var defaultCluster = Cluster{
	AWS: aws.AWS{
		Endpoints:    aws.EndpointsAll,
		Profile:      aws.DefaultProfile,
		Region:       aws.DefaultRegion,
		VPCCIDRBlock: aws.DefaultVPCCIDRBlock,
	},
	CA: CA{
		RootCAKeyAlg: "RSA",
	},
	Libvirt: libvirt.Libvirt{
		Network: libvirt.Network{
			IfName: libvirt.DefaultIfName,
		},
	},
	Networking: Networking{
		MTU:         "1480",
		PodCIDR:     "10.2.0.0/16",
		ServiceCIDR: "10.3.0.0/16",
		Type:        tectonicnetwork.NetworkFlannel,
	},
}

// Cluster defines the config for a cluster.
type Cluster struct {
	Admin      `json:",inline" yaml:"admin,omitempty"`
	aws.AWS    `json:",inline" yaml:"aws,omitempty"`
	BaseDomain string `json:"tectonic_base_domain,omitempty" yaml:"baseDomain,omitempty"`
	CA         `json:",inline" yaml:"CA,omitempty"`

	// Deprecated, will be removed soon.
	IgnitionMasterPaths []string `json:"tectonic_ignition_masters,omitempty" yaml:"-"`
	// Deprecated, will be removed soon.
	IgnitionWorkerPath string `json:"tectonic_ignition_worker,omitempty" yaml:"-"`

	IgnitionBootstrap string   `json:"openshift_ignition_bootstrap,omitempty" yaml:"-"`
	IgnitionMasters   []string `json:"openshift_ignition_master,omitempty" yaml:"-"`
	IgnitionWorker    string   `json:"openshift_ignition_worker,omitempty" yaml:"-"`

	Internal        `json:",inline" yaml:"-"`
	libvirt.Libvirt `json:",inline" yaml:"libvirt,omitempty"`
	Master          `json:",inline" yaml:"master,omitempty"`
	Name            string `json:"tectonic_cluster_name,omitempty" yaml:"name,omitempty"`
	Networking      `json:",inline" yaml:"networking,omitempty"`
	NodePools       `json:"-" yaml:"nodePools"`
	Platform        Platform `json:"tectonic_platform" yaml:"platform,omitempty"`
	PullSecret      string   `json:"tectonic_pull_secret,omitempty" yaml:"pullSecret,omitempty"`
	PullSecretPath  string   `json:"-" yaml:"pullSecretPath,omitempty"` // Deprecated: remove after openshift/release is ported to pullSecret
	Worker          `json:",inline" yaml:"worker,omitempty"`
}

// NodeCount will return the number of nodes specified in NodePools with matching names.
// If no matching NodePools are found, then 0 is returned.
func (c Cluster) NodeCount(names []string) int {
	var count int
	for _, name := range names {
		for _, n := range c.NodePools {
			if n.Name == name {
				count += n.Count
				break
			}
		}
	}
	return count
}

// TFVars will return the config for the cluster in tfvars format.
func (c *Cluster) TFVars() (string, error) {
	c.Master.Count = c.NodeCount(c.Master.NodePools)
	c.Worker.Count = c.NodeCount(c.Worker.NodePools)

	for i := 0; i < c.Master.Count; i++ {
		c.IgnitionMasterPaths = append(c.IgnitionMasterPaths, fmt.Sprintf(IgnitionPathMaster, i))
	}

	c.IgnitionWorkerPath = IgnitionPathWorker

	// fill in master ips
	if c.Platform == PlatformLibvirt {
		if err := c.Libvirt.TFVars(c.Master.Count, c.Worker.Count); err != nil {
			return "", err
		}
	}

	data, err := json.MarshalIndent(&c, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// YAML will return the config for the cluster in yaml format.
func (c *Cluster) YAML() (string, error) {
	c.NodePools = append(c.NodePools, NodePool{
		Count: c.Master.Count,
		Name:  "master",
	})
	c.Master.NodePools = []string{"master"}

	c.NodePools = append(c.NodePools, NodePool{
		Count: c.Worker.Count,
		Name:  "worker",
	})
	c.Worker.NodePools = []string{"worker"}

	yaml, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(yaml), nil
}

// ConvertInstallConfigToTFVar converts the installconfig to the Cluster struct
// that represents the terraform.tfvar file.
// TODO(yifan): Clean up the Cluster struct to trim unnecessary fields.
func ConvertInstallConfigToTFVar(cfg *types.InstallConfig, bootstrapIgn string, masterIgns []string, workerIgn string) (*Cluster, error) {
	cluster := &Cluster{
		Admin: Admin{
			Email:    cfg.Admin.Email,
			Password: cfg.Admin.Password,
			SSHKey:   cfg.Admin.SSHKey,
		},

		IgnitionMasters:   masterIgns,
		IgnitionWorker:    workerIgn,
		IgnitionBootstrap: bootstrapIgn,

		Internal: Internal{
			ClusterID: cfg.ClusterID,
		},

		Networking: Networking{
			Type:        tectonicnetwork.NetworkType(cfg.Networking.Type),
			ServiceCIDR: cfg.Networking.ServiceCIDR.String(),
			PodCIDR:     cfg.Networking.PodCIDR.String(),
		},
		BaseDomain: cfg.BaseDomain,
		Name:       cfg.Name,
		PullSecret: cfg.PullSecret,
	}

	if cfg.Platform.AWS != nil {
		cluster.Platform = PlatformAWS
		cluster.AWS = aws.AWS{
			Region:    cfg.Platform.AWS.Region,
			ExtraTags: cfg.Platform.AWS.UserTags,
			External: aws.External{
				VPCID: cfg.Platform.AWS.VPCID,
			},
			VPCCIDRBlock: cfg.Platform.AWS.VPCCIDRBlock,
		}
	} else if cfg.Platform.Libvirt != nil {
		cluster.Platform = PlatformLibvirt
		masterIPs := make([]string, len(cfg.Platform.Libvirt.MasterIPs))
		for i, ip := range cfg.Platform.Libvirt.MasterIPs {
			masterIPs[i] = ip.String()
		}
		cluster.Libvirt = libvirt.Libvirt{
			URI: cfg.Platform.Libvirt.URI,
			Network: libvirt.Network{
				Name:    cfg.Platform.Libvirt.Network.Name,
				IfName:  cfg.Platform.Libvirt.Network.IfName,
				IPRange: cfg.Platform.Libvirt.Network.IPRange,
			},
			MasterIPs: masterIPs,
		}
	}

	for _, m := range cfg.Machines {
		nodePool := NodePool{
			Name: m.Name,
		}
		if m.Replicas == nil {
			nodePool.Count = 1
		} else {
			nodePool.Count = int(*m.Replicas)
		}
		cluster.NodePools = append(cluster.NodePools, nodePool)

		switch m.Name {
		case "master":
			cluster.Master.Count += nodePool.Count
			cluster.Master.NodePools = append(cluster.Master.NodePools, m.Name)
			if m.Platform.AWS != nil {
				cluster.AWS.Master = aws.Master{
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
			cluster.Worker.Count += nodePool.Count
			cluster.Worker.NodePools = append(cluster.Worker.NodePools, m.Name)
			if m.Platform.AWS != nil {
				cluster.AWS.Worker = aws.Worker{
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
			return nil, fmt.Errorf("unrecognized machine pool %q", m.Name)
		}

	}

	return cluster, nil
}
