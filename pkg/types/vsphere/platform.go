package vsphere

// DiskType is a disk provisioning type for vsphere.
// +kubebuilder:validation:Enum="";thin;thick;eagerZeroedThick
type DiskType string

// FailureDomainType is the name of the failure domain type.
// There are two defined failure domains currently, Datacenter and ComputeCluster.
// Each represents a vCenter object type within a vSphere environment.
// +kubebuilder:validation:Enum=HostGroup;Datacenter;ComputeCluster
type FailureDomainType string

// DeploymentSuitable defines if a type of machine is suitable for a given DeploymentZone
// +kubebuilder:validation:Enum=Allowed;NotAllowed
type DeploymentSuitable string

const (
	// DiskTypeThin uses Thin disk provisioning type for vsphere in the cluster.
	DiskTypeThin DiskType = "thin"

	// DiskTypeThick uses Thick disk provisioning type for vsphere in the cluster.
	DiskTypeThick DiskType = "thick"

	// DiskTypeEagerZeroedThick uses EagerZeroedThick disk provisioning type for vsphere in the cluster.
	DiskTypeEagerZeroedThick DiskType = "eagerZeroedThick"

	// HostGroupFailureDomain as a type allows the use of a group of ESXi hosts
	// to be represented as a failure domain zone. When using this
	// case it is expected that region would be a cluster.
	// HostGroups within vCenter must be preconfigured and
	// assigned in the topology.
	HostGroupFailureDomain FailureDomainType = "HostGroup"

	// ComputeClusterFailureDomain failure domain can either be a zone or region.
	// The vCenter cluster is required to preconfigured and
	// assigned in the topology.
	ComputeClusterFailureDomain FailureDomainType = "ComputeCluster"

	// DatacenterFailureDomain failure domain can be only a region. The vcenter
	// datacenter is required to be preconfigred and assigned
	// in the topology. If used the zone would be of type ComputeCluster.
	DatacenterFailureDomain FailureDomainType = "Datacenter"

	// Allowed indicates that the deployment is suitable for
	// control plane nodes.
	Allowed DeploymentSuitable = "Allowed"

	// NotAllowed indicates that the deployment is not suitable for
	// control plane nodes.
	NotAllowed DeploymentSuitable = "NotAllowed"
)

