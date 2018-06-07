package kubecore

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Kind is the TypeMeta.Kind for the OperatorConfig.
	Kind = "KubeCoreOperatorConfig"
	// APIVersion is the TypeMeta.APIVersion for the OperatorConfig.
	APIVersion = "v1"
)

// OperatorConfig holds all configuration needed for the operator to make any install / upgrade time
// decisions.
type OperatorConfig struct {
	metav1.TypeMeta     `json:",inline"`
	ClusterConfig       `json:"clusterConfig,omitempty"`
	DNSConfig           `json:"dnsConfig,omitempty"`
	AuthConfig          `json:"authConfig,omitempty"`
	RoutingConfig       `json:"routingConfig,omitempty"`
	CloudProviderConfig `json:"cloudProviderConfig,omitempty"`
	NetworkConfig       `json:"networkConfig,omitempty"`
}

// AuthConfig holds Authentication related config values.
type AuthConfig struct {
	OIDCClientID      string `json:"oidc_client_id"`
	OIDCIssuerURL     string `json:"oidc_issuer_url"`
	OIDCGroupsClaim   string `json:"oidc_groups_claim"`
	OIDCUsernameClaim string `json:"oidc_username_claim"`
}

// CloudProviderConfig holds information on the cloud provider this cluster is operating in.
type CloudProviderConfig struct {
	CloudConfigPath      string `json:"cloud_config_path"`
	CloudProviderProfile string `json:"cloud_provider_profile"`
}

// NetworkConfig holds information on cluster networking.
type NetworkConfig struct {
	AdvertiseAddress string `json:"advertise_address"`
	ClusterCIDR      string `json:"cluster_cidr"`
	EtcdServers      string `json:"etcd_servers"`
	ServiceCIDR      string `json:"service_cidr"`
}

// ClusterConfig holds global/general information about the cluster.
type ClusterConfig struct {
	APIServerURL string `json:"apiserver_url"`
}

// DNSConfig options for the dns configuration
type DNSConfig struct {
	// ClusterIP ip address of the cluster
	ClusterIP string `json:"clusterIP"`
}

// RoutingConfig holds options for routes.
type RoutingConfig struct {
	// Subdomain is the suffix appended to $service.$namespace. to form the default route hostname
	// if empty, router.tectonic-ingress.cluster.local
	Subdomain string `json:"subdomain"`
}
