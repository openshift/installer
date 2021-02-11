package vsphere

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi/vim25/types"
)

var vmwareUplinkLacpPolicyModeAllowedValues = []string{
	string(types.VMwareUplinkLacpModeActive),
	string(types.VMwareUplinkLacpModePassive),
}

var vmwareUplinkPortTeamingPolicyModeAllowedValues = []string{
	string(types.DistributedVirtualSwitchNicTeamingPolicyModeLoadbalance_ip),
	string(types.DistributedVirtualSwitchNicTeamingPolicyModeLoadbalance_srcmac),
	string(types.DistributedVirtualSwitchNicTeamingPolicyModeLoadbalance_srcid),
	string(types.DistributedVirtualSwitchNicTeamingPolicyModeFailover_explicit),
	string(types.DistributedVirtualSwitchNicTeamingPolicyModeLoadbalance_loadbased),
}

// schemaVMwareDVSPortSetting returns schema items for resources that
// need to work with a VMwareDVSPortSetting.
func schemaVMwareDVSPortSetting() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// VmwareDistributedVirtualSwitchVlanIdSpec
		"vlan_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Computed:      true,
			Description:   "The VLAN ID for single VLAN mode. 0 denotes no VLAN.",
			ConflictsWith: []string{"vlan_range", "port_private_secondary_vlan_id"},
			ValidateFunc:  validation.IntBetween(0, 4094),
		},

		// VmwareDistributedVirtualSwitchTrunkVlanSpec
		"vlan_range": {
			Type:          schema.TypeSet,
			Optional:      true,
			Computed:      true,
			Description:   "The VLAN ID for single VLAN mode. 0 denotes no VLAN.",
			ConflictsWith: []string{"vlan_id", "port_private_secondary_vlan_id"},
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"min_vlan": {
						Type:         schema.TypeInt,
						Required:     true,
						Description:  "The minimum VLAN to use in the range.",
						ValidateFunc: validation.IntBetween(0, 4094),
					},
					"max_vlan": {
						Type:         schema.TypeInt,
						Required:     true,
						Description:  "The minimum VLAN to use in the range.",
						ValidateFunc: validation.IntBetween(0, 4094),
					},
				},
			},
		},

		// VmwareDistributedVirtualSwitchPvlanSpec
		"port_private_secondary_vlan_id": {
			Type:          schema.TypeInt,
			Optional:      true,
			Computed:      true,
			Description:   "The secondary VLAN ID for this port.",
			ConflictsWith: []string{"vlan_id", "vlan_range"},
			ValidateFunc:  validation.IntBetween(1, 4094),
		},

		// VmwareUplinkPortTeamingPolicy/DVSFailureCriteria
		"check_beacon": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable beacon probing on the ports this policy applies to.",
		},

		// VmwareUplinkPortTeamingPolicy/VMwareUplinkPortOrderPolicy
		"active_uplinks": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "List of active uplinks used for load balancing, matching the names of the uplinks assigned in the DVS.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"standby_uplinks": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "List of standby uplinks used for load balancing, matching the names of the uplinks assigned in the DVS.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},

		// VmwareUplinkPortTeamingPolicy
		"teaming_policy": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			Description:  "The network adapter teaming policy. Can be one of loadbalance_ip, loadbalance_srcmac, loadbalance_srcid, failover_explicit, or loadbalance_loadbased.",
			ValidateFunc: validation.StringInSlice(vmwareUplinkPortTeamingPolicyModeAllowedValues, false),
		},
		"notify_switches": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "If true, the teaming policy will notify the broadcast network of a NIC failover, triggering cache updates.",
		},
		"failback": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "If true, the teaming policy will re-activate failed interfaces higher in precedence when they come back up.",
		},

		// DVSSecurityPolicy
		"allow_promiscuous": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Enable promiscuous mode on the network. This flag indicates whether or not all traffic is seen on a given port.",
		},
		"allow_forged_transmits": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Controls whether or not the virtual network adapter is allowed to send network traffic with a different MAC address than that of its own.",
		},
		"allow_mac_changes": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Controls whether or not the Media Access Control (MAC) address can be changed.",
		},

		// VMwareUplinkLacpPolicy
		"lacp_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Whether or not to enable LACP on all uplink ports.",
		},
		"lacp_mode": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			Description:  "The uplink LACP mode to use. Can be one of active or passive.",
			ValidateFunc: validation.StringInSlice(vmwareUplinkLacpPolicyModeAllowedValues, false),
		},

		// DVSTrafficShapingPolicy - ingress
		"ingress_shaping_average_bandwidth": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The average ingress bandwidth in bits per second if ingress shaping is enabled on the port.",
		},
		"ingress_shaping_burst_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The maximum ingress burst size allowed in bytes if ingress shaping is enabled on the port.",
		},
		"ingress_shaping_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "True if the traffic shaper is enabled for ingress traffic on the port.",
		},
		"ingress_shaping_peak_bandwidth": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The peak ingress bandwidth during bursts in bits per second if ingress traffic shaping is enabled on the port.",
		},

		// DVSTrafficShapingPolicy - egress
		"egress_shaping_average_bandwidth": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The average egress bandwidth in bits per second if egress shaping is enabled on the port.",
		},
		"egress_shaping_burst_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The maximum egress burst size allowed in bytes if egress shaping is enabled on the port.",
		},
		"egress_shaping_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "True if the traffic shaper is enabled for egress traffic on the port.",
		},
		"egress_shaping_peak_bandwidth": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "The peak egress bandwidth during bursts in bits per second if egress traffic shaping is enabled on the port.",
		},

		// VMwareDVSPortSetting
		"block_all_ports": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Indicates whether to block all ports by default.",
		},
		"netflow_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Indicates whether to enable netflow on all ports.",
		},
		"tx_uplink": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "If true, a copy of packets sent to the switch will always be forwarded to an uplink in addition to the regular packet forwarded done by the switch.",
		},
		"directpath_gen2_allowed": {
			Type:        schema.TypeBool,
			Optional:    true,
			Computed:    true,
			Description: "Allow VMDirectPath Gen2 on the ports this policy applies to.",
		},
	}
}

