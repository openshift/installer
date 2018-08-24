package config

import "github.com/coreos/tectonic-config/config/tectonic-network"

// ContainerLinuxChannel indicates the selected Container Linux channel.
type ContainerLinuxChannel string

const (
	// ContainerLinuxChannelStable is the stable Container Linux channel.
	ContainerLinuxChannelStable ContainerLinuxChannel = "stable"
	// ContainerLinuxChannelBeta is the beta Container Linux channel.
	ContainerLinuxChannelBeta ContainerLinuxChannel = "beta"
	// ContainerLinuxChannelAlpha is the alpha Container Linux channel.
	ContainerLinuxChannelAlpha ContainerLinuxChannel = "alpha"
	// ContainerLinuxVersionLatest is the string to indicate the latest Container Linux version.
	ContainerLinuxVersionLatest = "latest"
)

// Admin converts admin related config.
type Admin struct {
	Email    string `json:"tectonic_admin_email" yaml:"email,omitempty"`
	Password string `json:"tectonic_admin_password" yaml:"password,omitempty"`
	SSHKey   string `json:"tectonic_admin_ssh_key,omitempty" yaml:"sshKey,omitempty"`
}

// CA related config
type CA struct {
	RootCACertPath string `json:"-" yaml:"rootCACertPath,omitempty"`
	RootCAKeyPath  string `json:"-" yaml:"rootCAKeyPath,omitempty"`
	RootCAKeyAlg   string `json:"-" yaml:"rootCAKeyAlg,omitempty"`
}

// ContainerLinux converts container linux related config.
type ContainerLinux struct {
	Channel ContainerLinuxChannel `json:"tectonic_container_linux_channel,omitempty" yaml:"channel,omitempty"`
	Version string                `json:"tectonic_container_linux_version,omitempty" yaml:"version,omitempty"`
}

// NodePool converts node pool related config.
type NodePool struct {
	Count        int    `json:"-" yaml:"count"`
	Name         string `json:"-" yaml:"name"`
	IgnitionFile string `json:"-" yaml:"ignitionFile"`
}

// NodePools converts node pools related config.
type NodePools []NodePool

// Map returns a nodePools' map equivalent.
func (n NodePools) Map() map[string]int {
	m := make(map[string]int)
	for i := range n {
		m[n[i].Name] = n[i].Count
	}
	return m
}

// Master converts master related config.
type Master struct {
	Count     int      `json:"tectonic_master_count,omitempty" yaml:"-"`
	NodePools []string `json:"-" yaml:"nodePools"`
}

// Networking converts networking related config.
type Networking struct {
	Type        tectonicnetwork.NetworkType `json:"tectonic_networking,omitempty" yaml:"type,omitempty"`
	MTU         string                      `json:"-" yaml:"mtu,omitempty"`
	ServiceCIDR string                      `json:"tectonic_service_cidr,omitempty" yaml:"serviceCIDR,omitempty"`
	PodCIDR     string                      `json:"tectonic_cluster_cidr,omitempty" yaml:"podCIDR,omitempty"`
}

// Worker converts worker related config.
type Worker struct {
	Count     int      `json:"tectonic_worker_count,omitempty" yaml:"-"`
	NodePools []string `json:"-" yaml:"nodePools"`
}

// Internal converts internal related config.
type Internal struct {
	ClusterID string `json:"tectonic_cluster_id,omitempty" yaml:"clusterId"`
}
