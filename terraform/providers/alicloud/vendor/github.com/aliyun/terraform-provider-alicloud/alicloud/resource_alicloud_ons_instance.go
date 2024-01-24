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

func resourceAlicloudOnsInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudOnsInstanceCreate,
		Read:   resourceAlicloudOnsInstanceRead,
		Update: resourceAlicloudOnsInstanceUpdate,
		Delete: resourceAlicloudOnsInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(6 * time.Minute),
			Delete: schema.DefaultTimeout(6 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"instance_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  validation.StringLenBetween(3, 64),
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				Deprecated:    "Field 'name' has been deprecated from version 1.97.0. Use 'instance_name' instead.",
				ConflictsWith: []string{"instance_name"},
				ValidateFunc:  validation.StringLenBetween(3, 64),
			},
			"instance_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"release_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"remark": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceAlicloudOnsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "OnsInstanceCreate"
	request := make(map[string]interface{})
	conn, err := client.NewOnsClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("instance_name"); ok {
		request["InstanceName"] = v
	} else if v, ok := d.GetOk("name"); ok {
		request["InstanceName"] = v
	} else {
		return WrapError(Error(`[ERROR] Argument "name" or "instance_name" must be set one!`))
	}

	if v, ok := d.GetOk("remark"); ok {
		request["Remark"] = v
	}

	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if IsExpectedErrors(err, []string{"Throttling.User"}) || NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_ons_instance", action, AlibabaCloudSdkGoERROR)
	}
	responseData := response["Data"].(map[string]interface{})
	d.SetId(fmt.Sprint(responseData["InstanceId"]))

	return resourceAlicloudOnsInstanceUpdate(d, meta)
}
func resourceAlicloudOnsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}
	object, err := onsService.DescribeOnsInstance(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_ons_instance onsService.DescribeOnsInstance Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	d.Set("instance_name", object["InstanceName"])
	d.Set("name", object["InstanceName"])
	d.Set("instance_type", formatInt(object["InstanceType"]))
	d.Set("release_time", fmt.Sprint(formatInt(object["ReleaseTime"])))
	d.Set("remark", object["Remark"])
	d.Set("status", formatInt(object["InstanceStatus"]))
	d.Set("instance_status", formatInt(object["InstanceStatus"]))

	listTagResourcesObject, err := onsService.ListTagResources(d.Id(), "INSTANCE")
	if err != nil {
		return WrapError(err)
	}
	d.Set("tags", tagsToMap(listTagResourcesObject))
	return nil
}
func resourceAlicloudOnsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	onsService := OnsService{client}
	var response map[string]interface{}
	d.Partial(true)

	if d.HasChange("tags") {
		if err := onsService.SetResourceTags(d, "INSTANCE"); err != nil {
			return WrapError(err)
		}
		d.SetPartial("tags")
	}
	update := false
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("instance_name") {
		update = true
		request["InstanceName"] = d.Get("instance_name")
	}
	if !d.IsNewResource() && d.HasChange("name") {
		update = true
		request["InstanceName"] = d.Get("name")
	}
	if !d.IsNewResource() && d.HasChange("remark") {
		update = true
		request["Remark"] = d.Get("remark")
	}
	if update {
		action := "OnsInstanceUpdate"
		conn, err := client.NewOnsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			addDebug(action, response, request)
			return nil
		})
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
		}
		d.SetPartial("name")
		d.SetPartial("instance_name")
		d.SetPartial("remark")
	}
	d.Partial(false)
	return resourceAlicloudOnsInstanceRead(d, meta)
}
func resourceAlicloudOnsInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "OnsInstanceDelete"
	var response map[string]interface{}
	conn, err := client.NewOnsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"InstanceId": d.Id(),
	}

	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	wait := incrementalWait(3*time.Second, 10*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-02-14"), StringPointer("AK"), nil, request, &runtime)
		if err != nil {
			if IsExpectedErrors(err, []string{"INSTANCE_NOT_EMPTY", "Throttling.User"}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		addDebug(action, response, request)
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
