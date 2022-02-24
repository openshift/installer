package vsphere

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi/vim25/types"
)

var lacpAPIVersionAllowedValues = []string{
	string(types.VMwareDvsLacpApiVersionSingleLag),
	string(types.VMwareDvsLacpApiVersionMultipleLag),
}

var multicastFilteringModeAllowedValues = []string{
	string(types.VMwareDvsMulticastFilteringModeLegacyFiltering),
	string(types.VMwareDvsMulticastFilteringModeSnooping),
}

var privateVLANTypeAllowedValues = []string{
	string(types.VmwareDistributedVirtualSwitchPvlanPortTypePromiscuous),
	string(types.VmwareDistributedVirtualSwitchPvlanPortTypeIsolated),
	string(types.VmwareDistributedVirtualSwitchPvlanPortTypeCommunity),
}

var networkResourceControlAllowedValues = []string{
	string(types.DistributedVirtualSwitchNetworkResourceControlVersionVersion2),
	string(types.DistributedVirtualSwitchNetworkResourceControlVersionVersion3),
}

var infrastructureTrafficClassValues = []string{
	string(types.DistributedVirtualSwitchHostInfrastructureTrafficClassManagement),
	string(types.DistributedVirtualSwitchHostInfrastructureTrafficClassFaultTolerance),
	string(types.DistributedVirtualSwitchHostInfrastructureTrafficClassVmotion),
	string(types.DistributedVirtualSwitchHostInfrastructureTrafficClassVirtualMachine),
	string(types.DistributedVirtualSwitchHostInfrastructureTrafficClassISCSI),
	string(types.DistributedVirtualSwitchHostInfrastructureTrafficClassNfs),
	string(types.DistributedVirtualSwitchHostInfrastructureTrafficClassHbr),
	string(types.DistributedVirtualSwitchHostInfrastructureTrafficClassVsan),
	string(types.DistributedVirtualSwitchHostInfrastructureTrafficClassVdp),
}

var sharesLevelAllowedValues = []string{
	string(types.SharesLevelLow),
	string(types.SharesLevelNormal),
	string(types.SharesLevelHigh),
	string(types.SharesLevelCustom),
}

