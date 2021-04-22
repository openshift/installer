package dns

import (
	"fmt"
	"net"
	"sort"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDnsNSRecordSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsNSRecordSetRead,
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"nameservers": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
	}
}

func dataSourceDnsNSRecordSetRead(d *schema.ResourceData, meta interface{}) error {
	host := d.Get("host").(string)

	nsRecords, err := net.LookupNS(host)
	if err != nil {
		return fmt.Errorf("error looking up NS records for %q: %s", host, err)
	}

	nameservers := make([]string, len(nsRecords))
	for i, record := range nsRecords {
		nameservers[i] = record.Host
	}
	sort.Strings(nameservers)

	err = d.Set("nameservers", nameservers)
	if err != nil {
		return err
	}
	d.SetId(host)

	return nil
}
