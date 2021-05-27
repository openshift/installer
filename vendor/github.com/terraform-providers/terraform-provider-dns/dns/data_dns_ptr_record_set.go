package dns

import (
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDnsPtrRecordSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsPtrRecordSetRead,
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ptr": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDnsPtrRecordSetRead(d *schema.ResourceData, meta interface{}) error {
	ipAddress := d.Get("ip_address").(string)
	names, err := net.LookupAddr(ipAddress)
	if err != nil {
		return fmt.Errorf("error looking up PTR records for %q: %s", ipAddress, err)
	}
	if len(names) == 0 {
		return fmt.Errorf("error looking up PTR records for %q: no records found", ipAddress)
	}

	d.Set("ptr", names[0])
	d.SetId(ipAddress)

	return nil
}
