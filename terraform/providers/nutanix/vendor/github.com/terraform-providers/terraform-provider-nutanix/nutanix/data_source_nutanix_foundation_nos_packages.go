package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/foundation"
)

// dataSourceFoundationNOSPackages datasource to pull list of nos packages file names available from foundation vm.
func dataSourceFoundationNOSPackages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFoundationNOSPackagesRead,
		Schema: map[string]*schema.Schema{
			//entities will have list(string) of nos packages
			"entities": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

// dataSourceFoundationNOSPackagesRead method to perform read operation on /enumerate_nos_packages api
func dataSourceFoundationNOSPackagesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// get Foundation api client
	conn := meta.(*Client).FoundationClientAPI

	var resp *foundation.ListNOSPackagesResponse
	resp, err := conn.FileManagement.ListNOSPackages(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if setErr := d.Set("entities", resp); setErr != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}
