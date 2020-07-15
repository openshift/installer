package packet

import (
	"fmt"
	"path"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/packethost/packngo"
)

func packetIPComputedFields() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"address_family": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"cidr": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"gateway": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"netmask": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"network": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"manageable": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"management": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}
}

func packetIPResourceComputedFields() map[string]*schema.Schema {
	s := packetIPComputedFields()
	s["address_family"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}
	s["public"] = &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}
	s["global"] = &schema.Schema{
		Type:     schema.TypeBool,
		Computed: true,
	}
	return s
}

func resourcePacketReservedIPBlock() *schema.Resource {
	reservedBlockSchema := packetIPResourceComputedFields()
	reservedBlockSchema["project_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
		ForceNew: true,
	}
	reservedBlockSchema["facility"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	}
	reservedBlockSchema["description"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		ForceNew: true,
	}
	reservedBlockSchema["quantity"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
		ForceNew: true,
	}
	reservedBlockSchema["type"] = &schema.Schema{
		Type:         schema.TypeString,
		ForceNew:     true,
		Default:      "public_ipv4",
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"public_ipv4", "global_ipv4"}, false),
	}
	reservedBlockSchema["cidr_notation"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}

	return &schema.Resource{
		Create: resourcePacketReservedIPBlockCreate,
		Read:   resourcePacketReservedIPBlockRead,
		Delete: resourcePacketReservedIPBlockDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: reservedBlockSchema,
	}
}

func resourcePacketReservedIPBlockCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	quantity := d.Get("quantity").(int)
	typ := d.Get("type").(string)

	req := packngo.IPReservationRequest{
		Type:     typ,
		Quantity: quantity,
	}
	f, ok := d.GetOk("facility")

	if ok && typ == "global_ipv4" {
		return fmt.Errorf("Facility can not be set for type == global_ipv4")
	}
	fs := f.(string)
	if typ == "public_ipv4" {
		req.Facility = &fs
	}
	desc, ok := d.GetOk("description")
	if ok {
		req.Description = desc.(string)
	}

	projectID := d.Get("project_id").(string)

	blockAddr, _, err := client.ProjectIPs.Request(projectID, &req)
	if err != nil {
		return fmt.Errorf("Error reserving IP address block: %s", err)
	}

	d.Set("project_id", projectID)
	d.SetId(blockAddr.ID)

	return resourcePacketReservedIPBlockRead(d, meta)
}

func getGlobalBool(r *packngo.IPAddressReservation) bool {
	if r.Global != nil {
		return *(r.Global)
	}
	return false
}

func getType(r *packngo.IPAddressReservation) (string, error) {
	globalBool := getGlobalBool(r)
	switch {
	case !r.Public:
		return fmt.Sprintf("private_ipv%d", r.AddressFamily), nil
	case r.Public && !globalBool:
		return fmt.Sprintf("public_ipv%d", r.AddressFamily), nil
	case r.Public && globalBool:
		return fmt.Sprintf("global_ipv%d", r.AddressFamily), nil
	}
	return "", fmt.Errorf("Unknown reservation type %+v", r)
}

func loadBlock(d *schema.ResourceData, reservedBlock *packngo.IPAddressReservation) error {
	ipv4CIDRToQuantity := map[int]int{32: 1, 31: 2, 30: 4, 29: 8, 28: 16, 27: 32, 26: 64, 25: 128, 24: 256}

	d.SetId(reservedBlock.ID)
	d.Set("address", reservedBlock.Address)
	if reservedBlock.Facility != nil {
		d.Set("facility", reservedBlock.Facility.Code)
	}
	d.Set("gateway", reservedBlock.Gateway)
	d.Set("network", reservedBlock.Network)
	d.Set("netmask", reservedBlock.Netmask)
	d.Set("address_family", reservedBlock.AddressFamily)
	d.Set("cidr", reservedBlock.CIDR)
	typ, err := getType(reservedBlock)
	if err != nil {
		return err
	}
	d.Set("type", typ)
	d.Set("public", reservedBlock.Public)
	d.Set("management", reservedBlock.Management)
	d.Set("manageable", reservedBlock.Manageable)
	if reservedBlock.AddressFamily == 4 {
		d.Set("quantity", ipv4CIDRToQuantity[reservedBlock.CIDR])
	} else {
		// In Packet, a reserved IPv6 block is allocated when a device is run in a project.
		// It's always /56, and it can't be created with Terraform, only imported.
		// The longest assignable prefix is /64, making it max 256 subnets per block.
		// The following logic will hold as long as /64 is the smallest assignable subnet size.
		bits := 64 - reservedBlock.CIDR
		if bits > 30 {
			return fmt.Errorf("Strange (too small) CIDR prefix: %d", reservedBlock.CIDR)
		}
		d.Set("quantity", 1<<uint(bits))
	}
	d.Set("project_id", path.Base(reservedBlock.Project.Href))
	d.Set("cidr_notation", fmt.Sprintf("%s/%d", reservedBlock.Network, reservedBlock.CIDR))
	return nil

}

func resourcePacketReservedIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	id := d.Id()

	reservedBlock, _, err := client.ProjectIPs.Get(id, nil)
	if err != nil {
		err = friendlyError(err)
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading IP address block with ID %s: %s", id, err)
	}
	err = loadBlock(d, reservedBlock)
	if (reservedBlock.Description != nil) && (*(reservedBlock.Description) != "") {
		d.Set("description", *(reservedBlock.Description))
	}
	d.Set("global", getGlobalBool(reservedBlock))
	if err != nil {
		return err
	}

	return nil
}

func resourcePacketReservedIPBlockDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)

	id := d.Id()

	_, err := client.ProjectIPs.Remove(id)

	if err != nil {
		return fmt.Errorf("Error deleting IP reservation block %s: %s", id, err)
	}

	d.SetId("")
	return nil
}
