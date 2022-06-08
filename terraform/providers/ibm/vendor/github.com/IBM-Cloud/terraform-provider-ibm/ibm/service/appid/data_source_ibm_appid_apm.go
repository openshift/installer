package appid

import (
	"context"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMAppIDAPM() *schema.Resource {
	return &schema.Resource{
		Description: "AppID advanced password management configuration (available for graduated tier only)",
		ReadContext: dataSourceIBMAppIDAPMRead,
		Schema: map[string]*schema.Schema{
			"tenant_id": {
				Description: "The AppID instance GUID",
				Type:        schema.TypeString,
				Required:    true,
			},
			"enabled": {
				Description: "`true` if APM is enabled",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"prevent_password_with_username": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"password_reuse": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max_password_reuse": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"password_expiration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"days_to_expire": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"lockout_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"lockout_time_sec": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"num_of_attempts": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"min_password_change_interval": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"min_hours_to_change_password": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMAppIDAPMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)

	apm, resp, err := appIDClient.GetCloudDirectoryAdvancedPasswordManagementWithContext(ctx, &appid.GetCloudDirectoryAdvancedPasswordManagementOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		return diag.Errorf("Error getting AppID APM configuration: %s\n%s", err, resp)
	}

	if apm.AdvancedPasswordManagement != nil {
		d.Set("enabled", *apm.AdvancedPasswordManagement.Enabled)

		if err := d.Set("password_reuse", flattenAppIDAPMPasswordReuse(apm.AdvancedPasswordManagement.PasswordReuse)); err != nil {
			return diag.Errorf("Failed setting AppID APM password_reuse: %s", err)
		}

		if apm.AdvancedPasswordManagement.PreventPasswordWithUsername != nil {
			d.Set("prevent_password_with_username", *apm.AdvancedPasswordManagement.PreventPasswordWithUsername.Enabled)
		}

		if err := d.Set("password_expiration", flattenAppIDAPMPasswordExpiration(apm.AdvancedPasswordManagement.PasswordExpiration)); err != nil {
			return diag.Errorf("Failed setting AppID APM password_expiration: %s", err)
		}

		if err := d.Set("lockout_policy", flattenAppIDAPMLockoutPolicy(apm.AdvancedPasswordManagement.LockOutPolicy)); err != nil {
			return diag.Errorf("Failed setting AppID APM lockout_policy: %s", err)
		}
		if err := d.Set("min_password_change_interval", flattenAppIDAPMPasswordChangeInterval(apm.AdvancedPasswordManagement.MinPasswordChangeInterval)); err != nil {
			return diag.Errorf("Failed setting AppID APM min_password_change_interval: %s", err)
		}

	}

	d.SetId(tenantID)
	return nil
}

func flattenAppIDAPMPasswordReuse(reuse *appid.ApmSchemaAdvancedPasswordManagementPasswordReuse) []interface{} {
	if reuse == nil {
		return []interface{}{}
	}

	mReuse := map[string]interface{}{}

	mReuse["enabled"] = *reuse.Enabled

	if reuse.Config != nil && reuse.Config.MaxPasswordReuse != nil {
		mReuse["max_password_reuse"] = *reuse.Config.MaxPasswordReuse
	}

	return []interface{}{mReuse}
}

func flattenAppIDAPMPasswordExpiration(exp *appid.ApmSchemaAdvancedPasswordManagementPasswordExpiration) []interface{} {
	if exp == nil {
		return []interface{}{}
	}

	mExp := map[string]interface{}{}

	mExp["enabled"] = *exp.Enabled

	if exp.Config != nil && exp.Config.DaysToExpire != nil {
		mExp["days_to_expire"] = *exp.Config.DaysToExpire
	}

	return []interface{}{mExp}
}

func flattenAppIDAPMLockoutPolicy(pol *appid.ApmSchemaAdvancedPasswordManagementLockOutPolicy) []interface{} {
	if pol == nil {
		return []interface{}{}
	}

	mPol := map[string]interface{}{}

	mPol["enabled"] = *pol.Enabled

	if pol.Config != nil && pol.Config.LockOutTimeSec != nil {
		mPol["lockout_time_sec"] = *pol.Config.LockOutTimeSec
	}

	if pol.Config != nil && pol.Config.NumOfAttempts != nil {
		mPol["num_of_attempts"] = *pol.Config.NumOfAttempts
	}

	return []interface{}{mPol}
}

func flattenAppIDAPMPasswordChangeInterval(in *appid.ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval) []interface{} {
	if in == nil {
		return []interface{}{}
	}

	mIn := map[string]interface{}{}

	mIn["enabled"] = *in.Enabled

	if in.Config != nil && in.Config.MinHoursToChangePassword != nil {
		mIn["min_hours_to_change_password"] = *in.Config.MinHoursToChangePassword
	}

	return []interface{}{mIn}
}
