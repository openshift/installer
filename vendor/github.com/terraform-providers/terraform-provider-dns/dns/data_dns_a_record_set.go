package dns

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDnsARecordSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsARecordSetRead,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addrs": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceDnsARecordSetRead(d *schema.ResourceData, meta interface{}) error {
	host := d.Get("host").(string)

	a, _, err := lookupIP(host)
	if err != nil {
		return fmt.Errorf("error looking up A records for %q: %s", host, err)
	}
	sort.Strings(a)

	d.Set("addrs", a)
	d.SetId(host)

	return nil
}
