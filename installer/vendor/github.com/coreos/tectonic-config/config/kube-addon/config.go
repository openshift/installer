package kubeaddon

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Kind is the TypeMeta.Kind for the OperatorConfig.
	Kind = "KubeAddonOperatorConfig"
	// APIVersion is the TypeMeta.APIVersion for the OperatorConfig.
	APIVersion = "v1"
)

// OperatorConfig contains configuration for KAO managed add-ons
type OperatorConfig struct {
	metav1.TypeMeta `json:",inline"`
	DNSConfig       `json:"dnsConfig,omitempty"`
	CloudProvider   string `json:"cloudProvider,omitempty"`
}

// DNSConfig options for the dns configuration
type DNSConfig struct {
	// ClusterIP ip address of the cluster
	ClusterIP string `json:"clusterIP"`
}
