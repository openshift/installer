package alicloud

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceAlicloudEcdSimpleOfficeSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudEcdSimpleOfficeSiteCreate,
		Read:   resourceAlicloudEcdSimpleOfficeSiteRead,
		Update: resourceAlicloudEcdSimpleOfficeSiteUpdate,
		Delete: resourceAlicloudEcdSimpleOfficeSiteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bandwidth": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 200),
				Deprecated:   "Field 'bandwidth' has been deprecated from provider version 1.142.0.",
			},
			"cen_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cen_owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"desktop_access_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{"Any", "Internet", "VPC"}, false),
			},
			"enable_admin_access": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"enable_cross_desktop_access": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"enable_internet_access": {
				Type:       schema.TypeBool,
				Computed:   true,
				Optional:   true,
				Deprecated: "Field 'enable_internet_access' has been deprecated from provider version 1.142.0.",
			},
			"mfa_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"office_site_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sso_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudEcdSimpleOfficeSiteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateSimpleOfficeSite"
	request := make(map[string]interface{})
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	request["CidrBlock"] = d.Get("cidr_block")
	request["RegionId"] = client.RegionId
	if v, ok := d.GetOkExists("enable_admin_access"); ok {
		request["EnableAdminAccess"] = v
	}
	if v, ok := d.GetOk("cen_owner_id"); ok {
		request["CenOwnerId"] = v
	}
	if v, ok := d.GetOk("cen_id"); ok {
		request["CenId"] = v
	}
	if v, ok := d.GetOk("office_site_name"); ok {
		request["OfficeSiteName"] = v
	}

	// todo: you can configure it using alicloud_ecd_network_package
	if v, ok := d.GetOkExists("enable_internet_access"); ok {
		request["EnableInternetAccess"] = v
	}
	if v, ok := d.GetOk("bandwidth"); ok {
		request["Bandwidth"] = v
	}

	if v, ok := d.GetOk("desktop_access_type"); ok {
		request["DesktopAccessType"] = v
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ecd_simple_office_site", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["OfficeSiteId"]))
	ecdService := EcdService{client}
	stateConf := BuildStateConf([]string{}, []string{"REGISTERED"}, d.Timeout(schema.TimeoutCreate), 5*time.Second, ecdService.EcdSimpleOfficeSiteStateRefreshFunc(d.Id()))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	return resourceAlicloudEcdSimpleOfficeSiteUpdate(d, meta)
}
func resourceAlicloudEcdSimpleOfficeSiteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	object, err := ecdService.DescribeEcdSimpleOfficeSite(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_simple_office_site ecdService.DescribeEcdSimpleOfficeSite Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	// todo:  bandwidth depends on network_package resource， you can find it in alicloud_ecd_network_package
	if v, ok := object["Bandwidth"]; ok && fmt.Sprint(v) != "0" {
		d.Set("bandwidth", formatInt(v))
	}
	d.Set("cen_id", object["CenId"])
	d.Set("cidr_block", object["CidrBlock"])
	d.Set("desktop_access_type", convertDesktopAccessType(object["DesktopAccessType"]))
	d.Set("enable_admin_access", object["EnableAdminAccess"])
	d.Set("enable_cross_desktop_access", object["EnableCrossDesktopAccess"])
	d.Set("office_site_name", object["Name"])
	d.Set("mfa_enabled", object["MfaEnabled"])
	d.Set("enable_internet_access", object["EnableInternetAccess"])
	d.Set("sso_enabled", object["SsoEnabled"])
	d.Set("status", object["Status"])
	return nil
}
func resourceAlicloudEcdSimpleOfficeSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	d.Partial(true)
	update := false
	request := map[string]interface{}{
		"OfficeSiteId": d.Id(),
	}
	if d.HasChange("enable_cross_desktop_access") {
		update = true
		if v, ok := d.GetOkExists("enable_cross_desktop_access"); ok {
			request["EnableCrossDesktopAccess"] = v
		}
	}
	if update {
		request["RegionId"] = client.RegionId
		action := "ModifyOfficeSiteCrossDesktopAccess"
		conn, err := client.NewGwsecdClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"DirectoryNotReady"}) {
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
		d.SetPartial("enable_cross_desktop_access")
	}

	update = false
	request = map[string]interface{}{
		"OfficeSiteId": d.Id(),
	}
	if d.HasChange("mfa_enabled") {
		update = true
		if v, ok := d.GetOkExists("mfa_enabled"); ok {
			request["MfaEnabled"] = v
		}
	}
	if update {
		request["RegionId"] = client.RegionId
		action := "ModifyOfficeSiteMfaEnabled"
		conn, err := client.NewGwsecdClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"DirectoryNotReady"}) {
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
		d.SetPartial("mfa_enabled")
	}

	update = false
	request = map[string]interface{}{
		"OfficeSiteId": d.Id(),
	}
	if d.HasChange("sso_enabled") {
		update = true
		request["RegionId"] = client.RegionId
		if v, ok := d.GetOkExists("sso_enabled"); ok {
			request["EnableSso"] = v
		}
	}
	if update {
		action := "SetOfficeSiteSsoStatus"
		conn, err := client.NewGwsecdClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"DirectoryNotReady"}) {
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
		d.SetPartial("sso_enabled")
	}

	update = false
	request = map[string]interface{}{
		"OfficeSiteId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("desktop_access_type") {
		update = true
	}
	if v, ok := d.GetOk("desktop_access_type"); ok {
		request["DesktopAccessType"] = v
	}
	if !d.IsNewResource() && d.HasChange("office_site_name") {
		update = true
	}
	if v, ok := d.GetOk("office_site_name"); ok {
		request["OfficeSiteName"] = v
	}
	if update {
		request["RegionId"] = client.RegionId
		action := "ModifyOfficeSiteAttribute"
		conn, err := client.NewGwsecdClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) || IsExpectedErrors(err, []string{"DirectoryNotReady"}) {
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
		d.SetPartial("desktop_access_type")
		d.SetPartial("office_site_name")
	}
	d.Partial(false)
	return resourceAlicloudEcdSimpleOfficeSiteRead(d, meta)
}
func resourceAlicloudEcdSimpleOfficeSiteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecdService := EcdService{client}
	object, err := ecdService.DescribeEcdSimpleOfficeSite(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ecd_simple_office_site ecdService.DescribeEcdSimpleOfficeSite Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if object["EnableInternetAccess"] == true {
		log.Printf("[WARN] Cannot destroy resource EcdSimpleOfficeSite. Terraform will remove this resource from the state file, however resources may remain.")
		return nil
	}
	action := "DeleteOfficeSites"
	var response map[string]interface{}
	conn, err := client.NewGwsecdClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"OfficeSiteId": []string{d.Id()},
	}

	request["RegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-09-30"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return nil
}

func convertDesktopAccessType(source interface{}) interface{} {
	switch source.(string) {
	case "ANY":
		return "Any"
	case "INTERNET":
		return "Internet"
	case "VPC":
		return "VPC"
	default:
		return source
	}
}
