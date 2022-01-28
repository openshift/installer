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

func resourceAlicloudCloudStorageGatewayGatewayCacheDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCloudStorageGatewayGatewayCacheDiskCreate,
		Read:   resourceAlicloudCloudStorageGatewayGatewayCacheDiskRead,
		Update: resourceAlicloudCloudStorageGatewayGatewayCacheDiskUpdate,
		Delete: resourceAlicloudCloudStorageGatewayGatewayCacheDiskDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cache_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"cloud_efficiency", "cloud_ssd"}, false),
			},
			"cache_disk_size_in_gb": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(20, 32768),
			},
			"cache_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gateway_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"local_file_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCloudStorageGatewayGatewayCacheDiskCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	var response map[string]interface{}
	action := "CreateGatewayCacheDisk"
	request := make(map[string]interface{})
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("cache_disk_category"); ok {
		request["CacheDiskCategory"] = v
	}
	request["CacheDiskSizeInGB"] = d.Get("cache_disk_size_in_gb")
	request["GatewayId"] = d.Get("gateway_id")
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		return WrapErrorf(err, DefaultErrorMsg, "alicloud_cloud_storage_gateway_gateway_cache_disk", action, AlibabaCloudSdkGoERROR)
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}

	stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}

	task, err := sgwService.DescribeTasks(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]))
	if err != nil {
		return nil
	}
	object, err := sgwService.DescribeCloudStorageGatewayGatewayCacheDisk(fmt.Sprint(request["GatewayId"], ":", task["RelatedResourceId"], ":"))
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_gateway_cache_disk sgwService.DescribeCloudStorageGatewayGatewayCacheDisk Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.SetId(fmt.Sprint(request["GatewayId"], ":", task["RelatedResourceId"], ":", object["LocalFilePath"]))

	return resourceAlicloudCloudStorageGatewayGatewayCacheDiskRead(d, meta)
}
func resourceAlicloudCloudStorageGatewayGatewayCacheDiskRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	object, err := sgwService.DescribeCloudStorageGatewayGatewayCacheDisk(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_cloud_storage_gateway_gateway_cache_disk sgwService.DescribeCloudStorageGatewayGatewayCacheDisk Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("gateway_id", parts[0])
	d.Set("local_file_path", object["LocalFilePath"])
	d.Set("cache_disk_category", object["CacheType"])
	d.Set("cache_disk_size_in_gb", formatInt(object["SizeInGB"]))
	d.Set("cache_id", object["CacheId"])
	d.Set("status", formatInt(object["ExpireStatus"]))
	return nil
}
func resourceAlicloudCloudStorageGatewayGatewayCacheDiskUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	if d.HasChange("cache_disk_size_in_gb") {
		request := map[string]interface{}{
			"GatewayId":     parts[0],
			"LocalFilePath": parts[2],
			"NewSizeInGB":   d.Get("cache_disk_size_in_gb"),
		}
		action := "ExpandCacheDisk"
		conn, err := client.NewHcsSgwClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
		if fmt.Sprint(response["Success"]) == "false" {
			return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
		}
		stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return WrapErrorf(err, IdMsg, d.Id())
		}
	}

	return resourceAlicloudCloudStorageGatewayGatewayCacheDiskRead(d, meta)
}
func resourceAlicloudCloudStorageGatewayGatewayCacheDiskDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	sgwService := SgwService{client}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteGatewayCacheDisk"
	var response map[string]interface{}
	conn, err := client.NewHcsSgwClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"GatewayId":     parts[0],
		"CacheId":       parts[1],
		"LocalFilePath": parts[2],
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2018-05-11"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
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
	if fmt.Sprint(response["Success"]) == "false" {
		return WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	stateConf := BuildStateConf([]string{}, []string{"task.state.completed"}, d.Timeout(schema.TimeoutCreate), 1*time.Second, sgwService.CloudStorageGatewayTaskStateRefreshFunc(fmt.Sprint(request["GatewayId"]), fmt.Sprint(response["TaskId"]), []string{"task.state.failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return WrapErrorf(err, IdMsg, d.Id())
	}
	return nil
}
