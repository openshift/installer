package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Important: Run hack/update-codegen.sh to regenerate code after modifying this file

// register our type with the k8s api scheme
func init() {
	SchemeBuilder.Register(&NetworkConfig{}, &NetworkConfigList{})
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkConfig describes the cluster's desired network configuration
// +k8s:openapi-gen=true
type NetworkConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkConfigSpec   `json:"spec,omitempty"`
	Status NetworkConfigStatus `json:"status,omitempty"`
}

// NetworkConfigStatus defines the observed state of NetworkConfig
type NetworkConfigStatus struct {
	// TODO
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NetworkConfigList contains a list of NetworkConfig
// We do not support more than one NetworkConfig, but the operator-sdk
// requires this
type NetworkConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkConfig `json:"items"`
}

// NetworkConfigSpec is the top-level network configuration object.
type NetworkConfigSpec struct {
	// IP address pool to use for pod IPs.
	// Some network providers, e.g. Openshift-sdn, support multiple ClusterNetworks.
	// Others only support one. This is equivalent to the cluster-cidr.
	ClusterNetworks []ClusterNetwork `json:"clusterNetworks"`

	// The CIDR to use for services
	ServiceNetwork string `json:"serviceNetwork"`

	// The "default" network that all pods will receive
	DefaultNetwork DefaultNetworkDefinition `json:"defaultNetwork"`

	// Additional networks to make available to pods. If they are specified,
	// pods can request them via annotations.
	//
	// Specifying any additionalNetworks will enable Multus across the cluster.
	AdditionalNetworks []AdditionalNetworkDefinition `json:"additionalNetworks"`

	// DeployKubeProxy specifies whether or not a standalone kube-proxy should
	// be deployed by the operator. Some network providers include kube-proxy
	// or similar functionality. If unset, the plugin will attempt to select
	// the correct value, which is false when Openshift-sdn and ovn-kubernetes are
	// used and true otherwise.
	// +optional
	DeployKubeProxy *bool `json:"deployKubeProxy,omitempty"`

	// KubeProxyConfig lets us configure desired proxy configuration.
	// If not specified, sensible defaults will be chosen by OpenShift directly.
	// Not consumed by all network providers - currently only openshift-sdn.
	KubeProxyConfig *ProxyConfig `json:"kubeProxyConfig,omitempty"`
}

// ClusterNetwork is a subnet from which to allocate PodIPs. A network of size
// 2^HostSubnetLength will be allocated when nodes join the cluster.
// Not all network providers support multiple ClusterNetworks
type ClusterNetwork struct {
	CIDR             string `json:"cidr"`
	HostSubnetLength uint32 `json:"hostSubnetLength"`
}

// NetworkDefinition represents a single network plugin's configuration.
// Kind must be specified, along with exactly one "Config" that matches
// the kind. Kinds that do not have a specific configuration parameter should
// use OtherConfig
type DefaultNetworkDefinition struct {
	// The type of network
	// All NetworkTypes are supported except for NetworkTypeRaw
	Type NetworkType `json:"type"`

	// OpenshiftSDNConfig configures the openshift-sdn plugin
	// +optional
	OpenshiftSDNConfig *OpenshiftSDNConfig `json:"openshiftSDNConfig,omitempty"`

	// OVNKubernetesConfig configures the ovn-kubernetes plugin
	// +optional
	OVNKubernetesConfig *OVNKubernetesConfig `json:"ovnKubernetesConfig,omitempty"`

	// OtherConfig is for network plugins that are supported by the operator
	// but do not need their own type. These values will be passed directly
	// to the manifest templates.
	// This is used by calico and kuryr
	// See the plugin-specific documentation for which values are required.
	// +optional
	OtherConfig map[string]string `json:"otherConfig,omitEmpty"`
}

// AdditionalNetworkDefinition is extra networks that are available but not
// created by default. Instead, pods must request them by name.
type AdditionalNetworkDefinition struct {
	// The type of network
	// The only supported value is NetworkTypeRaw
	Type NetworkType `json:"type"`

	// The name of the network. This will be populated in the resulting CRD
	Name string `json:"name"`

	// RawCNIConfig is the raw CNI configuration json to create in the
	// NetworkAttachmentDefinition CRD
	RawCNIConfig string `json:"rawCNIConfig"`
}

// OpenshiftSDNConfig configures the three openshift-sdn plugins
type OpenshiftSDNConfig struct {
	// Mode is one of "multitenant", "subnet", or "networkpolicy"
	Mode SDNMode `json:"mode"`

	// VXLANPort is the port to use for all vxlan packets. The default
	// is 4789
	// +optional
	VXLANPort *uint32 `json:"vxlanPort,omitempty"`

	// MTU is the mtu to use for the tunnel interface. Defaults to 1450 if unset.
	// +optional
	MTU *uint32 `json:"mtu,omitempty"`

	// If true, then assume the nodes already have a running openvswitch.
	// +optional
	UseExternalOpenvswitch *bool `json:"useExternalOpenvswitch,omitempty"`
}

// OVNKubernetesConfig is the configuration parameters for networks using the
// ovn-kubernetes netwok project
type OVNKubernetesConfig struct {
	// The UDP port to use for geneve
	// The default is 6081
	GenevePort *uint32 `json:"genevePort,omitempty"`

	// The MTU to use for the tunnel interface
	// Default is 1400
	MTU *uint32 `json:"mtu,omitempty"`
}

// NetworkType describes the network plugin type to configure
type NetworkType string

// ProxyConfig defines the configuration knobs for kubeproxy
// All of these are optional and have sensible defaults
type ProxyConfig struct {
	// The period that iptables rules are refreshed.
	// Default: 30s
	IptablesSyncPeriod string `json:"iptablesSyncPeriod,omitempty"`

	// The address to "bind" on
	// Defaults to 0.0.0.0
	BindAddress string

	// Any additional arguments to pass to the kubeproxy process
	ProxyArguments map[string][]string
}

const (
	// NetworkTypeOpenshiftSDN means the openshift-sdn plugin will be configured
	NetworkTypeOpenshiftSDN NetworkType = "OpenshiftSDN"

	// NetworkTypeOVNKubernetes means the ovn-kubernetes project will be configured
	NetworkTypeOVNKubernetes NetworkType = "OVNKubernetes"

	// NetworkTypeCalico means Calico will be configured
	NetworkTypeCalico NetworkType = "Calico"

	// NetworkType
	NetworkTypeKuryr NetworkType = "Kuryr"

	// NetworkTypeRaw
	NetworkTypeRaw NetworkType = "Raw"
)

// SDNMode is the Mode the openshift-sdn plugin is in
type SDNMode string

const (
	SDNModeMultitenant SDNMode = "Multitenant"

	SDNModeSubnet SDNMode = "Subnet"

	SDNModePolicy SDNMode = "Networkpolicy"
)
