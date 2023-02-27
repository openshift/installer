package ovirt

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ovirtclient "github.com/ovirt/go-ovirt-client/v3"
)

var tagSchema = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Computed: true,
	},
	"name": {
		Type:             schema.TypeString,
		Required:         true,
		Description:      "Name for the tag.",
		ForceNew:         true,
		ValidateDiagFunc: validateNonEmpty,
	},
	"description": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Description for the tag.",
		ForceNew:    true,
	},
}

func (p *provider) tagResource() *schema.Resource {
	return &schema.Resource{
		CreateContext: p.tagCreate,
		ReadContext:   p.tagRead,
		DeleteContext: p.tagDelete,
		Schema:        tagSchema,
		Description:   "The ovirt_tag resource creates tags for virtual machines to use.",
	}
}

func (p *provider) tagCreate(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	name := data.Get("name").(string)
	descriptionRaw, ok := data.GetOk("description")
	params := ovirtclient.NewCreateTagParams()
	if ok {
		var err error
		params, err = params.WithDescription(descriptionRaw.(string))
		if err != nil {
			return errorToDiags("set description", err)
		}
	}
	tag, err := client.CreateTag(name, params)
	if err != nil {
		return errorToDiags(fmt.Sprintf("create tag %s", name), err)
	}
	diags := diag.Diagnostics{}
	data.SetId(string(tag.ID()))
	diags = appendDiags(diags, "set name on ovirt_tag", data.Set("name", tag.Name()))
	diags = appendDiags(diags, "set description on ovirt_tag", data.Set("description", tag.Description()))
	return diags
}

func (p *provider) tagRead(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	tag, err := client.GetTag(ovirtclient.TagID(data.Id()))
	if err != nil {
		if isNotFound(err) {
			data.SetId("")
			return nil
		}
		return errorToDiags(fmt.Sprintf("get tag %s", data.Id()), err)
	}
	diags := diag.Diagnostics{}
	data.SetId(string(tag.ID()))
	diags = appendDiags(diags, "set name on ovirt_tag", data.Set("name", tag.Name()))
	diags = appendDiags(diags, "set description on ovirt_tag", data.Set("description", tag.Description()))
	return diags
}

func (p *provider) tagDelete(ctx context.Context, data *schema.ResourceData, i interface{}) diag.Diagnostics {
	client := p.client.WithContext(ctx)
	if err := client.RemoveTag(
		ovirtclient.TagID(data.Id()),
	); err != nil && ovirtclient.HasErrorCode(err, ovirtclient.ENotFound) {
		return errorToDiags(fmt.Sprintf("remove tag %s", data.Id()), err)
	}
	data.SetId("")
	return nil
}
