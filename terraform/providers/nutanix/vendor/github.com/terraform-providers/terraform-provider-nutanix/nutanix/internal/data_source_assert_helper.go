package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Data source func returning schema for assert helper
func DataSourceAssertHelper() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssertLogic,
		Schema: map[string]*schema.Schema{
			"checks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"error_message": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

// Data source logic for nutanix_assert_helper
func dataSourceAssertLogic(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// check all conditions results one by one and append erroes if found
	checks := d.Get("checks").([]interface{})
	var diags diag.Diagnostics
	for _, v := range checks {
		check := v.(map[string]interface{})
		if !(check["condition"]).(bool) {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  check["error_message"].(string),
			})
		}
	}

	// if any checks failed return the errors collected
	if len(diags) > 0 {
		return diags
	}
	d.SetId(resource.UniqueId())
	return nil
}
