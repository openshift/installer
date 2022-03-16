package vsphere

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vmware/govmomi/license"
)

func dataSourceVSphereLicense() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVSphereLicenseRead,

		Schema: map[string]*schema.Schema{
			"license_key": {
				Type:     schema.TypeString,
				Required: true,
			},

			// computed properties returned by the API
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"edition_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVSphereLicenseRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*Client).vimClient
	manager := license.NewManager(client.Client)
	licenseKey := d.Get("license_key").(string)
	if info := getLicenseInfoFromKey(d.Get("license_key").(string), manager); info != nil {
		log.Println("[INFO] Setting the values")
		d.Set("edition_key", info.EditionKey)
		d.Set("total", info.Total)
		d.Set("used", info.Used)
		d.Set("name", info.Name)
		if err := d.Set("labels", keyValuesToMap(info.Labels)); err != nil {
			return err
		}
		d.SetId(licenseKey)
		return nil
	} else {
		return ErrNoSuchKeyFound
	}
}
