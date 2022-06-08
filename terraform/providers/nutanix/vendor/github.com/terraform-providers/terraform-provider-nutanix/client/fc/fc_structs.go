package foundationcentral

type ErrorResponse struct {
	Code        *int32
	MessageList []*string
}

// Metadata for List Operations Input
type ListMetadataInput struct {
	Length *int `json:"length,omitempty"`
	Offset *int `json:"offset,omitempty"`
}

// Metadata for List Operations Output
type ListMetadataOutput struct {
	TotalMatches *int `json:"total_matches,omitempty"`
	Length       *int `json:"length,omitempty"`
	Offset       *int `json:"offset,omitempty"`
}

// CommonNetworkSetting ...
type CommonNetworkSettings struct {
	CvmDNSServers        []string `json:"cvm_dns_servers,omitempty"`
	HypervisorDNSServers []string `json:"hypervisor_dns_servers,omitempty"`
	CvmNtpServers        []string `json:"cvm_ntp_servers,omitempty"`
	HypervisorNtpServers []string `json:"hypervisor_ntp_servers,omitempty"`
}

// ImagedNodesList Filter
type ImagedNodeListFilter struct {
	NodeState *string `json:"node_state,omitempty"`
}

type HardwareAttribute struct {
}

// ImagedNodeDetails ...
type ImagedNodeDetails struct {
	CvmVlanID          *int                   `json:"cvm_vlan_id,omitempty"`
	NodeType           *string                `json:"node_type,omitempty"`
	CreatedTimestamp   *string                `json:"created_timestamp,omitempty"`
	Ipv6Interface      *string                `json:"ipv6_interface,omitempty"`
	APIKeyUUID         *string                `json:"api_key_uuid,omitempty"`
	FoundationVersion  *string                `json:"foundation_version,omitempty"`
	CurrentTime        *string                `json:"current_time,omitempty"`
	NodePosition       *string                `json:"node_position,omitempty"`
	CvmNetmask         *string                `json:"cvm_netmask,omitempty"`
	IpmiIP             *string                `json:"ipmi_ip,omitempty"`
	CvmUUID            *string                `json:"cvm_uuid,omitempty"`
	CvmIpv6            *string                `json:"cvm_ipv6,omitempty"`
	ImagedClusterUUID  *string                `json:"imaged_cluster_uuid,omitempty"`
	CvmUp              *bool                  `json:"cvm_up,omitempty"`
	Available          *bool                  `json:"available,omitempty"`
	ObjectVersion      *int                   `json:"object_version,omitempty"`
	IpmiNetmask        *string                `json:"ipmi_netmask,omitempty"`
	HypervisorHostname *string                `json:"hypervisor_hostname,omitempty"`
	NodeState          *string                `json:"node_state,omitempty"`
	HypervisorVersion  *string                `json:"hypervisor_version,omitempty"`
	HypervisorIP       *string                `json:"hypervisor_ip,omitempty"`
	Model              *string                `json:"model,omitempty"`
	IpmiGateway        *string                `json:"ipmi_gateway,omitempty"`
	HardwareAttributes map[string]interface{} `json:"hardware_attributes,omitempty"`
	CvmGateway         *string                `json:"cvm_gateway,omitempty"`
	NodeSerial         *string                `json:"node_serial,omitempty"`
	ImagedNodeUUID     *string                `json:"imaged_node_uuid,omitempty"`
	BlockSerial        *string                `json:"block_serial,omitempty"`
	HypervisorType     *string                `json:"hypervisor_type,omitempty"`
	LatestHbTSList     []*string              `json:"latest_hb_ts_list,omitempty"`
	HypervisorNetmask  *string                `json:"hypervisor_netmask,omitempty"`
	HypervisorGateway  *string                `json:"hypervisor_gateway,omitempty"`
	CvmIP              *string                `json:"cvm_ip,omitempty"`
	AosVersion         *string                `json:"aos_version,omitempty"`
	ClusterExternalIP  *string                `json:"cluster_external_ip,omitempty"`
	SupportedFeatures  []*string              `json:"supported_features,omitempty"`
}

