package ovirt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v3"
)

var templateSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "oVirt ID of this template.",
	},
	"vm_id": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "oVirt ID of the VM the template is based on.",
		ValidateDiagFunc: validateUUID,
	},
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		ForceNew:         true,
		Description:      "User-provided name for the template. Must only consist of lower- and uppercase letters, numbers, dash, underscore and dot.",
		ValidateDiagFunc: validateNonEmpty,
	},
	"description": {
		Type:             schema.TypeString,
		Optional:         true,
		ForceNew:         true,
		Description:      "User-provided description for the template.",
		ValidateDiagFunc: validateNonEmpty,
	},
}

func (p *provider) templateResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.templateCreate,
		ReadContext:   p.templateRead,
		DeleteContext: p.templateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: p.templateImport,
		},
		Schema:      templateSchema,
		Description: "The ovirt_template resource manages templates for virtual machines in oVirt.",
	}
}

func (p *provider) templateCreate(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	VMID := data.Get("vm_id").(string)
	templateName := data.Get("name").(string)

	params := ovirtclient.TemplateCreateParams()
	if templateDescription, ok := data.GetOk("description"); ok {
		_, err := params.WithDescription(templateDescription.(string))
		if err != nil {
			return diag.Diagnostics{
				diag.Diagnostic{
					Severity: diag.Error,
					Summary:  fmt.Sprintf("Invalid template description: %s", templateDescription),
					Detail:   err.Error(),
				},
			}
		}
	}

	template, err := client.CreateTemplate(ovirtclient.VMID(VMID), templateName, params)
	if err != nil {
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to create template for VM %s", VMID),
				Detail:   err.Error(),
			},
		}
	}

	data.SetId(string(template.ID()))
	diags := diag.Diagnostics{}
	diags = setResourceField(data, "vm_id", VMID, diags)
	diags = setResourceField(data, "name", template.Name(), diags)
	diags = setResourceField(data, "description", template.Description(), diags)

	return diags
}

func (p *provider) templateRead(
	ctx context.Context,
	data *schema.ResourceData,
	_ interface{},
) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	id := data.Id()
	template, err := client.GetTemplate(ovirtclient.TemplateID(id))
	if err != nil {
		if isNotFound(err) {
			data.SetId("")
			return nil
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to fetch template %s", id),
				Detail:   err.Error(),
			},
		}
	}

	data.SetId(string(template.ID()))
	diags := diag.Diagnostics{}
	diags = setResourceField(data, "name", template.Name(), diags)
	diags = setResourceField(data, "description", template.Description(), diags)

	return diags
}

func (p *provider) templateDelete(ctx context.Context, data *schema.ResourceData, _ interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)

	err := client.RemoveTemplate(ovirtclient.TemplateID(data.Id()))
	if err != nil {
		if isNotFound(err) {
			data.SetId("")
			return nil
		}
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("Failed to remove template %s", data.Id()),
				Detail:   err.Error(),
			},
		}
	}

	data.SetId("")

	return nil
}

func (p *provider) templateImport(ctx context.Context, data *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData,
	error,
) {
	client := p.client.WithContext(ctx)

	template, err := client.GetTemplate(ovirtclient.TemplateID(data.Id()))
	if err != nil {
		return nil, fmt.Errorf("failed to import template %s (%w)", data.Id(), err)
	}

	data.SetId(data.Id())
	diags := diag.Diagnostics{}
	diags = setResourceField(data, "name", template.Name(), diags)
	diags = setResourceField(data, "description", template.Description(), diags)
	if err := diagsToError(diags); err != nil {
		return nil, fmt.Errorf("failed to import template %s (%w)", data.Id(), err)
	}

	return []*schema.ResourceData{
		data,
	}, nil
}
