package vsphere

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/structure"
	"github.com/vmware/govmomi/vim25/types"
)

const hostVirtualSwitchIDPrefix = "tf-HostVirtualSwitch"

var linkDiscoveryProtocolConfigOperationAllowedValues = []string{
	string(types.LinkDiscoveryProtocolConfigOperationTypeNone),
	string(types.LinkDiscoveryProtocolConfigOperationTypeListen),
	string(types.LinkDiscoveryProtocolConfigOperationTypeAdvertise),
	string(types.LinkDiscoveryProtocolConfigOperationTypeBoth),
}

var linkDiscoveryProtocolConfigProtocolAllowedValues = []string{
	string(types.LinkDiscoveryProtocolConfigProtocolTypeCdp),
	string(types.LinkDiscoveryProtocolConfigProtocolTypeLldp),
}

// schemaHostVirtualSwitchBondBridge returns schema items for resources that
// need to work with a HostVirtualSwitchBondBridge, such as virtual switches.
func schemaHostVirtualSwitchBondBridge() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// HostVirtualSwitchBeaconConfig
		"beacon_interval": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "Determines how often, in seconds, a beacon should be sent to probe for the validity of a link.",
			Default:      1,
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

		// HostVirtualSwitchBondBridge
		"network_adapters": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "The list of network adapters to bind to this virtual switch.",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

// expandHostVirtualSwitchBeaconConfig reads certain ResourceData keys and
// returns a HostVirtualSwitchBeaconConfig.
func expandHostVirtualSwitchBeaconConfig(d *schema.ResourceData) *types.HostVirtualSwitchBeaconConfig {
	obj := &types.HostVirtualSwitchBeaconConfig{
		Interval: int32(d.Get("beacon_interval").(int)),
	}
	return obj
}

// flattenHostVirtualSwitchBeaconConfig reads various fields from a
// HostVirtualSwitchBeaconConfig into the passed in ResourceData.
func flattenHostVirtualSwitchBeaconConfig(d *schema.ResourceData, obj *types.HostVirtualSwitchBeaconConfig) error {
	d.Set("beacon_interval", obj.Interval)
	return nil
}

// expandLinkDiscoveryProtocolConfig reads certain ResourceData keys and
// returns a LinkDiscoveryProtocolConfig.
func expandLinkDiscoveryProtocolConfig(d *schema.ResourceData) *types.LinkDiscoveryProtocolConfig {
	obj := &types.LinkDiscoveryProtocolConfig{
		Operation: d.Get("link_discovery_operation").(string),
		Protocol:  d.Get("link_discovery_protocol").(string),
	}
	return obj
}

// flattenLinkDiscoveryProtocolConfig reads various fields from a
// LinkDiscoveryProtocolConfig into the passed in ResourceData.
func flattenLinkDiscoveryProtocolConfig(d *schema.ResourceData, obj *types.LinkDiscoveryProtocolConfig) error {
	d.Set("link_discovery_operation", obj.Operation)
	d.Set("link_discovery_protocol", obj.Protocol)
	return nil
}

// expandHostVirtualSwitchBondBridge reads certain ResourceData keys and
// returns a HostVirtualSwitchBondBridge.
func expandHostVirtualSwitchBondBridge(d *schema.ResourceData) *types.HostVirtualSwitchBondBridge {
	obj := &types.HostVirtualSwitchBondBridge{
		NicDevice: structure.SliceInterfacesToStrings(d.Get("network_adapters").([]interface{})),
	}
	obj.Beacon = expandHostVirtualSwitchBeaconConfig(d)
	obj.LinkDiscoveryProtocolConfig = expandLinkDiscoveryProtocolConfig(d)
	return obj
}

