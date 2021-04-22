package dns

import (
	"fmt"
	"sort"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/miekg/dns"
)

func resourceDnsNSRecordSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsNSRecordSetCreate,
		Read:   resourceDnsNSRecordSetRead,
		Update: resourceDnsNSRecordSetUpdate,
		Delete: resourceDnsNSRecordSetDelete,
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
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"nameservers": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateZone,
				},
				Set: schema.HashString,
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

func resourceDnsNSRecordSetCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(resourceFQDN(d))

	return resourceDnsNSRecordSetUpdate(d, meta)
}

func resourceDnsNSRecordSetRead(d *schema.ResourceData, meta interface{}) error {

	answers, err := resourceDnsRead(d, meta, dns.TypeNS)
	if err != nil {
		return err
	}

	if len(answers) > 0 {

		var ttl sort.IntSlice

		nameservers := schema.NewSet(schema.HashString, nil)
		for _, record := range answers {
			nameserver, t, err := getNSVal(record)
			if err != nil {
				return fmt.Errorf("Error querying DNS record: %s", err)
			}
			nameservers.Add(nameserver)
			ttl = append(ttl, t)
		}
		sort.Sort(ttl)

		d.Set("nameservers", nameservers)
		d.Set("ttl", ttl[0])
	} else {
		d.SetId("")
	}

	return nil

}

func resourceDnsNSRecordSetUpdate(d *schema.ResourceData, meta interface{}) error {

	if meta != nil {

		ttl := d.Get("ttl").(int)

		rec_fqdn := resourceFQDN(d)

		msg := new(dns.Msg)

		msg.SetUpdate(d.Get("zone").(string))

		if d.HasChange("nameservers") {
			o, n := d.GetChange("nameservers")
			os := o.(*schema.Set)
			ns := n.(*schema.Set)
			remove := os.Difference(ns).List()
			add := ns.Difference(os).List()

			// Loop through all the old nameservers and remove them
			for _, nameserver := range remove {
				rr_remove, _ := dns.NewRR(fmt.Sprintf("%s %d NS %s", rec_fqdn, ttl, nameserver.(string)))
				msg.Remove([]dns.RR{rr_remove})
			}
			// Loop through all the new nameservers and insert them
			for _, nameserver := range add {
				rr_insert, _ := dns.NewRR(fmt.Sprintf("%s %d NS %s", rec_fqdn, ttl, nameserver.(string)))
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

		return resourceDnsNSRecordSetRead(d, meta)
	} else {
		return fmt.Errorf("update server is not set")
	}
}

func resourceDnsNSRecordSetDelete(d *schema.ResourceData, meta interface{}) error {

	return resourceDnsDelete(d, meta, dns.TypeNS)
}
