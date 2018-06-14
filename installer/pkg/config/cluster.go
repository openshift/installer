package config

import (
	"encoding/json"
	"strings"

	"github.com/coreos/tectonic-config/config/tectonic-network"
	"gopkg.in/yaml.v2"

	"github.com/coreos/tectonic-installer/installer/pkg/config/aws"
	"github.com/coreos/tectonic-installer/installer/pkg/config/libvirt"
)

// Platform indicates the target platform of the cluster.
type Platform string

const (
	// IgnitionMaster is the relative path to the ign master cfg from the tf working directory
	IgnitionMaster = "ignition-master.ign"
	// IgnitionWorker is the relative path to the ign worker cfg from the tf working directory
	IgnitionWorker = "ignition-worker.ign"
	// IgnitionEtcd is the relative path to the ign etcd cfg from the tf working directory
	IgnitionEtcd = "ignition-etcd.ign"
	// PlatformAWS is the platform for a cluster launched on AWS.
	PlatformAWS Platform = "aws"
	// PlatformLibvirt is the platform for a cluster launched on libvirt.
	PlatformLibvirt Platform = "libvirt"
)

var defaultCluster = Cluster{
	AWS: aws.AWS{
		Endpoints:    aws.EndpointsAll,
		Profile:      aws.DefaultProfile,
		Region:       aws.DefaultRegion,
		VPCCIDRBlock: aws.DefaultVPCCIDRBlock,
	},
	ContainerLinux: ContainerLinux{
		Channel: ContainerLinuxChannelStable,
		Version: ContainerLinuxVersionLatest,
	},
	Libvirt: libvirt.Libvirt{
		Network: libvirt.Network{
			DNSServer: libvirt.DefaultDNSServer,
			IfName:    libvirt.DefaultIfName,
		},
	},
	Networking: Networking{
		MTU:         "1480",
		PodCIDR:     "10.2.0.0/16",
		ServiceCIDR: "10.3.0.0/16",
		Type:        tectonicnetwork.NetworkCanal,
	},
}

// Cluster defines the config for a cluster.
type Cluster struct {
	Admin           `json:",inline" yaml:"admin,omitempty"`
	BaseDomain      string `json:"tectonic_base_domain,omitempty" yaml:"baseDomain,omitempty"`
	CA              `json:",inline" yaml:"ca,omitempty"`
	ContainerLinux  `json:",inline" yaml:"containerLinux,omitempty"`
	Etcd            `json:",inline" yaml:"etcd,omitempty"`
	Internal        `json:",inline" yaml:"-"`
	LicensePath     string `json:"tectonic_license_path,omitempty" yaml:"licensePath,omitempty"`
	Master          `json:",inline" yaml:"master,omitempty"`
	Name            string `json:"tectonic_cluster_name,omitempty" yaml:"name,omitempty"`
	Networking      `json:",inline" yaml:"networking,omitempty"`
	NodePools       `json:"-" yaml:"nodePools"`
	Platform        Platform `json:"tectonic_platform" yaml:"platform,omitempty"`
	PullSecretPath  string   `json:"tectonic_pull_secret_path,omitempty" yaml:"pullSecretPath,omitempty"`
	Worker          `json:",inline" yaml:"worker,omitempty"`
	aws.AWS         `json:",inline" yaml:"aws,omitempty"`
	libvirt.Libvirt `json:",inline" yaml:"libvirt,omitempty"`
	IgnitionMaster  string `json:"tectonic_ignition_master,omitempty" yaml:"-"`
	IgnitionWorker  string `json:"tectonic_ignition_worker,omitempty" yaml:"-"`
	IgnitionEtcd    string `json:"tectonic_ignition_etcd,omitempty" yaml:"-"`
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
	c.Etcd.Count = c.NodeCount(c.Etcd.NodePools)
	c.Master.Count = c.NodeCount(c.Master.NodePools)
	c.Worker.Count = c.NodeCount(c.Worker.NodePools)

	c.IgnitionMaster = IgnitionMaster
	c.IgnitionWorker = IgnitionWorker
	c.IgnitionEtcd = IgnitionEtcd

	// fill in master ips
	if c.Platform == PlatformLibvirt {
		if err := c.Libvirt.TFVars(c.Master.Count); err != nil {
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
		Count: c.Etcd.Count,
		Name:  "etcd",
	})
	c.Etcd.NodePools = []string{"etcd"}

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

// String returns the string representation of a platform.
// This is the representation that should be used for
// consistent string comparison.
func (p Platform) String() string {
	return strings.ToLower(string(p))
}
