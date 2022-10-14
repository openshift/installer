package gcp

// CreateFirewallRules specifies if the installer should create firewall rules.
// +kubebuilder:validation:Enum="Enabled";"Disabled"
type CreateFirewallRules string

const (
	// CreateFirewallRulesEnabled is Enabled
	CreateFirewallRulesEnabled CreateFirewallRules = "Enabled"
	// CreateFirewallRulesDisabled is Disabled
	CreateFirewallRulesDisabled CreateFirewallRules = "Disabled"
)

// DNSZone stores the information common and required to create DNS zones including
// the project and id/name of the zone.
type DNSZone struct {
	// ID Technology Preview.
	// ID or name of the zone.
	// +optional
	ID string `json:"id,omitempty"`

	// ProjectID Technology Preview.
	// When the ProjectID is provided, the zone will exist in this project. When the ProjectID is
	// empty, the ProjectID defaults to the Service Project (GCP.ProjectID).
	// +optional
	ProjectID string `json:"project,omitempty"`
}

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// ProjectID is the the project that will be used for the cluster.
	ProjectID string `json:"projectID"`

	// Region specifies the GCP region where the cluster will be created.
	Region string `json:"region"`

	// CreateFirewallRules specifies if the installer should create the
	// cluster firewall rules in the gcp cloud network.
	// +optional
	CreateFirewallRules CreateFirewallRules `json:"createFirewallRules,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on GCP for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// Network specifies an existing VPC where the cluster should be created
	// rather than provisioning a new one.
	// +optional
	Network string `json:"network,omitempty"`

	// NetworkProjectID is currently TechPreview.
	// NetworkProjectID specifies which project the network and subnets exist in when
	// they are not in the main ProjectID.
	// +optional
	NetworkProjectID string `json:"networkProjectID,omitempty"`

	// ControlPlaneSubnet is an existing subnet where the control plane will be deployed.
	// The value should be the name of the subnet.
	// +optional
	ControlPlaneSubnet string `json:"controlPlaneSubnet,omitempty"`

	// ComputeSubnet is an existing subnet where the compute nodes will be deployed.
	// The value should be the name of the subnet.
	// +optional
	ComputeSubnet string `json:"computeSubnet,omitempty"`

	// Licenses is a list of licenses to apply to the compute images
	// The value should a list of strings (https URLs only) representing the license keys.
	// When set, this will cause the installer to copy the image into user's project.
	// This option is incompatible with any mechanism that makes use of pre-built images
	// such as the current env OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE
	// +optional
	Licenses []string `json:"licenses,omitempty"`

	// PrivateDNSZone Technology Preview.
	// PrivateDNSZone contains the zone ID and project where the Private DNS zone records will be created.
	// +optional
	PrivateDNSZone *DNSZone `json:"privateDNSZone,omitempty"`

	// PublicDNSZone Technology Preview.
	// PublicDNSZone contains the zone ID and project where the Public DNS zone records will be created.
	// +optional
	PublicDNSZone *DNSZone `json:"publicDNSZone,omitempty"`
}
