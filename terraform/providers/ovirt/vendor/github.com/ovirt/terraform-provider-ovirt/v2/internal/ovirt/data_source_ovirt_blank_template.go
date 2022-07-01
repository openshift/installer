package ovirt

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *provider) blankTemplateDataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: p.blankTemplateDataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "oVirt ID of the blank template.",
				Computed:    true,
			},
		},
		Description: `Returns the ID of the blank template, even when the blank template was deleted and re-created.`,
	}
}

func (p *provider) blankTemplateDataSourceRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	tpl, err := client.GetBlankTemplate()
	if err != nil {
		return errorToDiags("getting blank template", err)
	}

	data.SetId(string(tpl.ID()))

	return nil
}