// ImagedNodesInput ...
type ImagedNodesInput struct {
	CvmVlanID          *int               `json:"cvm_vlan_id,omitempty"`
	NodeType           *string            `json:"node_type,omitempty"`
	Ipv6Interface      *string            `json:"ipv6_interface,omitempty"`
	FoundationVersion  *string            `json:"foundation_version,omitempty"`
	IpmiNetmask        *string            `json:"ipmi_netmask,omitempty"`
	CvmNetmask         *string            `json:"cvm_netmask,omitempty"`
	IpmiIP             *string            `json:"ipmi_ip,omitempty"`
	CvmUUID            *string            `json:"cvm_uuid,omitempty"`
	CvmIpv6            *string            `json:"cvm_ipv6,omitempty"`
	CvmUp              *bool              `json:"cvm_up,omitempty"`
	NodePosition       *string            `json:"node_position,omitempty"`
	HypervisorHostname *string            `json:"hypervisor_hostname,omitempty"`
	HypervisorVersion  *string            `json:"hypervisor_version,omitempty"`
	HypervisorIP       *string            `json:"hypervisor_ip,omitempty"`
	CvmIP              *string            `json:"cvm_ip,omitempty"`
	IpmiGateway        *string            `json:"ipmi_gateway,omitempty"`
	HardwareAttributes *HardwareAttribute `json:"hardware_attributes,omitempty"`
	CvmGateway         *string            `json:"cvm_gateway,omitempty"`
	NodeSerial         *string            `json:"node_serial,omitempty"`
	BlockSerial        *string            `json:"block_serial,omitempty"`
	HypervisorType     *string            `json:"hypervisor_type,omitempty"`
	HypervisorNetmask  *string            `json:"hypervisor_netmask,omitempty"`
	HypervisorGateway  *string            `json:"hypervisor_gateway,omitempty"`
	Model              *string            `json:"model,omitempty"`
	AosVersion         *string            `json:"aos_version,omitempty"`
}

// ImagedNodeResponse ...
type ImagedNodesResponse struct {
	ObjectVersion  *int    `json:"object_version,omitempty"`
	ImagedNodeUUID *string `json:"imaged_node_uuid,omitempty"`
}

// Input for Imaged Nodes List
type ImagedNodesListInput struct {
	Length  *int                  `json:"length,omitempty"`
	Filters *ImagedNodeListFilter `json:"filters,omitempty"`
	Offset  *int                  `json:"offset,omitempty"`
}

// Response of Imaged Nodes List
type ImagedNodesListResponse struct {
	Metadata    *ListMetadataOutput  `json:"metadata,omitempty"`
	ImagedNodes []*ImagedNodeDetails `json:"imaged_nodes,omitempty"`
}

// Progress Details of a Node
type NodeProgressDetail struct {
	Status          *string   `json:"status,omitempty"`
	ImagedNodeUUID  *string   `json:"imaged_node_uuid,omitempty"`
	IntentPickedUp  *bool     `json:"intent_picked_up,omitempty"`
	ImagingStopped  *bool     `json:"imaging_stopped,omitempty"`
	PercentComplete *float64  `json:"percent_complete,omitempty"`
	MessageList     []*string `json:"message_list,omitempty"`
}

// Progress details of a Cluster
type ClusterProgressDetails struct {
	ClusterName     *string   `json:"cluster_name,omitempty"`
	Status          *string   `json:"status,omitempty"`
	PercentComplete *float64  `json:"percent_complete,omitempty"`
	MessageList     []*string `json:"message_list,omitempty"`
}

// Format of Cluster Status
type ClusterStatus struct {
	ClusterCreationStarted   *bool                   `json:"cluster_creation_started,omitempty"`
	IntentPickedUp           *bool                   `json:"intent_picked_up,omitempty"`
	ImagingStopped           *bool                   `json:"imaging_stopped,omitempty"`
	NodeProgressDetails      []*NodeProgressDetail   `json:"node_progress_details,omitempty"`
	AggregatePercentComplete *float64                `json:"aggregate_percent_complete,omitempty"`
	CurrentFoundationIP      *string                 `json:"current_foundation_ip,omitempty"`
	ClusterProgressDetails   *ClusterProgressDetails `json:"cluster_progress_details,omitempty"`
	FoundationSessionID      *string                 `json:"foundation_session_id,omitempty"`
}

