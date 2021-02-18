package vsphere

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi/vim25/types"
)

const (
	hostNetworkPolicyNicTeamingPolicyModeLoadbalanceIP     = "loadbalance_ip"
	hostNetworkPolicyNicTeamingPolicyModeLoadbalanceSrcMac = "loadbalance_srcmac"
	hostNetworkPolicyNicTeamingPolicyModeLoadbalanceSrcID  = "loadbalance_srcid"
	hostNetworkPolicyNicTeamingPolicyModeFailoverExplicit  = "failover_explicit"
)

var hostNetworkPolicyNicTeamingPolicyAllowedValues = []string{
	hostNetworkPolicyNicTeamingPolicyModeLoadbalanceIP,
	hostNetworkPolicyNicTeamingPolicyModeLoadbalanceSrcMac,
	hostNetworkPolicyNicTeamingPolicyModeLoadbalanceSrcID,
	hostNetworkPolicyNicTeamingPolicyModeFailoverExplicit,
}

// schemaHostNetworkPolicy returns schema items for resources that need to work
// with a HostNetworkPolicy, such as virtual switches and port groups.
func schemaHostNetworkPolicy() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// HostNicTeamingPolicy/HostNicFailureCriteria
		"check_beacon": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable beacon probing. Requires that the vSwitch has been configured to use a beacon. If disabled, link status is used only.",
		},

		// HostNicTeamingPolicy/HostNicOrderPolicy
		"active_nics": {
			Type:        schema.TypeList,
			Description: "List of active network adapters used for load balancing.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"standby_nics": {
			Type:        schema.TypeList,
			Description: "List of standby network adapters used for failover.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},

		// HostNicTeamingPolicy
		"teaming_policy": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "The network adapter teaming policy. Can be one of loadbalance_ip, loadbalance_srcmac, loadbalance_srcid, or failover_explicit.",
			ValidateFunc: validation.StringInSlice(hostNetworkPolicyNicTeamingPolicyAllowedValues, false),
		},
		"notify_switches": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If true, the teaming policy will notify the broadcast network of a NIC failover, triggering cache updates.",
		},
		"failback": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "If true, the teaming policy will re-activate failed interfaces higher in precedence when they come back up.",
		},

		// HostNetworkSecurityPolicy
		"allow_promiscuous": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable promiscuous mode on the network. This flag indicates whether or not all traffic is seen on a given port.",
		},
		"allow_forged_transmits": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Controls whether or not the virtual network adapter is allowed to send network traffic with a different MAC address than that of its own.",
		},
		"allow_mac_changes": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Controls whether or not the Media Access Control (MAC) address can be changed.",
		},

		// HostNetworkTrafficShapingPolicy
		"shaping_average_bandwidth": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The average bandwidth in bits per second if traffic shaping is enabled.",
		},
		"shaping_burst_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The maximum burst size allowed in bytes if traffic shaping is enabled.",
		},
		"shaping_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Enable traffic shaping on this virtual switch or port group.",
		},
		"shaping_peak_bandwidth": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "The peak bandwidth during bursts in bits per second if traffic shaping is enabled.",
		},
	}
}

// expandHostNicFailureCriteria reads certain ResourceData keys and returns a
// HostNicFailureCriteria.
func expandHostNicFailureCriteria(d *schema.ResourceData) *types.HostNicFailureCriteria {
	obj := &types.HostNicFailureCriteria{}

	if v, ok := d.GetOkExists("check_beacon"); ok {
		obj.CheckBeacon = structure.BoolPtr(v.(bool))
	}

	// These fields are deprecated and are set only to make things work. They are
	// not exposed to Terraform.
	obj.CheckSpeed = "minimum"
	obj.Speed = 10
	obj.CheckDuplex = structure.BoolPtr(false)
	obj.FullDuplex = structure.BoolPtr(false)
	obj.CheckErrorPercent = structure.BoolPtr(false)
	obj.Percentage = 0

	return obj
}

// flattenHostNicFailureCriteria reads various fields from a
// HostNicFailureCriteria into the passed in ResourceData.
func flattenHostNicFailureCriteria(d *schema.ResourceData, obj *types.HostNicFailureCriteria) error {
	if obj.CheckBeacon != nil {
		d.Set("check_beacon", obj.CheckBeacon)
	}
	return nil
}

// expandHostNicOrderPolicy reads certain ResourceData keys and returns a
// HostNicOrderPolicy.
func expandHostNicOrderPolicy(d *schema.ResourceData) *types.HostNicOrderPolicy {
	obj := &types.HostNicOrderPolicy{}
	activeNics, activeOk := d.GetOkExists("active_nics")
	standbyNics, standbyOk := d.GetOkExists("standby_nics")
	if !activeOk && !standbyOk {
		return nil
	}
	obj.ActiveNic = structure.SliceInterfacesToStrings(activeNics.([]interface{}))
	obj.StandbyNic = structure.SliceInterfacesToStrings(standbyNics.([]interface{}))
	return obj
}

// flattenHostNicOrderPolicy reads various fields from a HostNicOrderPolicy
// into the passed in ResourceData.
func flattenHostNicOrderPolicy(d *schema.ResourceData, obj *types.HostNicOrderPolicy) error {
	if obj == nil {
		return nil
	}
	if err := d.Set("active_nics", structure.SliceStringsToInterfaces(obj.ActiveNic)); err != nil {
		return err
	}
	if err := d.Set("standby_nics", structure.SliceStringsToInterfaces(obj.StandbyNic)); err != nil {
		return err
	}
	return nil
}

