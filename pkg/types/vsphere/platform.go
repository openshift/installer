package vsphere

import (
	configv1 "github.com/openshift/api/config/v1"
)

// DiskType is a disk provisioning type for vsphere.
// +kubebuilder:validation:Enum="";thin;thick;eagerZeroedThick
type DiskType string

// FailureDomainType is the string representation name of the failure domain type.
// There are three defined failure domains types currently, Datacenter, ComputeCluster and HostGroup.
// Each represents a vCenter object type within a vSphere environment.
// +kubebuilder:validation:Enum=HostGroup;Datacenter;ComputeCluster
type FailureDomainType string

const (
	// DiskTypeThin uses Thin disk provisioning type for vsphere in the cluster.
	DiskTypeThin DiskType = "thin"

	// DiskTypeThick uses Thick disk provisioning type for vsphere in the cluster.
	DiskTypeThick DiskType = "thick"

	// DiskTypeEagerZeroedThick uses EagerZeroedThick disk provisioning type for vsphere in the cluster.
	DiskTypeEagerZeroedThick DiskType = "eagerZeroedThick"

	// TagCategoryRegion the tag category associated with regions.
	TagCategoryRegion = "openshift-region"

	// TagCategoryZone the tag category associated with zones.
	TagCategoryZone = "openshift-zone"
)

const (
	// ControlPlaneRole represents control-plane nodes.
	ControlPlaneRole = "control-plane"
	// ComputeRole represents worker nodes.
	ComputeRole = "compute"
	// BootstrapRole represents bootstrap nodes.
	BootstrapRole = "bootstrap"
)

const (
	// HostGroupFailureDomain is a failure domain for a vCenter vm-host group.
	HostGroupFailureDomain FailureDomainType = "HostGroup"
	// ComputeClusterFailureDomain is a failure domain for a vCenter compute cluster.
	ComputeClusterFailureDomain FailureDomainType = "ComputeCluster"
	// DatacenterFailureDomain is a failure domain for a vCenter datacenter.
	DatacenterFailureDomain FailureDomainType = "Datacenter"
)

