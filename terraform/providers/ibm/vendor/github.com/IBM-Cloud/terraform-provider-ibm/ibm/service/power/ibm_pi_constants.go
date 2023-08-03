package power

import "time"

const (
	// used by all
	Arg_CloudInstanceID = "pi_cloud_instance_id"

	// Keys
	Arg_KeyName = "pi_key_name"
	Arg_Key     = "pi_ssh_key"

	Attr_KeyID           = "key_id"
	Attr_Keys            = "keys"
	Attr_KeyCreationDate = "creation_date"
	Attr_Key             = "ssh_key"
	Attr_KeyName         = "name"

	// SAP Profile
	PISAPProfiles         = "profiles"
	PISAPProfileCertified = "certified"
	PISAPProfileCores     = "cores"
	PISAPProfileMemory    = "memory"
	PISAPProfileID        = "profile_id"
	PISAPProfileType      = "type"

	// DHCP
	Arg_DhcpCidr              = "pi_cidr"
	Arg_DhcpID                = "pi_dhcp_id"
	Arg_DhcpCloudConnectionID = "pi_cloud_connection_id"
	Arg_DhcpDnsServer         = "pi_dns_server"
	Arg_DhcpName              = "pi_dhcp_name"
	Arg_DhcpSnatEnabled       = "pi_dhcp_snat_enabled"

	Attr_DhcpServers           = "servers"
	Attr_DhcpID                = "dhcp_id"
	Attr_DhcpLeases            = "leases"
	Attr_DhcpLeaseInstanceIP   = "instance_ip"
	Attr_DhcpLeaseInstanceMac  = "instance_mac"
	Attr_DhcpNetworkDeprecated = "network" // to deprecate
	Attr_DhcpNetworkID         = "network_id"
	Attr_DhcpNetworkName       = "network_name"
	Attr_DhcpStatus            = "status"

	// Instance
	Arg_PVMInstanceId           = "pi_instance_id"
	Arg_PVMInstanceActionType   = "pi_action"
	Arg_PVMInstanceHealthStatus = "pi_health_status"

	Attr_Status       = "status"
	Attr_Progress     = "progress"
	Attr_HealthStatus = "health_status"

	PVMInstanceHealthOk      = "OK"
	PVMInstanceHealthWarning = "WARNING"

	//Added timeout values for warning  and active status
	warningTimeOut = 60 * time.Second
	activeTimeOut  = 2 * time.Minute
	// power service instance capabilities
	CUSTOM_VIRTUAL_CORES                 = "custom-virtualcores"
	PIInstanceDeploymentType             = "pi_deployment_type"
	PIInstanceNetwork                    = "pi_network"
	PIInstanceStoragePool                = "pi_storage_pool"
	PISAPInstanceProfileID               = "pi_sap_profile_id"
	PISAPInstanceDeploymentType          = "pi_sap_deployment_type"
	PIInstanceStoragePoolAffinity        = "pi_storage_pool_affinity"
	Arg_PIInstanceSharedProcessorPool    = "pi_shared_processor_pool"
	Attr_PIInstanceSharedProcessorPool   = "shared_processor_pool"
	Attr_PIInstanceSharedProcessorPoolID = "shared_processor_pool_id"

	// Placement Group
	PIPlacementGroupID      = "placement_group_id"
	PIPlacementGroupMembers = "members"

	// Volume
	PIAffinityPolicy        = "pi_affinity_policy"
	PIAffinityVolume        = "pi_affinity_volume"
	PIAffinityInstance      = "pi_affinity_instance"
	PIAntiAffinityInstances = "pi_anti_affinity_instances"
	PIAntiAffinityVolumes   = "pi_anti_affinity_volumes"

	// IBM PI Volume Group
	PIVolumeGroupName                 = "pi_volume_group_name"
	PIVolumeGroupsVolumeIds           = "pi_volume_ids"
	PIVolumeGroupConsistencyGroupName = "pi_consistency_group_name"
	PIVolumeGroupID                   = "pi_volume_group_id"
	PIVolumeGroupAction               = "pi_volume_group_action"
	PIVolumeOnboardingID              = "pi_volume_onboarding_id"

	// Disaster Recovery Location
	PIDRLocation = "location"

	// VPN
	PIVPNConnectionId                         = "connection_id"
	PIVPNConnectionStatus                     = "connection_status"
	PIVPNConnectionDeadPeerDetection          = "dead_peer_detections"
	PIVPNConnectionDeadPeerDetectionAction    = "action"
	PIVPNConnectionDeadPeerDetectionInterval  = "interval"
	PIVPNConnectionDeadPeerDetectionThreshold = "threshold"
	PIVPNConnectionLocalGatewayAddress        = "local_gateway_address"
	PIVPNConnectionVpnGatewayAddress          = "gateway_address"

	// Cloud Connections
	PICloudConnectionTransitEnabled = "pi_cloud_connection_transit_enabled"

	// Shared Processor Pool
	Arg_SharedProcessorPoolName                      = "pi_shared_processor_pool_name"
	Arg_SharedProcessorPoolHostGroup                 = "pi_shared_processor_pool_host_group"
	Arg_SharedProcessorPoolPlacementGroupID          = "pi_shared_processor_pool_placement_group_id"
	Arg_SharedProcessorPoolReservedCores             = "pi_shared_processor_pool_reserved_cores"
	Arg_SharedProcessorPoolID                        = "pi_shared_processor_pool_id"
	Attr_SharedProcessorPoolID                       = "shared_processor_pool_id"
	Attr_SharedProcessorPoolName                     = "name"
	Attr_SharedProcessorPoolReservedCores            = "reserved_cores"
	Attr_SharedProcessorPoolAvailableCores           = "available_cores"
	Attr_SharedProcessorPoolAllocatedCores           = "allocated_cores"
	Attr_SharedProcessorPoolHostID                   = "host_id"
	Attr_SharedProcessorPoolStatus                   = "status"
	Attr_SharedProcessorPoolStatusDetail             = "status_detail"
	Attr_SharedProcessorPoolPlacementGroups          = "spp_placement_groups"
	Attr_SharedProcessorPoolInstances                = "instances"
	Attr_SharedProcessorPoolInstanceCpus             = "cpus"
	Attr_SharedProcessorPoolInstanceUncapped         = "uncapped"
	Attr_SharedProcessorPoolInstanceAvailabilityZone = "availability_zone"
	Attr_SharedProcessorPoolInstanceId               = "id"
	Attr_SharedProcessorPoolInstanceMemory           = "memory"
	Attr_SharedProcessorPoolInstanceName             = "name"
	Attr_SharedProcessorPoolInstanceStatus           = "status"
	Attr_SharedProcessorPoolInstanceVcpus            = "vcpus"

	// SPP Placement Group
	Arg_SPPPlacementGroupName     = "pi_spp_placement_group_name"
	Arg_SPPPlacementGroupPolicy   = "pi_spp_placement_group_policy"
	Attr_SPPPlacementGroupID      = "spp_placement_group_id"
	Attr_SPPPlacementGroupMembers = "members"
	Arg_SPPPlacementGroupID       = "pi_spp_placement_group_id"
	Attr_SPPPlacementGroupPolicy  = "policy"
	Attr_SPPPlacementGroupName    = "name"

	// status
	// common status states
	StatusShutoff = "SHUTOFF"
	StatusActive  = "ACTIVE"
	StatusResize  = "RESIZE"
	StatusError   = "ERROR"
	StatusBuild   = "BUILD"
	StatusPending = "PENDING"
	SctionStart   = "start"
	SctionStop    = "stop"
)
