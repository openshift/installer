package dns

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/miekg/dns"
)

func resourceDnsARecordSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsARecordSetCreate,
		Read:   resourceDnsARecordSetRead,
		Update: resourceDnsARecordSetUpdate,
		Delete: resourceDnsARecordSetDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDnsImport,
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateZone,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"addresses": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      hashIPString,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  3600,
			},
		},
	}
}

func resourceDnsARecordSetCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(resourceFQDN(d))

	return resourceDnsARecordSetUpdate(d, meta)
}

func resourceDnsARecordSetRead(d *schema.ResourceData, meta interface{}) error {

	answers, err := resourceDnsRead(d, meta, dns.TypeA)
	if err != nil {
		return err
	}

	if len(answers) > 0 {

		var ttl sort.IntSlice

		addresses := schema.NewSet(hashIPString, nil)
		for _, record := range answers {
			addr, t, err := getAVal(record)
			if err != nil {
				return fmt.Errorf("Error querying DNS record: %s", err)
			}
			addresses.Add(addr)
			ttl = append(ttl, t)
		}
		sort.Sort(ttl)

		d.Set("addresses", addresses)
		d.Set("ttl", ttl[0])
	} else {
		d.SetId("")
	}

	return nil
}

func resourceDnsARecordSetUpdate(d *schema.ResourceData, meta interface{}) error {

	if meta != nil {

		ttl := d.Get("ttl").(int)
		rec_fqdn := resourceFQDN(d)

		msg := new(dns.Msg)

		msg.SetUpdate(d.Get("zone").(string))

		if d.HasChange("addresses") {
			o, n := d.GetChange("addresses")
			os := o.(*schema.Set)
			ns := n.(*schema.Set)
			remove := os.Difference(ns).List()
			add := ns.Difference(os).List()

			// Loop through all the old addresses and remove them
			for _, addr := range remove {
				rr_remove, _ := dns.NewRR(fmt.Sprintf("%s %d A %s", rec_fqdn, ttl, addr.(string)))
				msg.Remove([]dns.RR{rr_remove})
			}
			// Loop through all the new addresses and insert them
			for _, addr := range add {
				rr_insert, _ := dns.NewRR(fmt.Sprintf("%s %d A %s", rec_fqdn, ttl, addr.(string)))
				msg.Insert([]dns.RR{rr_insert})
			}

			r, err := exchange(msg, true, meta)
			if err != nil {
				d.SetId("")
				return fmt.Errorf("Error updating DNS record: %s", err)
			}
			if r.Rcode != dns.RcodeSuccess {
				d.SetId("")
				return fmt.Errorf("Error updating DNS record: %v (%s)", r.Rcode, dns.RcodeToString[r.Rcode])
			}
		}

		return resourceDnsARecordSetRead(d, meta)
	} else {
		return fmt.Errorf("update server is not set")
	}
}

func resourceDnsARecordSetDelete(d *schema.ResourceData, meta interface{}) error {

	return resourceDnsDelete(d, meta, dns.TypeA)
}
