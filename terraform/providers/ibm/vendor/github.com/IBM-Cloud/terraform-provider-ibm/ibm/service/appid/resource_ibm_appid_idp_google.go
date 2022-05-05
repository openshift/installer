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

func ResourceIBMAppIDIDPGoogle() *schema.Resource {
	return &schema.Resource{
		Description:   "Update Google identity provider configuration.",
		CreateContext: resourceIBMAppIDIDPGoogleCreate,
		ReadContext:   resourceIBMAppIDIDPGoogleRead,
		DeleteContext: resourceIBMAppIDIDPGoogleDelete,
		UpdateContext: resourceIBMAppIDIDPGoogleUpdate,
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
				Description: "`true` if Google IDP configuration is active",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"config": {
				Description: "Google IDP configuration",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Description: "Google application id",
							Type:        schema.TypeString,
							Required:    true,
						},
						"application_secret": {
							Description: "Google application secret",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"redirect_url": {
				Description: "Paste the URI into the Authorized redirect URIs field in the Google Developer Console",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceIBMAppIDIDPGoogleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Id()

	gg, resp, err := appIDClient.GetGoogleIDPWithContext(ctx, &appid.GetGoogleIDPOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing Google IDP configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error loading AppID Google IDP: %s\n%s", err, resp)
	}

	d.Set("is_active", *gg.IsActive)

	if gg.RedirectURL != nil {
		d.Set("redirect_url", *gg.RedirectURL)
	}

	if gg.Config != nil {
		if err := d.Set("config", flattenIBMAppIDGoogleIDPConfig(gg.Config)); err != nil {
			return diag.Errorf("Failed setting AppID Google IDP config: %s", err)
		}
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDIDPGoogleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	config := &appid.SetGoogleIDPOptions{
		TenantID: &tenantID,
		IDP: &appid.FacebookGoogleConfigParams{
			IsActive: &isActive,
		},
	}

	if isActive {
		config.IDP.Config = expandAppIDGoogleIDPConfig(d.Get("config").([]interface{}))
	}

	_, resp, err := appIDClient.SetGoogleIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error applying AppID Google IDP configuration: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDIDPGoogleRead(ctx, d, meta)
}

func resourceIBMAppIDIDPGoogleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	config := appIDGoogleIDPConfigDefaults(tenantID)

	_, resp, err := appIDClient.SetGoogleIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error resetting AppID Google IDP configuration: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}

func resourceIBMAppIDIDPGoogleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// since this is configuration we can reuse create method
	return resourceIBMAppIDIDPGoogleCreate(ctx, d, m)
}

func expandAppIDGoogleIDPConfig(cfg []interface{}) *appid.FacebookGoogleConfigParamsConfig {
	config := &appid.FacebookGoogleConfigParamsConfig{}

	if len(cfg) == 0 || cfg[0] == nil {
		return nil
	}

	mCfg := cfg[0].(map[string]interface{})

	config.IDPID = helpers.String(mCfg["application_id"].(string))
	config.Secret = helpers.String(mCfg["application_secret"].(string))

	return config
}

func appIDGoogleIDPConfigDefaults(tenantID string) *appid.SetGoogleIDPOptions {
	return &appid.SetGoogleIDPOptions{
		TenantID: &tenantID,
		IDP: &appid.FacebookGoogleConfigParams{
			IsActive: helpers.Bool(false),
		},
	}
}