// schemaVMwareDVSConfigSpec returns schema items for resources that need to work
// with a VMwareDVSConfigSpec.
func schemaVMwareDVSConfigSpec() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		// DVSContactInfo
		"contact_detail": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The contact detail for this DVS.",
		},
		"contact_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The contact name for this DVS.",
		},

		// DistributedVirtualSwitchHostMemberConfigSpec
		"host": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "A host member specification.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// DistributedVirtualSwitchHostMemberPnicSpec
					"devices": {
						Type:        schema.TypeList,
						Description: "Name of the physical NIC to be added to the proxy switch.",
						Required:    true,
						MinItems:    1,
						Elem:        &schema.Schema{Type: schema.TypeString},
					},
					"host_system_id": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The managed object ID of the host this specification applies to.",
						ValidateFunc: validation.NoZeroValues,
					},
				},
			},
		},

		// VMwareDVSPvlanMapEntry
		"pvlan_mapping": {
			Type:        schema.TypeSet,
			Optional:    true,
			Description: "A private VLAN (PVLAN) mapping.",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"primary_vlan_id": {
						Type:         schema.TypeInt,
						Required:     true,
						Description:  "The primary VLAN ID. The VLAN IDs of 0 and 4095 are reserved and cannot be used in this property.",
						ValidateFunc: validation.IntBetween(1, 4094),
					},
					"secondary_vlan_id": {
						Type:         schema.TypeInt,
						Required:     true,
						Description:  "The secondary VLAN ID. The VLAN IDs of 0 and 4095 are reserved and cannot be used in this property.",
						ValidateFunc: validation.IntBetween(1, 4094),
					},
					"pvlan_type": {
						Type:         schema.TypeString,
						Required:     true,
						Description:  "The private VLAN type. Valid values are promiscuous, community and isolated.",
						ValidateFunc: validation.StringInSlice(privateVLANTypeAllowedValues, false),
					},
				},
			},
		},
		"ignore_other_pvlan_mappings": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to ignore existing PVLAN mappings not managed by this resource. Defaults to false.",
			Default:     false,
		},

		// VMwareIpfixConfig (Netflow)
		"netflow_active_flow_timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "The number of seconds after which active flows are forced to be exported to the collector.",
			Default:      60,
			ValidateFunc: validation.IntBetween(60, 3600),
		},
		"netflow_collector_ip_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "IP address for the netflow collector, using IPv4 or IPv6. IPv6 is supported in vSphere Distributed Switch Version 6.0 or later.",
		},
		"netflow_collector_port": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "The port for the netflow collector.",
			ValidateFunc: validation.IntBetween(0, 65535),
		},
		"netflow_idle_flow_timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "The number of seconds after which idle flows are forced to be exported to the collector.",
			Default:      15,
			ValidateFunc: validation.IntBetween(10, 600),
		},
		"netflow_internal_flows_only": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Whether to limit analysis to traffic that has both source and destination served by the same host.",
		},
		"netflow_observation_domain_id": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "The observation Domain ID for the netflow collector.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"netflow_sampling_rate": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "The ratio of total number of packets to the number of packets analyzed. Set to 0 to disable sampling, meaning that all packets are analyzed.",
			ValidateFunc: validation.IntAtLeast(0),
		},

		// LinkDiscoveryProtocolConfig
		"link_discovery_operation": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Whether to advertise or listen for link discovery. Valid values are advertise, both, listen, and none.",
			Default:      string(types.LinkDiscoveryProtocolConfigOperationTypeListen),
			ValidateFunc: validation.StringInSlice(linkDiscoveryProtocolConfigOperationAllowedValues, false),
		},
		"link_discovery_protocol": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "The discovery protocol type. Valid values are cdp and lldp.",
			Default:      string(types.LinkDiscoveryProtocolConfigProtocolTypeCdp),
			ValidateFunc: validation.StringInSlice(linkDiscoveryProtocolConfigProtocolAllowedValues, false),
		},

		// DVSNameArrayUplinkPortPolicy
		"uplinks": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "A list of uplink ports. The contents of this list control both the uplink count and names of the uplinks on the DVS across hosts.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},

		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name for the DVS. Must be unique in the folder that it is being created in.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the DVS.",
		},
		"ipv4_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The IPv4 address of the switch. This can be used to see the DVS as a unique device with NetFlow.",
		},
		"lacp_api_version": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			Description:  "The Link Aggregation Control Protocol group version in the switch. Can be one of singleLag or multipleLag.",
			ValidateFunc: validation.StringInSlice(lacpAPIVersionAllowedValues, false),
		},
		"max_mtu": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			Description:  "The maximum MTU on the switch.",
			ValidateFunc: validation.IntBetween(1, 9000),
		},
		"multicast_filtering_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			Description:  "The multicast filtering mode on the switch. Can be one of legacyFiltering, or snooping.",
			ValidateFunc: validation.StringInSlice(multicastFilteringModeAllowedValues, false),
		},
		"network_resource_control_version": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			Description:  "The network I/O control version to use. Can be one of version2 or version3.",
			ValidateFunc: validation.StringInSlice(networkResourceControlAllowedValues, false),
		},

		"config_version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "The version string of the configuration that this spec is trying to change.",
		},
	}

	structure.MergeSchema(s, schemaVMwareDVSPortSetting())
	structure.MergeSchema(s, schemaDvsHostInfrastructureTrafficResource())
	return s
}

// expandDVSContactInfo reads certain ResourceData keys and
// returns a DVSContactInfo.
func expandDVSContactInfo(d *schema.ResourceData) *types.DVSContactInfo {
	obj := &types.DVSContactInfo{
		Name:    d.Get("contact_name").(string),
		Contact: d.Get("contact_detail").(string),
	}
	return obj
}

