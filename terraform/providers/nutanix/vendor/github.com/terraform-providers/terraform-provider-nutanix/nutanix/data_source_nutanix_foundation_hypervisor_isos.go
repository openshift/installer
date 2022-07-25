package nutanix

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-providers/terraform-provider-nutanix/client/foundation"
)

func dataSourceFoundationHypervisorIsos() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFoundationHypervisorIsosRead,
		Schema: map[string]*schema.Schema{
			"hyperv": hypervisorSchema(),
			"kvm":    hypervisorSchema(),
			"linux":  hypervisorSchema(),
			"esx":    hypervisorSchema(),
			"xen":    hypervisorSchema(),
		},
	}
}

func dataSourceFoundationHypervisorIsosRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Get the client connection
	conn := meta.(*Client).FoundationClientAPI

	resp, err := conn.FileManagement.ListHypervisorISOs(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("hyperv", flattenHypervisorIsos(resp.Hyperv)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("kvm", flattenHypervisorIsos(resp.Kvm)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("esx", flattenHypervisorIsos(resp.Esx)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("linux", flattenHypervisorIsos(resp.Linux)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("xen", flattenHypervisorIsos(resp.Xen)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.UniqueId())
	return nil
}

func flattenHypervisorIsos(ref []*foundation.HypervisorISOReference) []map[string]interface{} {
	hyper := make([]map[string]interface{}, len(ref))
	for i, h := range ref {
		hyper[i] = map[string]interface{}{
			"filename":  h.Filename,
			"supported": h.Supported,
		}
	}
	return hyper
}

func hypervisorSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"filename": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"supported": {
					Type:     schema.TypeBool,
					Computed: true,
				},
			},
		},
	}
}