// Format of imaged cluster details
type ImagedClusterDetails struct {
	CurrentTime            *string                `json:"current_time,omitempty"`
	Archived               *bool                  `json:"archived,omitempty"`
	ClusterExternalIP      *string                `json:"cluster_external_ip,omitempty"`
	ImagedNodeUUIDList     []*string              `json:"imaged_node_uuid_list,omitempty"`
	CommonNetworkSettings  *CommonNetworkSettings `json:"common_network_settings,omitempty"`
	StorageNodeCount       *int                   `json:"storage_node_count,omitempty"`
	RedundancyFactor       *int                   `json:"redundancy_factor,omitempty"`
	FoundationInitNodeUUID *string                `json:"foundation_init_node_uuid,omitempty"`
	WorkflowType           *string                `json:"workflow_type,omitempty"`
	ClusterName            *string                `json:"cluster_name,omitempty"`
	FoundationInitConfig   *FoundationInitConfig  `json:"foundation_init_config,omitempty"`
	ClusterStatus          *ClusterStatus         `json:"cluster_status,omitempty"`
	ClusterSize            *int                   `json:"cluster_size,omitempty"`
	Destroyed              *bool                  `json:"destroyed,omitempty"`
	CreatedTimestamp       *string                `json:"created_timestamp,omitempty"`
	ImagedClusterUUID      *string                `json:"imaged_cluster_uuid,omitempty"`
}

// filter for Imaged Cluster List
type ImagedClustersListFilter struct {
	Archived *bool `json:"archived,omitempty"`
}

// Input for Imaged Cluster List
type ImagedClustersListInput struct {
	Length  *int                      `json:"length,omitempty"`
	Filters *ImagedClustersListFilter `json:"filters,omitempty"`
	Offset  *int                      `json:"offset,omitempty"`
}

// Response for Imaged Cluster List
type ImagedClustersListResponse struct {
	Metadata       *ListMetadataOutput     `json:"metadata,omitempty"`
	ImagedClusters []*ImagedClusterDetails `json:"imaged_clusters,omitempty"`
}

// Format of Hypervisor ISO Details
type HypervisorIsoDetails struct {
	HypervSku        *string `json:"hyperv_sku,omitempty"`
	URL              *string `json:"url,omitempty"`
	HypervProductKey *string `json:"hyperv_product_key,omitempty"`
	Sha256sum        *string `json:"sha256sum,omitempty"`
}

// format of Node to be Imaged
type Node struct {
	CvmGateway                 *string                `json:"cvm_gateway,omitempty"`
	IpmiNetmask                *string                `json:"ipmi_netmask,omitempty"`
	RdmaPassthrough            *bool                  `json:"rdma_passthrough,omitempty"`
	ImagedNodeUUID             *string                `json:"imaged_node_uuid,omitempty"`
	CvmVlanID                  *int                   `json:"cvm_vlan_id,omitempty"`
	HypervisorType             *string                `json:"hypervisor_type,omitempty"`
	ImageNow                   *bool                  `json:"image_now,omitempty"`
	HypervisorHostname         *string                `json:"hypervisor_hostname,omitempty"`
	HypervisorNetmask          *string                `json:"hypervisor_netmask,omitempty"`
	CvmNetmask                 *string                `json:"cvm_netmask,omitempty"`
	IpmiIP                     *string                `json:"ipmi_ip,omitempty"`
	HypervisorGateway          *string                `json:"hypervisor_gateway,omitempty"`
	HardwareAttributesOverride map[string]interface{} `json:"hardware_attributes_override,omitempty"`
	CvmRAMGb                   *int                   `json:"cvm_ram_gb,omitempty"`
	CvmIP                      *string                `json:"cvm_ip,omitempty"`
	HypervisorIP               *string                `json:"hypervisor_ip,omitempty"`
	UseExistingNetworkSettings *bool                  `json:"use_existing_network_settings,omitempty"`
	IpmiGateway                *string                `json:"ipmi_gateway,omitempty"`
}

// Input to Create a cluster
type CreateClusterInput struct {
	ClusterExternalIP     *string                `json:"cluster_external_ip,omitempty"`
	CommonNetworkSettings *CommonNetworkSettings `json:"common_network_settings,omitempty"`
	HypervisorIsoDetails  *HypervisorIsoDetails  `json:"hypervisor_iso_details,omitempty"`
	StorageNodeCount      *int                   `json:"storage_node_count,omitempty"`
	RedundancyFactor      *int                   `json:"redundancy_factor,omitempty"`
	ClusterName           *string                `json:"cluster_name,omitempty"`
	AosPackageURL         *string                `json:"aos_package_url,omitempty"`
	ClusterSize           *int                   `json:"cluster_size,omitempty"`
	AosPackageSha256sum   *string                `json:"aos_package_sha256sum,omitempty"`
	Timezone              *string                `json:"timezone,omitempty"`
	NodesList             []*Node                `json:"nodes_list,omitempty"`
	SkipClusterCreation   bool                   `json:"skip_cluster_creation,omitempty"`
}

// Response of cluster creation
type CreateClusterResponse struct {
	ImagedClusterUUID *string `json:"imaged_cluster_uuid,omitempty"`
}