// flattenDVSContactInfo reads various fields from a
// DVSContactInfo into the passed in ResourceData.
func flattenDVSContactInfo(d *schema.ResourceData, obj types.DVSContactInfo) error {
	d.Set("contact_name", obj.Name)
	d.Set("conatct_detail", obj.Contact)
	return nil
}

// expandDistributedVirtualSwitchHostMemberConfigSpec reads certain keys from a
// Set object map and returns a DistributedVirtualSwitchHostMemberConfigSpec.
func expandDistributedVirtualSwitchHostMemberConfigSpec(d map[string]interface{}) types.DistributedVirtualSwitchHostMemberConfigSpec {
	hostRef := &types.ManagedObjectReference{
		Type:  "HostSystem",
		Value: d["host_system_id"].(string),
	}

	var pnSpecs []types.DistributedVirtualSwitchHostMemberPnicSpec
	nics := structure.SliceInterfacesToStrings(d["devices"].([]interface{}))
	for _, nic := range nics {
		pnSpec := types.DistributedVirtualSwitchHostMemberPnicSpec{
			PnicDevice: nic,
		}
		pnSpecs = append(pnSpecs, pnSpec)
	}
	backing := types.DistributedVirtualSwitchHostMemberPnicBacking{
		PnicSpec: pnSpecs,
	}

	obj := types.DistributedVirtualSwitchHostMemberConfigSpec{
		Host:    *hostRef,
		Backing: &backing,
	}
	return obj
}

// flattenDistributedVirtualSwitchHostMemberConfigSpec reads various fields
// from a DistributedVirtualSwitchHostMemberConfigSpec and returns a Set object
// map.
//
// This is the flatten counterpart to
// expandDistributedVirtualSwitchHostMemberConfigSpec.
func flattenDistributedVirtualSwitchHostMember(obj types.DistributedVirtualSwitchHostMember) map[string]interface{} {
	d := make(map[string]interface{})
	d["host_system_id"] = obj.Config.Host.Value

	var devices []string
	backing := obj.Config.Backing.(*types.DistributedVirtualSwitchHostMemberPnicBacking)
	for _, spec := range backing.PnicSpec {
		devices = append(devices, spec.PnicDevice)
	}

	d["devices"] = devices

	return d
}

// expandSliceOfDistributedVirtualSwitchHostMemberConfigSpec expands all host
// entires for a VMware DVS, detecting if a host spec needs to be added,
// removed, or updated as well. The whole slice is returned.
func expandSliceOfDistributedVirtualSwitchHostMemberConfigSpec(d *schema.ResourceData) []types.DistributedVirtualSwitchHostMemberConfigSpec {
	var specs []types.DistributedVirtualSwitchHostMemberConfigSpec
	o, n := d.GetChange("host")
	os := o.(*schema.Set)
	ns := n.(*schema.Set)

	// Make an intersection set. These hosts have not changed so we don't bother
	// with them.
	is := os.Intersection(ns)
	os = os.Difference(is)
	ns = ns.Difference(is)

	// Our old and new sets now have an accurate description of hosts that may
	// have been added, removed, or changed. Add removed and modified hosts
	// first.
	for _, oe := range os.List() {
		om := oe.(map[string]interface{})
		var found bool
		for _, ne := range ns.List() {
			nm := ne.(map[string]interface{})
			if nm["host_system_id"] == om["host_system_id"] {
				found = true
			}
		}
		if !found {
			spec := expandDistributedVirtualSwitchHostMemberConfigSpec(om)
			spec.Operation = string(types.ConfigSpecOperationRemove)
			specs = append(specs, spec)
		}
	}

	// Process new hosts now. These are ones that are only present in the new
	// set.
	for _, ne := range ns.List() {
		nm := ne.(map[string]interface{})
		var found bool
		for _, oe := range os.List() {
			om := oe.(map[string]interface{})
			if om["host_system_id"] == nm["host_system_id"] {
				found = true
			}
		}
		spec := expandDistributedVirtualSwitchHostMemberConfigSpec(nm)
		if !found {
			spec.Operation = string(types.ConfigSpecOperationAdd)
		} else {
			spec.Operation = string(types.ConfigSpecOperationEdit)
		}
		specs = append(specs, spec)
	}

	// Done!
	return specs
}

