package packet

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/packethost/packngo"
)

func dataSourcePacketIPBlockRanges() *schema.Resource {
	return &schema.Resource{
		Read: dataSourcePacketIPBlockRangesRead,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"facility": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ipv4": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"global_ipv4": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"private_ipv4": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"ipv6": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func faclityMatch(ref, ipFacility string) bool {
	if ref == "" {
		return true
	}
	if ref == ipFacility {
		return true
	}
	return false
}

func dataSourcePacketIPBlockRangesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*packngo.Client)
	projectID := d.Get("project_id").(string)
	ips, _, err := client.ProjectIPs.List(projectID)
	if err != nil {
		return err
	}

	facility := d.Get("facility").(string)

	publicIPv4s := []string{}
	globalIPv4s := []string{}
	privateIPv4s := []string{}
	theIPv6s := []string{}
	var targetSlice *[]string

	for _, ip := range ips {
		targetSlice = nil
		cnStr := fmt.Sprintf("%s/%d", ip.Network, ip.CIDR)
		if ip.AddressFamily == 4 {
			if ip.Public {
				if getGlobalBool(&ip) {
					globalIPv4s = append(globalIPv4s, cnStr)
				} else {
					targetSlice = &publicIPv4s
				}
			} else {
				targetSlice = &privateIPv4s
			}
		} else {
			targetSlice = &theIPv6s
		}
		if targetSlice != nil && faclityMatch(facility, ip.Facility.Code) {
			*targetSlice = append(*targetSlice, cnStr)
		}
	}

	d.Set("public_ipv4", publicIPv4s)
	d.Set("global_ipv4", globalIPv4s)
	d.Set("private_ipv4", privateIPv4s)
	d.Set("ipv6", theIPv6s)
	if facility != "" {
		facility = "-" + facility
	}
	d.SetId(projectID + facility + "-IPs")
	return nil

}
