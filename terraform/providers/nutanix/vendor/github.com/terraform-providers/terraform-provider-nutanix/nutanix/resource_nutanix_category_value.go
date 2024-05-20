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

func resourceNutanixCategoryValue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNutanixCategoryValueCreateOrUpdate,
		ReadContext:   resourceNutanixCategoryValueRead,
		UpdateContext: resourceNutanixCategoryValueCreateOrUpdate,
		DeleteContext: resourceNutanixCategoryValueDelete,
		Schema: map[string]*schema.Schema{
			"value": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
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

func resourceNutanixCategoryValueCreateOrUpdate(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Creating CategoryValue: %s", resourceData.Get("value").(string))

	conn := meta.(*Client).API

	request := &v3.CategoryValue{}

	name, nameOK := resourceData.GetOk("name")

	value, valueOK := resourceData.GetOk("value")

	// Read Arguments and set request values
	if desc, ok := resourceData.GetOk("description"); ok {
		request.Description = utils.StringPtr(desc.(string))
	}

	// validate required fields
	if !nameOK || !valueOK {
		return diag.Errorf("please provide the required attributes name and value")
	}

	request.Value = utils.StringPtr(value.(string))

	// Make request to the API
	resp, err := conn.V3.CreateOrUpdateCategoryValue(name.(string), request)

	if err != nil {
		return diag.FromErr(err)
	}

	v := *resp.Value

	// set terraform state
	resourceData.SetId(v)

	return resourceNutanixCategoryValueRead(ctx, resourceData, meta)
}

func resourceNutanixCategoryValueRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Reading CategoryValue: %s", d.Id())

	name, nameOK := d.GetOk("name")

	if !nameOK {
		return diag.Errorf("please provide the required attributes name")
	}

	// Get client connection
	conn := meta.(*Client).API

	// Make request to the API
	resp, err := conn.V3.GetCategoryValue(name.(string), d.Id())

	if err != nil {
		if strings.Contains(fmt.Sprint(err), "ENTITY_NOT_FOUND") || strings.Contains(fmt.Sprint(err), "CATEGORY_NAME_VALUE_MISMATCH") {
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

func resourceNutanixCategoryValueDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*Client).API

	name, nameOK := d.GetOk("name")

	if !nameOK {
		return diag.Errorf("please provide the required attributes name")
	}

	log.Printf("[Debug] Destroying the category with the ID %s", d.Id())

	if err := conn.V3.DeleteCategoryValue(name.(string), d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