// expandVmwareDistributedVirtualSwitchVlanIDSpec reads certain ResourceData keys and
// returns a VmwareDistributedVirtualSwitchVlanIdSpec.
func expandVmwareDistributedVirtualSwitchVlanIDSpec(d *schema.ResourceData) *types.VmwareDistributedVirtualSwitchVlanIdSpec {
	obj := &types.VmwareDistributedVirtualSwitchVlanIdSpec{
		VlanId: int32(d.Get("vlan_id").(int)),
	}
	return obj
}

// flattenVmwareDistributedVirtualSwitchVlanIDSpec reads various fields from a
// VmwareDistributedVirtualSwitchVlanIdSpec into the passed in ResourceData.
func flattenVmwareDistributedVirtualSwitchVlanIDSpec(d *schema.ResourceData, obj *types.VmwareDistributedVirtualSwitchVlanIdSpec) error {
	d.Set("vlan_id", obj.VlanId)
	return nil
}

// expandVmwareDistributedVirtualSwitchTrunkVlanSpec reads certain ResourceData keys and
// returns a VmwareDistributedVirtualSwitchTrunkVlanSpec.
func expandVmwareDistributedVirtualSwitchTrunkVlanSpec(d *schema.ResourceData) *types.VmwareDistributedVirtualSwitchTrunkVlanSpec {
	var ranges []types.NumericRange
	data := d.Get("vlan_range").(*schema.Set).List()
	for _, v := range data {
		log.Printf("[DEBUG] processing range: %#v", v)
		r := v.(map[string]interface{})
		min := r["min_vlan"].(int)
		max := r["max_vlan"].(int)
		rng := types.NumericRange{
			Start: int32(min),
			End:   int32(max),
		}
		ranges = append(ranges, rng)
	}

	if len(ranges) < 1 {
		return nil
	}

	obj := &types.VmwareDistributedVirtualSwitchTrunkVlanSpec{
		VlanId: ranges,
	}
	return obj
}

