package ironic

import (
	"github.com/gophercloud/gophercloud/openstack/baremetal/v1/ports"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourcePortV1() *schema.Resource {
	return &schema.Resource{
		Create: resourcePortV1Create,
		Read:   resourcePortV1Read,
		Update: resourcePortV1Update,
		Delete: resourcePortV1Delete,

		Schema: map[string]*schema.Schema{
			"node_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"port_group_uuid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"local_link_connection": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"pxe_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"physical_network": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"extra": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"is_smart_nic": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourcePortV1Create(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	opts := portSchemaToCreateOpts(d)
	result, err := ports.Create(client, opts).Extract()
	if err != nil {
		return err
	}
	d.SetId(result.UUID)

	return resourcePortV1Read(d, meta)
}

func resourcePortV1Read(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Clients).GetIronicClient()
	if err != nil {
		return err
	}

	port, err := ports.Get(client, d.Id()).Extract()
	if err != nil {
		return err
	}

	d.Set("address", port.Address)
	d.Set("node_uuid", port.NodeUUID)
	d.Set("port_group_id", port.PortGroupUUID)
	d.Set("local_link_collection", port.LocalLinkConnection)
	d.Set("pxe_enabled", port.PXEEnabled)
	d.Set("physical_network", port.PhysicalNetwork)
	d.Set("extra", port.Extra)
	d.Set("is_smart_nic", port.IsSmartNIC)

	return nil
}

func resourcePortV1Update(d *schema.ResourceData, meta interface{}) error {
	return nil

}

func resourcePortV1Delete(d *schema.ResourceData, meta interface{}) error {
	return nil

}

func portSchemaToCreateOpts(d *schema.ResourceData) *ports.CreateOpts {
	pxeEnabled := d.Get("pxe_enabled").(bool)
	isSmartNic := d.Get("is_smart_nic").(bool)

	opts := ports.CreateOpts{
		NodeUUID:      d.Get("node_uuid").(string),
		Address:       d.Get("address").(string),
		PortGroupUUID: d.Get("port_group_uuid").(string),
		//LocalLinkConnection: d.Get("local_link_collection").(map[string]interface{}),
		PXEEnabled:      &pxeEnabled,
		PhysicalNetwork: d.Get("physical_network").(string),
		//Extra:               d.Get("extra").(map[string]interface{}),
		IsSmartNIC: &isSmartNic,
	}

	return &opts
}