// flattenHostVirtualSwitchBondBridge reads various fields from a
// HostVirtualSwitchBondBridge into the passed in ResourceData.
func flattenHostVirtualSwitchBondBridge(d *schema.ResourceData, obj *types.HostVirtualSwitchBondBridge) error {
	if err := d.Set("network_adapters", structure.SliceStringsToInterfaces(obj.NicDevice)); err != nil {
		return err
	}
	if err := flattenHostVirtualSwitchBeaconConfig(d, obj.Beacon); err != nil {
		return err
	}
	if err := flattenLinkDiscoveryProtocolConfig(d, obj.LinkDiscoveryProtocolConfig); err != nil {
		return err
	}
	return nil
}

// schemaHostVirtualSwitchSpec returns schema items for resources that need to
// work with a HostVirtualSwitchSpec, such as virtual switches.
func schemaHostVirtualSwitchSpec() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		// HostVirtualSwitchSpec
		"mtu": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "The maximum transmission unit (MTU) of the virtual switch in bytes.",
			Default:      1500,
			ValidateFunc: validation.IntBetween(1, 9000),
		},
		"number_of_ports": {
			Type:         schema.TypeInt,
			Optional:     true,
			Description:  "The number of ports that this virtual switch is configured to use.",
			Default:      128,
			ValidateFunc: validation.IntBetween(0, 1024),
		},
	}
	structure.MergeSchema(s, schemaHostVirtualSwitchBondBridge())
	structure.MergeSchema(s, schemaHostNetworkPolicy())
	return s
}

// expandHostVirtualSwitchSpec reads certain ResourceData keys and returns a
// HostVirtualSwitchSpec.
func expandHostVirtualSwitchSpec(d *schema.ResourceData) *types.HostVirtualSwitchSpec {
	obj := &types.HostVirtualSwitchSpec{
		Mtu:      int32(d.Get("mtu").(int)),
		NumPorts: int32(d.Get("number_of_ports").(int)),
		Bridge:   expandHostVirtualSwitchBondBridge(d),
		Policy:   expandHostNetworkPolicy(d),
	}
	// If there are no NICs, we need to nil-out Bridge, as a
	// HostVirtualSwitchBondBridge with no NICs is an invalid config.
	if len(obj.Bridge.(*types.HostVirtualSwitchBondBridge).NicDevice) < 1 {
		obj.Bridge = nil
	}
	return obj
}

// flattenHostVirtualSwitchSpec reads various fields from a
// HostVirtualSwitchSpec into the passed in ResourceData.
func flattenHostVirtualSwitchSpec(d *schema.ResourceData, obj *types.HostVirtualSwitchSpec) error {
	d.Set("mtu", obj.Mtu)
	d.Set("number_of_ports", obj.NumPorts)
	if obj.Bridge != nil {
		if err := flattenHostVirtualSwitchBondBridge(d, obj.Bridge.(*types.HostVirtualSwitchBondBridge)); err != nil {
			return err
		}
	}
	if err := flattenHostNetworkPolicy(d, obj.Policy); err != nil {
		return err
	}
	return nil
}

// saveHostVirtualSwitchID sets a special ID for a host virtual switch,
// composed of the MOID for the concerned HostSystem and the virtual switch's
// key.
func saveHostVirtualSwitchID(d *schema.ResourceData, hsID, name string) {
	d.SetId(fmt.Sprintf("%s:%s:%s", hostVirtualSwitchIDPrefix, hsID, name))
}

// splitHostVirtualSwitchID splits a vsphere_host_virtual_switch resource ID
// into its counterparts: the prefix, the HostSystem ID, and the virtual switch
// name.
func splitHostVirtualSwitchID(raw string) (string, string, error) {
	s := strings.SplitN(raw, ":", 3)
	if len(s) != 3 || s[0] != hostVirtualSwitchIDPrefix || s[1] == "" || s[2] == "" {
		return "", "", fmt.Errorf("corrupt ID: %s", raw)
	}
	return s[1], s[2], nil
}

// virtualSwitchIDsFromResourceID passes a resource's ID through
// splitHostVirtualSwitchID.
func virtualSwitchIDsFromResourceID(d *schema.ResourceData) (string, string, error) {
	return splitHostVirtualSwitchID(d.Id())
}
