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
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type ca struct {
	Cert   string `yaml:"cert"`
	Key    string `yaml:"key"`
	KeyAlg string `yaml:"keyAlg"`
}

// Cluster defines the config for a cluster.
type Cluster struct {
	AWS               aws.Config       `yaml:"aws,omitempty"`
	Admin             admin            `yaml:"admin"`
	Azure             azure.Config     `yaml:"azure,omitempty"`
	BaseDomain        string           `yaml:"baseDomain"`
	CA                ca               `yaml:"ca"`
	ContainerLinux    containerLinux   `yaml:"containerLinux"`
	CustomCAPEMList   string           `yaml:"customCAPEMList"`
	DDNS              ddns             `yaml:"ddns"`
	DNSName           string           `yaml:"dnsName"`
	Etcd              etcd             `yaml:"etcd"`
	GCP               gcp.Config       `yaml:"gcp,omitempty"`
	GovCloud          govcloud.Config  `yaml:"govcloud,omitempty"`
	ISCSI             iscsi            `yaml:"iscsi"`
	LicensePath       string           `yaml:"licensePath"`
	Master            master           `yaml:"master"`
	Metal             metal.Config     `yaml:"metal,omitempty"`
	Name              string           `yaml:"name"`
	Networking        networking       `yaml:"networking"`
	OpenStack         openstack.Config `yaml:"openstack,omitempty"`
	Platform          string           `yaml:"platform"`
	Proxy             proxy            `yaml:"proxy"`
	PullSecretPath    string           `yaml:"pullSecretPath"`
	TLSValidityPeriod int              `yaml:"tlsValidityPeriod"`
	VMware            vmware.Config    `yaml:"vmware,omitempty"`
	Worker            worker           `yaml:"worker"`
}

// Config defines the top level config for a configuration file.
type Config struct {
	Clusters []Cluster `yaml:"clusters"`
}

type containerLinux struct {
	Channel string `yaml:"channel"`
	Version string `yaml:"version"`
}

type ddns struct {
	Key    ddnsKey `yaml:"key"`
	Server string  `yaml:"secret"`
}

type ddnsKey struct {
	Algorithm string `yaml:"algorithm"`
	Name      string `yaml:"name"`
	Secret    string `yaml:"secret"`
}

type etcd struct {
	Count    int          `yaml:"count"`
	External etcdExternal `yaml:"external"`
}

type etcdExternal struct {
	CACertPath     string   `yaml:"caCertPath"`
	ClientCertPath string   `yaml:"clientCertPath"`
	ClientKeyPath  string   `yaml:"clientKeyPath"`
	Servers        []string `yaml:"servers"`
}

type iscsi struct {
	Enabled bool `yaml:"enabled"`
}

type master struct {
	Count int `yaml:"count"`
}

type networking struct {
	Type        string `yaml:"type"`
	MTU         string `yaml:"mtu"`
	ServiceCIDR string `yaml:"serviceCIDR"`
	PodCIDR     string `yaml:"podCIDR"`
}

type proxy struct {
	HTTP  string `yaml:"http"`
	HTTPS string `yaml:"https"`
	No    string `yaml:"no"`
}

type worker struct {
	Count int `yaml:"count"`
}
