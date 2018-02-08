package config

import (
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
	Console              console              `yaml:"Console"`
	ContainerLinux       containerLinux       `yaml:"ContainerLinux"`
	DNS                  dns                  `yaml:"DNS"`
	Etcd                 etcd                 `yaml:"Etcd"`
	ExternalTLSMaterials externalTLSMaterials `yaml:"ExternalTLSMaterials"`
	Masters              masters              `yaml:"Masters"`
	Name                 string               `yaml:"Name"`
	Networking           networking           `yaml:"Networking"`
	Platform             string               `yaml:"Platform"`
	Tectonic             tectonic             `yaml:"Tectonic"`
	Update               update               `yaml:"Update"`
	Workers              workers              `yaml:"Workers"`

	AWS       aws.Config       `yaml:"AWS,omitempty"`
	Azure     azure.Config     `yaml:"Azure,omitempty"`
	GCP       gcp.Config       `yaml:"GCP,omitempty"`
	GovCloud  govcloud.Config  `yaml:"GovCloud,omitempty"`
	Metal     metal.Config     `yaml:"Metal,omitempty"`
	OpenStack openstack.Config `yaml:"OpenStack,omitempty"`
	VMware    vmware.Config    `yaml:"VMware,omitempty"`
}

// Config defines the top level config for a configuration file.
type Config struct {
	Clusters []Cluster `yaml:"Clusters"`
}

type console struct {
	AdminEmail    string `yaml:"AdminEmail"`
	AdminPassword string `yaml:"AdminPassword"`
}

type containerLinux struct {
	Channel string `yaml:"Channel"`
	Version string `yaml:"Version"`
}

type dns struct {
	BaseDomain string `yaml:"BaseDomain"`
}

type etcd struct {
	NodeCount       int      `yaml:"NodeCount"`
	MachineType     string   `yaml:"MachineType"`
	ExternalServers []string `yaml:"ExternalServers"`
}

type externalTLSMaterials struct {
	ValidityPeriod int    `yaml:"ValidityPeriod"`
	EtcdCACertPath string `yaml:"EtcdCACertPath"`
}

type masters struct {
	NodeCount   int    `yaml:"NodeCount"`
	MachineType string `yaml:"MachineType"`
}

type networking struct {
	Type        string `yaml:"Type"`
	MTU         string `yaml:"MTU"`
	NodeCIDR    string `yaml:"NodeCIDR"`
	ServiceCIDR string `yaml:"ServiceCIDR"`
	PodCIDR     string `yaml:"PodCIDR"`
}

type tectonic struct {
	PullSecretPath string `yaml:"PullSecretPath"`
	LicensePath    string `yaml:"LicensePath"`
}

type update struct {
	Server  string `yaml:"Server"`
	Channel string `yaml:"Channel"`
	AppID   string `yaml:"AppID"`
}

type workers struct {
	NodeCount   int    `yaml:"NodeCount"`
	MachineType string `yaml:"MachineType"`
}
