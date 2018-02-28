package tectonicnetwork

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// Kind is the TypeMeta.Kind for the OperatorConfig.
	Kind = "TectonicNetworkOperatorConfig"
	// APIVersion is the TypeMeta.APIVersion for the OperatorConfig.
	APIVersion = "v1"
)

// NetworkType indicates the network configuration of the cluster.
//
// NOTE: only one of none, flannel, canal or calico can be enabled at a time.
type NetworkType string

const (
	// NetworkNone is the network profile for a cluster that does not use the TNO to configure
	// networking.
	NetworkNone NetworkType = "none"
	// NetworkFlannel is the network profile for a cluster that implements flannel.
	NetworkFlannel NetworkType = "flannel"
	// NetworkCanal is the network profile for a cluster that implements canal.
	NetworkCanal = "canal"
	// NetworkCalicoIPIP is the network profile for a cluster that implements calico.
	NetworkCalicoIPIP = "calico-ipip"
)

// OperatorConfig defines the configuration needed by the Tectonic Network Operator.
type OperatorConfig struct {
	metav1.TypeMeta `json:",inline"`

	// PodCIDR is an IP range from which pod IPs can be assigned.
	PodCIDR string `json:"podCIDR"`
	// NetworkProfile describes the network configuration for the cluster.
	NetworkProfile NetworkType `json:"networkProfile"`
	// CalicoConfig is used only when the networkType is `calico`.
	CalicoConfig `json:"calicoConfig"`
}

// CalicoConfig defines config values when the network profile supports `calico`.
type CalicoConfig struct {
	// MTU sets the MTU size for workload interfaces and the IP-in-IP tunnel device.
	MTU string `json:"mtu"`
}
