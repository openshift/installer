package terraformgenerator

import (
	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

// Tectonic defines all variables for this platform.
type Tectonic struct {
	AutoScalingGroupExtraTags string `json:"tectonic_autoscaling_group_extra_tags,omitempty"`
	BaseDomain                string `json:"tectonic_base_domain,omitempty"`
	CACert                    string `json:"tectonic_ca_cert,omitempty"`
	CAKey                     string `json:"tectonic_ca_key,omitempty"`
	CAKeyAlg                  string `json:"tectonic_ca_key_alg,omitempty"`
	ClusterCIDR               string `json:"tectonic_cluster_cidr,omitempty"`
	ClusterName               string `json:"tectonic_cluster_name,omitempty"`
	ContainerLinuxChannel     string `json:"tectonic_container_linux_channel,omitempty"`
	ContainerLinuxVersion     string `json:"tectonic_container_linux_version,omitempty"`
	CustomCAPEMList           string `json:"tectonic_custom_ca_pem_list,omitempty"`
	DDNSKeyAlgorithm          string `json:"tectonic_ddns_key_algorithm,omitempty"`
	DDNSKeyName               string `json:"tectonic_ddns_key_name,omitempty"`
	DDNSKeySecret             string `json:"tectonic_ddns_key_secret,omitempty"`
	DDNSServer                string `json:"tectonic_ddns_server,omitempty"`
	DNSName                   string `json:"tectonic_dns_name,omitempty"`
	EtcdCACertPath            string `json:"tectonic_etcd_ca_cert_path,omitempty"`
	EtcdClientCertPath        string `json:"tectonic_etcd_client_cert_path,omitempty"`
	EtcdClientKeyPath         string `json:"tectonic_etcd_client_key_path,omitempty"`
	EtcdCount                 string `json:"tectonic_etcd_count,omitempty"`
	EtcdServers               string `json:"tectonic_etcd_servers,omitempty"`
	EtcdTLSEnabled            string `json:"tectonic_etcd_tls_enabled,omitempty"`
	HTTPProxyAddress          string `json:"tectonic_http_proxy_address,omitempty"`
	HTTPSProxyAddress         string `json:"tectonic_https_proxy_address,omitempty"`
	ISCSIEnabled              string `json:"tectonic_iscsi_enabled,omitempty"`
	LicensePath               string `json:"tectonic_license_path,omitempty"`
	MasterCount               int    `json:"tectonic_master_count,omitempty"`
	Networking                string `json:"tectonic_networking,omitempty"`
	NoProxy                   string `json:"tectonic_no_proxy,omitempty"`
	PullSecretPath            string `json:"tectonic_pull_secret_path,omitempty"`
	ServiceCIDR               string `json:"tectonic_service_cidr,omitempty"`
	SSHAuthorizedKey          string `json:"tectonic_ssh_authorized_key,omitempty"`
	TLSValidityPeriod         string `json:"tectonic_tls_validity_period,omitempty"`
	WorkerCount               int    `json:"tectonic_worker_count,omitempty"`
}

// NewTectonic returns the config for Tectonic.
func NewTectonic(cluster config.Cluster) Tectonic {
	return Tectonic{
		// AutoScalingGroupExtraTags: "",
		BaseDomain: cluster.DNS.BaseDomain,
		// CACert:                    "",
		// CAKey:                     "",
		// CAKeyAlg:                  "",
		ClusterCIDR:           cluster.Networking.NodeCIDR,
		ClusterName:           cluster.Name,
		ContainerLinuxChannel: cluster.ContainerLinux.Channel,
		ContainerLinuxVersion: cluster.ContainerLinux.Version,
		// CustomCAPEMList:           "",
		// DDNSKeyAlgorithm:          "",
		// DDNSKeyName:               "",
		// DDNSKeySecret:             "",
		// DDNSServer:                "",
		// DNSName:                   "",
		// EtcdCACertPath:            "",
		// EtcdClientCertPath:        "",
		// EtcdClientKeyPath:         "",
		// EtcdCount:                 "",
		// EtcdServers:               "",
		// EtcdTLSEnabled:            "",
		// HTTPProxyAddress:          "",
		// HTTPSProxyAddress:         "",
		// ISCSIEnabled:              "",
		LicensePath: cluster.Tectonic.LicensePath,
		MasterCount: cluster.Masters.NodeCount,
		Networking:  cluster.Networking.Type,
		// NoProxy:                   "",
		PullSecretPath: cluster.Tectonic.PullSecretPath,
		ServiceCIDR:    cluster.Networking.ServiceCIDR,
		// SSHAuthorizedKey:          "",
		// TLSValidityPeriod:         "",
		WorkerCount: cluster.Workers.NodeCount,
	}
}
