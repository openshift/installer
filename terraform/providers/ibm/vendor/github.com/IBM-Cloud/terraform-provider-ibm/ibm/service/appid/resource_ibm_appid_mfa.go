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

func ResourceIBMAppIDMFA() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceIBMAppIDMFARead,
		CreateContext: resourceIBMAppIDMFACreate,
		UpdateContext: resourceIBMAppIDMFACreate,
		DeleteContext: resourceIBMAppIDMFADelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"is_active": {
				Description: "`true` if MFA is active",
				Type:        schema.TypeBool,
				Required:    true,
			},
		},
	}
}

func resourceIBMAppIDMFARead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Id()

	mfa, resp, err := appIDClient.GetMFAConfigWithContext(ctx, &appid.GetMFAConfigOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing AppID MFA configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error getting AppID MFA configuration: %s\n%s", err, resp)
	}

	if mfa.IsActive != nil {
		d.Set("is_active", *mfa.IsActive)
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDMFACreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	input := &appid.UpdateMFAConfigOptions{
		TenantID: &tenantID,
		IsActive: &isActive,
	}

	_, resp, err := appIDClient.UpdateMFAConfigWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error updating AppID MFA configuration: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDMFARead(ctx, d, meta)
}

func resourceIBMAppIDMFADelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.UpdateMFAConfigOptions{
		TenantID: &tenantID,
		IsActive: helpers.Bool(false),
	}

	_, resp, err := appIDClient.UpdateMFAConfigWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error resetting AppID MFA configuration: %s\n%s", err, resp)
	}

	d.SetId("")
	return nil
}
