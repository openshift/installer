// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package appid

import (
	"context"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIBMAppIDTokenConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMAppIDTokenConfigCreate,
		ReadContext:   resourceIBMAppIDTokenConfigRead,
		UpdateContext: resourceIBMAppIDTokenConfigUpdate,
		DeleteContext: resourceIBMAppIDTokenConfigDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The service `tenantId`",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"access_token_expires_in": {
				Description: "The length of time for which access tokens are valid in seconds",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
			"refresh_token_expires_in": {
				Description: "The length of time for which refresh tokens are valid in seconds",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     2592000,
			},
			"anonymous_token_expires_in": {
				Type:     schema.TypeInt,
				Default:  2592000,
				Optional: true,
			},
			"anonymous_access_enabled": {
				Description: "The length of time for which an anonymous token is valid in seconds",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"refresh_token_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"access_token_claim": {
				Description: "A set of objects that are created when claims that are related to access tokens are mapped",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Description:  "Defines the source of the claim. Options include: `saml`, `cloud_directory`, `facebook`, `google`, `appid_custom`, and `attributes`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"saml", "cloud_directory", "appid_custom", "facebook", "google", "ibmid", "attributes", "roles"}, false),
						},
						"source_claim": {
							Description: "Defines the claim as provided by the source. It can refer to the identity provider's user information or the user's App ID custom attributes.",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"destination_claim": {
							Description: "Optional: Defines the custom attribute that can override the current claim in token.",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"id_token_claim": {
				Description: "A set of objects that are created when claims that are related to identity tokens are mapped",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"saml", "cloud_directory", "appid_custom", "facebook", "google", "ibmid", "attributes", "roles"}, false),
						},
						"source_claim": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"destination_claim": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceIBMAppIDTokenConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appidClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	input := expandTokenConfig(d)

	_, resp, err := appidClient.PutTokensConfigWithContext(ctx, input)

	if err != nil {
		return diag.Errorf("Error updating AppID token configuration: %s\n%s", err, resp)
	}

	d.SetId(tenantID)

	return resourceIBMAppIDTokenConfigRead(ctx, d, meta)
}

func resourceIBMAppIDTokenConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	appidClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Id()

	tokenConfig, response, err := appidClient.GetTokensConfigWithContext(ctx, &appid.GetTokensConfigOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error reading AppID token configuration: %s\n%s", err, response)
	}

	if tokenConfig.Access != nil {
		d.Set("access_token_expires_in", *tokenConfig.Access.ExpiresIn)
	}

	if tokenConfig.Refresh != nil {
		if tokenConfig.Refresh.Enabled != nil {
			d.Set("refresh_token_enabled", *tokenConfig.Refresh.Enabled)
		} else {
			d.Set("refresh_token_enabled", nil)
		}

		d.Set("refresh_token_expires_in", *tokenConfig.Refresh.ExpiresIn)
	}

	if tokenConfig.AnonymousAccess != nil {
		if tokenConfig.AnonymousAccess.Enabled != nil {
			d.Set("anonymous_access_enabled", *tokenConfig.AnonymousAccess.Enabled)
		} else {
			d.Set("anonymous_access_enabled", nil)
		}

		d.Set("anonymous_token_expires_in", *tokenConfig.AnonymousAccess.ExpiresIn)
	}

	if tokenConfig.AccessTokenClaims != nil {
		if err := d.Set("access_token_claim", flattenTokenClaims(tokenConfig.AccessTokenClaims)); err != nil {
			return diag.FromErr(err)
		}
	}

	if tokenConfig.IDTokenClaims != nil {
		if err := d.Set("id_token_claim", flattenTokenClaims(tokenConfig.IDTokenClaims)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.Set("tenant_id", tenantID)

	return diags
}

func expandTokenClaims(l []interface{}) []appid.TokenClaimMapping {
	if len(l) == 0 {
		return nil
	}

	result := make([]appid.TokenClaimMapping, len(l))

	for i, item := range l {
		cMap := item.(map[string]interface{})

		claim := appid.TokenClaimMapping{
			Source: helpers.String(cMap["source"].(string)),
		}

		// source_claim and destination_claim are optional
		if sClaim, ok := cMap["source_claim"]; ok {
			claim.SourceClaim = helpers.String(sClaim.(string))
		}

		if dClaim, ok := cMap["destination_claim"]; ok {
			claim.DestinationClaim = helpers.String(dClaim.(string))
		}

		result[i] = claim
	}

	return result
}

func expandTokenConfig(d *schema.ResourceData) *appid.PutTokensConfigOptions {
	config := &appid.PutTokensConfigOptions{
		TenantID: helpers.String(d.Get("tenant_id").(string)),
	}

	if accessExpiresIn, ok := d.GetOk("access_token_expires_in"); ok {
		config.Access = &appid.AccessTokenConfigParams{
			ExpiresIn: core.Int64Ptr(int64(accessExpiresIn.(int))),
		}
	}

	if anonymousExpiresIn, ok := d.GetOk("anonymous_token_expires_in"); ok {
		config.AnonymousAccess = &appid.TokenConfigParams{
			ExpiresIn: core.Int64Ptr(int64(anonymousExpiresIn.(int))),
		}
	}

	if refreshExpiresIn, ok := d.GetOk("refresh_token_expires_in"); ok {
		config.Refresh = &appid.TokenConfigParams{
			ExpiresIn: core.Int64Ptr(int64(refreshExpiresIn.(int))),
		}
	}

	// can't really use GetOk with bool
	anonymousAccessEnabled := d.Get("anonymous_access_enabled")

	if anonymousAccessEnabled != nil {
		if config.AnonymousAccess == nil {
			config.AnonymousAccess = &appid.TokenConfigParams{}
		}

		config.AnonymousAccess.Enabled = helpers.Bool(anonymousAccessEnabled.(bool))
	}

	refreshTokenEnabled := d.Get("refresh_token_enabled")

	if refreshTokenEnabled != nil {
		if config.Refresh == nil {
			config.Refresh = &appid.TokenConfigParams{}
		}

		config.Refresh.Enabled = helpers.Bool(refreshTokenEnabled.(bool))
	}

	if accessClaims, ok := d.GetOk("access_token_claim"); ok {
		config.AccessTokenClaims = expandTokenClaims(accessClaims.(*schema.Set).List())
	}

	if idClaims, ok := d.GetOk("id_token_claim"); ok {
		config.IDTokenClaims = expandTokenClaims(idClaims.(*schema.Set).List())
	}

	return config
}

func resourceIBMAppIDTokenConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceIBMAppIDTokenConfigCreate(ctx, d, m)
}

func tokenConfigDefaults(tenantID string) *appid.PutTokensConfigOptions {
	return &appid.PutTokensConfigOptions{
		TenantID: helpers.String(tenantID),
		Access: &appid.AccessTokenConfigParams{
			ExpiresIn: core.Int64Ptr(3600),
		},
		Refresh: &appid.TokenConfigParams{
			Enabled:   helpers.Bool(false),
			ExpiresIn: core.Int64Ptr(2592000),
		},
		AnonymousAccess: &appid.TokenConfigParams{
			Enabled:   helpers.Bool(true),
			ExpiresIn: core.Int64Ptr(2592000),
		},
	}
}

func resourceIBMAppIDTokenConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appidClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	config := tokenConfigDefaults(tenantID)
	_, resp, err := appidClient.PutTokensConfigWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error resetting AppID token configuration: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}
