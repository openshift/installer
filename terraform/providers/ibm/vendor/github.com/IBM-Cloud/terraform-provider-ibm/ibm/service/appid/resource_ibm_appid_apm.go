package appid

import (
	"context"
	"log"

	"github.com/IBM-Cloud/bluemix-go/helpers"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	appid "github.com/IBM/appid-management-go-sdk/appidmanagementv4"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMAppIDAPM() *schema.Resource {
	return &schema.Resource{
		Description:   "AppID advanced password management configuration (available for graduated tier only)",
		ReadContext:   resourceIBMAppIDAPMRead,
		CreateContext: resourceIBMAppIDAPMCreate,
		UpdateContext: resourceIBMAppIDAPMCreate,
		DeleteContext: resourceIBMAppIDAPMDelete,
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
			"enabled": {
				Description: "`true` if APM is enabled",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"prevent_password_with_username": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"password_reuse": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"max_password_reuse": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  8,
						},
					},
				},
			},
			"password_expiration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"days_to_expire": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  30,
						},
					},
				},
			},
			"lockout_policy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"lockout_time_sec": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1800,
						},
						"num_of_attempts": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
					},
				},
			},
			"min_password_change_interval": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"min_hours_to_change_password": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
					},
				},
			},
		},
	}
}

func resourceIBMAppIDAPMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Id()

	apm, resp, err := appIDClient.GetCloudDirectoryAdvancedPasswordManagementWithContext(ctx, &appid.GetCloudDirectoryAdvancedPasswordManagementOptions{
		TenantID: &tenantID,
	})

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[WARN] AppID instance '%s' is not found, removing APM configuration from state", tenantID)
			d.SetId("")
			return nil
		}

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

	d.Set("tenant_id", tenantID)
	return nil
}

func resourceIBMAppIDAPMCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	enabled := d.Get("enabled").(bool)

	config := &appid.SetCloudDirectoryAdvancedPasswordManagementOptions{
		TenantID: &tenantID,
		AdvancedPasswordManagement: &appid.ApmSchemaAdvancedPasswordManagement{
			Enabled:                   &enabled,
			PasswordReuse:             expandAppIDAPMPasswordReuse(d.Get("password_reuse").([]interface{})),
			PasswordExpiration:        expandAppIDAPMPasswordExpiration(d.Get("password_expiration").([]interface{})),
			LockOutPolicy:             expandAppIDAPMLockoutPolicy(d.Get("lockout_policy").([]interface{})),
			MinPasswordChangeInterval: expandAppIDAPMMinPasswordChangeInterval(d.Get("min_password_change_interval").([]interface{})),
			PreventPasswordWithUsername: &appid.ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername{
				Enabled: helpers.Bool(d.Get("prevent_password_with_username").(bool)),
			},
		},
	}

	_, resp, err := appIDClient.SetCloudDirectoryAdvancedPasswordManagementWithContext(ctx, config)

	if err != nil {
		return diag.Errorf("Error updating AppID APM configuration: %s\n%s", err, resp)
	}

	d.SetId(tenantID)
	return resourceIBMAppIDAPMRead(ctx, d, meta)
}

func expandAppIDAPMPasswordReuse(reuse []interface{}) *appid.ApmSchemaAdvancedPasswordManagementPasswordReuse {
	if len(reuse) == 0 || reuse[0] == nil {
		return nil
	}

	mReuse := reuse[0].(map[string]interface{})

	result := &appid.ApmSchemaAdvancedPasswordManagementPasswordReuse{
		Enabled: helpers.Bool(mReuse["enabled"].(bool)),
		Config: &appid.ApmSchemaAdvancedPasswordManagementPasswordReuseConfig{
			MaxPasswordReuse: core.Int64Ptr(int64(mReuse["max_password_reuse"].(int))),
		},
	}

	return result
}

