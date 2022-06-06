package foundation

//Node imaging input
type ImageNodesInput struct {
	XsMasterLabel         string        `json:"xs_master_label,omitempty"`
	LayoutEggUUID         string        `json:"layout_egg_uuid,omitempty"`
	IpmiPassword          string        `json:"ipmi_password,omitempty"`
	CvmGateway            string        `json:"cvm_gateway"`
	HypervExternalVnic    string        `json:"hyperv_external_vnic,omitempty"`
	XenConfigType         string        `json:"xen_config_type,omitempty"`
	UcsmIP                string        `json:"ucsm_ip,omitempty"`
	UcsmPassword          string        `json:"ucsm_password,omitempty"`
	HypervisorIso         HypervisorIso `json:"hypervisor_iso"`
	UncPath               string        `json:"unc_path,omitempty"`
	HypervisorNetmask     string        `json:"hypervisor_netmask"`
	FcSettings            *FcSettings   `json:"fc_settings,omitempty"`
	XsMasterPassword      string        `json:"xs_master_password,omitempty"`
	SvmRescueArgs         []string      `json:"svm_rescue_args,omitempty"`
	CvmNetmask            string        `json:"cvm_netmask"`
	XsMasterIP            string        `json:"xs_master_ip,omitempty"`
	Clusters              []*Clusters   `json:"clusters,omitempty"`
	HypervExternalVswitch string        `json:"hyperv_external_vswitch,omitempty"`
	HypervisorNameserver  string        `json:"hypervisor_nameserver,omitempty"`
	HypervSku             string        `json:"hyperv_sku,omitempty"`
	EosMetadata           *EosMetadata  `json:"eos_metadata,omitempty"`
	Tests                 *Tests        `json:"tests,omitempty"`
	Blocks                []*Block      `json:"blocks"`
	HypervProductKey      string        `json:"hyperv_product_key,omitempty"`
	UncUsername           string        `json:"unc_username,omitempty"`
	InstallScript         string        `json:"install_script,omitempty"`
	IpmiUser              string        `json:"ipmi_user,omitempty"`
	HypervisorPassword    string        `json:"hypervisor_password,omitempty"`
	UncPassword           string        `json:"unc_password,omitempty"`
	XsMasterUsername      string        `json:"xs_master_username,omitempty"`
	SkipHypervisor        *bool         `json:"skip_hypervisor,omitempty"`
	HypervisorGateway     string        `json:"hypervisor_gateway"`
	NosPackage            string        `json:"nos_package"` //will be null for cluster creation
	UcsmUser              string        `json:"ucsm_user,omitempty"`
	IPMINetmask           string        `json:"ipmi_netmask,omitempty"`
	IPMIGateway           string        `json:"ipmi_gateway,omitempty"`
}

//Specific hypervisor definition for imaging
type Hypervisor struct {
	Checksum string `json:"checksum,omitempty"`
	Filename string `json:"filename"`
}

//Hypervisor ISO's for various kinds of hypervisor
type HypervisorIso struct {
	Hyperv *Hypervisor `json:"hyperv,omitempty"`
	Kvm    *Hypervisor `json:"kvm,omitempty"`
	Xen    *Hypervisor `json:"xen,omitempty"`
	Esx    *Hypervisor `json:"esx,omitempty"`
}

//Foundation Central Metadata
type FcMetadata struct {
	FcIP   string `json:"fc_ip"`
	APIKey string `json:"api_key"`
}

//Foundaton Central settings
type FcSettings struct {
	FcMetadata        *FcMetadata `json:"fc_metadata"`
	FoundationCentral *bool       `json:"foundation_central"`
}

