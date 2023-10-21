package alicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudOosParameter() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOosParameterCreate,
		Read:   resourceAlicloudOosParameterRead,
		Update: resourceAlicloudOosParameterUpdate,
		Delete: resourceAlicloudOosParameterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"constraints": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.ValidateJsonString,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := compareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(1, 200),
			},
			"parameter_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.All(validation.StringDoesNotMatch(regexp.MustCompile(`(^ALIYUN.*)|(^ACS.*)|(^ALIBABA.*)|(^ALICLOUD.*)|(^OOS.*)`), "It cannot start with `ALIYUN`, `ACS`, `ALIBABA`, `ALICLOUD`, or `OOS`"), validation.StringMatch(regexp.MustCompile(`^[a-zA-Z0-9_/-]{2,180}`), "The name must be `2` to `180` characters in length, and can contain letters, digits, hyphens (-), forward slashes (/) and underscores (_).")),
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"tags": tagsSchema(),
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"String", "StringList"}, false),
			},
			"value": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 4096),
			},
		},
	}
}

func resourceAlicloudOosParameterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateParameter"
	request := make(map[string]interface{})
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("constraints"); ok {
		request["Constraints"] = v
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	if v, ok := d.GetOk("resource_group_id"); ok {
		request["ResourceGroupId"] = v
	}
	if v, ok := d.GetOk("tags"); ok {
		respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
		if err != nil {
			return WrapError(err)
		}
		request["Tags"] = respJson
	}
	request["Name"] = d.Get("parameter_name")
	request["Type"] = d.Get("type")
	request["Value"] = d.Get("value")
	request["ClientToken"] = buildClientToken("CreateParameter")
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &runtime)
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_oos_parameter", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["Name"]))

	return resourceAlicloudOosParameterRead(d, meta)
}
func resourceAlicloudOosParameterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	oosService := OosService{client}
	object, err := oosService.DescribeOosParameter(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_oos_parameter oosService.DescribeOosParameter Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("parameter_name", d.Id())
	d.Set("constraints", object["Constraints"])
	d.Set("description", object["Description"])
	d.Set("resource_group_id", object["ResourceGroupId"])
	d.Set("tags", tagsToMap(object["Tags"]))
	d.Set("type", object["Type"])
	d.Set("value", object["Value"])
	return nil
}
func resourceAlicloudOosParameterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	update := false
	request := map[string]interface{}{
		"Name": d.Id(),
	}
	if d.HasChange("value") {
		update = true
	}
	request["Value"] = d.Get("value")
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("resource_group_id") {
		update = true
		if v, ok := d.GetOk("resource_group_id"); ok {
			request["ResourceGroupId"] = v
		}
	}
	if d.HasChange("tags") {
		update = true
		if v, ok := d.GetOk("tags"); ok {
			respJson, err := convertMaptoJsonString(v.(map[string]interface{}))
			if err != nil {
				return WrapError(err)
			}
			request["Tags"] = respJson
		}
	}
	if update {
		action := "UpdateParameter"
		conn, err := client.NewOosClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudOosParameterRead(d, meta)
}
func resourceAlicloudOosParameterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteParameter"
	var response map[string]interface{}
	conn, err := client.NewOosClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"Name": d.Id(),
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-06-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
