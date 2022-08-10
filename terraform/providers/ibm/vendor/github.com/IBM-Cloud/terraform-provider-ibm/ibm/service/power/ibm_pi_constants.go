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
	//Added timeout values for warning  and active status
	warningTimeOut = 60 * time.Second
	activeTimeOut  = 2 * time.Minute
	// power service instance capabilities
	CUSTOM_VIRTUAL_CORES          = "custom-virtualcores"
	PIInstanceNetwork             = "pi_network"
	PIInstanceStoragePool         = "pi_storage_pool"
	PISAPInstanceProfileID        = "pi_sap_profile_id"
	PISAPInstanceDeploymentType   = "pi_sap_deployment_type"
	PIInstanceStoragePoolAffinity = "pi_storage_pool_affinity"

	// Placement Group
	PIPlacementGroupID      = "placement_group_id"
	PIPlacementGroupMembers = "members"

	// Volume
	PIAffinityPolicy        = "pi_affinity_policy"
	PIAffinityVolume        = "pi_affinity_volume"
	PIAffinityInstance      = "pi_affinity_instance"
	PIAntiAffinityInstances = "pi_anti_affinity_instances"
	PIAntiAffinityVolumes   = "pi_anti_affinity_volumes"

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

	// status
	// common status states
	StatusShutoff = "SHUTOFF"
	StatusActive  = "ACTIVE"
	StatusResize  = "RESIZE"
	StatusError   = "ERROR"
	StatusBuild   = "BUILD"
	SctionStart   = "start"
	SctionStop    = "stop"
)
