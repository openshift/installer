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

func ResourceIBMAppIDIDPFacebook() *schema.Resource {
	return &schema.Resource{
		Description:   "Update Facebook identity provider configuration.",
		CreateContext: resourceIBMAppIDIDPFacebookCreate,
		ReadContext:   resourceIBMAppIDIDPFacebookRead,
		DeleteContext: resourceIBMAppIDIDPFacebookDelete,
		UpdateContext: resourceIBMAppIDIDPFacebookUpdate,
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
				Description: "`true` if Facebook IDP configuration is active",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"config": {
				Description: "Facebook IDP configuration",
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Description: "Facebook application id",
							Type:        schema.TypeString,
							Required:    true,
						},
						"application_secret": {
							Description: "Facebook application secret",
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"redirect_url": {
				Description: "Paste the URI into the Valid OAuth redirect URIs field in the Facebook Login section of the Facebook Developers Portal",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceIBMAppIDIDPFacebookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Id()

	fb, resp, err := appIDClient.GetFacebookIDPWithContext(ctx, &appid.GetFacebookIDPOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing Facebook IDP configuration from state", tenantID)
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error loading AppID Facebook IDP: %s\n%s", err, resp)
	}

	d.Set("is_active", *fb.IsActive)

	if fb.RedirectURL != nil {
		d.Set("redirect_url", *fb.RedirectURL)
	}

	if fb.Config != nil {
		if err := d.Set("config", flattenIBMAppIDFacebookIDPConfig(fb.Config)); err != nil {
			return diag.Errorf("Failed setting AppID Facebook IDP config: %s", err)
		}
	}

	d.Set("tenant_id", tenantID)

	return nil
}

func resourceIBMAppIDIDPFacebookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	isActive := d.Get("is_active").(bool)

	config := &appid.SetFacebookIDPOptions{
		TenantID: &tenantID,
		IDP: &appid.FacebookGoogleConfigParams{
			IsActive: &isActive,
		},
	}

	if isActive {
		config.IDP.Config = expandAppIDFBIDPConfig(d.Get("config").([]interface{}))
	}

	_, resp, err := appIDClient.SetFacebookIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error applying AppID Facebook IDP configuration: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDIDPFacebookRead(ctx, d, meta)
}

func resourceIBMAppIDIDPFacebookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	config := appIDFacebookIDPConfigDefaults(tenantID)

	_, resp, err := appIDClient.SetFacebookIDPWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error resetting AppID Facebook IDP configuration: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}

func resourceIBMAppIDIDPFacebookUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// since this is configuration we can reuse create method
	return resourceIBMAppIDIDPFacebookCreate(ctx, d, m)
}

func expandAppIDFBIDPConfig(cfg []interface{}) *appid.FacebookGoogleConfigParamsConfig {
	config := &appid.FacebookGoogleConfigParamsConfig{}

	if len(cfg) == 0 || cfg[0] == nil {
		return nil
	}

	mCfg := cfg[0].(map[string]interface{})

	config.IDPID = helpers.String(mCfg["application_id"].(string))
	config.Secret = helpers.String(mCfg["application_secret"].(string))

	return config
}

func appIDFacebookIDPConfigDefaults(tenantID string) *appid.SetFacebookIDPOptions {
	return &appid.SetFacebookIDPOptions{
		TenantID: &tenantID,
		IDP: &appid.FacebookGoogleConfigParams{
			IsActive: helpers.Bool(false),
		},
	}
}