// expandHostNicTeamingPolicy reads certain ResourceData keys and returns a
// HostNicTeamingPolicy.
func expandHostNicTeamingPolicy(d *schema.ResourceData) *types.HostNicTeamingPolicy {
	obj := &types.HostNicTeamingPolicy{
		Policy: d.Get("teaming_policy").(string),
	}
	if v, ok := d.GetOkExists("failback"); ok {
		obj.RollingOrder = structure.BoolPtr(!v.(bool))
	}
	if v, ok := d.GetOkExists("notify_switches"); ok {
		obj.NotifySwitches = structure.BoolPtr(v.(bool))
	}
	obj.FailureCriteria = expandHostNicFailureCriteria(d)
	obj.NicOrder = expandHostNicOrderPolicy(d)

	// These fields are deprecated and are set only to make things work. They are
	// not exposed to Terraform.
	obj.ReversePolicy = structure.BoolPtr(true)

	return obj
}

// flattenHostNicTeamingPolicy reads various fields from a HostNicTeamingPolicy
// into the passed in ResourceData.
func flattenHostNicTeamingPolicy(d *schema.ResourceData, obj *types.HostNicTeamingPolicy) error {
	if obj.RollingOrder != nil {
		v := *obj.RollingOrder
		d.Set("failback", !v)
	}
	if obj.NotifySwitches != nil {
		d.Set("notify_switches", obj.NotifySwitches)
	}
	d.Set("teaming_policy", obj.Policy)
	if obj.FailureCriteria != nil {
		if err := flattenHostNicFailureCriteria(d, obj.FailureCriteria); err != nil {
			return err
		}
	}
	if err := flattenHostNicOrderPolicy(d, obj.NicOrder); err != nil {
		return err
	}
	return nil
}

// expandHostNetworkSecurityPolicy reads certain ResourceData keys and returns
// a HostNetworkSecurityPolicy.
func expandHostNetworkSecurityPolicy(d *schema.ResourceData) *types.HostNetworkSecurityPolicy {
	obj := &types.HostNetworkSecurityPolicy{}
	if v, ok := d.GetOkExists("allow_promiscuous"); ok {
		obj.AllowPromiscuous = structure.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOkExists("allow_forged_transmits"); ok {
		obj.ForgedTransmits = structure.BoolPtr(v.(bool))
	}
	if v, ok := d.GetOkExists("allow_mac_changes"); ok {
		obj.MacChanges = structure.BoolPtr(v.(bool))
	}
	return obj
}

// flattenHostNetworkSecurityPolicy reads various fields from a
// HostNetworkSecurityPolicy into the passed in ResourceData.
func flattenHostNetworkSecurityPolicy(d *schema.ResourceData, obj *types.HostNetworkSecurityPolicy) error {
	if obj.AllowPromiscuous != nil {
		d.Set("allow_promiscuous", obj.AllowPromiscuous)
	}
	if obj.ForgedTransmits != nil {
		d.Set("allow_forged_transmits", obj.ForgedTransmits)
	}
	if obj.MacChanges != nil {
		d.Set("allow_mac_changes", obj.MacChanges)
	}
	return nil
}

// expandHostNetworkTrafficShapingPolicy reads certain ResourceData keys and
// returns a HostNetworkTrafficShapingPolicy.
func expandHostNetworkTrafficShapingPolicy(d *schema.ResourceData) *types.HostNetworkTrafficShapingPolicy {
	obj := &types.HostNetworkTrafficShapingPolicy{
		AverageBandwidth: int64(d.Get("shaping_average_bandwidth").(int)),
		BurstSize:        int64(d.Get("shaping_burst_size").(int)),
		PeakBandwidth:    int64(d.Get("shaping_peak_bandwidth").(int)),
	}
	if v, ok := d.GetOkExists("shaping_enabled"); ok {
		obj.Enabled = structure.BoolPtr(v.(bool))
	}
	return obj
}

// flattenHostNetworkTrafficShapingPolicy reads various fields from a
// HostNetworkTrafficShapingPolicy into the passed in ResourceData.
func flattenHostNetworkTrafficShapingPolicy(d *schema.ResourceData, obj *types.HostNetworkTrafficShapingPolicy) error {
	if obj.Enabled != nil {
		d.Set("shaping_enabled", obj.Enabled)
	}
	d.Set("shaping_average_bandwidth", obj.AverageBandwidth)
	d.Set("shaping_burst_size", obj.BurstSize)
	d.Set("shaping_peak_bandwidth", obj.PeakBandwidth)
	return nil
}

// expandHostNetworkPolicy reads certain ResourceData keys and returns a
// HostNetworkPolicy.
func expandHostNetworkPolicy(d *schema.ResourceData) *types.HostNetworkPolicy {
	obj := &types.HostNetworkPolicy{
		Security:      expandHostNetworkSecurityPolicy(d),
		NicTeaming:    expandHostNicTeamingPolicy(d),
		ShapingPolicy: expandHostNetworkTrafficShapingPolicy(d),
	}
	return obj
}

// flattenHostNetworkPolicy reads various fields from a HostNetworkPolicy into
// the passed in ResourceData.
func flattenHostNetworkPolicy(d *schema.ResourceData, obj *types.HostNetworkPolicy) error {
	if err := flattenHostNetworkSecurityPolicy(d, obj.Security); obj.Security != nil && err != nil {
		return err
	}
	if err := flattenHostNicTeamingPolicy(d, obj.NicTeaming); obj.NicTeaming != nil && err != nil {
		return err
	}
	if err := flattenHostNetworkTrafficShapingPolicy(d, obj.ShapingPolicy); obj.ShapingPolicy != nil && err != nil {
		return err
	}
	return nil
}
