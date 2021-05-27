package dns

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/miekg/dns"
)

func resourceDnsCnameRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsCnameRecordCreate,
		Read:   resourceDnsCnameRecordRead,
		Update: resourceDnsCnameRecordUpdate,
		Delete: resourceDnsCnameRecordDelete,
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
			"cname": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateZone,
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

func resourceDnsCnameRecordCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(resourceFQDN(d))

	return resourceDnsCnameRecordUpdate(d, meta)
}

func resourceDnsCnameRecordRead(d *schema.ResourceData, meta interface{}) error {

	answers, err := resourceDnsRead(d, meta, dns.TypeCNAME)
	if err != nil {
		return err
	}

	if len(answers) > 0 {

		if len(answers) > 1 {
			return fmt.Errorf("Error querying DNS record: multiple responses received")
		}
		record := answers[0]
		cname, ttl, err := getCnameVal(record)
		if err != nil {
			return fmt.Errorf("Error querying DNS record: %s", err)
		}
		d.Set("cname", cname)
		d.Set("ttl", ttl)
	} else {
		d.SetId("")
	}

	return nil
}

func resourceDnsCnameRecordUpdate(d *schema.ResourceData, meta interface{}) error {

	if meta != nil {

		ttl := d.Get("ttl").(int)

		rec_fqdn := resourceFQDN(d)

		msg := new(dns.Msg)

		msg.SetUpdate(d.Get("zone").(string))

		if d.HasChange("cname") {
			o, n := d.GetChange("cname")

			if o != "" {
				rr_remove, _ := dns.NewRR(fmt.Sprintf("%s %d CNAME %s", rec_fqdn, ttl, o))
				msg.Remove([]dns.RR{rr_remove})
			}
			if n != "" {
				rr_insert, _ := dns.NewRR(fmt.Sprintf("%s %d CNAME %s", rec_fqdn, ttl, n))
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

		return resourceDnsCnameRecordRead(d, meta)
	} else {
		return fmt.Errorf("update server is not set")
	}
}

func resourceDnsCnameRecordDelete(d *schema.ResourceData, meta interface{}) error {

	return resourceDnsDelete(d, meta, dns.TypeCNAME)
}
