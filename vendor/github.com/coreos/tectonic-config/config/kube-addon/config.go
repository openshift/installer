package kubeaddon

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	// Kind is the TypeMeta.Kind for the OperatorConfig.
	Kind = "KubeAddonOperatorConfig"
	// APIVersion is the TypeMeta.APIVersion for the OperatorConfig.
	APIVersion = "v1"
)

// OperatorConfig contains configuration for KAO managed add-ons
type OperatorConfig struct {
	metav1.TypeMeta    `json:",inline"`
	ClusterConfig      `json:"clusterConfig,omitempty"`
	CloudProvider      string `json:"cloudProvider,omitempty"`
	RegistryHTTPSecret string `json:"registryHTTPSecret,omitempty"`
}

// ClusterConfig holds global/general information about the cluster.
type ClusterConfig struct {
	APIServerURL string `json:"apiserver_url"`
}
