package packet

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourcePacketPreCreatedIPBlock() *schema.Resource {
	s := packetIPComputedFields()
	s["project_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	s["global"] = &schema.Schema{
		Type:     schema.TypeBool,
		Optional: true,
		Default:  false,
	}
	s["public"] = &schema.Schema{
		Type:     schema.TypeBool,
		Required: true,
	}

	s["facility"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
	}
	s["address_family"] = &schema.Schema{
		Type:     schema.TypeInt,
		Required: true,
	}
	s["cidr_notation"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	s["quantity"] = &schema.Schema{
		Type:     schema.TypeInt,
		Computed: true,
	}

	return &schema.Resource{
		Read:   dataSourcePacketReservedIPBlockRead,
		Schema: s,
	}
}

func dataSourcePacketReservedIPBlockRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	projectID := d.Get("project_id").(string)
	log.Println("[DEBUG] packet_precreated_ip_block - getting list of IPs in a project")
	ips, _, err := client.ProjectIPs.List(projectID)
	if err != nil {
		return err
	}
	ipv := d.Get("address_family").(int)
	public := d.Get("public").(bool)
	global := d.Get("global").(bool)

	if !public && global {
		return fmt.Errorf("Private (non-public) global IP address blocks are not supported in Packet")
	}

	fval, fok := d.GetOk("facility")
	if fok && global {
		return fmt.Errorf("You can't specify facility for global IP block - addresses from global blocks can be assigned to devices across several facilities")
	}

	if fok {
		// lookup of not-global block
		facility := fval.(string)
		for _, ip := range ips {
			if ip.Public == public && ip.AddressFamily == ipv && facility == ip.Facility.Code {
				loadBlock(d, &ip)
				break
			}
		}
	} else {
		// lookup of global block
		for _, ip := range ips {
			blockGlobal := getGlobalBool(&ip)
			if ip.Public == public && ip.AddressFamily == ipv && blockGlobal {
				loadBlock(d, &ip)
				break
			}
		}

	}
	if d.Get("cidr_notation") == "" {
		return fmt.Errorf("Could not find matching reserved block, all IPs were %v", ips)
	}
	return nil

}
