package dns

import (
	"fmt"
	"net"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDnsSRVRecordSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsSRVRecordSetRead,
		Schema: map[string]*schema.Schema{
			"service": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"srv": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}

func dataSourceDnsSRVRecordSetRead(d *schema.ResourceData, meta interface{}) error {
	service := d.Get("service").(string)

	_, records, err := net.LookupSRV("", "", service)
	if err != nil {
		return fmt.Errorf("error looking up SRV records for %q: %s", service, err)
	}

	// Sort by priority ascending, weight descending, target
	// alphabetically, and port ascending
	sort.Slice(records, func(i, j int) bool {
		if records[i].Priority < records[j].Priority {
			return true
		}
		if records[i].Priority > records[j].Priority {
			return false
		}
		if records[i].Weight > records[j].Weight {
			return true
		}
		if records[i].Weight < records[j].Weight {
			return false
		}
		if records[i].Target < records[j].Target {
			return true
		}
		if records[i].Target > records[j].Target {
			return false
		}
		return records[i].Port < records[j].Port
	})

	srv := make([]map[string]interface{}, len(records))
	for i, record := range records {
		srv[i] = map[string]interface{}{
			"priority": int(record.Priority),
			"weight":   int(record.Weight),
			"port":     int(record.Port),
			"target":   record.Target,
		}
	}

	err = d.Set("srv", srv)
	if err != nil {
		return err
	}
	d.SetId(service)

	return nil
}
