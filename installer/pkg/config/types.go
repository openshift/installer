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

type admin struct {
	Email    string `yaml:"Email"`
	Password string `yaml:"Password"`
}

type ca struct {
	Cert   string `yaml:"Cert"`
	Key    string `yaml:"Key"`
	KeyAlg string `yaml:"KeyAlg"`
}

// Cluster defines the config for a cluster.
type Cluster struct {
	AWS               aws.Config       `yaml:"AWS,omitempty"`
	Admin             admin            `yaml:"Admin"`
	Azure             azure.Config     `yaml:"Azure,omitempty"`
	BaseDomain        string           `yaml:"BaseDomain"`
	CA                ca               `yaml:"CA"`
	ContainerLinux    containerLinux   `yaml:"ContainerLinux"`
	CustomCAPEMList   string           `yaml:"CustomCAPEMList"`
	DDNS              ddns             `yaml:"DDNS"`
	DNSName           string           `yaml:"DNSName"`
	Etcd              etcd             `yaml:"Etcd"`
	GCP               gcp.Config       `yaml:"GCP,omitempty"`
	GovCloud          govcloud.Config  `yaml:"GovCloud,omitempty"`
	ISCSI             iscsi            `yaml:"ISCSI"`
	LicensePath       string           `yaml:"LicensePath"`
	Master            master           `yaml:"Master"`
	Metal             metal.Config     `yaml:"Metal,omitempty"`
	Name              string           `yaml:"Name"`
	Networking        networking       `yaml:"Networking"`
	OpenStack         openstack.Config `yaml:"OpenStack,omitempty"`
	Platform          string           `yaml:"Platform"`
	Proxy             proxy            `yaml:"Proxy"`
	PullSecretPath    string           `yaml:"PullSecretPath"`
	TLSValidityPeriod int              `yaml:"TLSValidityPeriod"`
	VMware            vmware.Config    `yaml:"VMware,omitempty"`
	Worker            worker           `yaml:"Worker"`
}

// Config defines the top level config for a configuration file.
type Config struct {
	Clusters []Cluster `yaml:"Clusters"`
}

type containerLinux struct {
	Channel string `yaml:"Channel"`
	Version string `yaml:"Version"`
}

type ddns struct {
	Key    ddnsKey `yaml:"Key"`
	Server string  `yaml:"Secret"`
}

type ddnsKey struct {
	Algorithm string `yaml:"Algorithm"`
	Name      string `yaml:"Name"`
	Secret    string `yaml:"Secret"`
}

type etcd struct {
	Count    int          `yaml:"Count"`
	External etcdExternal `yaml:"External"`
}

type etcdExternal struct {
	CACertPath     string   `yaml:"CACertPath"`
	ClientCertPath string   `yaml:"ClientCertPath"`
	ClientKeyPath  string   `yaml:"ClientKeyPath"`
	Servers        []string `yaml:"Servers"`
}

type iscsi struct {
	Enabled bool `yaml:"Enabled"`
}

type master struct {
	Count int `yaml:"Count"`
}

type networking struct {
	Type        string `yaml:"Type"`
	MTU         string `yaml:"MTU"`
	ServiceCIDR string `yaml:"ServiceCIDR"`
	PodCIDR     string `yaml:"PodCIDR"`
}

type proxy struct {
	HTTP  string `yaml:"HTTP"`
	HTTPS string `yaml:"HTTPS"`
	No    string `yaml:"No"`
}

type worker struct {
	Count int `yaml:"Count"`
}