// flattenVmwareDistributedVirtualSwitchTrunkVlanSpec reads various fields from a
// VmwareDistributedVirtualSwitchTrunkVlanSpec into the passed in ResourceData.
func flattenVmwareDistributedVirtualSwitchTrunkVlanSpec(d *schema.ResourceData, obj *types.VmwareDistributedVirtualSwitchTrunkVlanSpec) error {
	var s []interface{}
	for _, rng := range obj.VlanId {
		m := make(map[string]interface{})
		m["min_vlan"] = rng.Start
		m["max_vlan"] = rng.End
		s = append(s, m)
	}
	if err := d.Set("vlan_range", s); err != nil {
		return err
	}
	return nil
}

// expandVmwareDistributedVirtualSwitchPvlanSpec reads certain ResourceData keys and
// returns a VmwareDistributedVirtualSwitchPvlanSpec.
func expandVmwareDistributedVirtualSwitchPvlanSpec(d *schema.ResourceData) *types.VmwareDistributedVirtualSwitchPvlanSpec {
	obj := &types.VmwareDistributedVirtualSwitchPvlanSpec{
		PvlanId: int32(d.Get("port_private_secondary_vlan_id").(int)),
	}
	return obj
}

// flattenVmwareDistributedVirtualSwitchPvlanSpec reads various fields from a
// VmwareDistributedVirtualSwitchPvlanSpec into the passed in ResourceData.
func flattenVmwareDistributedVirtualSwitchPvlanSpec(d *schema.ResourceData, obj *types.VmwareDistributedVirtualSwitchPvlanSpec) error {
	d.Set("port_private_secondary_vlan_id", obj.PvlanId)
	return nil
}

// expandBaseVmwareDistributedVirtualSwitchVlanSpec reads certain ResourceData keys and
// returns a BaseVmwareDistributedVirtualSwitchVlanSpec.
func expandBaseVmwareDistributedVirtualSwitchVlanSpec(d *schema.ResourceData) types.BaseVmwareDistributedVirtualSwitchVlanSpec {
	var obj types.BaseVmwareDistributedVirtualSwitchVlanSpec

	_, ide := d.GetOkExists("vlan_id")
	_, pvid := d.GetOkExists("port_private_secondary_vlan_id")
	vteList, vteOK := d.GetOkExists("vlan_range")
	vte := vteOK && len(vteList.(*schema.Set).List()) > 0
	switch {
	case vte:
		obj = expandVmwareDistributedVirtualSwitchTrunkVlanSpec(d)
	case pvid:
		obj = expandVmwareDistributedVirtualSwitchPvlanSpec(d)
	case ide:
		obj = expandVmwareDistributedVirtualSwitchVlanIDSpec(d)
	}

	return obj
}

// flattenBaseVmwareDistributedVirtualSwitchVlanSpec reads various fields from a
// BaseVmwareDistributedVirtualSwitchVlanSpec into the passed in ResourceData.
func flattenBaseVmwareDistributedVirtualSwitchVlanSpec(d *schema.ResourceData, obj types.BaseVmwareDistributedVirtualSwitchVlanSpec) error {
	if obj == nil {
		return nil
	}

	var err error

	switch t := obj.(type) {
	case *types.VmwareDistributedVirtualSwitchVlanIdSpec:
		err = flattenVmwareDistributedVirtualSwitchVlanIDSpec(d, t)
	case *types.VmwareDistributedVirtualSwitchTrunkVlanSpec:
		err = flattenVmwareDistributedVirtualSwitchTrunkVlanSpec(d, t)
	case *types.VmwareDistributedVirtualSwitchPvlanSpec:
		err = flattenVmwareDistributedVirtualSwitchPvlanSpec(d, t)
	}

	return err
}

// expandDVSFailureCriteria reads certain ResourceData keys and
// returns a DVSFailureCriteria.
func expandDVSFailureCriteria(d *schema.ResourceData) *types.DVSFailureCriteria {
	obj := &types.DVSFailureCriteria{
		CheckBeacon: structure.GetBoolPolicy(d, "check_beacon"),
	}

	if structure.AllFieldsEmpty(obj) {
		return nil
	}
	return obj
}

