package config

import (
	"encoding/json"

	"gopkg.in/yaml.v2"

	"github.com/coreos/tectonic-installer/installer/pkg/config/aws"
	"github.com/coreos/tectonic-installer/installer/pkg/config/azure"
	"github.com/coreos/tectonic-installer/installer/pkg/config/gcp"
	"github.com/coreos/tectonic-installer/installer/pkg/config/govcloud"
	"github.com/coreos/tectonic-installer/installer/pkg/config/metal"
	"github.com/coreos/tectonic-installer/installer/pkg/config/openstack"
	"github.com/coreos/tectonic-installer/installer/pkg/config/vmware"
)

// Cluster defines the config for a cluster.
type Cluster struct {
	Admin               `json:",inline" yaml:"admin,omitempty"`
	BaseDomain          string `json:"tectonic_base_domain,omitempty" yaml:"baseDomain,omitempty"`
	CA                  `json:",inline" yaml:"ca,omitempty"`
	ContainerLinux      `json:",inline" yaml:"containerLinux,omitempty"`
	CustomCAPEMList     string `json:"tectonic_custom_ca_pem_list,omitempty" yaml:"customCAPEMList,omitempty"`
	DDNS                `json:",inline" yaml:"ddns,omitempty"`
	DNSName             string `json:"tectonic_dns_name,omitempty" yaml:"dnsName,omitempty"`
	Etcd                `json:",inline" yaml:"etcd,omitempty"`
	ISCSI               `json:",inline" yaml:"iscsi,omitempty"`
	Internal            `json:",inline" yaml:"-"`
	LicensePath         string `json:"tectonic_license_path,omitempty" yaml:"licensePath,omitempty"`
	Master              `json:",inline" yaml:"master,omitempty"`
	Name                string `json:"tectonic_cluster_name,omitempty" yaml:"name,omitempty"`
	Networking          `json:",inline" yaml:"networking,omitempty"`
	NodePools           `json:"-" yaml:"nodePools"`
	Platform            string `json:"-" yaml:"platform,omitempty"`
	Proxy               `json:",inline" yaml:"proxy,omitempty"`
	PullSecretPath      string `json:"tectonic_pull_secret_path,omitempty" yaml:"pullSecretPath,omitempty"`
	TLSValidityPeriod   int    `json:"tectonic_tls_validity_period,omitempty" yaml:"tlsValidityPeriod,omitempty"`
	Worker              `json:",inline" yaml:"worker,omitempty"`
	aws.AWS             `json:",inline" yaml:"aws,omitempty"`
	azure.Azure         `json:",inline" yaml:"azure,omitempty"`
	gcp.GCP             `json:",inline" yaml:"gcp,omitempty"`
	govcloud.GovCloud   `json:",inline" yaml:"govcloud,omitempty"`
	metal.Metal         `json:",inline" yaml:"metal,omitempty"`
	openstack.OpenStack `json:",inline" yaml:"openstack,omitempty"`
	vmware.VMware       `json:",inline" yaml:"vmware,omitempty"`
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
