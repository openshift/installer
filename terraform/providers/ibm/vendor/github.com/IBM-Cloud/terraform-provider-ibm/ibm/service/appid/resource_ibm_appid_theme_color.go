package appid

import (
	"context"
	"log"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const defaultHeaderColor = "#EEF2F5" // AppID default

func ResourceIBMAppIDThemeColor() *schema.Resource {
	return &schema.Resource{
		Description:   "Colors of the App ID login widget",
		CreateContext: resourceIBMAppIDThemeColorCreate,
		UpdateContext: resourceIBMAppIDThemeColorUpdate,
		ReadContext:   resourceIBMAppIDThemeColorRead,
		DeleteContext: resourceIBMAppIDThemeColorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The AppID instance GUID",
			},
			"header_color": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceIBMAppIDThemeColorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Id()

	colors, resp, err := appIDClient.GetThemeColorWithContext(ctx, &appid.GetThemeColorOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing AppID theme color configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error getting AppID theme colors: %s\n%s", err, resp)
	}

	if colors.HeaderColor != nil {
		d.Set("header_color", *colors.HeaderColor)
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDThemeColorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.PostThemeColorOptions{
		TenantID:    &tenantID,
		HeaderColor: helpers.String(d.Get("header_color").(string)),
	}

	resp, err := appIDClient.PostThemeColorWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error setting AppID theme color: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDThemeColorRead(ctx, d, meta)
}

func resourceIBMAppIDThemeColorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceIBMAppIDThemeColorCreate(ctx, d, meta)
}

func resourceIBMAppIDThemeColorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.PostThemeColorOptions{
		TenantID:    &tenantID,
		HeaderColor: helpers.String(defaultHeaderColor),
	}

	resp, err := appIDClient.PostThemeColorWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error resetting AppID theme color: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