// flattenDVSFailureCriteria reads various fields from a
// DVSFailureCriteria into the passed in ResourceData.
func flattenDVSFailureCriteria(d *schema.ResourceData, obj *types.DVSFailureCriteria) error {
	if obj == nil {
		return nil
	}

	structure.SetBoolPolicy(d, "check_beacon", obj.CheckBeacon)
	return nil
}

// expandVMwareUplinkPortOrderPolicy reads certain ResourceData keys and
// returns a VMwareUplinkPortOrderPolicy.
func expandVMwareUplinkPortOrderPolicy(d *schema.ResourceData) *types.VMwareUplinkPortOrderPolicy {
	obj := &types.VMwareUplinkPortOrderPolicy{
		ActiveUplinkPort:  structure.SliceInterfacesToStrings(d.Get("active_uplinks").([]interface{})),
		StandbyUplinkPort: structure.SliceInterfacesToStrings(d.Get("standby_uplinks").([]interface{})),
	}

	if structure.AllFieldsEmpty(obj) {
		return nil
	}
	return obj
}

// flattenVMwareUplinkPortOrderPolicy reads various fields from a
// VMwareUplinkPortOrderPolicy into the passed in ResourceData.
func flattenVMwareUplinkPortOrderPolicy(d *schema.ResourceData, obj *types.VMwareUplinkPortOrderPolicy) error {
	if obj == nil {
		return nil
	}

	if err := d.Set("active_uplinks", obj.ActiveUplinkPort); err != nil {
		return err
	}
	if err := d.Set("standby_uplinks", obj.StandbyUplinkPort); err != nil {
		return err
	}
	return nil
}

// expandVmwareUplinkPortTeamingPolicy reads certain ResourceData keys and
// returns a VmwareUplinkPortTeamingPolicy.
func expandVmwareUplinkPortTeamingPolicy(d *schema.ResourceData) *types.VmwareUplinkPortTeamingPolicy {
	obj := &types.VmwareUplinkPortTeamingPolicy{
		Policy:          structure.GetStringPolicy(d, "teaming_policy"),
		NotifySwitches:  structure.GetBoolPolicy(d, "notify_switches"),
		RollingOrder:    structure.GetBoolPolicyReverse(d, "failback"),
		FailureCriteria: expandDVSFailureCriteria(d),
		UplinkPortOrder: expandVMwareUplinkPortOrderPolicy(d),
	}

	if structure.AllFieldsEmpty(obj) {
		return nil
	}
	return obj
}

// flattenVmwareUplinkPortTeamingPolicy reads various fields from a
// VmwareUplinkPortTeamingPolicy into the passed in ResourceData.
func flattenVmwareUplinkPortTeamingPolicy(d *schema.ResourceData, obj *types.VmwareUplinkPortTeamingPolicy) error {
	if obj == nil {
		return nil
	}

	structure.SetStringPolicy(d, "teaming_policy", obj.Policy)
	structure.SetBoolPolicy(d, "notify_switches", obj.NotifySwitches)
	structure.SetBoolPolicyReverse(d, "failback", obj.RollingOrder)

	if err := flattenDVSFailureCriteria(d, obj.FailureCriteria); err != nil {
		return err
	}
	if err := flattenVMwareUplinkPortOrderPolicy(d, obj.UplinkPortOrder); err != nil {
		return err
	}
	return nil
}

// expandDVSSecurityPolicy reads certain ResourceData keys and
// returns a DVSSecurityPolicy.
func expandDVSSecurityPolicy(d *schema.ResourceData) *types.DVSSecurityPolicy {
	obj := &types.DVSSecurityPolicy{
		AllowPromiscuous: structure.GetBoolPolicy(d, "allow_promiscuous"),
		MacChanges:       structure.GetBoolPolicy(d, "allow_mac_changes"),
		ForgedTransmits:  structure.GetBoolPolicy(d, "allow_forged_transmits"),
	}

	if structure.AllFieldsEmpty(obj) {
		return nil
	}
	return obj
}

