package alicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlicloudCmsMonitorGroupInstances() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCmsMonitorGroupInstancesCreate,
		Read:   resourceAlicloudCmsMonitorGroupInstancesRead,
		Update: resourceAlicloudCmsMonitorGroupInstancesUpdate,
		Delete: resourceAlicloudCmsMonitorGroupInstancesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instances": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"category": {
							Type:     schema.TypeString,
							Required: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceAlicloudCmsMonitorGroupInstancesCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	action := "CreateMonitorGroupInstances"
	request := make(map[string]interface{})
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	request["GroupId"] = d.Get("group_id")
	Instances := make([]map[string]interface{}, len(d.Get("instances").(*schema.Set).List()))
	for i, InstancesValue := range d.Get("instances").(*schema.Set).List() {
		InstancesMap := InstancesValue.(map[string]interface{})
		Instances[i] = make(map[string]interface{})
		Instances[i]["Category"] = InstancesMap["category"]
		Instances[i]["InstanceId"] = InstancesMap["instance_id"]
		Instances[i]["InstanceName"] = InstancesMap["instance_name"]
		Instances[i]["RegionId"] = InstancesMap["region_id"]
	}
	request["Instances"] = Instances

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cms_monitor_group_instances", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprintf(`%v`, response["Code"]) != "200" {
		return WrapError(Error("CreateMonitorGroupInstances failed for " + response["Message"].(string)))
	}

	d.SetId(fmt.Sprint(request["GroupId"]))

	return resourceAlicloudCmsMonitorGroupInstancesRead(d, meta)
}
func resourceAlicloudCmsMonitorGroupInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	cmsService := CmsService{client}
	object, err := cmsService.DescribeCmsMonitorGroupInstances(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cms_monitor_group_instances cmsService.DescribeCmsMonitorGroupInstances Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("group_id", d.Id())
	resourceMap := make([]map[string]interface{}, 0)
	for _, v := range object {
		temp1 := map[string]interface{}{
			"category":      strings.ToLower(v["Category"].(string)),
			"instance_id":   v["InstanceId"],
			"instance_name": v["InstanceName"],
			"region_id":     v["RegionId"],
		}
		resourceMap = append(resourceMap, temp1)
	}
	if err := d.Set("instances", resourceMap); err != nil {
		return WrapError(err)
	}
	return nil
}
func resourceAlicloudCmsMonitorGroupInstancesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	var response map[string]interface{}
	if d.HasChange("instances") {
		request := map[string]interface{}{
			"GroupId": d.Id(),
		}
		Instances := make([]map[string]interface{}, len(d.Get("instances").(*schema.Set).List()))
		for i, InstancesValue := range d.Get("instances").(*schema.Set).List() {
			InstancesMap := InstancesValue.(map[string]interface{})
			Instances[i] = make(map[string]interface{})
			Instances[i]["Category"] = InstancesMap["category"]
			Instances[i]["InstanceId"] = InstancesMap["instance_id"]
			Instances[i]["InstanceName"] = InstancesMap["instance_name"]
			Instances[i]["RegionId"] = InstancesMap["region_id"]
		}
		request["Instances"] = Instances

		action := "ModifyMonitorGroupInstances"
		conn, err := client.NewCmsClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprintf(`%v`, response["Code"]) != "200" {
			return WrapError(Error("ModifyMonitorGroupInstances failed for " + response["Message"].(string)))
		}
	}
	return resourceAlicloudCmsMonitorGroupInstancesRead(d, meta)
}
func resourceAlicloudCmsMonitorGroupInstancesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	action := "DeleteMonitorGroupInstances"
	var response map[string]interface{}
	conn, err := client.NewCmsClient()
	if err != nil {
		return WrapError(err)
	}
	info := make(map[string]string)
	for _, InstancesValue := range d.Get("instances").(*schema.Set).List() {
		InstancesMap := InstancesValue.(map[string]interface{})
		if value, ok := info[InstancesMap["category"].(string)]; ok {
			info[InstancesMap["category"].(string)] = value + "," + InstancesMap["instance_id"].(string)
		} else {
			info[InstancesMap["category"].(string)] = InstancesMap["instance_id"].(string)
		}
	}
	for category, instanceIds := range info {
		request := map[string]interface{}{
			"GroupId":        d.Id(),
			"Category":       category,
			"InstanceIdList": instanceIds,
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-01-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if IsExpectedErrorCodes(fmt.Sprintf("%v", response["Code"]), []string{"ResourceNotFound"}) {
			return nil
		}
		if fmt.Sprintf("%v", response["Code"]) != "200" {
			return WrapError(Error("DeleteMonitorGroupInstances failed for " + response["Message"].(string)))
		}
	}
	return nil
}
