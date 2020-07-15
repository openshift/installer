package packet

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func resourcePacketPortVlanAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourcePacketPortVlanAttachmentCreate,
		Read:   resourcePacketPortVlanAttachmentRead,
		Delete: resourcePacketPortVlanAttachmentDelete,
		Update: resourcePacketPortVlanAttachmentUpdate,
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
			"native": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourcePacketPortVlanAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	deviceID := d.Get("device_id").(string)
	pName := d.Get("port_name").(string)
	vlanVNID := d.Get("vlan_vnid").(int)

	dev, _, err := client.Devices.Get(deviceID, &packngo.GetOptions{
		Includes: []string{"virtual_networks,project,native_virtual_network"},
	})
	if err != nil {
		return err
	}

	portFound := false
	vlanFound := false
	vlanID := ""
	var port packngo.Port
	for _, p := range dev.NetworkPorts {
		if p.Name == pName {
			portFound = true
			port = p
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
	if !portFound {
		return fmt.Errorf("Device %s doesn't have port %s", deviceID, pName)
	}

	par := &packngo.PortAssignRequest{PortID: port.ID}
	if vlanFound {
		log.Printf("Port %s already has VLAN %d assigned", pName, vlanVNID)
		par.VirtualNetworkID = vlanID
	} else {
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

		par.VirtualNetworkID = vlanID

		// Packet doesn't allow multiple VLANs to be assigned
		// to the same port at the same time
		lockId := "vlan-attachment-" + port.ID
		packetMutexKV.Lock(lockId)
		defer packetMutexKV.Unlock(lockId)

		_, _, err = client.DevicePorts.Assign(par)
		if err != nil {
			return err
		}
	}

	d.SetId(port.ID + ":" + vlanID)

	native := d.Get("native").(bool)
	if native {
		_, _, err = client.DevicePorts.AssignNative(par)
		if err != nil {
			return err
		}
	}

	return resourcePacketPortVlanAttachmentRead(d, meta)
}

func resourcePacketPortVlanAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	deviceID := d.Get("device_id").(string)
	pName := d.Get("port_name").(string)
	vlanVNID := d.Get("vlan_vnid").(int)

	dev, _, err := client.Devices.Get(deviceID, &packngo.GetOptions{Includes: []string{"virtual_networks,project,native_virtual_network"}})
	if err != nil {
		return err
	}
	portFound := false
	vlanFound := false
	portID := ""
	vlanID := ""
	vlanNative := false
	for _, p := range dev.NetworkPorts {
		if p.Name == pName {
			portFound = true
			portID = p.ID
			for _, n := range p.AttachedVirtualNetworks {
				if vlanVNID == n.VXLAN {
					vlanFound = true
					vlanID = n.ID
					if p.NativeVirtualNetwork != nil {
						vlanNative = vlanID == p.NativeVirtualNetwork.ID
					}
					break
				}
			}
			break
		}
	}
	if !portFound {
		return fmt.Errorf("Device %s doesn't have port %s", deviceID, pName)
	}
	if !vlanFound {
		d.SetId("")
	}
	d.Set("port_id", portID)
	d.Set("vlan_id", vlanID)
	d.Set("native", vlanNative)
	return nil
}

func resourcePacketPortVlanAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	if d.HasChange("native") {
		native := d.Get("native").(bool)
		portID := d.Get("port_id").(string)
		if native {
			vlanID := d.Get("vlan_id").(string)
			par := &packngo.PortAssignRequest{PortID: portID, VirtualNetworkID: vlanID}
			_, _, err := client.DevicePorts.AssignNative(par)
			if err != nil {
				return err
			}
		} else {
			_, _, err := client.DevicePorts.UnassignNative(portID)
			if err != nil {
				return err
			}
		}
	}
	return resourcePacketPortVlanAttachmentRead(d, meta)
}

func resourcePacketPortVlanAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	pID := d.Get("port_id").(string)
	vlanID := d.Get("vlan_id").(string)
	native := d.Get("native").(bool)
	if native {
		_, _, err := client.DevicePorts.UnassignNative(pID)
		if err != nil {
			return err
		}
	}
	par := &packngo.PortAssignRequest{PortID: pID, VirtualNetworkID: vlanID}
	lockId := "vlan-detachment-" + pID
	packetMutexKV.Lock(lockId)
	defer packetMutexKV.Unlock(lockId)
	portPtr, _, err := client.DevicePorts.Unassign(par)
	if err != nil {
		return err
	}
	forceBond := d.Get("force_bond").(bool)
	if forceBond && (len(portPtr.AttachedVirtualNetworks) == 0) {
		deviceID := d.Get("device_id").(string)
		portName := d.Get("port_name").(string)
		port, err := client.DevicePorts.GetPortByName(deviceID, portName)
		if err != nil {
			return friendlyError(err)
		}
		_, _, err = client.DevicePorts.Bond(port, false)
		if err != nil {
			return friendlyError(err)
		}
	}
	return nil
}
