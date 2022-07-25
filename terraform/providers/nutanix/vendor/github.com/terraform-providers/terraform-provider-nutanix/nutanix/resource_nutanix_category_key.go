package nutanix

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
	"github.com/terraform-providers/terraform-provider-nutanix/utils"
)

func resourceNutanixCategoryKey() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixCategoryKeyCreateOrUpdate,
		ReadContext:   resourceNutanixCategoryKeyRead,
		UpdateContext: resourceNutanixCategoryKeyCreateOrUpdate,
		DeleteContext: resourceNutanixCategoryKeyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"system_defined": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"api_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceNutanixCategoryKeyCreateOrUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating CategoryKey: %s", resourceData.Get("name").(string))

	conn := meta.(*Client).API

	request := &v3.CategoryKey{}

	name, nameOK := resourceData.GetOk("name")

	// Read Arguments and set request values
	if desc, ok := resourceData.GetOk("description"); ok {
		request.Description = utils.StringPtr(desc.(string))
	}

	// validate required fields
	if !nameOK {
		return diag.Errorf("please provide the required attribute name")
	}

	request.Name = utils.StringPtr(name.(string))

	// Make request to the API
	resp, err := conn.V3.CreateOrUpdateCategoryKey(request)

	if err != nil {
		return diag.FromErr(err)
	}

	n := *resp.Name

	// set terraform state
	resourceData.SetId(n)

	return resourceNutanixCategoryKeyRead(ctx, resourceData, meta)
}

func resourceNutanixCategoryKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading CategoryKey: %s", d.Get("name").(string))

	// Get client connection
	conn := meta.(*Client).API

	// Make request to the API
	resp, err := conn.V3.GetCategoryKey(d.Id())

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("api_version", utils.StringValue(resp.APIVersion))
	d.Set("name", utils.StringValue(resp.Name))
	d.Set("description", utils.StringValue(resp.Description))

	return diag.FromErr(d.Set("system_defined", utils.BoolValue(resp.SystemDefined)))
}

func resourceNutanixCategoryKeyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

	log.Printf("[Debug] Destroying the category with the ID %s", d.Id())

	if err := conn.V3.DeleteCategoryKey(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
