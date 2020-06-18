package packet

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketPortVlanAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketPortVlanAttachmentCreate,
		Read:   resourcePacketPortVlanAttachmentRead,
		Delete: resourcePacketPortVlanAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"force_bond": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"port_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vlan_vnid": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"vlan_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePacketPortVlanAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	deviceID := d.Get("device_id").(string)
	pName := d.Get("port_name").(string)
	vlanVNID := d.Get("vlan_vnid").(int)

	dev, _, err := client.Devices.Get(deviceID, &packngo.GetOptions{Includes: []string{"virtual_networks,project"}})
	if err != nil {
		return err
	}

	portFound := false
	vlanFound := false
	var port packngo.Port
	for _, p := range dev.NetworkPorts {
		if p.Name == pName {
			portFound = true
			port = p
			for _, n := range p.AttachedVirtualNetworks {
				if vlanVNID == n.VXLAN {
					vlanFound = true
					break
				}
			}
			break
		}
	}
	if !portFound {
		return fmt.Errorf("Device %s doesn't have port %s", deviceID, pName)
	}
	if vlanFound {
		log.Printf("Port %s already has VLAN %d assigned", pName, vlanVNID)
		return nil
	}

	vlanID := ""
	facility := dev.Facility.Code
	vlans, _, err := client.ProjectVirtualNetworks.List(dev.Project.ID, nil)
	if err != nil {
		return err
	}
	for _, n := range vlans.VirtualNetworks {
		if (n.VXLAN == vlanVNID) && (n.FacilityCode == facility) {
			vlanID = n.ID
		}
	}
	if len(vlanID) == 0 {
		return fmt.Errorf("VLAN with VNID %d doesn't exist in facilty %s", vlanVNID, facility)
	}

	par := &packngo.PortAssignRequest{PortID: port.ID, VirtualNetworkID: vlanID}

	_, _, err = client.DevicePorts.Assign(par)
	if err != nil {
		return err
	}

	d.SetId(port.ID + ":" + vlanID)
	return resourcePacketPortVlanAttachmentRead(d, meta)
}

func resourcePacketPortVlanAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	deviceID := d.Get("device_id").(string)
	pName := d.Get("port_name").(string)
	vlanVNID := d.Get("vlan_vnid").(int)

	dev, _, err := client.Devices.Get(deviceID, &packngo.GetOptions{Includes: []string{"virtual_networks,project"}})
	if err != nil {
		return err
	}
	portFound := false
	vlanFound := false
	portID := ""
	vlanID := ""
	for _, p := range dev.NetworkPorts {
		if p.Name == pName {
			portFound = true
			portID = p.ID
			for _, n := range p.AttachedVirtualNetworks {
				if vlanVNID == n.VXLAN {
					vlanFound = true
					vlanID = n.ID
					break
				}
			}
			break
		}
	}
	d.Set("port_id", portID)
	d.Set("vlan_id", vlanID)
	if !portFound {
		return fmt.Errorf("Device %s doesn't have port %s", deviceID, pName)
	}
	if !vlanFound {
		d.SetId(portID)
	}
	return nil
}

func resourcePacketPortVlanAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	pID := d.Get("port_id").(string)
	vlanID := d.Get("vlan_id").(string)
	par := &packngo.PortAssignRequest{PortID: pID, VirtualNetworkID: vlanID}
	client := meta.(*packngo.Client)
	portPtr, _, err := client.DevicePorts.Unassign(par)
	if err != nil {
		return err
	}
	forceBond := d.Get("force_bond").(bool)
	if forceBond && (len(portPtr.AttachedVirtualNetworks) == 0) {
		_, _, err = client.DevicePorts.Bond(&packngo.BondRequest{PortID: pID, BulkEnable: false})
		if err != nil {
			return friendlyError(err)
		}
	}
	return nil
}
