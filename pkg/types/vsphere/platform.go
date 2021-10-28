package vsphere

// DiskType is a disk provisioning type for vsphere.
// +kubebuilder:validation:Enum="";thin;thick;eagerZeroedThick
type DiskType string

const (
	// DiskTypeThin uses Thin disk provisioning type for vsphere in the cluster.
	DiskTypeThin DiskType = "thin"

	// DiskTypeThick uses Thick disk provisioning type for vsphere in the cluster.
	DiskTypeThick DiskType = "thick"

	// DiskTypeEagerZeroedThick uses EagerZeroedThick disk provisioning type for vsphere in the cluster.
	DiskTypeEagerZeroedThick DiskType = "eagerZeroedThick"
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

	// VCenters slice is the configuration for the use of Regions and Zones. Also
	// deploying openshift cluster virtual machines in multiple datacenters and clusters.
	// +optional
	VCenters []VCenter `json:"vcenters,omitempty"`
}

// VCenter stores the vCenter connection fields and the Region slice
// that is used to create virtual machines in various datacenters and clusters
type VCenter struct {
	// Server is the domain name or IP address of the vCenter.
	// +required
	Server string `json:"server"`

	// User is the name of the user to use to connect to the vCenter.
	// +required
	User string `json:"user"`

	// Password is the password for the user to use to connect to the vCenter.
	// +required
	Password string `json:"password"`

	// Port is the vCenter API TCP port that the vCenter SDK clients connect to.
	// +optional
	Port uint `json:"port,omitempty"`

	// Regions slice is the configuration of the vCenter datacenter and Zones that
	// will be used in deployment of virtual machines, cloud provider and tags.
	// +required
	Regions []Region `json:"regions"`
}

// Region stores the name of the region, the datacenter that defines the region and the Zones
// that are a part of a particular region.
type Region struct {

	// Name is the region name that will be used for cloud provider and tag
	// +required
	Name string `json:"name"`

	// Datacenter is the name of the datacenter to use in the vCenter.
	// +required
	Datacenter string `json:"datacenter"`

	// Zones slice is the configuration of the vCenter Resource Pool (optionally),
	// Cluster, Network and Datastore to deploy virtual machines.
	// +required
	Zones []Zone `json:"zones"`
}

// Zone stores the name of the zone, the cluster associated with the zone.
// ResourcePool, Network and Datastore is provided for virtual machine deployment.
type Zone struct {
	// Name is the zone name that will be used for the cloud provider and tag
	// +required
	Name string `json:"name"`

	// ResourcePool is the name of the already configured Resource Pool to use
	// when deploying virtual machines.
	// +optional
	ResourcePool string `json:"resourcePool,omitempty"`

	// Cluster is the name of the cluster virtual machines will be cloned into.
	// +required
	Cluster string `json:"cluster"`

	// Network specifies the name of the network to be used by the cluster.
	// +required
	Network string `json:"network"`

	// Datastore is the default datastore to use for provisioning volumes.
	// +required
	Datastore string `json:"datastore"`
}