func expandAppIDAPMPasswordExpiration(exp []interface{}) *appid.ApmSchemaAdvancedPasswordManagementPasswordExpiration {
	if len(exp) == 0 || exp[0] == nil {
		return nil
	}

	mExp := exp[0].(map[string]interface{})

	result := &appid.ApmSchemaAdvancedPasswordManagementPasswordExpiration{
		Enabled: helpers.Bool(mExp["enabled"].(bool)),
		Config: &appid.ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig{
			DaysToExpire: core.Int64Ptr(int64(mExp["days_to_expire"].(int))),
		},
	}

	return result
}

func expandAppIDAPMLockoutPolicy(loc []interface{}) *appid.ApmSchemaAdvancedPasswordManagementLockOutPolicy {
	if len(loc) == 0 || loc[0] == nil {
		return nil
	}

	mLock := loc[0].(map[string]interface{})

	result := &appid.ApmSchemaAdvancedPasswordManagementLockOutPolicy{
		Enabled: helpers.Bool(mLock["enabled"].(bool)),
		Config: &appid.ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig{
			LockOutTimeSec: core.Int64Ptr(int64(mLock["lockout_time_sec"].(int))),
			NumOfAttempts:  core.Int64Ptr(int64(mLock["num_of_attempts"].(int))),
		},
	}

	return result
}

func expandAppIDAPMMinPasswordChangeInterval(chg []interface{}) *appid.ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval {
	if len(chg) == 0 || chg[0] == nil {
		return nil
	}

	mChg := chg[0].(map[string]interface{})

	result := &appid.ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval{
		Enabled: helpers.Bool(mChg["enabled"].(bool)),
		Config: &appid.ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig{
			MinHoursToChangePassword: core.Int64Ptr(int64(mChg["min_hours_to_change_password"].(int))),
		},
	}

	return result
}

func resourceIBMAppIDAPMDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	appIDClient, err := meta.(conns.ClientSession).AppIDAPI()

	if err != nil {
		return diag.FromErr(err)
	}

	tenantID := d.Get("tenant_id").(string)
	config := getDefaultAppIDAPMConfig()

	_, resp, err := appIDClient.SetCloudDirectoryAdvancedPasswordManagementWithContext(ctx, &appid.SetCloudDirectoryAdvancedPasswordManagementOptions{
		TenantID:                   &tenantID,
		AdvancedPasswordManagement: config,
	})

	if err != nil {
		return diag.Errorf("Error resetting AppID APM configuration: %s\n%s", err, resp)
	}

	d.SetId("")

	return nil
}

func getDefaultAppIDAPMConfig() *appid.ApmSchemaAdvancedPasswordManagement {
	return &appid.ApmSchemaAdvancedPasswordManagement{
		Enabled: helpers.Bool(false),
		PasswordReuse: &appid.ApmSchemaAdvancedPasswordManagementPasswordReuse{
			Enabled: helpers.Bool(false),
			Config: &appid.ApmSchemaAdvancedPasswordManagementPasswordReuseConfig{
				MaxPasswordReuse: core.Int64Ptr(8),
			},
		},
		PasswordExpiration: &appid.ApmSchemaAdvancedPasswordManagementPasswordExpiration{
			Enabled: helpers.Bool(false),
			Config: &appid.ApmSchemaAdvancedPasswordManagementPasswordExpirationConfig{
				DaysToExpire: core.Int64Ptr(30),
			},
		},
		MinPasswordChangeInterval: &appid.ApmSchemaAdvancedPasswordManagementMinPasswordChangeInterval{
			Enabled: helpers.Bool(false),
			Config: &appid.ApmSchemaAdvancedPasswordManagementMinPasswordChangeIntervalConfig{
				MinHoursToChangePassword: core.Int64Ptr(0),
			},
		},
		LockOutPolicy: &appid.ApmSchemaAdvancedPasswordManagementLockOutPolicy{
			Enabled: helpers.Bool(false),
			Config: &appid.ApmSchemaAdvancedPasswordManagementLockOutPolicyConfig{
				LockOutTimeSec: core.Int64Ptr(1800),
				NumOfAttempts:  core.Int64Ptr(3),
			},
		},
		PreventPasswordWithUsername: &appid.ApmSchemaAdvancedPasswordManagementPreventPasswordWithUsername{
			Enabled: helpers.Bool(false),
		},
	}
}