// flattenSliceOfDistributedVirtualSwitchHostMember creates a set of all host
// entries for a supplied slice of DistributedVirtualSwitchHostMember.
//
// This is the flatten counterpart to
// expandSliceOfDistributedVirtualSwitchHostMemberConfigSpec.
func flattenSliceOfDistributedVirtualSwitchHostMember(d *schema.ResourceData, members []types.DistributedVirtualSwitchHostMember) error {
	var hosts []map[string]interface{}
	for _, m := range members {
		hosts = append(hosts, flattenDistributedVirtualSwitchHostMember(m))
	}
	if err := d.Set("host", hosts); err != nil {
		return err
	}
	return nil
}

// expandVMwareDVSPvlanConfigSpec reads certain keys from a Set object map
// representing a VMwareDVSPvlanMapEntry and returns a
// VMwareDVSPvlanConfigSpec.
func expandVMwareDVSPvlanConfigSpec(d map[string]interface{}) types.VMwareDVSPvlanConfigSpec {
	mapEntry := types.VMwareDVSPvlanMapEntry{
		PrimaryVlanId:   int32(d["primary_vlan_id"].(int)),
		SecondaryVlanId: int32(d["secondary_vlan_id"].(int)),
		PvlanType:       d["pvlan_type"].(string),
	}

	obj := types.VMwareDVSPvlanConfigSpec{
		PvlanEntry: mapEntry,
	}
	return obj
}

// flattenVMwareDVSPvlanMapEntry reads various fields
// from a VMwareDVSPvlanMapEntry and returns a Set object
// map.
//
// This is the flatten counterpart to
// expandVMwareDVSPvlanConfigSpec.
func flattenVMwareDVSPvlanMapEntry(obj types.VMwareDVSPvlanMapEntry) map[string]interface{} {
	d := make(map[string]interface{})
	d["primary_vlan_id"] = int(obj.PrimaryVlanId)
	d["secondary_vlan_id"] = int(obj.SecondaryVlanId)
	d["pvlan_type"] = obj.PvlanType
	return d
}

// expandSliceOfVMwareDVSPvlanConfigSpec expands all pvlan mapping
// entries for a VMware DVS, detecting if a pvlan mapping needs to be added,
// removed, or updated as well. The whole slice is returned.
func expandSliceOfVMwareDVSPvlanConfigSpec(d *schema.ResourceData) []types.VMwareDVSPvlanConfigSpec {
	var specs []types.VMwareDVSPvlanConfigSpec
	o, n := d.GetChange("pvlan_mapping")
	os := o.(*schema.Set)
	ns := n.(*schema.Set)

	// Make an intersection set. These mappings have not changed so we don't bother
	// with them.
	is := os.Intersection(ns)
	os = os.Difference(is)
	ns = ns.Difference(is)

	// Our old and new sets now have an accurate description of mappings that may
	// have been added, removed, or changed. Add removed and modified mappings
	// first.
	for _, oe := range os.List() {
		om := oe.(map[string]interface{})
		var found bool
		for _, ne := range ns.List() {
			nm := ne.(map[string]interface{})
			if nm["secondary_vlan_id"] == om["secondary_vlan_id"] && om["pvlan_type"] != "promiscuous" {
				found = true
			}
		}
		if !found {
			spec := expandVMwareDVSPvlanConfigSpec(om)
			spec.Operation = string(types.ConfigSpecOperationRemove)
			specs = append(specs, spec)
		}
	}

	// Process new mappings now. These are ones that are only present in the new
	// set.
	for _, ne := range ns.List() {
		nm := ne.(map[string]interface{})
		var found bool
		for _, oe := range os.List() {
			om := oe.(map[string]interface{})
			if nm["secondary_vlan_id"] == om["secondary_vlan_id"] && om["pvlan_type"] != "promiscuous" {
				found = true
			}
		}
		spec := expandVMwareDVSPvlanConfigSpec(nm)
		if !found {
			spec.Operation = string(types.ConfigSpecOperationAdd)
		} else {
			spec.Operation = string(types.ConfigSpecOperationEdit)
		}
		specs = append(specs, spec)
	}

	// Done!
	return specs
}

