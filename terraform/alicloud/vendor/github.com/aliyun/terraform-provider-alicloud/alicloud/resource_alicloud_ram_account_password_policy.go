package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
)

// defaults from API docs https://help.aliyun.com/document_detail/28739.html
// when resource is created it sets not specifies value in the resource block to defaults
// also during deletion it rollbacks changes to defaults (API has only a Set method)
var (
	default_minimum_password_length      = 12
	default_require_lowercase_characters = true
	default_require_uppercase_characters = true
	default_require_numbers              = true
	default_require_symbols              = true
	default_hard_expiry                  = false
	default_max_password_age             = 0 // means disable
	default_password_reuse_prevention    = 0 // means disable
	default_max_login_attempts           = 5
)

func resourceAlicloudRamAccountPasswordPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamAccountPasswordPolicyUpdate,
		Read:   resourceAlicloudRamAccountPasswordPolicyRead,
		Update: resourceAlicloudRamAccountPasswordPolicyUpdate,
		Delete: resourceAlicloudRamAccountPasswordPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"minimum_password_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      default_minimum_password_length,
				ValidateFunc: intBetween(8, 32),
			},
			"require_lowercase_characters": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  default_require_lowercase_characters,
			},
			"require_uppercase_characters": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  default_require_uppercase_characters,
			},
			"require_numbers": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  default_require_numbers,
			},
			"require_symbols": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  default_require_symbols,
			},
			"hard_expiry": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  default_hard_expiry,
			},
			"max_password_age": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      default_max_password_age,
				ValidateFunc: intBetween(0, 1095),
			},
			"password_reuse_prevention": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      default_password_reuse_prevention,
				ValidateFunc: intBetween(0, 24),
			},
			"max_login_attempts": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      default_max_login_attempts,
				ValidateFunc: intBetween(0, 32),
			},
		},
	}
}

func resourceAlicloudRamAccountPasswordPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateSetPasswordPolicyRequest()
	request.RegionId = client.RegionId
	request.MinimumPasswordLength = requests.NewInteger(d.Get("minimum_password_length").(int))
	request.RequireLowercaseCharacters = requests.NewBoolean(d.Get("require_lowercase_characters").(bool))
	request.RequireUppercaseCharacters = requests.NewBoolean(d.Get("require_uppercase_characters").(bool))
	request.RequireNumbers = requests.NewBoolean(d.Get("require_numbers").(bool))
	request.RequireSymbols = requests.NewBoolean(d.Get("require_symbols").(bool))
	request.MaxLoginAttemps = requests.NewInteger(d.Get("max_login_attempts").(int))
	request.HardExpiry = requests.NewBoolean(d.Get("hard_expiry").(bool))
	request.MaxPasswordAge = requests.NewInteger(d.Get("max_password_age").(int))
	request.PasswordReusePrevention = requests.NewInteger(d.Get("password_reuse_prevention").(int))
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.SetPasswordPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_account_password_policy", request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	d.SetId("ram-account-password-policy")

	return resourceAlicloudRamAccountPasswordPolicyRead(d, meta)
}

func resourceAlicloudRamAccountPasswordPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	object, err := ramService.DescribeRamAccountPasswordPolicy(d.Id())
	if err != nil {
		return WrapError(err)
	}
	passwordPolicy := object.PasswordPolicy
	d.Set("minimum_password_length", passwordPolicy.MinimumPasswordLength)
	d.Set("require_lowercase_characters", passwordPolicy.RequireLowercaseCharacters)
	d.Set("require_uppercase_characters", passwordPolicy.RequireUppercaseCharacters)
	d.Set("require_numbers", passwordPolicy.RequireNumbers)
	d.Set("require_symbols", passwordPolicy.RequireSymbols)
	d.Set("hard_expiry", passwordPolicy.HardExpiry)
	d.Set("max_password_age", passwordPolicy.MaxPasswordAge)
	d.Set("password_reuse_prevention", passwordPolicy.PasswordReusePrevention)
	d.Set("max_login_attempts", passwordPolicy.MaxLoginAttemps)
	return nil
}

func resourceAlicloudRamAccountPasswordPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	request := ram.CreateSetPasswordPolicyRequest()
	request.RegionId = client.RegionId
	request.MinimumPasswordLength = requests.NewInteger(default_minimum_password_length)
	request.RequireLowercaseCharacters = requests.NewBoolean(default_require_lowercase_characters)
	request.RequireUppercaseCharacters = requests.NewBoolean(default_require_uppercase_characters)
	request.RequireNumbers = requests.NewBoolean(default_require_numbers)
	request.RequireSymbols = requests.NewBoolean(default_require_symbols)
	request.HardExpiry = requests.NewBoolean(default_hard_expiry)
	request.MaxPasswordAge = requests.NewInteger(default_max_password_age)
	request.PasswordReusePrevention = requests.NewInteger(default_password_reuse_prevention)
	request.MaxLoginAttemps = requests.NewInteger(default_max_login_attempts)
	raw, err := client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.SetPasswordPolicy(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), request.GetActionName(), AlibabaCloudSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return nil
}
