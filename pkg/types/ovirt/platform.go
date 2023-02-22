package ovirt

import (
	configv1 "github.com/openshift/api/config/v1"
)

// Platform stores all the global configuration that all
// machinesets use.
type Platform struct {
	// The target cluster under which all VMs will run
	ClusterID string `json:"ovirt_cluster_id"`

	// The target storage domain under which all VM disk would be created.
	StorageDomainID string `json:"ovirt_storage_domain_id"`

	// NetworkName is the target network of all the network interfaces of the nodes.
	// When no ovirt_network_name is provided it defaults to `ovirtmgmt` network, which is a default network for every ovirt cluster.
	// +optional
	NetworkName string `json:"ovirt_network_name,omitempty"`

	// VNICProfileID defines the VNIC profile ID to use the the VM network interfaces.
	// When no vnicProfileID is provided it will be set to the profile of the network. If there are multiple
	// profiles for the network, the installer requires you to explicitly set the vnicProfileID.
	// +optional
	VNICProfileID string `json:"vnicProfileID,omitempty"`

	// DeprecatedAPIVIP is an IP which will be served by bootstrap and then pivoted masters, using keepalived
	// Deprecated: Use APIVIPs
	//
	// +kubebuilder:validation:Format=ip
	// +optional
	DeprecatedAPIVIP string `json:"api_vip,omitempty"`

	// APIVIPs contains the VIP(s) which will be served by bootstrap and then
	// pivoted masters, using keepalived. In dual stack clusters it contains an
	// IPv4 and IPv6 address, otherwise only one VIP
	//
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:UniqueItems=true
	// +kubebuilder:validation:Format=ip
	// +optional
	APIVIPs []string `json:"api_vips,omitempty"`

	// IngressIP is an external IP which routes to the default ingress controller.
	// The IP is a suitable target of a wildcard DNS record used to resolve default route host names.
	// Deprecated: Use IngressVIPs
	//
	// +kubebuilder:validation:Format=ip
	// +optional
	DeprecatedIngressVIP string `json:"ingress_vip,omitempty"`

	// IngressVIPs are external IP(s) which route to the default ingress
	// controller. The VIPs are suitable targets of wildcard DNS records used to
	// resolve default route host names. In dual stack clusters it contains an
	// IPv4 and IPv6 address, otherwise only one VIP
	//
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:UniqueItems=true
	// +kubebuilder:validation:Format=ip
	// +optional
	IngressVIPs []string `json:"ingress_vips,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on ovirt for machine pools which do not define their
	// own platform configuration.
	// Default will set the image field to the latest RHCOS image.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// AffinityGroups contains the RHV affinity groups that the installer will create.
	// +optional
	AffinityGroups []AffinityGroup `json:"affinityGroups"`

	// LoadBalancer defines how the load balancer used by the cluster is configured.
	// LoadBalancer is available in TechPreview.
	// +optional
	LoadBalancer *configv1.OvirtPlatformLoadBalancer `json:"loadBalancer,omitempty"`
}

// AffinityGroup defines the affinity group that the installer will create
type AffinityGroup struct {
	// Name name of the affinity group
	Name string `json:"name"`
	// Priority of the affinity group, needs to be between 1 (lowest) - 5 (highest)
	Priority int `json:"priority"`
	// Description of the affinity group
	// +optional
	Description string `json:"description,omitempty"`
	// Enforcing whether to create a hard affinity rule, default is false
	// +optional
	Enforcing bool `json:"enforcing"`
}