// flattenSliceOfVMwareDVSPvlanMapEntry creates a set of all host
// entries for a supplied slice of VMwareDVSPvlanMapEntry.
//
// This is the flatten counterpart to
// expandSliceOfVMwareDVSPvlanConfigSpec.
func flattenSliceOfVMwareDVSPvlanMapEntry(d *schema.ResourceData, entries []types.VMwareDVSPvlanMapEntry) error {
	oldPvlanMappings := d.Get("pvlan_mapping").(*schema.Set)
	var mappings []map[string]interface{}
	for _, entry := range entries {
		// TODO: Should the ignore_other_pvlan_mappings immediately affect the way it treats existing resources or should it have to be applied first?
		if flattened := flattenVMwareDVSPvlanMapEntry(entry); d.Get("ignore_other_pvlan_mappings").(bool) && !oldPvlanMappings.Contains(flattened) {
			log.Printf("[DEBUG] Found unmanaged pvlan_mapping (%v) and ignore_other_pvlan_mappings is true. Not reading into state.", flattened)
		} else {
			mappings = append(mappings, flattened)
		}
	}
	if err := d.Set("pvlan_mapping", mappings); err != nil {
		return err
	}
	return nil
}

// expandVMwareIpfixConfig reads certain ResourceData keys and
// returns a VMwareIpfixConfig.
func expandVMwareIpfixConfig(d *schema.ResourceData) *types.VMwareIpfixConfig {
	obj := &types.VMwareIpfixConfig{
		ActiveFlowTimeout:   int32(d.Get("netflow_active_flow_timeout").(int)),
		CollectorIpAddress:  d.Get("netflow_collector_ip_address").(string),
		CollectorPort:       int32(d.Get("netflow_collector_port").(int)),
		IdleFlowTimeout:     int32(d.Get("netflow_idle_flow_timeout").(int)),
		InternalFlowsOnly:   d.Get("netflow_internal_flows_only").(bool),
		ObservationDomainId: int64(d.Get("netflow_observation_domain_id").(int)),
		SamplingRate:        int32(d.Get("netflow_sampling_rate").(int)),
	}
	return obj
}

// flattenVMwareIpfixConfig reads various fields from a
// VMwareIpfixConfig into the passed in ResourceData.
func flattenVMwareIpfixConfig(d *schema.ResourceData, obj *types.VMwareIpfixConfig) error {
	d.Set("netflow_active_flow_timeout", obj.ActiveFlowTimeout)
	d.Set("netflow_collector_ip_address", obj.CollectorIpAddress)
	d.Set("netflow_collector_port", obj.CollectorPort)
	d.Set("netflow_idle_flow_timeout", obj.IdleFlowTimeout)
	d.Set("netflow_internal_flows_only", obj.InternalFlowsOnly)
	d.Set("netflow_observation_domain_id", obj.ObservationDomainId)
	d.Set("netflow_sampling_rate", obj.SamplingRate)
	return nil
}

