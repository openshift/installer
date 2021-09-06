package kubevirt

// Platform stores all the global configuration used by the kubevirt platform installation.
type Platform struct {
	// Namespace is the namespace in the infra cluster, which the control plane (master vms)
	// and the compute (worker vms) are installed in.
	Namespace string `json:"namespace"`

	// StorageClass is the Storage Class used in the infra cluster.
	// +optional
	StorageClass string `json:"storageClass,omitempty"`

	// NetworkName is the target network of all the network interfaces of the nodes.
	NetworkName string `json:"networkName"`

	// InterfaceBindingMethod is the the interface binding method of the nodes of the tenantcluster (Bridge | SRIOV).
	// +optional
	InterfaceBindingMethod string `json:"interfaceBindingMethod"`

	// APIVIP is the virtual IP address for the api endpoint.
	// +kubebuilder:validation:Format=ip
	APIVIP string `json:"apiVIP"`

	// IngressIP is an external IP which routes to the default ingress controller.
	// +kubebuilder:validation:Format=ip
	IngressVIP string `json:"ingressVIP"`

	// PersistentVolumeAccessMode is the access mode should be use with the persistent volumes.
	// +kubebuilder:default="ReadWriteMany"
	// +optional
	PersistentVolumeAccessMode string `json:"persistentVolumeAccessMode,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on Kubevirt for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
}
