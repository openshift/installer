package tnc

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// TNCConfigKind is the TypeMeta.Kind for the TNC config.
	TNCConfigKind = "TectonicNodeControllerConfig"
	// TNCConfigAPIVersion is the TypeMeta.APIVersion for the TNC config.
	TNCConfigAPIVersion = "v1"

	// TNCOConfigKind is the TypeMeta.Kind for the TNC Operator's config.
	TNCOConfigKind = "TectonicNodeControllerOperatorConfig"
	// TNCOConfigAPIVersion is the TypeMeta.APIVersion for the TNC Operator's config.
	TNCOConfigAPIVersion = "v1"
)

// ControllerConfig is the config structure TNC to generate NodeConfigs
type ControllerConfig struct {
	metav1.TypeMeta     `json:",inline"`
	ClusterDNSIP        string `json:"clusterDNSIP"`
	CloudProvider       string `json:"cloudProvider"`
	CloudProviderConfig string `json:"cloudProviderConfig"`
	ClusterName         string `json:"clusterName"`

	BaseDomain string `json:"baseDomain"`

	// Size of the initial etcd cluster.
	EtcdInitialCount int `json:"etcdInitialCount"`

	// User customizations, list of node configs in the cluster to collect and apply as part of the final config
	AdditionalConfigs []string `json:"additionalConfigs"`

	// Pool update limits allowing individual roles to have a specific percentage or constant number of nodes that can be upgrading at a given time. If an empty string is provided only a single node will be updated at a given time
	NodePoolUpdateLimit map[string]string `json:"nodePoolUpdateLimit"`
}

// OperatorConfig defines the configuration needed by the Tectonic Node Controller Operator.
type OperatorConfig struct {
	metav1.TypeMeta `json:",inline"`

	// ControllerConfig is the config that will be consumed by the TNC.
	ControllerConfig `json:"controllerConfig"`
}
