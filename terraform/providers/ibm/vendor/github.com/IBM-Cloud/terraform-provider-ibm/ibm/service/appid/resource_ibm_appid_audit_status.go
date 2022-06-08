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

func ResourceIBMAppIDAuditStatus() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMAppIDAuditStatusCreate,
		ReadContext:   resourceIBMAppIDAuditStatusRead,
		DeleteContext: resourceIBMAppIDAuditStatusDelete,
		UpdateContext: resourceIBMAppIDAuditStatusUpdate,
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
			"is_active": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "The auditing status of the tenant.",
			},
		},
	}
}

func resourceIBMAppIDAuditStatusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Id()

	auditStatus, resp, err := appIDClient.GetAuditStatusWithContext(ctx, &appid.GetAuditStatusOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing audit status configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error getting AppID audit status: %s\n%s", err, resp)
	}

	d.Set("is_active", *auditStatus.IsActive)
	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDAuditStatusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	resp, err := appIDClient.SetAuditStatusWithContext(ctx, &appid.SetAuditStatusOptions{
		TenantID: &tenantID,
		IsActive: &isActive,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing audit status from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error setting AppID audit status: %s\n%s", err, resp)
	}

	d.SetId(tenantID)
	return resourceIBMAppIDAuditStatusRead(ctx, d, meta)
}

func resourceIBMAppIDAuditStatusDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	resp, err := appIDClient.SetAuditStatusWithContext(ctx, &appid.SetAuditStatusOptions{
		TenantID: &tenantID,
		IsActive: helpers.Bool(false),
	})

	if err != nil {
		return diag.Errorf("Error resetting AppID audit status: %s\n%s", err, resp)
	}

	d.SetId("")
	return nil
}

func resourceIBMAppIDAuditStatusUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceIBMAppIDAuditStatusCreate(ctx, d, m)
}
