package vsphere

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi/vim25/types"
)

var distributedVirtualPortgroupPortgroupTypeAllowedValues = []string{
	string(types.DistributedVirtualPortgroupPortgroupTypeEarlyBinding),
	string(types.DistributedVirtualPortgroupPortgroupTypeEphemeral),
}

// schemaDVPortgroupConfigSpec returns schema items for resources that
// need to work with a DVPortgroupConfigSpec.
func schemaDVPortgroupConfigSpec() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		// VMwareDVSPortgroupPolicy
		"block_override_allowed": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow the blocked setting of an individual port to override the setting in the portgroup.",
		},
		"live_port_moving_allowed": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow a live port to be moved in and out of the portgroup.",
		},
		"network_resource_pool_override_allowed": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow the network resource pool of an individual port to override the setting in the portgroup.",
		},
		"port_config_reset_at_disconnect": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Reset the setting of any ports in this portgroup back to the default setting when the port disconnects.",
		},
		"shaping_override_allowed": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow the traffic shaping policies of an individual port to override the settings in the portgroup.",
		},
		"traffic_filter_override_allowed": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow any filter policies set on the individual port to override those in the portgroup.",
		},
		"netflow_override_allowed": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow the enabling or disabling of Netflow on a port, contrary to the policy in the portgroup.",
		},
		"security_policy_override_allowed": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow security policy settings on a port to override those on the portgroup.",
		},
		"uplink_teaming_override_allowed": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow the uplink teaming policies on a port to override those on the portgroup.",
		},
		"vlan_override_allowed": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Allow the VLAN configuration on a port to override those on the portgroup.",
		},

		// DVPortgroupConfigSpec
		"auto_expand": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Auto-expands the port group beyond the port count configured in number_of_ports when necessary.",
		},
		"config_version": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Version string of the configuration that this spec is trying to change.",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The description of the portgroup.",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the portgroup.",
		},
		"number_of_ports": {
			Type:         schema.TypeInt,
			Optional:     true,
			Computed:     true,
			Description:  "The number of ports in this portgroup. The DVS will expand and shrink by modifying this setting.",
			ValidateFunc: validation.IntAtLeast(0),
		},
		"port_name_format": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A template string to use when creating ports in the portgroup.",
		},
		"type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      string(types.DistributedVirtualPortgroupPortgroupTypeEarlyBinding),
			Description:  "The type of portgroup. Can be one of earlyBinding (static) or ephemeral.",
			ValidateFunc: validation.StringInSlice(distributedVirtualPortgroupPortgroupTypeAllowedValues, false),
		},
		"network_resource_pool_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "-1",
			Description: "The key of a network resource pool to associate with this portgroup.",
		},
	}
	structure.MergeSchema(s, schemaVMwareDVSPortSetting())
	return s
}

// expandVMwareDVSPortgroupPolicy reads certain ResourceData keys and
// returns a VMwareDVSPortgroupPolicy.
func expandVMwareDVSPortgroupPolicy(d *schema.ResourceData) *types.VMwareDVSPortgroupPolicy {
	obj := &types.VMwareDVSPortgroupPolicy{
		DVPortgroupPolicy: types.DVPortgroupPolicy{
			BlockOverrideAllowed:               d.Get("block_override_allowed").(bool),
			ShapingOverrideAllowed:             d.Get("shaping_override_allowed").(bool),
			LivePortMovingAllowed:              d.Get("live_port_moving_allowed").(bool),
			PortConfigResetAtDisconnect:        d.Get("port_config_reset_at_disconnect").(bool),
			NetworkResourcePoolOverrideAllowed: structure.GetBoolPtr(d, "network_resource_pool_override_allowed"),
			TrafficFilterOverrideAllowed:       structure.GetBoolPtr(d, "traffic_filter_override_allowed"),
		},
		VlanOverrideAllowed:           d.Get("vlan_override_allowed").(bool),
		UplinkTeamingOverrideAllowed:  d.Get("uplink_teaming_override_allowed").(bool),
		SecurityPolicyOverrideAllowed: d.Get("security_policy_override_allowed").(bool),
		IpfixOverrideAllowed:          structure.GetBoolPtr(d, "netflow_override_allowed"),
	}
	return obj
}

// flattenVMwareDVSPortgroupPolicy reads various fields from a
// VMwareDVSPortgroupPolicy into the passed in ResourceData.
func flattenVMwareDVSPortgroupPolicy(d *schema.ResourceData, obj *types.VMwareDVSPortgroupPolicy) error {
	d.Set("block_override_allowed", obj.BlockOverrideAllowed)
	d.Set("shaping_override_allowed", obj.ShapingOverrideAllowed)
	d.Set("live_port_moving_allowed", obj.LivePortMovingAllowed)
	d.Set("port_config_reset_at_disconnect", obj.PortConfigResetAtDisconnect)
	d.Set("vlan_override_allowed", obj.VlanOverrideAllowed)
	d.Set("uplink_teaming_override_allowed", obj.UplinkTeamingOverrideAllowed)
	d.Set("security_policy_override_allowed", obj.SecurityPolicyOverrideAllowed)

	structure.SetBoolPtr(d, "network_resource_pool_override_allowed", obj.NetworkResourcePoolOverrideAllowed)
	structure.SetBoolPtr(d, "traffic_filter_override_allowed", obj.TrafficFilterOverrideAllowed)
	structure.SetBoolPtr(d, "netflow_override_allowed", obj.IpfixOverrideAllowed)
	return nil
}

// expandDVPortgroupConfigSpec reads certain ResourceData keys and
// returns a DVPortgroupConfigSpec.
func expandDVPortgroupConfigSpec(d *schema.ResourceData) types.DVPortgroupConfigSpec {
	obj := types.DVPortgroupConfigSpec{
		ConfigVersion:                d.Get("config_version").(string),
		Name:                         d.Get("name").(string),
		NumPorts:                     int32(d.Get("number_of_ports").(int)),
		PortNameFormat:               d.Get("port_name_format").(string),
		DefaultPortConfig:            expandVMwareDVSPortSetting(d, "distributed_port_group"),
		Description:                  d.Get("description").(string),
		Type:                         d.Get("type").(string),
		Policy:                       expandVMwareDVSPortgroupPolicy(d),
		AutoExpand:                   structure.GetBoolPtr(d, "auto_expand"),
		VmVnicNetworkResourcePoolKey: d.Get("network_resource_pool_key").(string),
	}
	return obj
}

// flattenDVPortgroupConfigInfo reads various fields from a
// DVPortgroupConfigInfo into the passed in ResourceData.
//
// This is the flatten counterpart to expandDVPortgroupConfigSpec.
func flattenDVPortgroupConfigInfo(d *schema.ResourceData, obj types.DVPortgroupConfigInfo) error {
	d.Set("config_version", obj.ConfigVersion)
	d.Set("name", obj.Name)
	d.Set("number_of_ports", obj.NumPorts)
	d.Set("port_name_format", obj.PortNameFormat)
	d.Set("description", obj.Description)
	d.Set("type", obj.Type)
	structure.SetBoolPtr(d, "auto_expand", obj.AutoExpand)
	d.Set("network_resource_pool_key", obj.VmVnicNetworkResourcePoolKey)

	if err := flattenVMwareDVSPortSetting(d, obj.DefaultPortConfig.(*types.VMwareDVSPortSetting)); err != nil {
		return err
	}
	if err := flattenVMwareDVSPortgroupPolicy(d, obj.Policy.(*types.VMwareDVSPortgroupPolicy)); err != nil {
		return err
	}
	return nil
}