// Update cluster data
type UpdateClusterData struct {
	Archived *bool `json:"archived,omitempty"`
}

// Input to create API KEY
type CreateAPIKeysInput struct {
	Alias string `json:"alias"`
}

// Response of API KEY creation
type CreateAPIKeysResponse struct {
	CreatedTimestamp string `json:"created_timestamp,omitempty"`
	Alias            string `json:"alias,omitempty"`
	KeyUUID          string `json:"key_uuid,omitempty"`
	APIKey           string `json:"api_key,omitempty"`
	CurrentTime      string `json:"current_time,omitempty"`
}

// format to list all the API key
type ListAPIKeysResponse struct {
	Metadata *ListMetadataOutput      `json:"metadata,omitempty" mapstructure:"metadata,omitempty"`
	APIKeys  []*CreateAPIKeysResponse `json:"api_keys,omitempty"`
}

// Format of Foundation Init Config
type FoundationInitConfig struct {
	Blocks            []*Blocks        `json:"blocks,omitempty"`
	Clusters          []*Clusters      `json:"clusters,omitempty"`
	CvmGateway        string           `json:"cvm_gateway,omitempty"`
	CvmNetmask        string           `json:"cvm_netmask,omitempty"`
	DNSServers        string           `json:"dns_servers,omitempty"`
	HypervProductKey  string           `json:"hyperv_product_key,omitempty"`
	HypervSku         string           `json:"hyperv_sku,omitempty"`
	HypervisorGateway string           `json:"hypervisor_gateway,omitempty"`
	HypervisorNetmask string           `json:"hypervisor_netmask,omitempty"`
	IpmiGateway       string           `json:"ipmi_gateway,omitempty"`
	IpmiNetmask       string           `json:"ipmi_netmask,omitempty"`
	HypervisorIsoURL  *HypervisorIso   `json:"hypervisor_iso_url"`
	HypervisorIsos    []*HypervisorIso `json:"hypervisor_isos,omitempty"`
	NosPackageURL     *NosPackageURL   `json:"nos_package_url,omitempty"`
}

// format for Hypervisor ISO
type HypervisorIso struct {
	HypervisorType *string `json:"hypervisor_type,omitempty"`
	Sha256sum      *string `json:"sha256sum,omitempty"`
	URL            *string `json:"url,omitempty"`
}

// NOS Package URL
type NosPackageURL struct {
	Sha256sum string `json:"sha256sum,omitempty"`
	URL       string `json:"url,omitempty"`
}

// format of blocks
type Blocks struct {
	BlockID string   `json:"block_id,omitempty"`
	Nodes   []*Nodes `json:"nodes,omitempty"`
}

// format of nodes inside block field
type Nodes struct {
	CvmIP                      string                 `json:"cvm_ip,omitempty"`
	CvmVlanID                  int                    `json:"cvm_vlan_id,omitempty"`
	FcImagedNodeUUID           string                 `json:"fc_imaged_node_uuid,omitempty"`
	Hypervisor                 string                 `json:"hypervisor,omitempty"`
	HypervisorHostname         string                 `json:"hypervisor_hostname,omitempty"`
	HypervisorIP               string                 `json:"hypervisor_ip,omitempty"`
	ImageNow                   bool                   `json:"image_now,omitempty"`
	IpmiIP                     string                 `json:"ipmi_ip,omitempty"`
	IPv6Address                string                 `json:"ipv6_address,omitempty"`
	NodePosition               string                 `json:"node_position,omitempty"`
	NodeSerial                 string                 `json:"node_serial,omitempty"`
	HardwareAttributesOverride map[string]interface{} `json:"hardware_attributes_override,omitempty"`
}

// format of clusters
type Clusters struct {
	ClusterExternalIP     string    `json:"cluster_external_ip,omitempty"`
	ClusterInitNow        bool      `json:"cluster_init_now,omitempty"`
	ClusterInitSuccessful bool      `json:"cluster_init_successful,omitempty"`
	ClusterMembers        []*string `json:"cluster_members,omitempty"`
	ClusterName           string    `json:"cluster_name,omitempty"`
	CvmDNSServers         string    `json:"cvm_dns_servers,omitempty"`
	CvmNtpServers         string    `json:"cvm_ntp_servers,omitempty"`
	RedundancyFactor      int       `json:"redundancy_factor,omitempty"`
	TimeZone              string    `json:"timezone,omitempty"`
}

// input of Imaged node details
type ImagedNodeDetailsInput struct {
	ImagedNodeUUID string `json:"imaged_node_uuid"`
}