// flattenDVSSecurityPolicy reads various fields from a
// DVSSecurityPolicy into the passed in ResourceData.
func flattenDVSSecurityPolicy(d *schema.ResourceData, obj *types.DVSSecurityPolicy) error {
	if obj == nil {
		return nil
	}

	structure.SetBoolPolicy(d, "allow_promiscuous", obj.AllowPromiscuous)
	structure.SetBoolPolicy(d, "allow_mac_changes", obj.MacChanges)
	structure.SetBoolPolicy(d, "allow_forged_transmits", obj.ForgedTransmits)
	return nil
}

// expandVMwareUplinkLacpPolicy reads certain ResourceData keys and
// returns a VMwareUplinkLacpPolicy.
func expandVMwareUplinkLacpPolicy(d *schema.ResourceData) *types.VMwareUplinkLacpPolicy {
	obj := &types.VMwareUplinkLacpPolicy{
		Enable: structure.GetBoolPolicy(d, "lacp_enabled"),
		Mode:   structure.GetStringPolicy(d, "lacp_mode"),
	}

	if structure.AllFieldsEmpty(obj) {
		return nil
	}
	return obj
}

// flattenVMwareUplinkLacpPolicy reads various fields from a
// VMwareUplinkLacpPolicy into the passed in ResourceData.
func flattenVMwareUplinkLacpPolicy(d *schema.ResourceData, obj *types.VMwareUplinkLacpPolicy) error {
	if obj == nil {
		return nil
	}

	structure.SetBoolPolicy(d, "lacp_enabled", obj.Enable)
	structure.SetStringPolicy(d, "lacp_mode", obj.Mode)
	return nil
}

// expandDVSTrafficShapingPolicyIngress reads certain ResourceData keys and
// returns a DVSTrafficShapingPolicy for ingress traffic.
func expandDVSTrafficShapingPolicyIngress(d *schema.ResourceData) *types.DVSTrafficShapingPolicy {
	obj := &types.DVSTrafficShapingPolicy{
		Enabled:          structure.GetBoolPolicy(d, "ingress_shaping_enabled"),
		AverageBandwidth: structure.GetLongPolicy(d, "ingress_shaping_average_bandwidth"),
		PeakBandwidth:    structure.GetLongPolicy(d, "ingress_shaping_peak_bandwidth"),
		BurstSize:        structure.GetLongPolicy(d, "ingress_shaping_burst_size"),
	}

	if structure.AllFieldsEmpty(obj) {
		return nil
	}
	return obj
}

// flattenDVSTrafficShapingPolicyIngress reads various fields from the
// DVSTrafficShapingPolicy ingress policy into the passed in ResourceData.
func flattenDVSTrafficShapingPolicyIngress(d *schema.ResourceData, obj *types.DVSTrafficShapingPolicy) error {
	if obj == nil {
		return nil
	}

	structure.SetBoolPolicy(d, "ingress_shaping_enabled", obj.Enabled)
	structure.SetLongPolicy(d, "ingress_shaping_average_bandwidth", obj.AverageBandwidth)
	structure.SetLongPolicy(d, "ingress_shaping_peak_bandwidth", obj.PeakBandwidth)
	structure.SetLongPolicy(d, "ingress_shaping_burst_size", obj.BurstSize)

	return nil
}

// expandDVSTrafficShapingPolicyEgress reads certain ResourceData keys and
// returns a DVSTrafficShapingPolicy for egress traffic.
func expandDVSTrafficShapingPolicyEgress(d *schema.ResourceData) *types.DVSTrafficShapingPolicy {
	obj := &types.DVSTrafficShapingPolicy{
		Enabled:          structure.GetBoolPolicy(d, "egress_shaping_enabled"),
		AverageBandwidth: structure.GetLongPolicy(d, "egress_shaping_average_bandwidth"),
		PeakBandwidth:    structure.GetLongPolicy(d, "egress_shaping_peak_bandwidth"),
		BurstSize:        structure.GetLongPolicy(d, "egress_shaping_burst_size"),
	}

	if structure.AllFieldsEmpty(obj) {
		return nil
	}
	return obj
}

