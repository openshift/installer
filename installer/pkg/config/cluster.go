package config

import (
	"encoding/json"
	"fmt"

	"github.com/coreos/tectonic-config/config/tectonic-network"
	"gopkg.in/yaml.v2"

	"github.com/openshift/installer/installer/pkg/config/aws"
	"github.com/openshift/installer/installer/pkg/config/libvirt"
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
	ContainerLinux: ContainerLinux{
		Channel: ContainerLinuxChannelStable,
		Version: ContainerLinuxVersionLatest,
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
	Admin           `json:",inline" yaml:"admin,omitempty"`
	aws.AWS         `json:",inline" yaml:"aws,omitempty"`
	BaseDomain      string `json:"tectonic_base_domain,omitempty" yaml:"baseDomain,omitempty"`
	CA              `json:",inline" yaml:"CA,omitempty"`
	ContainerLinux  `json:",inline" yaml:"containerLinux,omitempty"`
	IgnitionMasters []string `json:"tectonic_ignition_masters,omitempty" yaml:"-"`
	IgnitionWorker  string   `json:"tectonic_ignition_worker,omitempty" yaml:"-"`
	Internal        `json:",inline" yaml:"-"`
	libvirt.Libvirt `json:",inline" yaml:"libvirt,omitempty"`
	LicensePath     string `json:"tectonic_license_path,omitempty" yaml:"licensePath,omitempty"`
	Master          `json:",inline" yaml:"master,omitempty"`
	Name            string `json:"tectonic_cluster_name,omitempty" yaml:"name,omitempty"`
	Networking      `json:",inline" yaml:"networking,omitempty"`
	NodePools       `json:"-" yaml:"nodePools"`
	Platform        Platform `json:"tectonic_platform" yaml:"platform,omitempty"`
	PullSecretPath  string   `json:"tectonic_pull_secret_path,omitempty" yaml:"pullSecretPath,omitempty"`
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
		c.IgnitionMasters = append(c.IgnitionMasters, fmt.Sprintf(IgnitionPathMaster, i))
	}

	c.IgnitionWorker = IgnitionPathWorker

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