//Clusters creation related information
type Clusters struct {
	EnableNs              *bool    `json:"enable_ns,omitempty"`
	BackplaneSubnet       string   `json:"backplane_subnet,omitempty"`
	ClusterInitSuccessful *bool    `json:"cluster_init_successful,omitempty"`
	BackplaneNetmask      string   `json:"backplane_netmask,omitempty"`
	RedundancyFactor      *int64   `json:"redundancy_factor"`
	BackplaneVlan         string   `json:"backplane_vlan,omitempty"`
	ClusterName           string   `json:"cluster_name"`
	ClusterExternalIP     *string  `json:"cluster_external_ip"`
	CvmNtpServers         string   `json:"cvm_ntp_servers,omitempty"`
	SingleNodeCluster     *bool    `json:"single_node_cluster,omitempty"`
	ClusterMembers        []string `json:"cluster_members"`
	CvmDNSServers         string   `json:"cvm_dns_servers,omitempty"`
	ClusterInitNow        *bool    `json:"cluster_init_now"`
	HypervisorNtpServers  string   `json:"hypervisor_ntp_servers,omitempty"`
}

// EOS metadata struct used in node imaging when importing cluster config got from install.nutanix.com
type EosMetadata struct {
	ConfigID    string   `json:"config_id"`
	AccountName []string `json:"account_name"`
	Email       string   `json:"email"`
}

// Tests struct for enabling sys check and ncc cecks
type Tests struct {
	RunSyscheck *bool `json:"run_syscheck"`
	RunNcc      *bool `json:"run_ncc"`
}

// UCSM prams for cisco nodes configuration
type UcsmParams struct {
	NativeVlan       *bool  `json:"native_vlan"`
	KeepUcsmSettings *bool  `json:"keep_ucsm_settings"`
	MacPool          string `json:"mac_pool"`
	VlanName         string `json:"vlan_name"`
}

// Vswitches struct for defining virtual switches
type Vswitches struct {
	Lacp        string   `json:"lacp"`
	BondMode    string   `json:"bond_mode"`
	Name        string   `json:"name"`
	Uplinks     []string `json:"uplinks"`
	OtherConfig []string `json:"other_config"`
	Mtu         *int64   `json:"mtu"`
}

//Single node definition
type Node struct {
	IPMINetmask             string       `json:"ipmi_netmask,omitempty"`
	IPMIGateway             string       `json:"ipmi_gateway,omitempty"`
	Ipv6Address             string       `json:"ipv6_address,omitempty"`
	NodePosition            string       `json:"node_position"`
	ImageDelay              *int64       `json:"image_delay,omitempty"`
	UcsmParams              *UcsmParams  `json:"ucsm_params,omitempty"`
	HypervisorHostname      string       `json:"hypervisor_hostname"`
	CvmGbRAM                *int64       `json:"cvm_gb_ram,omitempty"`
	DeviceHint              string       `json:"device_hint,omitempty"`
	BondMode                string       `json:"bond_mode,omitempty"`
	RdmaPassthrough         *bool        `json:"rdma_passthrough,omitempty"`
	ClusterID               string       `json:"cluster_id,omitempty"`
	UcsmNodeSerial          string       `json:"ucsm_node_serial,omitempty"`
	HypervisorIP            string       `json:"hypervisor_ip"`
	NodeSerial              string       `json:"node_serial,omitempty"`
	IpmiConfigureNow        *bool        `json:"ipmi_configure_now,omitempty"`
	ImageSuccessful         *bool        `json:"image_successful,omitempty"`
	Ipv6Interface           string       `json:"ipv6_interface,omitempty"`
	CvmNumVcpus             *int64       `json:"cvm_num_vcpus,omitempty"`
	IpmiMac                 string       `json:"ipmi_mac,omitempty"`
	RdmaMacAddr             string       `json:"rdma_mac_addr,omitempty"`
	BondUplinks             []string     `json:"bond_uplinks,omitempty"`
	CurrentNetworkInterface string       `json:"current_network_interface,omitempty"`
	Hypervisor              string       `json:"hypervisor"`
	Vswitches               []*Vswitches `json:"vswitches,omitempty"`
	BondLacpRate            string       `json:"bond_lacp_rate,omitempty"`
	ImageNow                *bool        `json:"image_now"`
	UcsmManagedMode         string       `json:"ucsm_managed_mode,omitempty"`
	IpmiIP                  string       `json:"ipmi_ip"`
	CurrentCvmVlanTag       *int64       `json:"current_cvm_vlan_tag,omitempty"`
	CvmIP                   string       `json:"cvm_ip"`
	ExludeBootSerial        string       `json:"exlude_boot_serial,omitempty"`
	MitigateLowBootSpace    *bool        `json:"mitigate_low_boot_space,omitempty"`
	IpmiUser                string       `json:"ipmi_user,omitempty"`
	IpmiPassword            string       `json:"ipmi_password,omitempty"`
}

