package ovirt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) templatesDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.templatesDataSourceRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				Description:      "Name of the template to look for",
				ValidateDiagFunc: validateNonEmpty,
			},
			"fail_on_empty": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Fail if no templates with the given name were found.",
			},
			"templates": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "oVirt identifier for the template",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "User-provided description for the template.",
						},
					},
				},
			},
		},
		Description: `Search oVirt templates by name.`,
	}
}

func (p *provider) templatesDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	templates, err := client.ListTemplates()
	if err != nil {
		return errorToDiags("list templates", err)
	}
	name := data.Get("name").(string)
	var result []map[string]interface{}
	for _, template := range templates {
		if template.Name() == name {
			result = append(
				result, map[string]interface{}{
					"id":          template.ID(),
					"description": template.Description(),
				},
			)
		}
	}
	data.SetId(name)
	if err := data.Set("templates", result); err != nil {
		return errorToDiags("set templates", err)
	}
	if data.Get("fail_on_empty").(bool) && len(result) == 0 {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "No template found",
				Detail:   fmt.Sprintf("No template with the name %s found.", name),
			},
		}
	}
	return nil
}
