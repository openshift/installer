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

func ResourceIBMAppIDThemeText() *schema.Resource {
	return &schema.Resource{
		Description:   "Update theme texts of the App ID login widget",
		CreateContext: resourceIBMAppIDThemeTextCreate,
		ReadContext:   resourceIBMAppIDThemeTextRead,
		UpdateContext: resourceIBMAppIDThemeTextUpdate,
		DeleteContext: resourceIBMAppIDThemeTextDelete,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The AppID instance GUID",
			},
			"tab_title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"footnote": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIBMAppIDThemeTextRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Id()

	text, resp, err := appIDClient.GetThemeTextWithContext(ctx, &appid.GetThemeTextOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing AppID theme text configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error getting AppID theme text: %s\n%s", err, resp)
	}

	if text.TabTitle != nil {
		d.Set("tab_title", *text.TabTitle)
	}

	if text.Footnote != nil {
		d.Set("footnote", *text.Footnote)
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDThemeTextCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.PostThemeTextOptions{
		TenantID: &tenantID,
		TabTitle: helpers.String(d.Get("tab_title").(string)),
		Footnote: helpers.String(d.Get("footnote").(string)),
	}

	resp, err := appIDClient.PostThemeTextWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error setting AppID theme text: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDThemeTextRead(ctx, d, meta)
}

func resourceIBMAppIDThemeTextUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceIBMAppIDThemeTextCreate(ctx, d, meta)
}

func resourceIBMAppIDThemeTextDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	input := &appid.PostThemeTextOptions{
		TenantID: &tenantID,
		TabTitle: helpers.String("Login"),
		Footnote: helpers.String("Powered by App ID"),
	}

	resp, err := appIDClient.PostThemeTextWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error resetting AppID theme text: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