// flattenDVSTrafficShapingPolicyEgress reads various fields from the
// DVSTrafficShapingPolicy egress policy into the passed in ResourceData.
func flattenDVSTrafficShapingPolicyEgress(d *schema.ResourceData, obj *types.DVSTrafficShapingPolicy) error {
	if obj == nil {
		return nil
	}

	structure.SetBoolPolicy(d, "egress_shaping_enabled", obj.Enabled)
	structure.SetLongPolicy(d, "egress_shaping_average_bandwidth", obj.AverageBandwidth)
	structure.SetLongPolicy(d, "egress_shaping_peak_bandwidth", obj.PeakBandwidth)
	structure.SetLongPolicy(d, "egress_shaping_burst_size", obj.BurstSize)
	return nil
}

// expandVMwareDVSPortSetting reads certain ResourceData keys and
// returns a VMwareDVSPortSetting.
func expandVMwareDVSPortSetting(d *schema.ResourceData, resourceType string) *types.VMwareDVSPortSetting {
	var lacpPolicy *types.VMwareUplinkLacpPolicy = nil
	if resourceType == "distributed_virtual_switch" {
		lacpPolicy = expandVMwareUplinkLacpPolicy(d)
	}
	obj := &types.VMwareDVSPortSetting{
		DVPortSetting: types.DVPortSetting{
			Blocked:                 structure.GetBoolPolicy(d, "block_all_ports"),
			InShapingPolicy:         expandDVSTrafficShapingPolicyIngress(d),
			OutShapingPolicy:        expandDVSTrafficShapingPolicyEgress(d),
			VmDirectPathGen2Allowed: structure.GetBoolPolicy(d, "directpath_gen2_allowed"),
		},
		Vlan:                expandBaseVmwareDistributedVirtualSwitchVlanSpec(d),
		UplinkTeamingPolicy: expandVmwareUplinkPortTeamingPolicy(d),
		SecurityPolicy:      expandDVSSecurityPolicy(d),
		IpfixEnabled:        structure.GetBoolPolicy(d, "netflow_enabled"),
		TxUplink:            structure.GetBoolPolicy(d, "tx_uplink"),
		LacpPolicy:          lacpPolicy,
	}

	if structure.AllFieldsEmpty(obj) {
		return nil
	}
	return obj
}

// flattenVMwareDVSPortSetting reads various fields from a
// VMwareDVSPortSetting into the passed in ResourceData.
func flattenVMwareDVSPortSetting(d *schema.ResourceData, obj *types.VMwareDVSPortSetting) error {
	if obj == nil {
		return nil
	}

	structure.SetBoolPolicy(d, "block_all_ports", obj.Blocked)
	structure.SetBoolPolicy(d, "netflow_enabled", obj.IpfixEnabled)
	structure.SetBoolPolicy(d, "tx_uplink", obj.TxUplink)
	structure.SetBoolPolicy(d, "directpath_gen2_allowed", obj.VmDirectPathGen2Allowed)

	if err := flattenDVSTrafficShapingPolicyIngress(d, obj.InShapingPolicy); err != nil {
		return err
	}
	if err := flattenDVSTrafficShapingPolicyEgress(d, obj.OutShapingPolicy); err != nil {
		return err
	}
	if err := flattenBaseVmwareDistributedVirtualSwitchVlanSpec(d, obj.Vlan); err != nil {
		return err
	}
	if err := flattenVmwareUplinkPortTeamingPolicy(d, obj.UplinkTeamingPolicy); err != nil {
		return err
	}
	if err := flattenDVSSecurityPolicy(d, obj.SecurityPolicy); err != nil {
		return err
	}
	if err := flattenVMwareUplinkLacpPolicy(d, obj.LacpPolicy); err != nil {
		return err
	}
	return nil
}
