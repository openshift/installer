package cognitoidp

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	"github.com/hashicorp/terraform-provider-aws/internal/flex"
)

// @SDKDataSource("aws_cognito_user_pool_client")
func DataSourceUserPoolClient() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceUserPoolClientRead,

		Schema: map[string]*schema.Schema{
			"access_token_validity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"allowed_oauth_flows": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"allowed_oauth_flows_user_pool_client": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allowed_oauth_scopes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"analytics_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"application_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"application_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_arn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_data_shared": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
			"callback_urls": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"default_redirect_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_token_revocation": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"enable_propagate_additional_user_context_data": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"explicit_auth_flows": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"generate_secret": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"id_token_validity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"logout_urls": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"prevent_user_existence_errors": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"read_attributes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"refresh_token_validity": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"supported_identity_providers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"token_validity_units": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_token": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id_token": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"refresh_token": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"user_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"write_attributes": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceUserPoolClientRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).CognitoIDPConn(ctx)

	clientId := d.Get("client_id").(string)
	d.SetId(clientId)

	userPoolClient, err := FindCognitoUserPoolClientByID(ctx, conn, d.Get("user_pool_id").(string), d.Id())

	if err != nil {
		return sdkdiag.AppendErrorf(diags, "reading Cognito User Pool Client (%s): %s", clientId, err)
	}

	d.Set("user_pool_id", userPoolClient.UserPoolId)
	d.Set("name", userPoolClient.ClientName)
	d.Set("explicit_auth_flows", flex.FlattenStringSet(userPoolClient.ExplicitAuthFlows))
	d.Set("read_attributes", flex.FlattenStringSet(userPoolClient.ReadAttributes))
	d.Set("write_attributes", flex.FlattenStringSet(userPoolClient.WriteAttributes))
	d.Set("refresh_token_validity", userPoolClient.RefreshTokenValidity)
	d.Set("access_token_validity", userPoolClient.AccessTokenValidity)
	d.Set("id_token_validity", userPoolClient.IdTokenValidity)
	d.Set("client_secret", userPoolClient.ClientSecret)
	d.Set("allowed_oauth_flows", flex.FlattenStringSet(userPoolClient.AllowedOAuthFlows))
	d.Set("allowed_oauth_flows_user_pool_client", userPoolClient.AllowedOAuthFlowsUserPoolClient)
	d.Set("allowed_oauth_scopes", flex.FlattenStringSet(userPoolClient.AllowedOAuthScopes))
	d.Set("callback_urls", flex.FlattenStringSet(userPoolClient.CallbackURLs))
	d.Set("default_redirect_uri", userPoolClient.DefaultRedirectURI)
	d.Set("logout_urls", flex.FlattenStringSet(userPoolClient.LogoutURLs))
	d.Set("prevent_user_existence_errors", userPoolClient.PreventUserExistenceErrors)
	d.Set("supported_identity_providers", flex.FlattenStringSet(userPoolClient.SupportedIdentityProviders))
	d.Set("enable_token_revocation", userPoolClient.EnableTokenRevocation)
	d.Set("enable_propagate_additional_user_context_data", userPoolClient.EnablePropagateAdditionalUserContextData)

	if err := d.Set("analytics_configuration", flattenUserPoolClientAnalyticsConfig(userPoolClient.AnalyticsConfiguration)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting analytics_configuration: %s", err)
	}

	if err := d.Set("token_validity_units", flattenUserPoolClientTokenValidityUnitsType(userPoolClient.TokenValidityUnits)); err != nil {
		return sdkdiag.AppendErrorf(diags, "setting token_validity_units: %s", err)
	}

	return diags
}

func flattenUserPoolClientAnalyticsConfig(analyticsConfig *cognitoidentityprovider.AnalyticsConfigurationType) []interface{} {
	if analyticsConfig == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"user_data_shared": aws.BoolValue(analyticsConfig.UserDataShared),
	}

	if analyticsConfig.ExternalId != nil {
		m["external_id"] = aws.StringValue(analyticsConfig.ExternalId)
	}

	if analyticsConfig.RoleArn != nil {
		m["role_arn"] = aws.StringValue(analyticsConfig.RoleArn)
	}

	if analyticsConfig.ApplicationId != nil {
		m["application_id"] = aws.StringValue(analyticsConfig.ApplicationId)
	}

	if analyticsConfig.ApplicationArn != nil {
		m["application_arn"] = aws.StringValue(analyticsConfig.ApplicationArn)
	}

	return []interface{}{m}
}

func flattenUserPoolClientTokenValidityUnitsType(tokenValidityConfig *cognitoidentityprovider.TokenValidityUnitsType) []interface{} {
	if tokenValidityConfig == nil {
		return nil
	}

	//tokenValidityConfig is never nil and if everything is empty it causes diffs
	if tokenValidityConfig.IdToken == nil && tokenValidityConfig.AccessToken == nil && tokenValidityConfig.RefreshToken == nil {
		return nil
	}

	m := map[string]interface{}{}

	if tokenValidityConfig.IdToken != nil {
		m["id_token"] = aws.StringValue(tokenValidityConfig.IdToken)
	}

	if tokenValidityConfig.AccessToken != nil {
		m["access_token"] = aws.StringValue(tokenValidityConfig.AccessToken)
	}

	if tokenValidityConfig.RefreshToken != nil {
		m["refresh_token"] = aws.StringValue(tokenValidityConfig.RefreshToken)
	}

	return []interface{}{m}
}