// Platform stores any global configuration used for vsphere platforms.
type Platform struct {
	// VCenter is the domain name or IP address of the vCenter.
	// Deprecated: Use VCenters.Server
	DeprecatedVCenter string `json:"vCenter,omitempty"`
	// Username is the name of the user to use to connect to the vCenter.
	// Deprecated: Use VCenters.Username
	DeprecatedUsername string `json:"username,omitempty"`
	// Password is the password for the user to use to connect to the vCenter.
	// Deprecated: Use VCenters.Password
	DeprecatedPassword string `json:"password,omitempty"`
	// Datacenter is the name of the datacenter to use in the vCenter.
	// Deprecated: Use FailureDomains.Topology.Datacenter
	DeprecatedDatacenter string `json:"datacenter,omitempty"`
	// DefaultDatastore is the default datastore to use for provisioning volumes.
	// Deprecated: Use FailureDomains.Topology.Datastore
	DeprecatedDefaultDatastore string `json:"defaultDatastore,omitempty"`
	// Folder is the absolute path of the folder that will be used and/or created for
	// virtual machines. The absolute path is of the form /<datacenter>/vm/<folder>/<subfolder>.
	// +kubebuilder:validation:Pattern=`^/.*?/vm/.*?`
	// +optional
	// Deprecated: Use FailureDomains.Topology.Folder
	DeprecatedFolder string `json:"folder,omitempty"`
	// Cluster is the name of the cluster virtual machines will be cloned into.
	// Deprecated: Use FailureDomains.Topology.Cluster
	DeprecatedCluster string `json:"cluster,omitempty"`
	// ResourcePool is the absolute path of the resource pool where virtual machines will be
	// created. The absolute path is of the form /<datacenter>/host/<cluster>/Resources/<resourcepool>.
	// Deprecated: Use FailureDomains.Topology.ResourcePool
	DeprecatedResourcePool string `json:"resourcePool,omitempty"`
	// ClusterOSImage overrides the url provided in rhcos.json to download the RHCOS OVA
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

	// DeprecatedAPIVIP is the virtual IP address for the api endpoint
	// Deprecated: Use APIVIPs
	//
	// +kubebuilder:validation:format=ip
	// +optional
	DeprecatedAPIVIP string `json:"apiVIP,omitempty"`

	// APIVIPs contains the VIP(s) for the api endpoint. In dual stack clusters
	// it contains an IPv4 and IPv6 address, otherwise only one VIP
	//
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:Format=ip
	// +optional
	APIVIPs []string `json:"apiVIPs,omitempty"`

	// DeprecatedIngressVIP is the virtual IP address for ingress
	// Deprecated: Use IngressVIPs
	//
	// +kubebuilder:validation:format=ip
	// +optional
	DeprecatedIngressVIP string `json:"ingressVIP,omitempty"`

	// IngressVIPs contains the VIP(s) for ingress. In dual stack clusters it
	// contains an IPv4 and IPv6 address, otherwise only one VIP
	//
	// +kubebuilder:validation:MaxItems=2
	// +kubebuilder:validation:Format=ip
	// +optional
	IngressVIPs []string `json:"ingressVIPs,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on VSphere for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`
	// Network specifies the name of the network to be used by the cluster.
	// Deprecated: Use FailureDomains.Topology.Network
	DeprecatedNetwork string `json:"network,omitempty"`
	// DiskType is the name of the disk provisioning type,
	// valid values are thin, thick, and eagerZeroedThick. When not
	// specified, it will be set according to the default storage policy
	// of vsphere.
	DiskType DiskType `json:"diskType,omitempty"`
	// VCenters holds the connection details for services to communicate with vCenter.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:MaxItems=3
	// +kubebuilder:validation:MinItems=1
	VCenters []VCenter `json:"vcenters,omitempty"`
	// FailureDomains holds the VSpherePlatformFailureDomainSpec which contains
	// the definition of region, zone and the vCenter topology.
	// If this is omitted failure domains (regions and zones) will not be used.
	// +kubebuilder:validation:Optional
	FailureDomains []FailureDomain `json:"failureDomains,omitempty"`

	// nodeNetworking contains the definition of internal and external network constraints for
	// assigning the node's networking.
	// If this field is omitted, networking defaults to the legacy
	// address selection behavior which is to only support a single address and
	// return the first one found.
	// +optional
	NodeNetworking *configv1.VSpherePlatformNodeNetworking `json:"nodeNetworking,omitempty"`

	// LoadBalancer defines how the load balancer used by the cluster is configured.
	// LoadBalancer is available in TechPreview.
	// +optional
	LoadBalancer *configv1.VSpherePlatformLoadBalancer `json:"loadBalancer,omitempty"`

	// dnsRecordsType determines whether records for api, api-int, and ingress
	// are provided by the internal DNS service or externally.
	// Allowed values are `Internal`, `External`, and omitted.
	// When set to `Internal`, records are provided by the internal infrastructure and
	// no additional user configuration is required for the cluster to function.
	// When set to `External`, records are not provided by the internal infrastructure
	// and must be configured by the user on a DNS server outside the cluster.
	// Cluster nodes must use this external server for their upstream DNS requests.
	// This value may only be set when loadBalancer.type is set to UserManaged.
	// When omitted, this means the user has no opinion and the platform is left
	// to choose reasonable defaults. These defaults are subject to change over time.
	// The current default is `Internal`.
	// +openshift:enable:FeatureGate=OnPremDNSRecords
	// +optional
	DNSRecordsType configv1.DNSRecordsType `json:"dnsRecordsType,omitempty"`

	// Hosts defines network configurations to be applied by the installer. Hosts is available in TechPreview.
	Hosts []*Host `json:"hosts,omitempty"`
}

