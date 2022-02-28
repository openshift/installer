package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudMhubApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudMhubAppCreate,
		Read:   resourceAlicloudMhubAppRead,
		Update: resourceAlicloudMhubAppUpdate,
		Delete: resourceAlicloudMhubAppDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bundle_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"package_name"},
			},
			"encoded_icon": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"industry_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"package_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"bundle_id"},
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Android", "iOS"}, false),
			},
		},
	}
}

func resourceAlicloudMhubAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateApp"
	request := make(map[string]interface{})
	conn, err := client.NewMhubClient()
	if err != nil {
		return WrapError(err)
	}
	request["Name"] = d.Get("app_name")
	request["ProductId"] = d.Get("product_id")

	if v, ok := d.GetOk("package_name"); ok {
		request["PackageName"] = v.(string)
	} else if v, ok := d.GetOk("bundle_id"); ok {
		request["BundleId"] = v.(string)
	} else {
		return WrapError(fmt.Errorf("either bundle_id or package_name must be set"))
	}

	request["Type"] = convertTypeRequest(d.Get("type").(string))
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-08-25"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_mhub_app", action, AlibabaCloudSdkGoERROR)
	}
	responseAppInfo := response["AppInfo"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseAppInfo["ProductId"], ":", responseAppInfo["AppKey"]))

	return resourceAlicloudMhubAppUpdate(d, meta)
}
func resourceAlicloudMhubAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	mhubService := MhubService{client}
	object, err := mhubService.DescribeMhubApp(d.Id())
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_mhub_app mhubService.DescribeMhubApp Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("app_name", object["Name"])
	d.Set("bundle_id", object["BundleId"])
	d.Set("product_id", parts[0])

	d.Set("type", convertTypeResponse(formatInt(object["Type"])))
	d.Set("encoded_icon", object["EncodedIcon"])
	d.Set("industry_id", fmt.Sprint(formatInt(object["IndustryId"])))
	d.Set("package_name", object["PackageName"])
	return nil
}
func resourceAlicloudMhubAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AppKey": parts[1],
	}
	if !d.IsNewResource() && d.HasChange("app_name") {
		update = true
		request["Name"] = d.Get("app_name")
	}
	if d.HasChange("encoded_icon") {
		update = true
		if v, ok := d.GetOk("encoded_icon"); ok {
			request["EncodedIcon"] = v
		}
	}
	if d.HasChange("industry_id") {
		update = true
		if v, ok := d.GetOk("industry_id"); ok {
			request["IndustryId"] = v
		}
	}
	if update {
		action := "ModifyApp"
		conn, err := client.NewMhubClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-08-25"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudMhubAppRead(d, meta)
}
func resourceAlicloudMhubAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteApp"
	var response map[string]interface{}
	conn, err := client.NewMhubClient()
	if err != nil {
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	request := map[string]interface{}{
		"AppKey": parts[1],
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-08-25"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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

func convertTypeResponse(source int) string {
	switch source {
	case 1:
		return "iOS"
	case 2:
		return "Android"
	}
	return ""
}
func convertTypeRequest(source string) string {
	switch source {
	case "iOS":
		return "1"
	case "Android":
		return "2"
	}
	return source
}
