package terraformgenerator

import (
	"github.com/coreos/tectonic-installer/installer/pkg/config"
)

// Tectonic defines all variables for this platform.
type Tectonic struct {
	BaseDomain            string   `json:"tectonic_base_domain,omitempty"`
	CACert                string   `json:"tectonic_ca_cert,omitempty"`
	CAKey                 string   `json:"tectonic_ca_key,omitempty"`
	CAKeyAlg              string   `json:"tectonic_ca_key_alg,omitempty"`
	ClusterCIDR           string   `json:"tectonic_cluster_cidr,omitempty"`
	ClusterName           string   `json:"tectonic_cluster_name,omitempty"`
	ContainerLinuxChannel string   `json:"tectonic_container_linux_channel,omitempty"`
	ContainerLinuxVersion string   `json:"tectonic_container_linux_version,omitempty"`
	CustomCAPEMList       string   `json:"tectonic_custom_ca_pem_list,omitempty"`
	DDNSKeyAlgorithm      string   `json:"tectonic_ddns_key_algorithm,omitempty"`
	DDNSKeyName           string   `json:"tectonic_ddns_key_name,omitempty"`
	DDNSKeySecret         string   `json:"tectonic_ddns_key_secret,omitempty"`
	DDNSServer            string   `json:"tectonic_ddns_server,omitempty"`
	DNSName               string   `json:"tectonic_dns_name,omitempty"`
	EtcdCACertPath        string   `json:"tectonic_etcd_ca_cert_path,omitempty"`
	EtcdClientCertPath    string   `json:"tectonic_etcd_client_cert_path,omitempty"`
	EtcdClientKeyPath     string   `json:"tectonic_etcd_client_key_path,omitempty"`
	EtcdCount             int      `json:"tectonic_etcd_count,omitempty"`
	EtcdServers           []string `json:"tectonic_etcd_servers,omitempty"`
	HTTPProxyAddress      string   `json:"tectonic_http_proxy_address,omitempty"`
	HTTPSProxyAddress     string   `json:"tectonic_https_proxy_address,omitempty"`
	ISCSIEnabled          bool     `json:"tectonic_iscsi_enabled,omitempty"`
	LicensePath           string   `json:"tectonic_license_path,omitempty"`
	MasterCount           int      `json:"tectonic_master_count,omitempty"`
	Networking            string   `json:"tectonic_networking,omitempty"`
	NoProxy               string   `json:"tectonic_no_proxy,omitempty"`
	PullSecretPath        string   `json:"tectonic_pull_secret_path,omitempty"`
	ServiceCIDR           string   `json:"tectonic_service_cidr,omitempty"`
	TLSValidityPeriod     int      `json:"tectonic_tls_validity_period,omitempty"`
	WorkerCount           int      `json:"tectonic_worker_count,omitempty"`
}

// NewTectonic returns the config for Tectonic.
func NewTectonic(cluster config.Cluster) Tectonic {
	return Tectonic{
		BaseDomain:            cluster.BaseDomain,
		CACert:                cluster.CA.Cert,
		CAKey:                 cluster.CA.Key,
		CAKeyAlg:              cluster.CA.KeyAlg,
		ClusterCIDR:           cluster.Networking.PodCIDR,
		ClusterName:           cluster.Name,
		ContainerLinuxChannel: cluster.ContainerLinux.Channel,
		ContainerLinuxVersion: cluster.ContainerLinux.Version,
		CustomCAPEMList:       cluster.CustomCAPEMList,
		DDNSKeyAlgorithm:      cluster.DDNS.Key.Algorithm,
		DDNSKeyName:           cluster.DDNS.Key.Name,
		DDNSKeySecret:         cluster.DDNS.Key.Secret,
		DDNSServer:            cluster.DDNS.Server,
		DNSName:               cluster.DNSName,
		EtcdCACertPath:        cluster.Etcd.External.CACertPath,
		EtcdClientCertPath:    cluster.Etcd.External.ClientCertPath,
		EtcdClientKeyPath:     cluster.Etcd.External.ClientKeyPath,
		EtcdCount:             cluster.Etcd.Count,
		EtcdServers:           cluster.Etcd.External.Servers,
		HTTPProxyAddress:      cluster.Proxy.HTTP,
		HTTPSProxyAddress:     cluster.Proxy.HTTP,
		ISCSIEnabled:          cluster.ISCSI.Enabled,
		LicensePath:           cluster.LicensePath,
		MasterCount:           cluster.Master.Count,
		Networking:            cluster.Networking.Type,
		NoProxy:               cluster.Proxy.No,
		PullSecretPath:        cluster.PullSecretPath,
		ServiceCIDR:           cluster.Networking.ServiceCIDR,
		TLSValidityPeriod:     cluster.TLSValidityPeriod,
		WorkerCount:           cluster.Worker.Count,
	}
}