//Block containing multiple nodes
type Block struct {
	Nodes   []*Node `json:"nodes"`
	BlockID string  `json:"block_id,omitempty"`
}

//Response from /image_nodes API call
//Its union of successful and errored response
type ImageNodesAPIResponse struct {
	SessionID string `json:"session_id"`

	//Error details incase of errored responses
	Error *Error `json:"error"`
}

//Error details for image nodes errored response
type Error struct {
	Message   string            `json:"message"`
	Details   map[string]string `json:"details"`
	SessionID string            `json:"session_id"`
}

//Node Imaging progress response
type ImageNodesProgressResponse struct {
	AbortSession             *bool              `json:"abort_session"`
	Results                  []string           `json:"results"`
	SessionID                string             `json:"session_id"`
	ImagingStopped           *bool              `json:"imaging_stopped"`
	AggregatePercentComplete *float64           `json:"aggregate_percent_complete"`
	Action                   string             `json:"action"`
	Clusters                 []*ClusterProgress `json:"clusters"`
	Nodes                    []*NodeProgress    `json:"nodes"`

	//Message in case of errored response
	Message string `json:"message"`
}

//Individual cluster progress
type ClusterProgress struct {
	Category        []string `json:"category,omitempty"`
	Status          string   `json:"status"`
	Messages        []string `json:"messages"`
	ClusterName     string   `json:"cluster_name"`
	TimeElapsed     *float64 `json:"time_elapsed"`
	ClusterMembers  []string `json:"cluster_members"`
	PercentComplete *float64 `json:"percent_complete"`
	TimeTotal       *float64 `json:"time_total"`
}

//Individual Node progress
type NodeProgress struct {
	Category        []string `json:"category,omitempty"`
	Status          string   `json:"status"`
	Messages        []string `json:"messages"`
	TimeElapsed     *float64 `json:"time_elapsed"`
	CvmIP           string   `json:"cvm_ip"`
	PercentComplete *float64 `json:"percent_complete"`
	HypervisorIP    string   `json:"hypervisor_ip"`
	TimeTotal       *float64 `json:"time_total"`
}

//Response from /enumerate_nos_packages api
type ListNOSPackagesResponse []string

//Reference to hypervisor for ListHypervisorISOsResponse
type HypervisorISOReference struct {
	Supported *bool  `json:"supported"`
	Filename  string `json:"filename"`
}

//Response from /enumerate_hypervisor_isos api
type ListHypervisorISOsResponse struct {
	Hyperv []*HypervisorISOReference `json:"hyperv"`
	Kvm    []*HypervisorISOReference `json:"kvm"`
	Esx    []*HypervisorISOReference `json:"esx"`
	Linux  []*HypervisorISOReference `json:"linux"`
	Xen    []*HypervisorISOReference `json:"xen"`
}

//Response struct format for /discover_nodes api
type DiscoverNodesAPIResponse []DiscoveredBlock

//struct format for discovered block
type DiscoveredBlock struct {
	Model    string           `json:"model"`
	Nodes    []DiscoveredNode `json:"nodes"`
	ChassisN *int64           `json:"chassis_n"`
	BlockID  string           `json:"block_id"`
}

//Format for single discovered node
type DiscoveredNode struct {
	FoundationVersion       string      `json:"foundation_version"`
	Ipv6Address             string      `json:"ipv6_address"`
	NodeUUID                string      `json:"node_uuid"`
	CurrentNetworkInterface string      `json:"current_network_interface"`
	NodePosition            string      `json:"node_position"`
	Hypervisor              string      `json:"hypervisor"`
	Configured              *bool       `json:"configured"`
	NosVersion              string      `json:"nos_version"`
	ClusterID               interface{} `json:"cluster_id"`
	CurrentCvmVlanTag       interface{} `json:"current_cvm_vlan_tag"`
	HypervisorVersion       string      `json:"hypervisor_version"`
	SvmIP                   string      `json:"svm_ip"`
	Model                   string      `json:"model"`
	NodeSerial              string      `json:"node_serial"`
}

