package equinixmetal

type Platform struct {
	// Metro represents the Equinix Metal metro code for the location where your devices will be provisioned
	// (https://metal.equinix.com/developers/docs/getting-started/metros/)
	Metro string `json:"metro"`

	// Facility represents the Equinix Metal facility code for the region and
	// datacenter where your devices will be provisioned
	// (https://metal.equinix.com/developers/docs/getting-started/facilities/)
	Facility string `json:"facility"`

	// ProjectID represents the Equinix Metal project used for logical grouping and invoicing (https://metal.equinix.com/developers/docs/API/getting-started/)
	ProjectID string `json:"project_id"`

	// APIVIP is the static IP that was provisioned as the persistant endpoint for the cluster API and ignition
	// +kubebuilder:validation:Format=ip
	APIVIP string `json:"apivip"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on Equinix Metal for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// ClusterOSImage is a URL to override the default OS image
	// for cluster nodes. The URL must contain a sha256 hash of the image
	// e.g https://mirror.example.com/images/metal.qcow2.gz?sha256=3b5a8...
	//
	// +optional
	ClusterOSImage string `json:"clusterOSImage,omitempty" validate:"omitempty,osimageuri,urlexist"`
}