// schemaDvsHostInfrastructureTrafficResource returns the respective schema
// keys for the various kinds of network I/O control traffic classes. The
// schema items are generated dynamically off of the list of available traffic
// classes for the currently supported vSphere API. Not all traffic classes may
// be supported across all DVS and network I/O control versions.
func schemaDvsHostInfrastructureTrafficResource() map[string]*schema.Schema {
	s := make(map[string]*schema.Schema)
	shareLevelFmt := "The allocation level for the %s traffic class. Can be one of high, low, normal, or custom."
	shareCountFmt := "The amount of shares to allocate to the %s traffic class for a custom share level."
	maxMbitFmt := "The maximum allowed usage for the %s traffic class, in Mbits/sec."
	resMbitFmt := "The amount of guaranteed bandwidth for the %s traffic class, in Mbits/sec."

	for _, class := range infrastructureTrafficClassValues {
		shareLevelKey := fmt.Sprintf("%s_share_level", strings.ToLower(class))
		shareCountKey := fmt.Sprintf("%s_share_count", strings.ToLower(class))
		maxMbitKey := fmt.Sprintf("%s_maximum_mbit", strings.ToLower(class))
		resMbitKey := fmt.Sprintf("%s_reservation_mbit", strings.ToLower(class))

		s[shareLevelKey] = &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			Description:  fmt.Sprintf(shareLevelFmt, class),
			ValidateFunc: validation.StringInSlice(sharesLevelAllowedValues, false),
		}
		s[shareCountKey] = &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			Description:  fmt.Sprintf(shareCountFmt, class),
			ValidateFunc: validation.IntAtLeast(0),
		}
		s[maxMbitKey] = &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			Description:  fmt.Sprintf(maxMbitFmt, class),
			ValidateFunc: validation.IntAtLeast(-1),
		}
		s[resMbitKey] = &schema.Schema{
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			Description:  fmt.Sprintf(resMbitFmt, class),
			ValidateFunc: validation.IntAtLeast(-1),
		}
	}

	return s
}

// expandDvsHostInfrastructureTrafficResource reads the network I/O control
// resource data keys for the traffic class supplied by key and returns an
// appropriate types.DvsHostInfrastructureTrafficResource reference. This
// should be checked for nil to see if it should be added to the slice in the
// config.
func expandDvsHostInfrastructureTrafficResource(d *schema.ResourceData, key string) *types.DvsHostInfrastructureTrafficResource {
	shareLevelKey := fmt.Sprintf("%s_share_level", strings.ToLower(key))
	shareCountKey := fmt.Sprintf("%s_share_count", strings.ToLower(key))
	maxMbitKey := fmt.Sprintf("%s_maximum_mbit", strings.ToLower(key))
	resMbitKey := fmt.Sprintf("%s_reservation_mbit", strings.ToLower(key))

	obj := &types.DvsHostInfrastructureTrafficResource{
		AllocationInfo: types.DvsHostInfrastructureTrafficResourceAllocation{
			Limit:       structure.GetInt64Ptr(d, maxMbitKey),
			Reservation: structure.GetInt64Ptr(d, resMbitKey),
		},
	}
	shares := &types.SharesInfo{
		Level:  types.SharesLevel(d.Get(shareLevelKey).(string)),
		Shares: int32(d.Get(shareCountKey).(int)),
	}
	if !structure.AllFieldsEmpty(shares) {
		obj.AllocationInfo.Shares = shares
	}

	if structure.AllFieldsEmpty(obj) {
		return nil
	}
	obj.Key = key
	return obj
}

// flattenDvsHostInfrastructureTrafficResource reads various fields from a
// DvsHostInfrastructureTrafficResource and sets appropriate keys in the
// supplied ResourceData.
func flattenDvsHostInfrastructureTrafficResource(d *schema.ResourceData, obj types.DvsHostInfrastructureTrafficResource, key string) error {
	shareLevelKey := fmt.Sprintf("%s_share_level", strings.ToLower(key))
	shareCountKey := fmt.Sprintf("%s_share_count", strings.ToLower(key))
	maxMbitKey := fmt.Sprintf("%s_maximum_mbit", strings.ToLower(key))
	resMbitKey := fmt.Sprintf("%s_reservation_mbit", strings.ToLower(key))

	structure.SetInt64Ptr(d, maxMbitKey, obj.AllocationInfo.Limit)
	structure.SetInt64Ptr(d, resMbitKey, obj.AllocationInfo.Reservation)
	if obj.AllocationInfo.Shares != nil {
		d.Set(shareLevelKey, obj.AllocationInfo.Shares.Level)
		d.Set(shareCountKey, obj.AllocationInfo.Shares.Shares)
	}
	return nil
}

