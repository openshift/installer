package dns

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/miekg/dns"
)

func resourceDnsTXTRecordSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsTXTRecordSetCreate,
		Read:   resourceDnsTXTRecordSetRead,
		Update: resourceDnsTXTRecordSetUpdate,
		Delete: resourceDnsTXTRecordSetDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDnsImport,
		},

		Schema: map[string]*schema.Schema{
			"zone": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateZone,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateName,
			},
			"txt": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"ttl": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  3600,
			},
		},
	}
}

func resourceDnsTXTRecordSetCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(resourceFQDN(d))

	return resourceDnsTXTRecordSetUpdate(d, meta)
}

func resourceDnsTXTRecordSetRead(d *schema.ResourceData, meta interface{}) error {

	answers, err := resourceDnsRead(d, meta, dns.TypeTXT)
	if err != nil {
		return err
	}

	if len(answers) > 0 {

		var ttl sort.IntSlice

		txt := schema.NewSet(schema.HashString, nil)
		for _, record := range answers {
			switch r := record.(type) {
			case *dns.TXT:
				txt.Add(strings.Join(r.Txt, ""))
				ttl = append(ttl, int(r.Hdr.Ttl))
			default:
				return fmt.Errorf("didn't get an TXT record")
			}
		}
		sort.Sort(ttl)

		d.Set("txt", txt)
		d.Set("ttl", ttl[0])
	} else {
		d.SetId("")
	}

	return nil
}

func resourceDnsTXTRecordSetUpdate(d *schema.ResourceData, meta interface{}) error {

	if meta != nil {

		ttl := d.Get("ttl").(int)
		fqdn := resourceFQDN(d)

		msg := new(dns.Msg)

		msg.SetUpdate(d.Get("zone").(string))

		if d.HasChange("txt") {
			o, n := d.GetChange("txt")
			os := o.(*schema.Set)
			ns := n.(*schema.Set)
			remove := os.Difference(ns).List()
			add := ns.Difference(os).List()

			// Loop through all the old addresses and remove them
			for _, txt := range remove {
				rr_remove, _ := dns.NewRR(fmt.Sprintf("%s %d TXT \"%s\"", fqdn, ttl, txt))
				msg.Remove([]dns.RR{rr_remove})
			}
			// Loop through all the new addresses and insert them
			for _, txt := range add {
				rr_insert, _ := dns.NewRR(fmt.Sprintf("%s %d TXT \"%s\"", fqdn, ttl, txt))
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

		return resourceDnsTXTRecordSetRead(d, meta)
	} else {
		return fmt.Errorf("update server is not set")
	}
}

func resourceDnsTXTRecordSetDelete(d *schema.ResourceData, meta interface{}) error {

	return resourceDnsDelete(d, meta, dns.TypeTXT)
}
