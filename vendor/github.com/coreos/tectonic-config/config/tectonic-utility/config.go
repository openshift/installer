package tectonicutility

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Kind is the TypeMeta.Kind for the OperatorConfig.
	Kind = "TectonicUtilityOperatorConfig"
	// APIVersion is the TypeMeta.APIVersion for the OperatorConfig.
	APIVersion = "v1"
)

// OperatorConfig defines the config for Tectonic Utility Operator.
type OperatorConfig struct {
	metav1.TypeMeta         `json:",inline"`
	IdentityConfig          `json:"identityConfig"`
	IngressConfig           `json:"ingressConfig"`
	StatsEmitterConfig      `json:"statsEmitterConfig"`
	TectonicConfigMapConfig `json:"tectonicConfigMap"`
	NetworkConfig           `json:"networkConfig"`
}

// IdentityConfig defines the config for Tectonic Identity.
type IdentityConfig struct {
	AdminEmail        string `json:"adminEmail"`
	AdminPasswordHash string `json:"adminPasswordHash"`
	AdminUserID       string `json:"adminUserID"`
	ConsoleClientID   string `json:"consoleClientID"`
	ConsoleSecret     string `json:"consoleSecret"`
	KubectlClientID   string `json:"kubectlClientID"`
	KubectlSecret     string `json:"kubectlSecret"`
}

// IngressConfig defines the config for Tectonic Ingress.
type IngressConfig struct {
	ConsoleBaseHost string `json:"consoleBaseHost"`
	IngressKind     string `json:"ingressKind"`
}

// StatsEmitterConfig defines the config for Tectonic Stats Emitter.
type StatsEmitterConfig struct {
	StatsURL string `json:"statsURL"`
}

// TectonicConfigMapConfig defines the variables that will be used by the Tectonic ConfigMap.
type TectonicConfigMapConfig struct {
	BaseAddress          string `json:"baseAddress"`
	CertificatesStrategy string `json:"certificatesStrategy"`
	ClusterID            string `json:"clusterID"`
	ClusterName          string `json:"clusterName"`
	IdentityAPIService   string `json:"identityAPIService"`
	InstallerPlatform    string `json:"installerPlatform"`
	KubeAPIServerURL     string `json:"kubeAPIserverURL"`
	TectonicVersion      string `json:"tectonicVersion"`
}

// NetworkConfig holds information on cluster networking.
// Copied from kube-core
type NetworkConfig struct {
	AdvertiseAddress string `json:"advertise_address"`
	ClusterCIDR      string `json:"cluster_cidr"`
	EtcdServers      string `json:"etcd_servers"`
	ServiceCIDR      string `json:"service_cidr"`
}
