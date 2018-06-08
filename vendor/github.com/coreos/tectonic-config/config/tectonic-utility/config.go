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
	StatsEmitterConfig      `json:"statsEmitterConfig"`
	TectonicConfigMapConfig `json:"tectonicConfigMap"`
}

// StatsEmitterConfig defines the config for Tectonic Stats Emitter.
type StatsEmitterConfig struct {
	StatsURL string `json:"statsURL"`
}

// TectonicConfigMapConfig defines the variables that will be used by the Tectonic ConfigMap.
type TectonicConfigMapConfig struct {
	ClusterID            string `json:"clusterID"`
	ClusterName          string `json:"clusterName"`
	CertificatesStrategy string `json:"certificatesStrategy"`
	InstallerPlatform    string `json:"installerPlatform"`
	KubeAPIServerURL     string `json:"kubeAPIserverURL"`
	TectonicVersion      string `json:"tectonicVersion"`
}