//Response struct for /node_network_details api
type NodeNetworkDetailsResponse struct {
	Nodes []NodeNetworkDetail `json:"nodes"`
	Error string              `json:"error"`
}

//Network detail of particular node
type NodeNetworkDetail struct {
	CvmGateway         string `json:"cvm_gateway"`
	IpmiNetmask        string `json:"ipmi_netmask"`
	Ipv6Address        string `json:"ipv6_address"`
	CvmVlanID          string `json:"cvm_vlan_id"`
	HypervisorHostname string `json:"hypervisor_hostname"`
	HypervisorNetmask  string `json:"hypervisor_netmask"`
	CvmNetmask         string `json:"cvm_netmask"`
	IpmiIP             string `json:"ipmi_ip"`
	HypervisorGateway  string `json:"hypervisor_gateway"`
	Error              string `json:"error"`
	CvmIP              string `json:"cvm_ip"`
	HypervisorIP       string `json:"hypervisor_ip"`
	IpmiGateway        string `json:"ipmi_gateway"`
	NodeSerial         string `json:"node_serial"`
}

//Node network details api input
type NodeNetworkDetailsInput struct {
	Nodes   []NodeIpv6Input `json:"nodes"`
	Timeout string          `json:"timeout,omitempty"`
}

//Node ipv6 details for node network details api
type NodeIpv6Input struct {
	Ipv6Address string `json:"ipv6_address"`
}

// Input struct for IPMI Configuration using /ipmi_config api in foundation
type IPMIConfigAPIInput struct {
	IpmiUser     string                 `json:"ipmi_user"`
	IpmiNetmask  string                 `json:"ipmi_netmask"`
	Blocks       []IPMIConfigBlockInput `json:"blocks"`
	IpmiGateway  string                 `json:"ipmi_gateway"`
	IpmiPassword string                 `json:"ipmi_password"`
}

// IPMI config info for every node
type IPMIConfigNodeInput struct {
	IpmiMac          string `json:"ipmi_mac"`
	IpmiConfigureNow bool   `json:"ipmi_configure_now"`
	IpmiIP           string `json:"ipmi_ip"`
}

// Block Info for IPMI Config
type IPMIConfigBlockInput struct {
	Nodes   []IPMIConfigNodeInput `json:"nodes"`
	BlockID string                `json:"block_id,omitempty"`
}

// Response stuct for ipmi config using /ipmi_config api in foundation
type IPMIConfigAPIResponse struct {
	IpmiUser     string                    `json:"ipmi_user"`
	IpmiNetmask  string                    `json:"ipmi_netmask"`
	Blocks       []IPMIConfigBlockResponse `json:"blocks"`
	IpmiGateway  string                    `json:"ipmi_gateway"`
	IpmiPassword string                    `json:"ipmi_password"`
	Error        *Error                    `json:"error"`
}

// IPMI config info response for every node along with success flag
type IPMIConfigNodeResponse struct {
	IpmiConfigureSuccessful bool   `json:"ipmi_configure_successful"`
	IpmiConfigureNow        bool   `json:"ipmi_configure_now"`
	IpmiIP                  string `json:"ipmi_ip"`
	IpmiMac                 string `json:"ipmi_mac"`
	IpmiMessage             string `json:"ipmi_message"`
}

// Block Info for IPMI Config response
type IPMIConfigBlockResponse struct {
	Nodes   []IPMIConfigNodeResponse `json:"nodes"`
	BlockID string                   `json:"block_id"`
}

// Response struct for upload image
type UploadImageResponse struct {
	Md5Sum      string `json:"md5sum"`
	Name        string `json:"name"`
	InWhitelist bool   `json:"in_whitelist"`
}
