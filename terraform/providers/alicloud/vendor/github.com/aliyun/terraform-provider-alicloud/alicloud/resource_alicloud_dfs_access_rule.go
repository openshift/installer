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

func resourceAlicloudDfsAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudDfsAccessRuleCreate,
		Read:   resourceAlicloudDfsAccessRuleRead,
		Update: resourceAlicloudDfsAccessRuleUpdate,
		Delete: resourceAlicloudDfsAccessRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"access_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"access_rule_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_segment": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"priority": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"rw_access_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"RDONLY", "RDWR"}, false),
			},
		},
	}
}

func resourceAlicloudDfsAccessRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateAccessRule"
	request := make(map[string]interface{})
	conn, err := client.NewAlidfsClient()
	if err != nil {
		return WrapError(err)
	}
	request["AccessGroupId"] = d.Get("access_group_id")
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v
	}
	request["NetworkSegment"] = d.Get("network_segment")
	request["Priority"] = d.Get("priority")
	request["InputRegionId"] = client.RegionId
	request["RWAccessType"] = d.Get("rw_access_type")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_dfs_access_rule", action, AlibabaCloudSdkGoERROR)
	}

	d.SetId(fmt.Sprint(request["AccessGroupId"], ":", response["AccessRuleId"]))

	return resourceAlicloudDfsAccessRuleRead(d, meta)
}
func resourceAlicloudDfsAccessRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	dfsService := DfsService{client}
	object, err := dfsService.DescribeDfsAccessRule(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_dfs_access_rule dfsService.DescribeDfsAccessRule Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("access_group_id", parts[0])
	d.Set("access_rule_id", parts[1])
	d.Set("description", object["Description"])
	d.Set("network_segment", object["NetworkSegment"])
	if v, ok := object["Priority"]; ok && fmt.Sprint(v) != "0" {
		d.Set("priority", formatInt(v))
	}
	d.Set("rw_access_type", object["RWAccessType"])
	return nil
}
func resourceAlicloudDfsAccessRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	update := false
	request := map[string]interface{}{
		"AccessGroupId": parts[0],
		"AccessRuleId":  parts[1],
	}
	request["InputRegionId"] = client.RegionId
	if d.HasChange("description") {
		update = true
		if v, ok := d.GetOk("description"); ok {
			request["Description"] = v
		}
	}
	if d.HasChange("priority") {
		update = true
		request["Priority"] = d.Get("priority")
	}
	if d.HasChange("rw_access_type") {
		update = true
		request["RWAccessType"] = d.Get("rw_access_type")
	}
	if update {
		action := "ModifyAccessRule"
		conn, err := client.NewAlidfsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	return resourceAlicloudDfsAccessRuleRead(d, meta)
}
func resourceAlicloudDfsAccessRuleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteAccessRule"
	var response map[string]interface{}
	conn, err := client.NewAlidfsClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"AccessGroupId": parts[0],
		"AccessRuleId":  parts[1],
	}

	request["InputRegionId"] = client.RegionId
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-06-20"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrors(err, []string{"InvalidParameter.AccessRuleNotFound"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabaCloudSdkGoERROR)
	}
	return nil
}
