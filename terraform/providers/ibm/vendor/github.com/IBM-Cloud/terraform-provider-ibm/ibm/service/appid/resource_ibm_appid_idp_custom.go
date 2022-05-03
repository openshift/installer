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

func ResourceIBMAppIDIDPCustom() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMAppIDIDPCustomCreate,
		ReadContext:   resourceIBMAppIDIDPCustomRead,
		DeleteContext: resourceIBMAppIDIDPCustomDelete,
		UpdateContext: resourceIBMAppIDIDPCustomUpdate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"public_key": {
				Description: "This is the public key used to validate your signed JWT. It is required to be a PEM in the RS256 or greater format.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceIBMAppIDIDPCustomRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Id()

	config, resp, err := appIDClient.GetCustomIDPWithContext(ctx, &appid.GetCustomIDPOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing AppID custom IDP configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error loading AppID custom IDP: %s\n%s", err, resp)
	}

	d.Set("is_active", *config.IsActive)

	if config.Config != nil && config.Config.PublicKey != nil {
		if err := d.Set("public_key", *config.Config.PublicKey); err != nil {
			return diag.Errorf("Failed setting AppID custom IDP public_key: %s", err)
		}
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDIDPCustomCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	config := &appid.SetCustomIDPOptions{
		TenantID: &tenantID,
		IsActive: &isActive,
	}

	if isActive {
		config.Config = &appid.CustomIDPConfigParamsConfig{}

		if pKey, ok := d.GetOk("public_key"); ok {
			config.Config.PublicKey = helpers.String(pKey.(string))
		}
	}

	_, resp, err := appIDClient.SetCustomIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error applying AppID custom IDP configuration: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDIDPCustomRead(ctx, d, meta)
}

func appIDCustomIDPDefaults(tenantID string) *appid.SetCustomIDPOptions {
	return &appid.SetCustomIDPOptions{
		TenantID: &tenantID,
		IsActive: helpers.Bool(false),
	}
}

func resourceIBMAppIDIDPCustomDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	config := appIDCustomIDPDefaults(tenantID)

	_, resp, err := appIDClient.SetCustomIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error resetting AppID custom IDP configuration: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}

func resourceIBMAppIDIDPCustomUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// since this is configuration we can reuse create method
	return resourceIBMAppIDIDPCustomCreate(ctx, d, m)
}