// expandSliceOfDvsHostInfrastructureTrafficResource expands all network I/O
// control resource entries that are currently supported in API, and returns a
// slice of DvsHostInfrastructureTrafficResource.
func expandSliceOfDvsHostInfrastructureTrafficResource(d *schema.ResourceData) []types.DvsHostInfrastructureTrafficResource {
	var s []types.DvsHostInfrastructureTrafficResource
	for _, key := range infrastructureTrafficClassValues {
		v := expandDvsHostInfrastructureTrafficResource(d, key)
		if v != nil {
			s = append(s, *v)
		}
	}
	return s
}

// flattenSliceOfDvsHostInfrastructureTrafficResource reads in the supplied network I/O control allocation entries supplied via a respective DVSConfigInfo field and sets the appropriate keys in the supplied ResourceData.
func flattenSliceOfDvsHostInfrastructureTrafficResource(d *schema.ResourceData, s []types.DvsHostInfrastructureTrafficResource) error {
	for _, v := range s {
		if err := flattenDvsHostInfrastructureTrafficResource(d, v, v.Key); err != nil {
			return err
		}
	}
	return nil
}

// expandDVSNameArrayUplinkPortPolicy reads certain ResourceData keys and
// returns a DVSNameArrayUplinkPortPolicy.
func expandDVSNameArrayUplinkPortPolicy(d *schema.ResourceData) *types.DVSNameArrayUplinkPortPolicy {
	obj := &types.DVSNameArrayUplinkPortPolicy{
		UplinkPortName: structure.SliceInterfacesToStrings(d.Get("uplinks").([]interface{})),
	}
	if structure.AllFieldsEmpty(obj) {
		return nil
	}
	return obj
}

// flattenDVSNameArrayUplinkPortPolicy reads various fields from a
// DVSNameArrayUplinkPortPolicy into the passed in ResourceData.
func flattenDVSNameArrayUplinkPortPolicy(d *schema.ResourceData, obj *types.DVSNameArrayUplinkPortPolicy) error {
	if err := d.Set("uplinks", obj.UplinkPortName); err != nil {
		return err
	}
	return nil
}

// expandVMwareDVSConfigSpec reads certain ResourceData keys and
// returns a VMwareDVSConfigSpec.
func expandVMwareDVSConfigSpec(d *schema.ResourceData) *types.VMwareDVSConfigSpec {
	obj := &types.VMwareDVSConfigSpec{
		DVSConfigSpec: types.DVSConfigSpec{
			Name:                                d.Get("name").(string),
			ConfigVersion:                       d.Get("config_version").(string),
			DefaultPortConfig:                   expandVMwareDVSPortSetting(d, "distributed_virtual_switch"),
			Host:                                expandSliceOfDistributedVirtualSwitchHostMemberConfigSpec(d),
			Description:                         d.Get("description").(string),
			Contact:                             expandDVSContactInfo(d),
			SwitchIpAddress:                     d.Get("ipv4_address").(string),
			InfrastructureTrafficResourceConfig: expandSliceOfDvsHostInfrastructureTrafficResource(d),
			NetworkResourceControlVersion:       d.Get("network_resource_control_version").(string),
			UplinkPortPolicy:                    expandDVSNameArrayUplinkPortPolicy(d),
		},
		PvlanConfigSpec:             expandSliceOfVMwareDVSPvlanConfigSpec(d),
		MaxMtu:                      int32(d.Get("max_mtu").(int)),
		LinkDiscoveryProtocolConfig: expandLinkDiscoveryProtocolConfig(d),
		IpfixConfig:                 expandVMwareIpfixConfig(d),
		LacpApiVersion:              d.Get("lacp_api_version").(string),
		MulticastFilteringMode:      d.Get("multicast_filtering_mode").(string),
	}
	return obj
}