// FailureDomain holds the region and zone failure domain and
// the vCenter topology of that failure domain.
type FailureDomain struct {
	// name defines the name of the FailureDomain
	// This name is arbitrary but will be used
	// in VSpherePlatformDeploymentZone for association.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=256
	Name string `json:"name"`
	// region defines a FailureDomainCoordinate which
	// includes the name of the vCenter tag, the failure domain type
	// and the name of the vCenter tag category.

	// The region is the name of the tag in vCenter that is associated with the
	// tag category `openshift-region`. The region name must match the tag name
	// and must exist prior to installation. When the regionType is Datacenter
	// the tag must be attached to the toplogy.datacenter object in vCenter.
	// When the regionType is ComputeCluster the tag must be attached to the topology.computeCluster
	// object in vCenter.
	// +kubebuilder:validation:Required
	Region string `json:"region"`
	// The zone is the name of the tag in vCenter that is associated with
	// the tag category `openshift-zone`. The zone name must match the tag name
	// and must exist prior to installation. When zoneType is HostGroup the
	// ESXi hosts defined in the provided in the topology.hostGroup field must be tagged.
	// When the zoneType is ComputeCluster the tag must be attached to the topology.computeCluster
	// object in vCenter.
	// +kubebuilder:validation:Required
	Zone string `json:"zone"`
	// server is the fully-qualified domain name or the IP address of the vCenter server.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=255
	Server string `json:"server"`
	// Topology describes a given failure domain using vSphere constructs
	// +kubebuilder:validation:Required
	Topology Topology `json:"topology"`

	// regionType is the type of failure domain region, the current values are "Datacenter" and "ComputeCluster"
	// +kubebuilder:validation:Enum=Datacenter;ComputeCluster
	// When regionType is Datacenter the zoneType must be ComputeCluster.
	// When regionType is ComputeCluster the zoneType must be HostGroup
	// +optional
	RegionType FailureDomainType `json:"regionType,omitempty"`
	// zoneType is the type of the failure domain zone, the current values are "ComputeCluster" and "HostGroup"

	// When zoneType is ComputeCluster the regionType must be Datacenter
	// When zoneType is HostGroup the regionType must be ComputeCluster
	// If the zoneType is HostGroup topology.hostGroup must be defined and exist in vCenter
	// prior to installation.
	// +kubebuilder:validation:Enum=ComputeCluster;HostGroup
	// +optional
	ZoneType FailureDomainType `json:"zoneType,omitempty"`
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
	// networks is the list of networks within this failure domain
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=10
	Networks []string `json:"networks,omitempty"`
	// datastore is the name or inventory path of the datastore in which the
	// virtual machine is created/located.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	Datastore string `json:"datastore"`
	// resourcePool is the absolute path of the resource pool where virtual machines will be
	// created. The absolute path is of the form /<datacenter>/host/<cluster>/Resources/<resourcepool>.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	// +kubebuilder:validation:Pattern=`^/.*?/host/.*?/Resources.*`
	// +optional
	ResourcePool string `json:"resourcePool,omitempty"`
	// folder is the inventory path of the folder in which the
	// virtual machine is created/located.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	// +kubebuilder:validation:Pattern=`^/.*?/vm/.*?`
	// +optional
	Folder string `json:"folder,omitempty"`
	// template is the inventory path of the virtual machine or template
	// that will be used for cloning.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=2048
	// +kubebuilder:validation:Pattern=`^/.*?/vm/.*?`
	// +optional
	Template string `json:"template,omitempty"`
	// tagIDs is an optional set of tags to add to an instance. Specified tagIDs
	// must use URN-notation instead of display names. A maximum of 10 tag IDs may be specified.
	// +kubebuilder:example=`urn:vmomi:InventoryServiceTag:5736bf56-49f5-4667-b38c-b97e09dc9578:GLOBAL`
	// +optional
	TagIDs []string `json:"tagIDs,omitempty"`

	// hostGroup is the name of the vm-host group of type host within vCenter for this failure domain.
	// hostGroup is limited to 80 characters.
	// This field is required when the ZoneType is HostGroup
	// +kubebuilder:validation:MaxLength=80
	// +optional
	HostGroup string `json:"hostGroup,omitempty"`
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
	Port int32 `json:"port,omitempty"`
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

// Host defines host VMs to generate as part of the installation.
type Host struct {
	// FailureDomain refers to the name of a FailureDomain as described in https://github.com/openshift/enhancements/blob/master/enhancements/installer/vsphere-ipi-zonal.md
	// +optional
	FailureDomain string `json:"failureDomain"`
	// NetworkDeviceSpec to be applied to the host
	// +kubebuilder:validation:Required
	NetworkDevice *NetworkDeviceSpec `json:"networkDevice"`
	// Role defines the role of the node
	// +kubebuilder:validation:Enum="";bootstrap;control-plane;compute
	// +kubebuilder:validation:Required
	Role string `json:"role"`
}

// NetworkDeviceSpec defines network config for static IP assignment.
type NetworkDeviceSpec struct {
	// gateway is an IPv4 or IPv6 address which represents the subnet gateway,
	// for example, 192.168.1.1.
	// +kubebuilder:validation:Format=ipv4
	// +kubebuilder:validation:Format=ipv6
	Gateway string `json:"gateway,omitempty"`

	// ipAddrs is a list of one or more IPv4 and/or IPv6 addresses and CIDR to assign to
	// this device, for example, 192.168.1.100/24. IP addresses provided via ipAddrs are
	// intended to allow explicit assignment of a machine's IP address.
	// +kubebuilder:validation:Format=ipv4
	// +kubebuilder:validation:Format=ipv6
	// +kubebuilder:example=`192.168.1.100/24`
	// +kubebuilder:example=`2001:DB8:0000:0000:244:17FF:FEB6:D37D/64`
	// +kubebuilder:validation:Required
	IPAddrs []string `json:"ipAddrs"`

	// nameservers is a list of IPv4 and/or IPv6 addresses used as DNS nameservers, for example,
	// 8.8.8.8. a nameserver is not provided by a fulfilled IPAddressClaim. If DHCP is not the
	// source of IP addresses for this network device, nameservers should include a valid nameserver.
	// +kubebuilder:validation:Format=ipv4
	// +kubebuilder:validation:Format=ipv6
	// +kubebuilder:example=`8.8.8.8`
	Nameservers []string `json:"nameservers,omitempty"`
}

// IsControlPlane checks if the current host is a master.
func (h *Host) IsControlPlane() bool {
	return h.Role == ControlPlaneRole
}

// IsCompute checks if the current host is a worker.
func (h *Host) IsCompute() bool {
	return h.Role == ComputeRole
}

// IsBootstrap checks if the current host is a bootstrap.
func (h *Host) IsBootstrap() bool {
	return h.Role == BootstrapRole
}
