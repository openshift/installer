package config

import "gopkg.in/yaml.v2"

// Config defines the top level config for a configuration file.
type Config struct {
	Clusters []Cluster `json:",inline" yaml:"clusters,omitempty"`
}

// YAML will return the config for the cluster in yaml format.
func (c *Config) YAML() (string, error) {
	for _, cluster := range c.Clusters {
		cluster.NodePools = append(cluster.NodePools, NodePool{
			Count: cluster.Etcd.Count,
			Name:  "etcd",
		})
		cluster.Etcd.NodePools = []string{"etcd"}

		cluster.NodePools = append(cluster.NodePools, NodePool{
			Count: cluster.Master.Count,
			Name:  "master",
		})
		cluster.Master.NodePools = []string{"master"}

		cluster.NodePools = append(cluster.NodePools, NodePool{
			Count: cluster.Worker.Count,
			Name:  "worker",
		})
		cluster.Worker.NodePools = []string{"worker"}
	}

	yaml, err := yaml.Marshal(c)
	if err != nil {
		return "", err
	}

	return string(yaml), nil
}
