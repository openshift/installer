package appid

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDIDPSAML() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMAppIDIDPSAMLRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_active": {
				Description: "SAML IDP activation",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"config": {
				Description: "SAML IDP configuration",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"entity_id": {
							Description: "Unique name for an Identity Provider",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"sign_in_url": {
							Description: "SAML SSO url",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"certificates": {
							Description: "List of certificates, primary and optional secondary",
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Computed: true,
						},
						"display_name": {
							Description: "Provider name",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"encrypt_response": {
							Description: "`true` if SAML responses should be encrypted",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"sign_request": {
							Description: "`true` if SAML requests should be signed",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"include_scoping": {
							Description: "`true` if scopes are included",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"authn_context": {
							Description: "SAML authNContext configuration",
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"class": {
										Description: "List of `authnContext` classes",
										Type:        schema.TypeList,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed: true,
									},
									"comparison": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAppIDIDPSAMLRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	saml, resp, err := appIDClient.GetSAMLIDPWithContext(ctx, &appid.GetSAMLIDPOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error loading SAML IDP: %s\n%s", err, resp)
	}

	d.Set("is_active", *saml.IsActive)

	if saml.Config != nil {
		if err := d.Set("config", flattenAppIDIDPSAMLConfig(saml.Config)); err != nil {
			return diag.Errorf("Failed setting AppID IDP SAML config: %s", err)
		}
	}

	d.SetId(tenantID)

	return nil
}

func flattenAppIDIDPSAMLConfig(config *appid.SAMLConfigParams) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	mConfig := map[string]interface{}{}

	if config.EntityID != nil {
		mConfig["entity_id"] = *config.EntityID
	}

	if config.SignInURL != nil {
		mConfig["sign_in_url"] = *config.SignInURL
	}

	mConfig["certificates"] = flex.FlattenStringList(config.Certificates)

	if config.DisplayName != nil {
		mConfig["display_name"] = *config.DisplayName
	}

	if config.SignRequest != nil {
		mConfig["sign_request"] = *config.SignRequest
	}

	if config.EncryptResponse != nil {
		mConfig["encrypt_response"] = *config.EncryptResponse
	}

	if config.IncludeScoping != nil {
		mConfig["include_scoping"] = *config.IncludeScoping
	}

	if config.AuthnContext != nil {
		mConfig["authn_context"] = flattenAuthNContext(config.AuthnContext)
	}

	return []interface{}{mConfig}
}

func flattenAuthNContext(context *appid.SAMLConfigParamsAuthnContext) []interface{} {
	if context == nil {
		return []interface{}{}
	}

	mContext := map[string]interface{}{}

	if context.Class != nil {
		var class []interface{}

		for _, c := range context.Class {
			class = append(class, c)
		}

		mContext["class"] = class
	}

	if context.Comparison != nil {
		mContext["comparison"] = *context.Comparison
	}

	return []interface{}{mContext}
}