// Platform stores any global configuration used for vsphere platforms
type Platform struct {
	// VCenter is the domain name or IP address of the vCenter.
	VCenter string `json:"vCenter"`

	// Username is the name of the user to use to connect to the vCenter.
	Username string `json:"username"`

	// Password is the password for the user to use to connect to the vCenter.
	Password string `json:"password"`

	// Datacenter is the name of the datacenter to use in the vCenter.
	Datacenter string `json:"datacenter"`

	// DefaultDatastore is the default datastore to use for provisioning volumes.
	DefaultDatastore string `json:"defaultDatastore"`

	// Folder is the absolute path of the folder that will be used and/or created for
	// virtual machines. The absolute path is of the form /<datacenter>/vm/<folder>/<subfolder>.
	Folder string `json:"folder,omitempty"`

	// Cluster is the name of the cluster virtual machines will be cloned into.
	Cluster string `json:"cluster,omitempty"`

	// ResourcePool is the absolute path of the resource pool where virtual machines will be
	// created. The absolute path is of the form /<datacenter>/host/<cluster>/Resources/<resourcepool>.
	ResourcePool string `json:"resourcePool,omitempty"`

	// ClusterOSImage overrides the url provided in rhcos.json to download the RHCOS OVA
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

	// APIVIP is the virtual IP address for the api endpoint
	//
	// +kubebuilder:validation:format=ip
	// +optional
	APIVIP string `json:"apiVIP,omitempty"`

	// IngressVIP is the virtual IP address for ingress
	//
	// +kubebuilder:validation:format=ip
	// +optional
	IngressVIP string `json:"ingressVIP,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on VSphere for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// Network specifies the name of the network to be used by the cluster.
	Network string `json:"network,omitempty"`

	// DiskType is the name of the disk provisioning type,
	// valid values are thin, thick, and eagerZeroedThick. When not
	// specified, it will be set according to the default storage policy
	// of vsphere.
	DiskType DiskType `json:"diskType,omitempty"`

	// vcenters holds the connection details for services to communicate with vCenter.
	// Currently only a single vCenter is supported.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=1
	// +kubebuilder:validation:MinItems=1
	VCenters []VCenter `json:"vcenters,omitempty"`

	// vSphere location where openshift rhcos virtual machines will be deployed
	// based on vCenter and failureDomain
	// If this is omitted failure domains (regions and zones) will not be used.
	// +kubebuilder:validation:Optional
	DeploymentZones []DeploymentZone `json:"deploymentZones,omitempty"`

	// failureDomains holds the VSpherePlatformFailureDomainSpec which contains
	// the definition of region, zone and the vCenter topology.
	// If this is omitted failure domains (regions and zones) will not be used.
	// +kubebuilder:validation:Optional
	FailureDomains []FailureDomain `json:"failureDomains,omitempty"`
}

// FailureDomainCoordinate holds the name of the associated tag, the type
// of the failure domain, and the vCenter tag category associated with this
// failure domain.
type FailureDomainCoordinate struct {
	// name is the name of the vCenter tag that represents this failure domain
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=80
	Name string `json:"name"`

	// type is the name of the failure domain type, which includes
	// Datacenter, ComputeCluster and HostGroup
	// +kubebuilder:validation:Enum=HostGroup;Datacenter;ComputeCluster
	// +kubebuilder:validation:Required
	Type FailureDomainType `json:"type"`

	// tagCategory is the category used for the tag
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=80
	TagCategory string `json:"tagCategory"`
}

// FailureDomain holds the region and zone failure domain and
// the vCenter topology of that failure domain.
type FailureDomain struct {
	// name defines the name of the FailureDomain
	// This name is abritrary but will be used
	// in VSpherePlatformDeploymentZone for association.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	Name string `json:"name"`

	// region defines a FailureDomainCoordinate which
	// includes the name of the vCenter tag, the failure domain type
	// and the name of the vCenter tag category.
	// +kubebuilder:validation:Required
	Region FailureDomainCoordinate `json:"region"`

	// zone defines a VSpherePlatformFailureDomain which
	// includes the name of the vCenter tag, the failure domain type
	// and the name of the vCenter tag category.
	// +kubebuilder:validation:Required
	Zone FailureDomainCoordinate `json:"zone"`

	// Topology describes a given failure domain using vSphere constructs
	// +kubebuilder:validation:Required
	Topology Topology `json:"topology"`
}

// Topology holds the required and optional vCenter objects - datacenter,
// computeCluster, networks, datastore and resourcePool - to provision virtual machines.
type Topology struct {
	// datacenter is the vCenter datacenter in which virtual machines will be located
	// and defined as the failure domain.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=80
	Datacenter string `json:"datacenter"`

	// computeCluster as the failure domain
	// This is required to be a path
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	ComputeCluster string `json:"computeCluster"`

	// Hosts has information required for placement of machines on VSphere hosts.
	// +optional
	Hosts *FailureDomainHosts `json:"hosts,omitempty"`

	// networks is the list of networks within this failure domain
	// +kubebuilder:validation:Optional
	Networks []string `json:"networks,omitempty"`

	// datastore is the name or inventory path of the datastore in which the
	// virtual machine is created/located.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	Datastore string `json:"datastore"`
}

// FailureDomainHosts defines the attributes of a host group failure domain
type FailureDomainHosts struct {
	// vmGroupName is the Virtual Machine Group name configured
	// within a vCenter cluster that is associated with
	// the corresponding Host Group.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=80
	VMGroupName string `json:"vmGroupName"`

	// hostGroupName is the Host Gorup name configured
	// within a vCenter cluster defining a group
	// of ESXi hosts.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=80
	HostGroupName string `json:"hostGroupName"`
}

// VCenter stores the vCenter connection fields
// https://github.com/kubernetes/cloud-provider-vsphere/blob/master/pkg/common/config/types_yaml.go
type VCenter struct {

	// server is the fully-qualified domain name or the IP address of the vCenter server.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=255
	Server string `json:"server"`

	// port is the TCP port that will be used to communicate to
	// the vCenter endpoint. This is typically unchanged from
	// the default of HTTPS TCP/443.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=32767
	// +kubebuilder:default=443
	Port uint `json:"port,omitempty"`

	// Username is the username that will be used to connect to vCenter
	// +kubebuilder:validation:Required
	Username string `json:"user"`

	// Password is the password for the user to use to connect to the vCenter.
	// +kubebuilder:validation:Required
	Password string `json:"password"`

	// Datacenter in which VMs are located.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	Datacenters []string `json:"datacenters"`
}

// PlacementConstraint is the context information for VM placements within a failure domain
type PlacementConstraint struct {
	// resourcePool is the absolute path of the resource pool where virtual machines will be
	// created. The absolute path is of the form /<datacenter>/host/<cluster>/Resources/<resourcepool>.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	ResourcePool string `json:"resourcePool,omitempty"`

	// folder is the name or inventory path of the folder in which the
	// virtual machine is created/located.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	Folder string `json:"folder,omitempty"`
}

// DeploymentZone holds the association between a
// vCenter, failure domain and the virtual machine placementConstraints
type DeploymentZone struct {

	// name is the abritary name of the DeploymentZone
	// This name will be used in the install MachinePool
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// server is the fully-qualified domain name or the IP address of the vCenter server.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=255
	Server string `json:"server,omitempty"`

	// failureDomain is the name of FailureDomain used for this DeploymentZone
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	FailureDomain string `json:"failureDomain,omitempty"`

	// ControlPlane determines if this failure domain is suitable for use by control plane machines.
	// +kubebuilder:validation:Optional
	ControlPlane DeploymentSuitable `json:"controlPlane,omitempty"`

	// PlacementConstraint encapsulates the placement constraints
	// used within this deployment zone.
	// +kubebuilder:validation:Required
	PlacementConstraint PlacementConstraint `json:"placementConstraint"`
}
