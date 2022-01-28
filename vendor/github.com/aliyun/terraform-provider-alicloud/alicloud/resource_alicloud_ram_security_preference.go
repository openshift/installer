package alicloud

import (
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudRamSecurityPreference() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamSecurityPreferenceCreate,
		Read:   resourceAlicloudRamSecurityPreferenceRead,
		Update: resourceAlicloudRamSecurityPreferenceUpdate,
		Delete: resourceAlicloudRamSecurityPreferenceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"enable_save_mfa_ticket": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"allow_user_to_change_password": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"allow_user_to_manage_access_keys": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"allow_user_to_manage_mfa_devices": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"login_session_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"login_network_masks": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enforce_mfa_for_login": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamSecurityPreferenceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "SetSecurityPreference"
	request := make(map[string]interface{})
	conn, err := client.NewImsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOkExists("allow_user_to_change_password"); ok {
		request["AllowUserToChangePassword"] = v
	}
	if v, ok := d.GetOkExists("allow_user_to_manage_access_keys"); ok {
		request["AllowUserToManageAccessKeys"] = v
	}
	if v, ok := d.GetOkExists("allow_user_to_manage_mfa_devices"); ok {
		request["AllowUserToManageMFADevices"] = v
	}
	if v, ok := d.GetOkExists("enable_save_mfa_ticket"); ok {
		request["EnableSaveMFATicket"] = v
	}
	if v, ok := d.GetOk("login_network_masks"); ok {
		request["LoginNetworkMasks"] = v
	}
	if v, ok := d.GetOk("login_session_duration"); ok {
		request["LoginSessionDuration"] = v
	}
	if v, ok := d.GetOk("enforce_mfa_for_login"); ok {
		request["EnforceMFAForLogin"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ram_security_preference", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId("RamSecurityPreference")

	return resourceAlicloudRamSecurityPreferenceRead(d, meta)
}
func resourceAlicloudRamSecurityPreferenceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ramService := RamService{client}
	object, err := ramService.DescribeRamSecurityPreference(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ram_security_preference ramService.DescribeRamSecurityPreference Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if v, exist := object["AccessKeyPreference"].(map[string]interface{}); exist {
		d.Set("allow_user_to_manage_access_keys", v["AllowUserToManageAccessKeys"])
	}
	if v, exist := object["LoginProfilePreference"].(map[string]interface{}); exist {
		d.Set("allow_user_to_change_password", v["AllowUserToChangePassword"])
		d.Set("enable_save_mfa_ticket", v["EnableSaveMFATicket"])
		d.Set("login_network_masks", v["LoginNetworkMasks"])
		d.Set("login_session_duration", v["LoginSessionDuration"])
		d.Set("enforce_mfa_for_login", v["EnforceMFAForLogin"])
	}

	if v, exist := object["MFAPreference"].(map[string]interface{}); exist {
		d.Set("allow_user_to_manage_mfa_devices", v["AllowUserToManageMFADevices"])
	}

	return nil
}
func resourceAlicloudRamSecurityPreferenceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	conn, err := client.NewImsClient()
	if err != nil {
		return WrapError(err)
	}
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{}

	if d.HasChange("allow_user_to_change_password") {
		update = true
	}
	request["AllowUserToChangePassword"] = d.Get("allow_user_to_change_password")

	if d.HasChange("allow_user_to_manage_access_keys") {
		update = true
	}
	request["AllowUserToManageAccessKeys"] = d.Get("allow_user_to_manage_access_keys")

	if d.HasChange("allow_user_to_manage_mfa_devices") {
		update = true
	}
	request["AllowUserToManageMFADevices"] = d.Get("allow_user_to_manage_mfa_devices")

	if d.HasChange("enable_save_mfa_ticket") {
		update = true
	}
	request["EnableSaveMFATicket"] = d.Get("enable_save_mfa_ticket")

	if d.HasChange("login_network_masks") {
		update = true
	}
	request["LoginNetworkMasks"] = d.Get("login_network_masks")
	if d.HasChange("login_session_duration") {
		update = true
	}
	request["LoginSessionDuration"] = d.Get("login_session_duration")

	if d.HasChange("enforce_mfa_for_login") {
		update = true
	}
	request["EnforceMFAForLogin"] = d.Get("enforce_mfa_for_login")

	if update {
		action := "SetSecurityPreference"
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-08-15"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
	}
	return resourceAlicloudRamSecurityPreferenceRead(d, meta)
}
func resourceAlicloudRamSecurityPreferenceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[WARN] Cannot destroy resourceAlicloudRamSecurityPreference. Terraform will remove this resource from the state file, however resources may remain.")
	return nil
}
