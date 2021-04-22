package dns

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/miekg/dns"
)

func resourceDnsMXRecordSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsMXRecordSetCreate,
		Read:   resourceDnsMXRecordSetRead,
		Update: resourceDnsMXRecordSetUpdate,
		Delete: resourceDnsMXRecordSetDelete,
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
			"mx": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"preference": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"exchange": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validateZone,
						},
					},
				},
				Set: resourceDnsMXRecordSetHash,
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

func resourceDnsMXRecordSetCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(resourceFQDN(d))

	return resourceDnsMXRecordSetUpdate(d, meta)
}

func resourceDnsMXRecordSetRead(d *schema.ResourceData, meta interface{}) error {

	answers, err := resourceDnsRead(d, meta, dns.TypeMX)
	if err != nil {
		return err
	}

	if len(answers) > 0 {

		var ttl sort.IntSlice

		mx := schema.NewSet(resourceDnsMXRecordSetHash, nil)
		for _, record := range answers {
			switch r := record.(type) {
			case *dns.MX:
				m := map[string]interface{}{
					"preference": int(r.Preference),
					"exchange":   r.Mx,
				}
				mx.Add(m)
				ttl = append(ttl, int(r.Hdr.Ttl))
			default:
				return fmt.Errorf("didn't get an MX record")
			}
		}
		sort.Sort(ttl)

		d.Set("mx", mx)
		d.Set("ttl", ttl[0])
	} else {
		d.SetId("")
	}

	return nil
}

func resourceDnsMXRecordSetUpdate(d *schema.ResourceData, meta interface{}) error {

	if meta != nil {

		ttl := d.Get("ttl").(int)
		fqdn := resourceFQDN(d)

		msg := new(dns.Msg)

		msg.SetUpdate(d.Get("zone").(string))

		if d.HasChange("mx") {
			o, n := d.GetChange("mx")
			os := o.(*schema.Set)
			ns := n.(*schema.Set)
			remove := os.Difference(ns).List()
			add := ns.Difference(os).List()

			// Loop through all the old addresses and remove them
			for _, mx := range remove {
				m := mx.(map[string]interface{})
				rr_remove, _ := dns.NewRR(fmt.Sprintf("%s %d MX %d %s", fqdn, ttl, m["preference"], m["exchange"]))
				msg.Remove([]dns.RR{rr_remove})
			}
			// Loop through all the new addresses and insert them
			for _, mx := range add {
				m := mx.(map[string]interface{})
				rr_insert, _ := dns.NewRR(fmt.Sprintf("%s %d MX %d %s", fqdn, ttl, m["preference"], m["exchange"]))
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

		return resourceDnsMXRecordSetRead(d, meta)
	} else {
		return fmt.Errorf("update server is not set")
	}
}

func resourceDnsMXRecordSetDelete(d *schema.ResourceData, meta interface{}) error {

	return resourceDnsDelete(d, meta, dns.TypeMX)
}

func resourceDnsMXRecordSetHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["preference"].(int)))
	buf.WriteString(fmt.Sprintf("%s-", m["exchange"].(string)))

	return hashcode.String(buf.String())
}