// flattenVMwareDVSConfigInfo reads various fields from a
// VMwareDVSConfigInfo into the passed in ResourceData.
//
// This is the flatten counterpart to expandVMwareDVSConfigSpec, as the
// configuration info from a DVS comes back as this type instead of a specific
// ConfigSpec.
func flattenVMwareDVSConfigInfo(d *schema.ResourceData, obj *types.VMwareDVSConfigInfo) error {
	d.Set("name", obj.Name)
	d.Set("config_version", obj.ConfigVersion)
	d.Set("description", obj.Description)
	d.Set("ipv4_address", obj.SwitchIpAddress)
	d.Set("max_mtu", obj.MaxMtu)
	d.Set("lacp_api_version", obj.LacpApiVersion)
	d.Set("multicast_filtering_mode", obj.MulticastFilteringMode)
	d.Set("network_resource_control_version", obj.NetworkResourceControlVersion)
	// This is not available in ConfigSpec but is available in ConfigInfo, so
	// flatten it here.
	d.Set("network_resource_control_enabled", obj.NetworkResourceManagementEnabled)

	// Version is set in this object too as ConfigInfo has the productInfo
	// property that is outside of this ConfigSpec structure.
	d.Set("version", obj.ProductInfo.Version)

	if err := flattenDVSNameArrayUplinkPortPolicy(d, obj.UplinkPortPolicy.(*types.DVSNameArrayUplinkPortPolicy)); err != nil {
		return err
	}
	if err := flattenVMwareDVSPortSetting(d, obj.DefaultPortConfig.(*types.VMwareDVSPortSetting)); err != nil {
		return err
	}
	if err := flattenSliceOfDistributedVirtualSwitchHostMember(d, obj.Host); err != nil {
		return err
	}
	if err := flattenSliceOfVMwareDVSPvlanMapEntry(d, obj.PvlanConfig); err != nil {
		return err
	}
	if err := flattenSliceOfDvsHostInfrastructureTrafficResource(d, obj.InfrastructureTrafficResourceConfig); err != nil {
		return err
	}
	if err := flattenDVSContactInfo(d, obj.Contact); err != nil {
		return err
	}
	if err := flattenLinkDiscoveryProtocolConfig(d, obj.LinkDiscoveryProtocolConfig); err != nil {
		return err
	}
	if err := flattenVMwareIpfixConfig(d, obj.IpfixConfig); err != nil {
		return err
	}
	return nil
}

// schemaDVSCreateSpec returns schema items for resources that
// need to work with a DVSCreateSpec.
func schemaDVSCreateSpec() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		// DistributedVirtualSwitchProductSpec
		"version": {
			Type:         schema.TypeString,
			Computed:     true,
			Description:  "The version of this virtual switch. Allowed versions are 6.5.0, 6.0.0, 5.5.0, 5.1.0, and 5.0.0.",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(dvsVersions, false),
		},
	}
	structure.MergeSchema(s, schemaVMwareDVSConfigSpec())

	return s
}

// expandDVSCreateSpec reads certain ResourceData keys and
// returns a DVSCreateSpec.
func expandDVSCreateSpec(d *schema.ResourceData) types.DVSCreateSpec {
	// Since we are only working with the version string from the product spec,
	// we don't have a separate expander/flattener for it. Just do that here.
	obj := types.DVSCreateSpec{
		ProductInfo: &types.DistributedVirtualSwitchProductSpec{
			Version: d.Get("version").(string),
		},
		ConfigSpec: expandVMwareDVSConfigSpec(d),
	}
	return obj
}
